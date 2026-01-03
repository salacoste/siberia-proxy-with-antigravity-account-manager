package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/config"
)

func TestProxyService(t *testing.T) {
	// 1. Setup Config
	cfg := &config.AppConfig{
		ProxyPort: 0, // 0 lets the OS choose a port, but our service binds to this.
		// Wait, our service implementation binds to cfg.ProxyPort directly.
		// Let's explicitly choose a non-standard port for testing
		Theme: "system",
	}
	cfg.ProxyPort = 19999

	// 2. Start Proxy
	svc := NewService(cfg)
	if err := svc.Start(); err != nil {
		t.Fatalf("Failed to start proxy: %v", err)
	}
	defer svc.Stop(nil)

	// Give it a moment to bind
	time.Sleep(100 * time.Millisecond)

	// 3. Create a test upstream server (the "Internet")
	expectedBody := "Hello from Internet"
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedBody)
	}))
	defer upstream.Close()

	// 4. Configure Client to use our Proxy
	proxyURL, _ := url.Parse(fmt.Sprintf("http://localhost:%d", cfg.ProxyPort))
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// 5. Make Request
	resp, err := client.Get(upstream.URL)
	if err != nil {
		t.Fatalf("Failed to make request through proxy: %v", err)
	}
	defer resp.Body.Close()

	// 6. Verify Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if string(body) != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, string(body))
	}
}
