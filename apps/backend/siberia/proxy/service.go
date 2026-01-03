package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/salacoste/siberia/siberia/config"
)

type Service struct {
	server *http.Server
	proxy  *goproxy.ProxyHttpServer
	config *config.AppConfig
}

func NewService(cfg *config.AppConfig) *Service {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	return &Service{
		proxy:  proxy,
		config: cfg,
	}
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
			if r.Method == "CONNECT" {
				w.WriteHeader(http.StatusProxyAuthRequired)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
			w.Write([]byte("Proxy Authentication Required"))
			return
		}

		// Strip Auth headers so they don't leak upstream (unless logic dictates otherwise)
		r.Header.Del("Proxy-Authorization")
		r.Header.Del("Authorization")
	}

	// 1. Check if z.ai is enabled and if request is "Reverse Proxy" style (relative path or targeting localhost)

	if s.config.ZaiEnabled {
		// A simple heuristic: if scheme is empty, it's likely a relative request to us acting as endpoint
		if r.URL.Scheme == "" || r.URL.Host == "" {
			s.handleZaiForward(w, r)
			return
		}
	}

	// 2. Otherwise default to standard goproxy (Forward Proxy)
	s.proxy.ServeHTTP(w, r)
}

func (s *Service) handleZaiForward(w http.ResponseWriter, r *http.Request) {
	targetBase := s.config.ZaiBaseURL
	// Ensure base has no trailing slash
	targetBase = strings.TrimRight(targetBase, "/")

	// Construct downstream URL
	destURL := fmt.Sprintf("%s%s", targetBase, r.URL.Path)
	if r.URL.RawQuery != "" {
		destURL += "?" + r.URL.RawQuery
	}

	// Create new request
	req, err := http.NewRequest(r.Method, destURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// defer r.Body.Close() - handled by server or client

	// Copy headers
	for k, vv := range r.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	// Inject Auth
	req.Header.Set("Authorization", "Bearer "+s.config.ZaiApiKey)

	// Remove Hop-by-hop headers
	req.Header.Del("Connection")
	req.Header.Del("Keep-Alive")
	req.Header.Del("Proxy-Authenticate")
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("Te")
	req.Header.Del("Trailers")
	req.Header.Del("Transfer-Encoding")
	req.Header.Del("Upgrade")

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
		http.Error(w, "Prody Error: "+err.Error(), http.StatusBadGateway)
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

func (s *Service) Start() error {
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
