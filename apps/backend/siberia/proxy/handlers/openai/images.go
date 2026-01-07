package openai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/middleware"
)

// ImageGenerations handles POST /v1/images/generations
func (h *Handler) ImageGenerations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse into generic/internal struct directly or generic map
	// Since we defined mappers.ImageRequest to match OpenAI almost 1:1, we can decode directly into it
	// provided the JSON keys match.
	var req mappers.ImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Validate / Default
	if req.Prompt == "" {
		http.Error(w, "Missing prompt", http.StatusBadRequest)
		return
	}
	if req.N == 0 {
		req.N = 1
	}
	if req.Size == "" {
		req.Size = "1024x1024"
	}

	// 3. Upstream Call
	imgResp, identity, err := h.UpstreamClient.GenerateImage(r.Context(), &req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upstream Error: %v", err), http.StatusBadGateway)
		return
	}

	// 4. Attribution & Response
	maskedIdentity := "unknown"
	if len(identity) > 3 {
		maskedIdentity = identity[:3] + "***"
	} else if identity != "" {
		maskedIdentity = "***"
	}
	middleware.SetAttribution(w, "openai-image", "dall-e-3", maskedIdentity)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(imgResp); err != nil {
		// Too late to change status likely
		return
	}
}
