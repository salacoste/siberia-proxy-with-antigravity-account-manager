package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
)

func benchmarkProxy(b *testing.B, sampleRate int) {
	// Setup
	cfg := &config.AppConfig{
		ProxyPort:           0,
		AccessLogSampleRate: sampleRate,
		AccessLogEnabled:    true,
	}

	// Create Service with minimal deps
	svc := NewService(cfg, nil, nil, nil)
	// Mock TelemetryManager to avoid DB/File IO
	svc.TelemetryManager = NewTelemetryManager(10000, nil)
	svc.SkipWailsEvents = true

	// Upstream Server (The "Internet")
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return 10KB of data
		w.Header().Set("Content-Type", "text/plain")
		w.Write(bytes.Repeat([]byte("A"), 1024*10))
	}))
	defer upstream.Close()

	// Prepare Request
	// For forward proxy, URL must be absolute
	req, _ := http.NewRequest("GET", upstream.URL, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		svc.ServeHTTP(w, req)
	}
}

func BenchmarkProxy_NoSampling(b *testing.B) {
	benchmarkProxy(b, 1)
}

func BenchmarkProxy_Sampling10(b *testing.B) {
	benchmarkProxy(b, 10)
}

func BenchmarkProxy_Sampling100(b *testing.B) {
	benchmarkProxy(b, 100)
}
