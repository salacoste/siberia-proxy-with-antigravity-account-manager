package mcp

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/salacoste/siberia/siberia/config"
)

func TestMcpServer_Initialize(t *testing.T) {
	srv := NewServer(&config.Manager{})

	reqBody := `{"jsonrpc": "2.0", "method": "initialize", "id": 1}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	var jsonResp JsonRpcResponse
	json.Unmarshal(body, &jsonResp)

	if jsonResp.ID.(float64) != 1 {
		t.Errorf("Expected ID 1, got %v", jsonResp.ID)
	}
}

func TestMcpServer_ToolsList(t *testing.T) {
	srv := NewServer(&config.Manager{})
	srv.Registry.Register("test_tool", "A test", nil, nil) // Register dummy

	reqBody := `{"jsonrpc": "2.0", "method": "tools/list", "id": 2}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	var jsonResp JsonRpcResponse
	json.NewDecoder(w.Body).Decode(&jsonResp)

	resMap, _ := jsonResp.Result.(map[string]interface{})
	tools, _ := resMap["tools"].([]interface{})

	// Should be at least 1
	if len(tools) < 1 {
		t.Errorf("Expected at least 1 tool, got %d", len(tools))
	}
}

func TestMcpServer_ToolsCall(t *testing.T) {
	srv := NewServer(&config.Manager{})

	// Register a tool that echoes input
	srv.Registry.Register("echo", "Echo", nil, func(args map[string]interface{}) (string, error) {
		msg, _ := args["msg"].(string)
		return msg, nil
	})

	reqBody := `{"jsonrpc": "2.0", "method": "tools/call", "params": {"name": "echo", "arguments": {"msg": "hello"}}, "id": 3}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(reqBody))
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	var jsonResp JsonRpcResponse
	json.NewDecoder(w.Body).Decode(&jsonResp)

	resMap, _ := jsonResp.Result.(map[string]interface{})
	content, _ := resMap["content"].([]interface{})
	item, _ := content[0].(map[string]interface{})

	if item["text"] != "hello" {
		t.Errorf("Expected 'hello', got %v", item["text"])
	}
}
