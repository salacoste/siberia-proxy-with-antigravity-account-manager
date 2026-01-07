package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/salacoste/siberia/siberia/config"
)

// JSON-RPC Types
type JsonRpcRequest struct {
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

type JsonRpcResponse struct {
	JsonRpc string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
	ID      interface{}   `json:"id"`
}

type JsonRpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type McpServer struct {
	Config   *config.Manager
	Registry *Registry
}

func NewServer(cfgMgr *config.Manager) *McpServer {
	return &McpServer{
		Config:   cfgMgr,
		Registry: NewRegistry(cfgMgr),
	}
}

func (s *McpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JsonRpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, nil, -32700, "Parse error")
		return
	}

	var result interface{}
	var err error

	switch req.Method {
	case "initialize":
		// Return server capabilities
		result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]string{
				"name":    "siberia-internal-mcp",
				"version": "1.0.0",
			},
		}

	case "tools/list":
		tools := s.Registry.ListTools()
		result = map[string]interface{}{
			"tools": tools,
		}

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if json.Unmarshal(req.Params, &params) != nil {
			s.writeError(w, req.ID, -32602, "Invalid params")
			return
		}
		result, err = s.callTool(params.Name, params.Arguments)
		if err != nil {
			s.writeError(w, req.ID, -32000, err.Error())
			return
		}
		// Wrap result in content array as per MCP spec
		result = map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("%v", result),
				},
			},
		}

	default:
		s.writeError(w, req.ID, -32601, "Method not found")
		return
	}

	s.writeResponse(w, req.ID, result)
}

func (s *McpServer) callTool(name string, args map[string]interface{}) (interface{}, error) {
	return s.Registry.Call(name, args)
}

func (s *McpServer) writeResponse(w http.ResponseWriter, id interface{}, result interface{}) {
	resp := JsonRpcResponse{
		JsonRpc: "2.0",
		ID:      id,
		Result:  result,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *McpServer) writeError(w http.ResponseWriter, id interface{}, code int, message string) {
	resp := JsonRpcResponse{
		JsonRpc: "2.0",
		ID:      id,
		Error: &JsonRpcError{
			Code:    code,
			Message: message,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
