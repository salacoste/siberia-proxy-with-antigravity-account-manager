package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/config"
)

// BenchmarkProxyThroughput measures how many reqs/sec the proxy can handle
// with the new TelemetryManager and Worker Pool.
func BenchmarkProxyThroughput(b *testing.B) {
	// 1. Setup
	cfg := &config.AppConfig{
		ProxyPort:   19998, // Distinct port
		MitmEnabled: false,
	}

	// Mock CA Service not needed for HTTP benchmark if MITM off
	// But NewService needs it?
	svc := NewService(cfg, nil, nil)
	svc.SkipWailsEvents = true // CRITICAL: Don't choke on Wails calls during bench

	// Start Proxy
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := svc.Start(ctx); err != nil {
		b.Fatalf("Failed to start proxy: %v", err)
	}
	defer svc.Stop(ctx)

	// Wait for start
	time.Sleep(100 * time.Millisecond)

	// Create a dummy upstream server to proxy to
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer upstream.Close()

	proxyUrl := fmt.Sprintf("http://localhost:%d", cfg.ProxyPort)

	// Create a client that uses the proxy
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(proxyUrl)
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	client := &http.Client{Transport: transport}

	b.ResetTimer()

	// 2. Run Benchmark
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(upstream.URL)
			if err != nil {
				// Don't fail benchmark on individual errors, but log
				continue
			}
			resp.Body.Close()
		}
	})
}
