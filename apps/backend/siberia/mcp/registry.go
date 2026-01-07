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

	// 3. Vision: UI to Artifact
	r.Register("ui_to_artifact", "Generate React/Tailwind code from a UI screenshot.",
		schemaImageTool(),
		func(args map[string]interface{}) (string, error) {
			path, _ := args["image_path"].(string)
			if path == "" {
				return "", fmt.Errorf("image_path required")
			}
			cfg := r.cfgMgr.Get()
			return zaitools.CallUiToArtifact(cfg.ZaiApiKey, cfg.ZaiBaseURL, path)
		})

	// 4. Vision: Extract Text
	r.Register("extract_text_from_screenshot", "Extract text from an image (OCR).",
		schemaImageTool(),
		func(args map[string]interface{}) (string, error) {
			path, _ := args["image_path"].(string)
			if path == "" {
				return "", fmt.Errorf("image_path required")
			}
			cfg := r.cfgMgr.Get()
			return zaitools.CallExtractText(cfg.ZaiApiKey, cfg.ZaiBaseURL, path)
		})

	// 5. Vision: Diagnose Error
	r.Register("diagnose_error_screenshot", "Analyze an error dialog or traceback screenshot.",
		schemaImageTool(),
		func(args map[string]interface{}) (string, error) {
			path, _ := args["image_path"].(string)
			if path == "" {
				return "", fmt.Errorf("image_path required")
			}
			cfg := r.cfgMgr.Get()
			return zaitools.CallDiagnoseError(cfg.ZaiApiKey, cfg.ZaiBaseURL, path)
		})

	// 6. Vision: Understand Diagram
	r.Register("understand_technical_diagram", "Explain a technical architecture or workflow diagram.",
		schemaImageTool(),
		func(args map[string]interface{}) (string, error) {
			path, _ := args["image_path"].(string)
			if path == "" {
				return "", fmt.Errorf("image_path required")
			}
			cfg := r.cfgMgr.Get()
			return zaitools.CallUnderstandDiagram(cfg.ZaiApiKey, cfg.ZaiBaseURL, path)
		})
}

func schemaImageTool() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"image_path": { "type": "string", "description": "Local path or URL to the image" }
		},
		"required": ["image_path"]
	}`)
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
