package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/google/uuid"
	"github.com/salacoste/siberia/siberia/accounts"
	"github.com/salacoste/siberia/siberia/analytics"
	"github.com/salacoste/siberia/siberia/ca"
	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/proxy/handlers/claude"
	"github.com/salacoste/siberia/siberia/proxy/handlers/gemini"
	"github.com/salacoste/siberia/siberia/proxy/handlers/legacy"
	"github.com/salacoste/siberia/siberia/proxy/handlers/openai"
	"github.com/salacoste/siberia/siberia/proxy/middleware"
	"github.com/salacoste/siberia/siberia/proxy/providers"
	"github.com/salacoste/siberia/siberia/proxy/scripting"
	"github.com/salacoste/siberia/siberia/proxy/upstream"

	"github.com/salacoste/siberia/siberia/types"
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
	openaiHandler     *openai.Handler
	legacyHandler     http.Handler

	claudeHandler http.Handler
	geminiHandler http.Handler

	// Providers (Story-70)
	zaiProvider *providers.ZaiProvider

	MapLocal      *middleware.MapLocalMiddleware
	ScriptEngine  *scripting.ScriptEngine
	sampleCounter uint64
}

func NewService(cfg *config.AppConfig, caSvc *ca.Service, analyticsEngine *analytics.AnalyticsEngine, accSvc *accounts.Service) *Service {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	// Initialize Upstream Client
	geminiClient := upstream.NewGeminiClient(accSvc)

	// Initialize Handlers
	oaHandler := openai.NewHandler(geminiClient)
	clHandler := claude.NewHandler(geminiClient)
	gmHandler := gemini.NewHandler(geminiClient)
	lgHandler := legacy.NewHandler(geminiClient)

	svc := &Service{
		proxy:             proxy,
		config:            cfg,
		caService:         caSvc,
		CaptureBody:       true,
		BreakpointManager: NewBreakpointManager(),

		TelemetryManager: NewTelemetryManager(1000, analyticsEngine), // Buffer 1000 events
		openaiHandler:    oaHandler,

		legacyHandler: lgHandler,

		claudeHandler: clHandler,
		geminiHandler: gmHandler,
		zaiProvider:   providers.NewZaiProvider(cfg),
		MapLocal:      middleware.NewMapLocalMiddleware(),
		ScriptEngine:  scripting.NewScriptEngine(),
	}

	// Configure MitM if enabled
	if caSvc != nil {
		caPair, err := caSvc.GetCAPair()
		if err == nil {
			proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
				if svc.config.MitmEnabled {
					// Use ConnectHijack to fully control the Tunnel
					return &goproxy.ConnectAction{
						Action: goproxy.ConnectHijack,
						Hijack: func(req *http.Request, clientConn net.Conn, ctx *goproxy.ProxyCtx) {
							// 1. Send OK to Client to establish tunnel
							clientConn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

							// 2. Wrap in TLS Server
							// We need to generate a cert for the host?
							// goproxy's logic normally handles this dynamic cert generation.
							// We can use goproxy.TLSConfigFromCA to get a config that "Signs" on demand?
							// correct. TLSConfigFromCA returns a func that generates certs.
							conf, err := goproxy.TLSConfigFromCA(caPair)(host, ctx)
							if err != nil {
								log.Printf("[MitM] Failed to generate cert for %s: %v", host, err)
								clientConn.Close()
								return
							}

							tlsConn := tls.Server(clientConn, conf)

							// 3. Serve via OUR Service handler (s)
							// This routes the decrypted HTTP2/1.1 traffic back to s.ServeHTTP
							// where we can intercept WebSockets and Log everything.
							// We create a custom listener to serve this single connection?
							// http.Serve serves on a listener.
							// We can use a "SingleConnectionListener".

							// Or simpler:
							// http.Server{Handler: svc}.Serve(SingleConnListener{tlsConn})

							// Define a throwaway listener
							l := &singleConnListener{conn: tlsConn}
							server := &http.Server{Handler: svc}
							// Serve blocks until connection closes
							server.Serve(l)
						},
					}, host
				}
				return goproxy.OkConnect, host
			}))

			// Remove erroneous assignment
			// proxy.MitmHandler = svc

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
		connID := uuid.New().String()
		ctx.UserData = map[string]interface{}{
			"start":  time.Now(),
			"connID": connID,
		}

		// Sampling Logic
		shouldLog := true
		if s.config.AccessLogSampleRate > 1 {
			c := atomic.AddUint64(&s.sampleCounter, 1)
			if c%uint64(s.config.AccessLogSampleRate) != 0 {
				shouldLog = false
			}
		}
		if data, ok := ctx.UserData.(map[string]interface{}); ok {
			data["shouldLog"] = shouldLog
		}

		// Capture Request Body if enabled or if Breakpoint logic needs it
		var reqBody string
		// Optimization: Skip body read if not logging AND no breakpoints
		needsInspection := (s.CaptureBody && shouldLog) || s.BreakpointManager.HasActiveBreakpoints() || (s.ScriptEngine != nil && s.ScriptEngine.Active)

		if r.Body != nil && needsInspection {
			var err error
			reqBody, err = s.peekAndRestoreBody(r)
			if err != nil {
				log.Printf("[Proxy] Error capturing req body: %v", err)
			}
			if data, ok := ctx.UserData.(map[string]interface{}); ok {
				data["reqBody"] = reqBody
			}
		}

		// === Scripting Logic (OnRequest) ===
		if s.ScriptEngine != nil && s.ScriptEngine.Active {
			modReq, err := s.ScriptEngine.RunOnRequest(r)
			if err != nil {
				log.Printf("[Scripting] Error in onRequest: %v", err)
			} else {
				r = modReq
			}
		}
		// ===================================

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

		shouldLog := true
		if data, ok := ctx.UserData.(map[string]interface{}); ok {
			if v, ok := data["shouldLog"].(bool); ok {
				shouldLog = v
			}
		}

		// Optimization: Skip body read if not logging AND no scripting
		needsInspection := (s.CaptureBody && shouldLog) || (s.ScriptEngine != nil && s.ScriptEngine.Active)

		respBody := ""
		if resp != nil && resp.Body != nil && needsInspection {
			var err error
			respBody, err = s.peekAndRestoreResponse(resp)
			if err != nil {
				log.Printf("[Proxy] Error capturing resp body: %v", err)
			}
		}

		// === Scripting Logic (OnResponse) ===
		if s.ScriptEngine != nil && s.ScriptEngine.Active && resp != nil {
			modResp, err := s.ScriptEngine.RunOnResponse(resp)
			if err != nil {
				log.Printf("[Scripting] Error in onResponse: %v", err)
			} else {
				resp = modResp
			}
		}
		// ===================================

		// Emit Log
		if resp != nil {

			connID := ""
			if data, ok := ctx.UserData.(map[string]interface{}); ok {
				if id, ok := data["connID"].(string); ok {
					connID = id
				}
			}
			if shouldLog {
				s.emitFullEvent(resp.Request, resp, start, reqBody, respBody, connID)
			}
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

func (s *Service) emitFullEvent(req *http.Request, resp *http.Response, start time.Time, reqBody, respBody, connID string) {
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

	event := types.ProxyRequestEvent{
		Method:       req.Method,
		URL:          req.URL.String(),
		Status:       status,
		Duration:     duration,
		Time:         start.Format("15:04:05"),
		Size:         size,
		ReqHeaders:   reqHeaders,
		RespHeaders:  respHeaders,
		ReqBody:      reqBody,
		RespBody:     respBody,
		ConnectionID: connID,
	}

	// Non-blocking Emit via Manager
	s.TelemetryManager.Emit(event)

	// Concise Access Log (Story-67)
	if s.config.AccessLogEnabled {
		log.Printf("[Access] %s %s %d %dms", req.Method, req.URL.Path, status, duration)
	}
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

	// 0.5 Map Local Middleware
	// Check if this request should be served from a local file
	if s.MapLocal != nil {
		req, resp := s.MapLocal.HandleRequest(r, nil) // ctx can be nil for now or we plumb it
		if resp != nil {
			// Middleware served the response
			for k, vv := range resp.Header {
				for _, v := range vv {
					w.Header().Add(k, v)
				}
			}
			w.WriteHeader(resp.StatusCode)
			if resp.Body != nil {
				io.Copy(w, resp.Body)
				resp.Body.Close()
			}
			// Emit Log for Local Map
			s.emitFullEvent(r, resp, time.Now(), "[Map Local]", "[Mapped Content]", uuid.New().String())
			return
		}
		// If req modified, use it
		r = req
	}

	// 1. Handler Specific Routes (Story-43, Story-63, Story-70)
	// OpenAI Chat Completions
	if r.URL.Path == "/v1/chat/completions" && r.Method == "POST" {
		s.openaiHandler.ServeHTTP(w, r)
		return
	}

	// Legacy Codex Completions
	if r.URL.Path == "/v1/completions" && r.Method == "POST" {
		s.legacyHandler.ServeHTTP(w, r)
		return
	}

	// OpenAI Image Generations (Story-63)
	if r.URL.Path == "/v1/images/generations" && r.Method == "POST" {
		s.openaiHandler.ImageGenerations(w, r)
		return
	}

	// Claude Messages
	if r.URL.Path == "/v1/messages" && r.Method == "POST" {
		// Story-70: Z.ai Dispatch Check
		shouldUseZai, restoredReq := s.checkZaiDispatch(r)
		if shouldUseZai {
			s.zaiProvider.ForwardAnthropicJSON(w, restoredReq)
			return
		}
		// Use restored request just in case checkZaiDispatch consumed it
		s.claudeHandler.ServeHTTP(w, restoredReq)
		return
	}

	// Native Gemini
	if strings.Contains(r.URL.Path, "/v1beta/models/") && strings.Contains(r.URL.Path, ":") {
		s.geminiHandler.ServeHTTP(w, r)
		return
	}

	// 2. Check if z.ai is enabled and if request is "Reverse Proxy" style (Catch-All for Z.ai)
	if s.config.ZaiEnabled {
		if r.URL.Scheme == "" || r.URL.Host == "" {
			rw := &captureResponseWriter{ResponseWriter: w, statusCode: 200}
			s.handleZaiForward(rw, r)
			return
		}
	}

	// 3. Otherwise default to standard goproxy (Forward Proxy)
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

	if s.ctx != nil {
		runtime.LogDebugf(s.ctx, "[z.ai] Forwarding %s to %s", r.URL.Path, destURL)
	}

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
	s.emitFullEvent(r, resp, time.Now(), "[Z.AI Forward]", string(truncate(bodyBytes, 4096)), "zai-forward")
}

func (s *Service) peekAndRestoreBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}

	const limit = 4096
	bufPtr := GetBuffer()
	buf := *bufPtr
	defer PutBuffer(bufPtr)

	readBuf := buf
	if len(readBuf) > limit {
		readBuf = readBuf[:limit]
	}

	n, err := io.ReadFull(r.Body, readBuf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}

	chunk := make([]byte, n)
	copy(chunk, buf[:n])

	r.Body = &readCloserWrapper{
		Reader: io.MultiReader(bytes.NewReader(chunk), r.Body),
		Closer: r.Body,
	}

	return string(chunk), nil
}

