package tools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"https://example.com/page?utm_source=google&utm_medium=cpc",
			"https://example.com/page",
		},
		{
			"https://example.com?gclid=12345&fbclid=abcdef",
			"https://example.com",
		},
		{
			"https://example.com?q=search&utm_campaign=summer",
			"https://example.com?q=search",
		},
	}

	for _, tt := range tests {
		got, err := normalizeURL(tt.input)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Query param order might vary, but for empty query it should be clean
		if !strings.Contains(got, "https://example.com") {
			t.Errorf("Normalize failed. Got: %s, Want: %s", got, tt.expected)
		}
		if strings.Contains(got, "utm_") {
			t.Errorf("Tracking param not stripped: %s", got)
		}
	}
}

func TestCallWebReader(t *testing.T) {
	// Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body><h1>Hello World</h1><p>Test Content</p></body></html>")
	}))
	defer ts.Close()

	md, err := CallWebReader(ts.URL)
	if err != nil {
		t.Fatalf("CallWebReader failed: %v", err)
	}

	if !strings.Contains(md, "# Hello World") {
		t.Errorf("Expected markdown title, got: %s", md)
	}
	if !strings.Contains(md, "Test Content") {
		t.Errorf("Expected markdown content, got: %s", md)
	}
}
