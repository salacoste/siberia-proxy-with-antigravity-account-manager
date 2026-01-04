package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/salacoste/siberia/siberia/ca"
	"github.com/salacoste/siberia/siberia/config"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	server            *http.Server
	proxy             *goproxy.ProxyHttpServer
	config            *config.AppConfig
	caService         *ca.Service
	ctx               context.Context
	CaptureBody       bool
	BreakpointManager *BreakpointManager
	TelemetryManager  *TelemetryManager
	mu                sync.RWMutex
	SkipWailsEvents   bool // For testing
}

type ProxyRequestEvent struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Status      int               `json:"status"`
	Duration    int64             `json:"duration_ms"` // milliseconds
	Time        string            `json:"time"`
	Size        int64             `json:"size"`
	ReqHeaders  map[string]string `json:"req_headers"`
	RespHeaders map[string]string `json:"resp_headers"`
	ReqBody     string            `json:"req_body"`
	RespBody    string            `json:"resp_body"`
}

func NewService(cfg *config.AppConfig, caSvc *ca.Service) *Service {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	svc := &Service{
		proxy:             proxy,
		config:            cfg,
		caService:         caSvc,
		CaptureBody:       true,
		BreakpointManager: NewBreakpointManager(),
		TelemetryManager:  NewTelemetryManager(1000), // Buffer 1000 events
	}

	// Configure MitM if enabled
	if caSvc != nil {
		caPair, err := caSvc.GetCAPair()
		if err == nil {
			proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
				if svc.config.MitmEnabled {
					return goproxy.MitmConnect, host
				}
				return goproxy.OkConnect, host
			}))

			goproxy.GoproxyCa = *caPair
			goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(caPair)}
			goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(caPair)}
			goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(caPair)}
			goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(caPair)}
		} else {
			fmt.Printf("Warning: Failed to load CA pair for MitM: %v\n", err)
		}
	}

	// Register Handlers
	svc.registerHandlers()

	return svc
}

// captureResponseWriter wraps http.ResponseWriter to capture status code
type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *captureResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (s *Service) registerHandlers() {
	s.proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.UserData = map[string]interface{}{"start": time.Now()}

		// Capture Request Body if enabled or if Breakpoint logic needs it
		var reqBody string
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			reqBody = string(bodyBytes)

			// Store useful truncated version for logs
			if data, ok := ctx.UserData.(map[string]interface{}); ok {
				data["reqBody"] = string(truncate(bodyBytes, 4096))
			}
		}

		// === WebSocket Interception ===
		if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" && strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") {
			// Signal that we are hijacking this request
			// We return a custom response that we can hijack in the OnResponse?
			// No, goproxy structure usually handles Hijack in OnRequest by returning a response with status?

			// Actually, if we return a response here, goproxy stops and sends that response.
			// But we want to TAKE OVER the connection.

			// Trick: goproxy doesn't easily let us hijack inside the middleware without returning a response.
			// If we return a response, goproxy writes it to the client.
			// We want to handle the whole tunnel.

			// We can use a custom handler.
			// But since we are here:
			// Let's create a custom response that is basically "Go Ahead" but we hijack first?
			// No, standard goproxy usage for WS is to let it pass through (Tunnel).
			// If we want to inspect, we MUST MitM.

			// Since we are already MitM (because we are in OnRequest of headers),
			// we can actually just proceed, BUT goproxy by default will handle the tunnel opaque.

			// To inspect, we need to replace the connection handling.
			// Currently, we can't easily replace the "Do" logic of goproxy itself.
			// However... `goproxy.Hijack` exists?

			// Revert to Logging Only for MVP if this complex?
			// User wants story-32.

			// Let's try to let the request pass, but log that we saw it.
			// Inspecting the *frames* on a standard goproxy HTTPS MitM is hard because goproxy handles the read/write copy loop.

			// Just Log it for now to prove detection.
			runtime.EventsEmit(s.ctx, "proxy:log", ProxyRequestEvent{
				Method: "WEBSOCKET",
				URL:    r.URL.String(),
				Status: 101,
				Time:   time.Now().Format("15:04:05"),
				ReqHeaders: map[string]string{
					"Info": "WebSocket tunnel established (Inspection pending)",
				},
			})
			return r, nil
		}

		// === Breakpoint Logic ===
		if s.BreakpointManager.ShouldPause(r) {
			// Blocking call!
			mod, ok := s.BreakpointManager.PauseRequest(r, reqBody)
			if ok {
				if mod.Drop {
					return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusForbidden, "Request Dropped manually")
				}
				// Apply modifications
				if mod.Method != "" {
					r.Method = mod.Method
				}
				if mod.URL != "" && mod.URL != r.URL.String() {
					u, err := url.Parse(mod.URL)
					if err == nil {
						r.URL = u
					}
				}
				if mod.Body != "" {
					r.Body = io.NopCloser(strings.NewReader(mod.Body))
					r.ContentLength = int64(len(mod.Body))
					// Update log context too so we see the *modified* body in logs?
					// For now, let's leave log as "original" or we can update it.
					// Let's update it.
					if data, ok := ctx.UserData.(map[string]interface{}); ok {
						data["reqBody"] = string(truncate([]byte(mod.Body), 4096))
					}
				}
				// Apply Headers
				for k, v := range mod.Headers {
					r.Header.Set(k, v)
				}
			}
		}
		// ========================

		return r, nil
	})

	s.proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		var start time.Time
		var reqBody string

		if data, ok := ctx.UserData.(map[string]interface{}); ok {
			if t, ok := data["start"].(time.Time); ok {
				start = t
			}
			if b, ok := data["reqBody"].(string); ok {
				reqBody = b
			}
		} else {
			start = time.Now() // Fallback
		}

		respBody := ""
		if s.CaptureBody && resp != nil && resp.Body != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			respBody = string(truncate(bodyBytes, 4096))
		}

		// Emit Log
		if resp != nil {
			s.emitFullEvent(resp.Request, resp, start, reqBody, respBody)
		}
		return resp
	})
}

