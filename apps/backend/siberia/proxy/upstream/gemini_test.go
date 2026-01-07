package upstream

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/session"
)

// Mock Account Provider
type MockAccountProvider struct {
	TokenCount      int
	LastFingerprint string
}

func (m *MockAccountProvider) GetRotatingToken(fingerprint string) (string, string, error) {
	m.TokenCount++
	m.LastFingerprint = fingerprint
	return "mock-token", "mock-identity", nil
}

func (m *MockAccountProvider) GetSchedulingMode() string {
	return "PerformanceFirst"
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

	// Since we hardcoded sleep in gemini.go, this test might be slow (100ms). That's fine.

	_, _, err := client.GenerateContent(context.Background(), "gemini-1.5-pro", &mappers.GeminiRequest{})
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

	_, _, err := client.GenerateContent(context.Background(), "gemini-1.5-pro", &mappers.GeminiRequest{})
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

func TestGeminiClient_StickySession(t *testing.T) {
	mockAuth := &MockAccountProvider{}
	client := NewGeminiClient(mockAuth)

	// Mock Server for success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(mappers.GeminiResponse{})
	}))
	client.endpoint = server.URL
	defer server.Close()

	// 1. Request with content
	req := &mappers.GeminiRequest{
		Contents: []mappers.GeminiContent{
			{Role: "user", Parts: []mappers.GeminiPart{{Text: "Hello Sticky"}}},
		},
	}
	_, _, _ = client.GenerateContent(context.Background(), "gemini-pro", req)

	// Verify Fingerprint was generated and passed
	// Hash of "gemini-pro:Hello Sticky"
	// We trust session.GenerateFingerprint works (unit tested).
	// We just check that LastFingerprint is NOT empty.
	if mockAuth.LastFingerprint == "" {
		t.Error("Expected LastFingerprint to be set, got empty")
	}

	firstFingerprint := mockAuth.LastFingerprint

	// 2. Exact same request
	_, _, _ = client.GenerateContent(context.Background(), "gemini-pro", req)
	if mockAuth.LastFingerprint != firstFingerprint {
		t.Errorf("Expected same fingerprint for same content, got %s vs %s", mockAuth.LastFingerprint, firstFingerprint)
	}

	// 3. Different content
	reqDiff := &mappers.GeminiRequest{
		Contents: []mappers.GeminiContent{
			{Role: "user", Parts: []mappers.GeminiPart{{Text: "Hello Different"}}},
		},
	}
	_, _, _ = client.GenerateContent(context.Background(), "gemini-pro", reqDiff)
	if mockAuth.LastFingerprint == firstFingerprint {
		t.Error("Expected different fingerprint for different content")
	}
}

func TestGeminiClient_SessionHeader_Override(t *testing.T) {
	mockAuth := &MockAccountProvider{}
	client := NewGeminiClient(mockAuth)

	// Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(mappers.GeminiResponse{})
	}))
	client.endpoint = server.URL
	defer server.Close()

	// 1. Explicit Session ID
	sessionID := "my-explicit-session-123"
	ctx := context.WithValue(context.Background(), session.SessionIDKey, sessionID)

	req := &mappers.GeminiRequest{
		Contents: []mappers.GeminiContent{
			{Role: "user", Parts: []mappers.GeminiPart{{Text: "Prompt A"}}},
		},
	}

	_, _, _ = client.GenerateContent(ctx, "gemini-pro", req)

	if mockAuth.LastFingerprint != sessionID {
		t.Errorf("Expected fingerprint to be %s, got %s", sessionID, mockAuth.LastFingerprint)
	}

	// 2. Exact same session ID with DIFFERENT prompt
	reqDiff := &mappers.GeminiRequest{
		Contents: []mappers.GeminiContent{
			{Role: "user", Parts: []mappers.GeminiPart{{Text: "Prompt B"}}},
		},
	}
	_, _, _ = client.GenerateContent(ctx, "gemini-pro", reqDiff)

	if mockAuth.LastFingerprint != sessionID {
		t.Errorf("Expected fingerprint to persist as %s despite content change, got %s", sessionID, mockAuth.LastFingerprint)
	}
}
