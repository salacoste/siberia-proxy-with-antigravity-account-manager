package tools

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallWebSearchPrime(t *testing.T) {
	// Mock Z.ai Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search/prime" {
			t.Errorf("Expected path /search/prime, got %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("Expected Bearer test-key, got %s", auth)
		}

		// Return Mock Response
		resp := SearchResponse{
			Results: []SearchResult{
				{Title: "Test Result", Link: "https://example.com", Snippet: "A snippet"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	req := SearchRequest{Query: "test query"}
	results, err := CallWebSearchPrime("test-key", ts.URL, req)
	if err != nil {
		t.Fatalf("CallWebSearchPrime failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Test Result" {
		t.Errorf("Unexpected title: %s", results[0].Title)
	}
}