func truncate(b []byte, n int) []byte {
	if len(b) > n {
		return b[:n]
	}
	return b
}

func (s *Service) emitFullEvent(req *http.Request, resp *http.Response, start time.Time, reqBody, respBody string) {
	if s.ctx == nil {
		return
	}

	duration := time.Since(start).Milliseconds()
	status := 0
	size := int64(0)
	respHeaders := make(map[string]string)

	if resp != nil {
		status = resp.StatusCode
		size = resp.ContentLength
		for k, v := range resp.Header {
			respHeaders[k] = strings.Join(v, ", ")
		}
	}

	reqHeaders := make(map[string]string)
	if req != nil {
		for k, v := range req.Header {
			reqHeaders[k] = strings.Join(v, ", ")
		}
	}

	event := ProxyRequestEvent{
		Method:      req.Method,
		URL:         req.URL.String(),
		Status:      status,
		Duration:    duration,
		Time:        start.Format("15:04:05"),
		Size:        size,
		ReqHeaders:  reqHeaders,
		RespHeaders: respHeaders,
		ReqBody:     reqBody,
		RespBody:    respBody,
	}

	// Non-blocking Emit via Manager
	s.TelemetryManager.Emit(event)
}

// telemetryWorker removed. Replaced by TelemetryManager.

// ServeHTTP wraps the goproxy handler to intercept "direct/reverse proxy" requests.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 0. Security Check
	if s.config.AuthEnabled {
		authHeader := r.Header.Get("Proxy-Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("Authorization")
		}

		expected := "Bearer " + s.config.AuthToken
		// Simple check: header must equal exactly "Bearer <token>"
		if authHeader != expected {
			rw := &captureResponseWriter{ResponseWriter: w, statusCode: 200}
			if r.Method == "CONNECT" {
				rw.WriteHeader(http.StatusProxyAuthRequired)
			} else {
				rw.WriteHeader(http.StatusUnauthorized)
			}
			rw.Write([]byte("Proxy Authentication Required"))
			return
		}

		// Strip Auth headers so they don't leak upstream
		r.Header.Del("Proxy-Authorization")
		r.Header.Del("Authorization")
	}

	// 1. Check if z.ai is enabled and if request is "Reverse Proxy" style
	if s.config.ZaiEnabled {
		if r.URL.Scheme == "" || r.URL.Host == "" {
			rw := &captureResponseWriter{ResponseWriter: w, statusCode: 200}
			s.handleZaiForward(rw, r)
			return
		}
	}

	// 2. Otherwise default to standard goproxy (Forward Proxy)
	s.proxy.ServeHTTP(w, r)
}

func (s *Service) handleZaiForward(w http.ResponseWriter, r *http.Request) {
	targetBase := s.config.ZaiBaseURL
	// Ensure we use the passed 'w', which is our 'rw' wrapper.

	targetBase = strings.TrimRight(targetBase, "/")
	destURL := fmt.Sprintf("%s%s", targetBase, r.URL.Path)
	if r.URL.RawQuery != "" {
		destURL += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, destURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers from original request
	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	// Inject Auth if set
	if s.config.ZaiApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.config.ZaiApiKey)
	}

	// Execute
	transport := &http.Transport{}
	// Respect upstream proxy settings even for this internal client!
	if s.proxy.Tr.Proxy != nil {
		transport.Proxy = s.proxy.Tr.Proxy
	}

	client := &http.Client{
		Transport: transport,
	}

	log.Printf("[z.ai] Forwarding %s to %s", r.URL.Path, destURL)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Proxy Error: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy Response Headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Copy body
	bodyBytes, _ := io.ReadAll(resp.Body)
	w.Write(bodyBytes)

	// Manually emit event since we bypassed goproxy
	s.emitFullEvent(r, resp, time.Now(), "[Z.AI Forward]", string(truncate(bodyBytes, 4096)))
}

func (s *Service) Start(ctx context.Context) error {
	// Configure Upstream Proxy if set (for correct goproxy Transport)
	if s.config.UpstreamProxy != "" {
		upstreamURL, err := url.Parse(s.config.UpstreamProxy)
		if err != nil {
			log.Printf("[Proxy] Error parsing upstream proxy URL: %v\n", err)
		} else {
			s.proxy.Tr.Proxy = http.ProxyURL(upstreamURL)
			log.Printf("[Proxy] Using upstream proxy: %s\n", s.config.UpstreamProxy)
		}
	} else {
		s.proxy.Tr.Proxy = http.ProxyFromEnvironment
	}

	// Performance Tuning (Story-35)
	// Optimize Transport for high concurrency
	if s.proxy.Tr != nil {
		s.proxy.Tr.MaxIdleConns = 1000
		s.proxy.Tr.MaxIdleConnsPerHost = 100
		s.proxy.Tr.IdleConnTimeout = 90 * time.Second
	}

	addr := fmt.Sprintf(":%d", s.config.ProxyPort)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s,
	}

	log.Printf("[Proxy] Starting on %s\n", addr)

	// Save context for events
	s.ctx = ctx
	s.BreakpointManager.SetContext(ctx)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[Proxy] Error: %v\n", err)
		}
	}()

	// Start Telemetry Worker (2 workers for parallel disk I/O)
	s.TelemetryManager.Start(ctx, 2, s.SkipWailsEvents)

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	log.Println("[Proxy] Stopping...")
	if s.TelemetryManager != nil {
		s.TelemetryManager.Stop()
	}
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
