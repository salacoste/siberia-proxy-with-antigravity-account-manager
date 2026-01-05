package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
)

func TestProxyAuth(t *testing.T) {
	cfg := &config.AppConfig{
		AuthEnabled: true,
		AuthToken:   "secret-token",
		ZaiEnabled:  false,
	}
	svc := NewService(cfg, nil, nil, nil)

	// 1. Test Without Header -> 401
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Result().StatusCode)
	}

	// 2. Test With Invalid Header -> 401
	req = httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer wrong-token")
	w = httptest.NewRecorder()
	svc.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Result().StatusCode)
	}

	// 3. Test With Valid Authorization -> 200 (or whatever goproxy returns for nil upstream)
	// Actually goproxy without handler might fail or return something, but certainly NOT 401.
	// Since we haven't set up the full server+client loop, goproxy.ServeHTTP might panic or error if not fully mocked,
	// but let's see if we pass the auth check.
	req = httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	w = httptest.NewRecorder()

	// We expect the auth code to strip the header and pass control to s.proxy.ServeHTTP
	// To verify this without relying on goproxy internals, we can inspect if headers were deleted?
	// But ServeHTTP modifies the request object passed in.

	// Let's just run it. If it returns 404 or 500 that's fine, as long as it's NOT 401.
	// goproxy default behavior for unknown request?

	svc.ServeHTTP(w, req)
	if w.Result().StatusCode == http.StatusUnauthorized || w.Result().StatusCode == http.StatusProxyAuthRequired {
		t.Errorf("Should have passed auth check, but got %d", w.Result().StatusCode)
	}
}
