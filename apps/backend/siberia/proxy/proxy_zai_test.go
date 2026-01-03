package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
)

func TestZaiReverseProxy(t *testing.T) {
	// 1. Setup Mock z.ai Upstream
	mockZai := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Rewrite
		if !strings.HasPrefix(r.URL.Path, "/v1/models") {
			t.Errorf("Expected path /v1/models, got %s", r.URL.Path)
		}
		// Verify Auth Injection
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key-123" {
			t.Errorf("Expected Authorization header 'Bearer test-key-123', got '%s'", auth)
		}
		fmt.Fprint(w, `{"data": [{"id": "gpt-4"}]}`)
	}))
	defer mockZai.Close()

	// 2. Configure Proxy Service
	cfg := &config.AppConfig{
		ProxyPort:  0,
		ZaiEnabled: true,
		ZaiBaseURL: mockZai.URL, // Point to our mock instead of real z.ai
		ZaiApiKey:  "test-key-123",
	}
	svc := NewService(cfg)

	// 3. Create Request to the Proxy (acting as Reverse Proxy)
	// Note: Request URL is relative (/v1/models), simulating a client hitting the proxy endpoint directly
	req := httptest.NewRequest("GET", "/v1/models", nil)
	w := httptest.NewRecorder()

	// 4. Exec ServeHTTP (which contains the routing logic)
	svc.ServeHTTP(w, req)

	// 5. Verify Response
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "gpt-4") {
		t.Errorf("Expected response to contain 'gpt-4', got %s", string(body))
	}
}
