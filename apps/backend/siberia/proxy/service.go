package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"time"

	"github.com/elazarl/goproxy"
	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	server *http.Server
	proxy  *goproxy.ProxyHttpServer
	config *config.AppConfig
	ctx    context.Context
}

type ProxyRequestEvent struct {
	Method   string `json:"method"`
	URL      string `json:"url"`
	Status   int    `json:"status"`
	Duration int64  `json:"duration_ms"` // milliseconds
	Time     string `json:"time"`
}

func NewService(cfg *config.AppConfig) *Service {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	return &Service{
		proxy:  proxy,
		config: cfg,
	}
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

// ServeHTTP wraps the goproxy handler to intercept "direct/reverse proxy" requests.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Wrap writer to capture status
	rw := &captureResponseWriter{ResponseWriter: w, statusCode: 200} // Default 200

	// 0. Security Check
	if s.config.AuthEnabled {
		authHeader := r.Header.Get("Proxy-Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("Authorization")
		}

		expected := "Bearer " + s.config.AuthToken
		// Simple check: header must equal exactly "Bearer <token>"
		if authHeader != expected {
			if r.Method == "CONNECT" {
				rw.WriteHeader(http.StatusProxyAuthRequired)
			} else {
				rw.WriteHeader(http.StatusUnauthorized)
			}
			rw.Write([]byte("Proxy Authentication Required"))
			s.emitEvent(r, rw.statusCode, start)
			return
		}

		// Strip Auth headers so they don't leak upstream
		r.Header.Del("Proxy-Authorization")
		r.Header.Del("Authorization")
	}

	// 1. Check if z.ai is enabled and if request is "Reverse Proxy" style
	if s.config.ZaiEnabled {
		if r.URL.Scheme == "" || r.URL.Host == "" {
			s.handleZaiForward(rw, r)
			s.emitEvent(r, rw.statusCode, start)
			return
		}
	}

	// 2. Otherwise default to standard goproxy (Forward Proxy)
	// We need to use 'rw' here too? goproxy might not be happy with a wrapped writer if it tries to hijack.
	// Hijacker interface check:
	// If the underlying w supports Hijack, our wrapper must too.
	// For simple HTTP/HTTPS proxying via CONNECT, goproxy hijacks the connection.
	// Creating a robust wrapper that supports Hijack is complex.
	// For CONNECT method, we might lose status code capture because goproxy takes over the raw net.Conn.
	// For now, let's pass original 'w' to goproxy to avoid breaking CONNECT/SSL,
	// and accept that we might not capture the *exact* status code for CONNECT requests,
	// or we can try to guess/emit "200" if it returns successfully?
	//
	// Improved approach: Only wrap for non-CONNECT methods or handle Hijack.
	// For MVP simplicity: Pass 'w' (original) to s.proxy.ServeHTTP.
	// We won't see the status code for standard proxy requests in this iteration.
	//
	// Wait, if we want to monitor traffic, we really want that status code.
	// Let's implement a minimal event emit *after* the call.
	// If we can't capture status easily without Hijack support, sending 0 or 200 is a fallback.
	// Let's stick to emitting "Request Started" info or minimal info.

	s.proxy.ServeHTTP(w, r)

	// For standard proxy, since we didn't wrap, we assume success or rely on goproxy logging.
	// We can still emit the event that a request happened.
	s.emitEvent(r, 0, start) // Status 0 indicates "Unknown/Tunnel"
}

func (s *Service) emitEvent(r *http.Request, status int, start time.Time) {
	if s.ctx != nil {
		duration := time.Since(start).Milliseconds()
		event := ProxyRequestEvent{
			Method:   r.Method,
			URL:      r.URL.String(),
			Status:   status,
			Duration: duration,
			Time:     start.Format("15:04:05"),
		}
		runtime.EventsEmit(s.ctx, "proxy:request", event)

		// Persistent Logging
		logger.LogAccess(logger.AccessEntry{
			Time:       event.Time,
			Timestamp:  start.Unix(),
			Method:     r.Method,
			URL:        r.URL.String(),
			Status:     status,
			DurationMs: duration,
			ClientIP:   r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			Size:       0, // We aren't capturing size yet in captureResponseWriter without more work
		})
	}
}

func (s *Service) handleZaiForward(w http.ResponseWriter, r *http.Request) {
	// ... (implementation of handleZaiForward needs to write to 'w' which is our capture wrapper)
	// The implementation below uses 'w' which is the argument.
	targetBase := s.config.ZaiBaseURL
	// ...
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
	io.Copy(w, resp.Body)
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
		Handler: s, // Use 's' (Service) as the handler, which wraps ServeHTTP
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
