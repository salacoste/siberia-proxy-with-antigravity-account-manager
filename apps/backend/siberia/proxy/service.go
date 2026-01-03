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
	"github.com/salacoste/siberia/siberia/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	server      *http.Server
	proxy       *goproxy.ProxyHttpServer
	config      *config.AppConfig
	caService   *ca.Service
	ctx         context.Context
	CaptureBody bool
	mu          sync.RWMutex
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
		proxy:       proxy,
		config:      cfg,
		caService:   caSvc,
		CaptureBody: true,
	}

	// Configure MitM if enabled (or always configure it, but toggle in handler)
	// We need to load the CA pair to sign certs.
	if caSvc != nil {
		caPair, err := caSvc.GetCAPair()
		if err == nil {
			proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
				// Toggle MitM based on config
				if svc.config.MitmEnabled {
					return goproxy.MitmConnect, host
				}
				return goproxy.OkConnect, host
			}))

			// Set the MitM config
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

		// Capture Request Body if enabled
		if s.CaptureBody {
			// We need to read and restore the body
			if r.Body != nil {
				bodyBytes, _ := io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Store in context for Response hook
				if data, ok := ctx.UserData.(map[string]interface{}); ok {
					data["reqBody"] = string(truncate(bodyBytes, 4096)) // Cap at 4KB
				}
			}
		}
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

	runtime.EventsEmit(s.ctx, "proxy:log", event)

	// Also log to disk (minimal version)
	logger.LogAccess(logger.AccessEntry{
		Time:       event.Time,
		Timestamp:  start.Unix(),
		Method:     event.Method,
		URL:        event.URL,
		Status:     status,
		DurationMs: duration,
		Size:       size,
	})
}

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

	addr := fmt.Sprintf(":%d", s.config.ProxyPort)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s,
	}

	log.Printf("[Proxy] Starting on %s\n", addr)

	// Save context for events
	s.ctx = ctx

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[Proxy] Error: %v\n", err)
		}
	}()

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	log.Println("[Proxy] Stopping...")
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
