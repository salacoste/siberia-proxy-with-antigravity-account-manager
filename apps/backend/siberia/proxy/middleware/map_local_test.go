package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMapLocalMiddleware(t *testing.T) {
	// Setup temporary file
	content := "Hello, Map Local!"
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	mw := NewMapLocalMiddleware()

	// 1. Add Rule
	err := mw.AddRule(MapLocalRule{
		ID:        "rule1",
		Enabled:   true,
		UrlRegex:  `.*\/example\.com\/api\/test`,
		LocalPath: tmpFile,
	})
	if err != nil {
		t.Fatalf("failed to add rule: %v", err)
	}

	// 2. Test Match
	req := httptest.NewRequest("GET", "https://example.com/api/test", nil)
	_, resp := mw.HandleRequest(req, nil)

	if resp == nil {
		t.Fatal("expected match, got nil response")
	}

	// Verify status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// Verify header
	if val := resp.Header.Get("X-Siberia-Map-Local"); val != "rule1" {
		t.Errorf("expected header X-Siberia-Map-Local=rule1, got %s", val)
	}

	// Verify Content
	bodyBytes, _ := io.ReadAll(resp.Body)
	if string(bodyBytes) != content {
		t.Errorf("expected body %q, got %q", content, string(bodyBytes))
	}

	// 3. Test No Match
	req2 := httptest.NewRequest("GET", "https://example.com/other", nil)
	newReq2, resp2 := mw.HandleRequest(req2, nil)

	if resp2 != nil {
		t.Fatal("expected no match, got response")
	}
	if newReq2 != req2 {
		t.Fatal("expected request pass-through")
	}

	// 4. Test Disabled Rule
	mw.AddRule(MapLocalRule{
		ID:        "rule1",
		Enabled:   false,
		UrlRegex:  `.*\/example\.com\/api\/test`,
		LocalPath: tmpFile,
	})

	_, resp3 := mw.HandleRequest(req, nil)

	if resp3 != nil {
		t.Fatal("expected disabled rule to be ignored")
	}
}