func (s *Service) peekAndRestoreResponse(resp *http.Response) (string, error) {
	if resp.Body == nil {
		return "", nil
	}

	const limit = 4096
	bufPtr := GetBuffer()
	buf := *bufPtr
	defer PutBuffer(bufPtr)

	readBuf := buf
	if len(readBuf) > limit {
		readBuf = readBuf[:limit]
	}

	n, err := io.ReadFull(resp.Body, readBuf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}

	chunk := make([]byte, n)
	copy(chunk, buf[:n])

	resp.Body = &readCloserWrapper{
		Reader: io.MultiReader(bytes.NewReader(chunk), resp.Body),
		Closer: resp.Body,
	}

	return string(chunk), nil
}

type readCloserWrapper struct {
	io.Reader
	io.Closer
}

func (s *Service) Start(ctx context.Context) error {
	// Configure Upstream Proxy if set (for correct goproxy Transport)
	if s.config.UpstreamProxy != "" {
		upstreamURL, err := url.Parse(s.config.UpstreamProxy)
		if err != nil {
			log.Printf("[Proxy] Error parsing upstream proxy URL: %v\n", err)
		} else {
			s.proxy.Tr.Proxy = http.ProxyURL(upstreamURL)
			// log.Printf("[Proxy] Using upstream proxy: %s\n", s.config.UpstreamProxy)
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

// Helper for single connection serving
type singleConnListener struct {
	conn net.Conn
	once sync.Once
	done chan struct{}
}

func (l *singleConnListener) Accept() (net.Conn, error) {
	// We only return the connection once
	var c net.Conn
	l.once.Do(func() {
		c = l.conn
	})
	if c != nil {
		return c, nil
	}
	// Block forever after
	select {}
}

func (l *singleConnListener) Close() error {
	return l.conn.Close()
}

func (l *singleConnListener) Addr() net.Addr {
	return l.conn.LocalAddr()
}

func (s *Service) checkZaiDispatch(r *http.Request) (bool, *http.Request) {
	cfg := s.config
	if cfg.ZaiDispatchMode == "off" || cfg.ZaiDispatchMode == "" {
		return false, r
	}

	// 1. Peek Body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		// If error, assume standard handler can handle or fail
		return false, r
	}
	r.Body.Close()
	// Restore body immediately for the returned request
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 2. Parse Minimal to get Model
	var partial struct {
		Model string `json:"model"`
	}
	if err := json.Unmarshal(bodyBytes, &partial); err != nil {
		return false, r
	}

	// 3. Logic
	model := partial.Model

	// Explicit Z.ai prefix
	if strings.HasPrefix(model, "zai:") {
		return true, r
	}

	// Exclusive Mode
	if cfg.ZaiDispatchMode == "exclusive" {
		return true, r
	}

	// Pooled/Mapped check
	if _, ok := cfg.ZaiModelMapping[model]; ok {
		return true, r
	}

	return false, r
}
