package claude

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers/claude"
	"github.com/salacoste/siberia/siberia/proxy/middleware"
	"github.com/salacoste/siberia/siberia/proxy/upstream"
)

type Handler struct {
	UpstreamClient upstream.Client
}

func NewHandler(client upstream.Client) *Handler {
	return &Handler{
		UpstreamClient: client,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse
	var cReq claude.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&cReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Map
	gReq, err := claude.MapRequest(&cReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Mapping Error: %v", err), http.StatusBadRequest)
		return
	}

	// 3. Routing & Downgrade Logic
	// Default to Pro
	targetModel := "models/gemini-1.5-pro-latest"

	// Smart Downgrade: If it's a background agent loop or specifically requesting fast model
	// Or if detected "Thinking" usage without signature (handled in mapper), but here we might switch model too.
	// Simple Heuristic: If model name contains "haiku", map to flash.
	if strings.Contains(cReq.Model, "haiku") || strings.Contains(cReq.Model, "flash") {
		targetModel = "models/gemini-1.5-flash-latest"
	}

	// 4. Streaming
	if cReq.Stream {
		h.handleStream(w, r, targetModel, gReq, cReq.Model)
		return
	}

	gResp, identity, err := h.UpstreamClient.GenerateContent(r.Context(), targetModel, gReq)
	if err != nil {
		// Story-54: Fallback Retry on 400 (Invalid Argument)
		// Often caused by "Thinking" blocks not being supported by the specific upstream model/region
		// or strict validation failures. We strip ALL thinking blocks and retry.
		if strings.Contains(err.Error(), "400") {
			fmt.Println("[ClaudeHandler] 400 Error detected, retrying without thinking blocks...")
			cReq.Messages = stripAllThinking(cReq.Messages)

			// Remap
			if gReqRetry, errMap := claude.MapRequest(&cReq); errMap == nil {
				if gRespRetry, identRetry, errRetry := h.UpstreamClient.GenerateContent(r.Context(), targetModel, gReqRetry); errRetry == nil {
					gResp = gRespRetry
					identity = identRetry
					err = nil
				} else {
					// Return original error if retry fails
					fmt.Printf("[ClaudeHandler] Retry failed: %v\n", errRetry)
				}
			}
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Upstream Error: %v", err), http.StatusBadGateway)
			return
		}
	}

	cResp, err := claude.MapResponse(gResp, cReq.Model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Response Mapping Error: %v", err), http.StatusInternalServerError)
		return
	}

	middleware.SetAttribution(w, "anthropic-shim", targetModel, identity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cResp)
}

// sanitizeMessages removes invalid/empty thinking blocks and strips trailing thinking blocks
// to prevent upstream 400 errors (Vertex AI/Gemini).
func sanitizeMessages(msgs []claude.Message) []claude.Message {
	if len(msgs) == 0 {
		return msgs
	}

	out := make([]claude.Message, 0, len(msgs))

	for i, msg := range msgs {
		isLastMsg := (i == len(msgs)-1)

		var newContent any

		switch v := msg.Content.(type) {
		case string:
			newContent = v
		case []interface{}: // JSON Array
			var newBlocks []interface{}
			for _, item := range v {
				m, ok := item.(map[string]interface{})
				if !ok {
					newBlocks = append(newBlocks, item)
					continue
				}

				typ, _ := m["type"].(string)
				if typ == "thinking" {
					// Check content
					val, _ := m["thinking"].(string)
					val = strings.TrimSpace(val)

					// Rule 1: Drop empty thinking blocks
					if val == "" {
						continue
					}
					// Rule 2: If malformed (e.g. very short "..." or just signature), convert to text or drop?
					// For now, strict check: if < 5 chars, likely garbage.
					if len(val) < 5 {
						continue
					}
				}
				newBlocks = append(newBlocks, item)
			}

			// Rule 3: Remove Trailing Thinking Block (only if it's the very last block of last message)
			if isLastMsg && len(newBlocks) > 0 {
				lastIdx := len(newBlocks) - 1
				if lastMap, ok := newBlocks[lastIdx].(map[string]interface{}); ok {
					if t, _ := lastMap["type"].(string); t == "thinking" {
						// Drop it
						newBlocks = newBlocks[:lastIdx]
					}
				}
			}

			// If all blocks removed, default to empty string or skip message?
			// Claude API requires non-empty content usually.
			// But if it was ONLY thinking and we stripped it, we might have an issue.
			// Let's leave it empty slice for now, mapper might handle or upstream fails cleanly.
			newContent = newBlocks

		default:
			newContent = v
		}

		msg.Content = newContent
		out = append(out, msg)
	}

	return out
}

// stripAllThinking removes ALL thinking blocks for fallback retry
func stripAllThinking(msgs []claude.Message) []claude.Message {
	out := make([]claude.Message, 0, len(msgs))
	for _, msg := range msgs {
		var newContent any
		switch v := msg.Content.(type) {
		case []interface{}:
			var newBlocks []interface{}
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					if t, _ := m["type"].(string); t != "thinking" {
						newBlocks = append(newBlocks, item)
					}
				} else {
					newBlocks = append(newBlocks, item)
				}
			}
			newContent = newBlocks
		default:
			newContent = v
		}
		msg.Content = newContent
		out = append(out, msg)
	}
	return out
}
