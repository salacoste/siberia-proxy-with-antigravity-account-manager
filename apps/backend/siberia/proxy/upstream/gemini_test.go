package upstream

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// Mock Account Provider
type MockAccountProvider struct {
	TokenCount int
}

func (m *MockAccountProvider) GetRotatingToken() (string, error) {
	m.TokenCount++
	return "mock-token", nil
}

func TestGeminiClient_Retry429(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls == 1 {
			w.WriteHeader(429) // Quota
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(mappers.GeminiResponse{})
	}))
	defer server.Close()

	mockAuth := &MockAccountProvider{}
	client := NewGeminiClient(mockAuth)
	client.endpoint = server.URL // Override for test

	// Speed up backoff
	// Since we hardcoded sleep in gemini.go, this test might be slow (100ms). That's fine.

	_, err := client.GenerateContent(context.Background(), "gemini-1.5-pro", &mappers.GeminiRequest{})
	if err != nil {
		t.Fatalf("Expected success after retry, got: %v", err)
	}

	if calls != 2 {
		t.Errorf("Expected 2 calls (1 retry), got %d", calls)
	}
	if mockAuth.TokenCount != 2 {
		t.Errorf("Expected 2 token fetches (rotation), got %d", mockAuth.TokenCount)
	}
}

func TestGeminiClient_Failover500(t *testing.T) {
	primaryCalls := 0
	fallbackCalls := 0

	// Primary: Always 500
	srvPrimary := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		primaryCalls++
		w.WriteHeader(500)
	}))
	defer srvPrimary.Close()

	// Fallback: 200 OK
	srvFallback := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fallbackCalls++
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(mappers.GeminiResponse{})
	}))
	defer srvFallback.Close()

	mockAuth := &MockAccountProvider{}
	client := NewGeminiClient(mockAuth)

	// Override URLs
	client.primaryURL = srvPrimary.URL
	client.fallbackURL = srvFallback.URL
	client.endpoint = srvPrimary.URL

	_, err := client.GenerateContent(context.Background(), "gemini-1.5-pro", &mappers.GeminiRequest{})
	if err != nil {
		t.Fatalf("Expected failover success, got error: %v", err)
	}

	if primaryCalls < 1 {
		t.Errorf("Expected primary to be called")
	}
	if fallbackCalls < 1 {
		t.Errorf("Expected fallback to be called")
	}
	if client.endpoint != srvFallback.URL {
		t.Errorf("Expected endpoint to stick to fallback")
	}
}
