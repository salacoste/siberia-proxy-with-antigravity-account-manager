package mcp

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/salacoste/siberia/siberia/config"
	zaitools "github.com/salacoste/siberia/siberia/zai/tools" // Import the new tools package
)

type ToolHandler func(args map[string]interface{}) (string, error)

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type Registry struct {
	mu       sync.RWMutex
	tools    []Tool
	handlers map[string]ToolHandler
	cfgMgr   *config.Manager // Need config to access API Keys
}

func NewRegistry(cfgMgr *config.Manager) *Registry {
	r := &Registry{
		tools:    make([]Tool, 0),
		handlers: make(map[string]ToolHandler),
		cfgMgr:   cfgMgr,
	}
	r.registerDefaults()
	return r
}

func (r *Registry) registerDefaults() {
	// 1. Web Search
	r.Register("search_web", "Search the web for information using Z.ai Prime.",
		schemaSearch(),
		func(args map[string]interface{}) (string, error) {
			// Extract Args
			query, _ := args["query"].(string)
			recency, _ := args["recency"].(string)
			// Handle domain_filter array... simplified for now

			req := zaitools.SearchRequest{
				Query:   query,
				Recency: recency,
			}

			cfg := r.cfgMgr.Get()
			results, err := zaitools.CallWebSearchPrime(cfg.ZaiApiKey, cfg.ZaiBaseURL, req)
			if err != nil {
				return "", err
			}

			// Return as JSON string
			resBytes, _ := json.Marshal(results)
			return string(resBytes), nil
		})

	// 2. Web Reader
	r.Register("read_web_page", "Read the content of a URL and convert it to clean Markdown.",
		schemaReader(),
		func(args map[string]interface{}) (string, error) {
			url, ok := args["url"].(string)
			if !ok || url == "" {
				return "", fmt.Errorf("url required")
			}
			return zaitools.CallWebReader(url)
		})
}

// Helpers for Schema Definition (JSON Schema)
func schemaSearch() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"query": { "type": "string", "description": "Search query" },
			"recency": { "type": "string", "enum": ["past_24h", "past_week", "past_month", "past_year"], "description": "Filter by time" }
		},
		"required": ["query"]
	}`)
}

func schemaReader() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"url": { "type": "string", "description": "URL to read" }
		},
		"required": ["url"]
	}`)
}

func (r *Registry) Register(name, description string, schema json.RawMessage, handler ToolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tools = append(r.tools, Tool{
		Name:        name,
		Description: description,
		InputSchema: schema,
	})
	r.handlers[name] = handler
}

func (r *Registry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Return copy
	list := make([]Tool, len(r.tools))
	copy(list, r.tools)
	return list
}

func (r *Registry) Call(name string, args map[string]interface{}) (string, error) {
	r.mu.RLock()
	handler, ok := r.handlers[name]
	r.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	return handler(args)
}
