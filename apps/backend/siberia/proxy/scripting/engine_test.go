package scripting

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScriptEngine_OnRequest_ModifyHeader(t *testing.T) {
	engine := NewScriptEngine()
	script := `
	function onRequest(req) {
		req.Headers["X-Test"] = ["Scripted"];
		return req;
	}
	`
	engine.UpdateScript(script, true)

	req := httptest.NewRequest("GET", "http://example.com/api", nil)
	modReq, err := engine.RunOnRequest(req)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if modReq.Header.Get("X-Test") != "Scripted" {
		t.Errorf("Expected X-Test header 'Scripted', got '%s'", modReq.Header.Get("X-Test"))
	}
}

func TestScriptEngine_OnResponse_ModifyStatus(t *testing.T) {
	engine := NewScriptEngine()
	script := `
	function onResponse(res) {
		res.Status = 418;
		return res;
	}
	`
	// Javascript object properties for our struct are TitleCase because of Export?
	// Goja maps struct fields. Struct tags `json:"status"` might NOT be used by default logic unless using specific mapper?
	// Goja default: Field names.
	// Wait, in engine.go I defined structs with JSON tags.
	// Does goja use them?
	// Allow for FieldName usage or Tag usage.
	// Let's check common behavior. Usually it's FieldName unless mapped.
	// IMPORTANT: My engine.go uses `vm.ExportTo(val, &newJsResp)`.
	// For that to work, the JS object must match.
	// Let's test what field names work.

	// Updating script to use PascalCase if JSON tags aren't picked up automatically by vm.ToValue/ExportTo without configuration.
	// However, vm.ToValue uses current value.
	// Let's try standard field names first (Go struct fields).

	script = `
	function onResponse(res) {
		// Goja maps struct fields directly by name usually
		res.StatusCode = 418;
		return res;
	}
	`
	engine.UpdateScript(script, true)

	resp := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("OK")),
	}

	modResp, err := engine.RunOnResponse(resp)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if modResp.StatusCode != 418 {
		t.Errorf("Expected status 418, got %d", modResp.StatusCode)
	}
}

func TestScriptEngine_SyntaxError(t *testing.T) {
	engine := NewScriptEngine()
	engine.UpdateScript("var x = ;", true)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	_, err := engine.RunOnRequest(req)
	if err == nil {
		t.Error("Expected syntax error, got nil")
	}
}
