package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

type ZaiProvider struct {
	Config *config.AppConfig
}

func NewZaiProvider(cfg *config.AppConfig) *ZaiProvider {
	return &ZaiProvider{Config: cfg}
}

// AnthropicRequest represents the structure of a message request
type AnthropicRequest struct {
	Model     string `json:"model"`
	Messages  any    `json:"messages"` // We keep this generic to pass through
	MaxTokens int    `json:"max_tokens,omitempty"`
	Stream    bool   `json:"stream,omitempty"`
	// Add other fields as map to preserve them?
	// For now, let's decode to map[string]interface{} to modify model safely
	// keeping everything else.
}

func (p *ZaiProvider) ForwardAnthropicJSON(w http.ResponseWriter, r *http.Request) {
	cfg := p.Config
	if cfg.ZaiApiKey == "" {
		http.Error(w, "Z.ai API Key not configured", http.StatusUnauthorized)
		return
	}

	// 1. Parse Body to Map to replace Model
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rawPayload map[string]interface{}
	if err := json.Unmarshal(body, &rawPayload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Apply Robustness Cleaning (Story-71)
	cleaned := mappers.DeepClean(rawPayload)
	cleaned = mappers.SanitizeSchema(cleaned) // Also sanitize schema for Z.ai/Anthropic
	cleaned = mappers.FilterWebSearch(cleaned)

	payload, ok := cleaned.(map[string]interface{})
	if !ok {
		// Should not happen if input was map
		payload = rawPayload
	}

	// 2. Map Model
	originalModel, _ := payload["model"].(string)
	targetModel := originalModel

	// Strip "zai:" prefix if present
	targetModel = strings.TrimPrefix(targetModel, "zai:")

	// Check mapping
	if mapped, ok := cfg.ZaiModelMapping[targetModel]; ok {
		targetModel = mapped
	}
	payload["model"] = targetModel

	// 3. Re-encode
	newBody, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to re-encode body", http.StatusInternalServerError)
		return
	}

	// 4. Create Upstream Request
	// Determine URL. Z.ai usually uses /v1/messages for Anthropic compat
	// But let's verify if ZaiBaseURL includes /v1 or not.
	// Default is "https://api.z.ai/v1", so we append "/messages"
	upstreamURL := strings.TrimRight(cfg.ZaiBaseURL, "/") + "/messages"

	proxyReq, err := http.NewRequest(r.Method, upstreamURL, bytes.NewBuffer(newBody))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 5. Headers
	// Z.ai expects x-api-key or Authorization?
	// Usually strict Anthropic compat means x-api-key.
	proxyReq.Header.Set("x-api-key", cfg.ZaiApiKey)
	proxyReq.Header.Set("anthropic-version", r.Header.Get("anthropic-version"))
	proxyReq.Header.Set("content-type", "application/json")
	// Pass through others if needed?

	// 6. Execute
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upstream error: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 7. Stream Response
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}
