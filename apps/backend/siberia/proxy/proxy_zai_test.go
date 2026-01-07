package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/proxy/providers"
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
	svc := NewService(cfg, nil, nil, nil)

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

// Story-70: Test Z.ai Provider Dispatch Logic
func TestZaiAnthropicProvider(t *testing.T) {
	// 1. Setup Mock Z.ai API (Chat Completions)
	mockZai := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log.Printf("Mock Z.ai received: %s %s", r.Method, r.URL.Path)

		// Verify Path is appended correctly
		if r.URL.Path != "/v1/messages" {
			t.Errorf("Expected upstream path /v1/messages, got %s", r.URL.Path)
		}

		// Verify Auth header (x-api-key)
		key := r.Header.Get("x-api-key")
		if key != "zai-key-abc" {
			t.Errorf("Expected x-api-key 'zai-key-abc', got '%s'", key)
		}

		// Verify Body (Model Mapping)
		body, _ := io.ReadAll(r.Body)
		var payload providers.AnthropicRequest
		json.Unmarshal(body, &payload)

		if payload.Model != "metis" {
			t.Errorf("Expected mapped model 'metis', got '%s'", payload.Model)
		}

		w.WriteHeader(200)
		fmt.Fprint(w, `{"id":"msg_123","content":[{"text":"Hello from Z.ai"}]}`)
	}))
	defer mockZai.Close()

	// 2. Setup Config
	cfg := &config.AppConfig{
		ZaiEnabled:      true,
		ZaiBaseURL:      mockZai.URL + "/v1", // Mock URL
		ZaiApiKey:       "zai-key-abc",
		ZaiDispatchMode: "pooled",
		ZaiModelMapping: map[string]string{
			"claude-3-opus": "metis",
		},
	}

	svc := NewService(cfg, nil, nil, nil)

	// 3. Create Request (Standard Anthropic)
	reqBody := `{"model": "claude-3-opus", "messages": [{"role": "user", "content": "hi"}]}`
	req := httptest.NewRequest("POST", "/v1/messages", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	// 4. Exec
	svc.ServeHTTP(w, req)

	// 5. Verify
	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(respBody), "Hello from Z.ai") {
		t.Errorf("Expected response from Z.ai mock, got: %s", string(respBody))
	}
}

// Helper to test "OFF" mode
func TestZaiAnthropicProvider_Off(t *testing.T) {
	cfg := &config.AppConfig{
		ZaiEnabled:      true,
		ZaiDispatchMode: "off",
	}
	svc := NewService(cfg, nil, nil, nil)

	// Use a fake claude handler to verify fallback
	svc.claudeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("Fallback to Claude"))
	})

	reqBody := `{"model": "claude-3-opus", "messages": []}`
	req := httptest.NewRequest("POST", "/v1/messages", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	svc.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected fallback (201), got %d", w.Code)
	}
}
