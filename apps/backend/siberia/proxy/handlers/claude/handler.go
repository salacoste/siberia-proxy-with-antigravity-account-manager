package claude

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers/claude"
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

	// 5. Unary
	gResp, err := h.UpstreamClient.GenerateContent(r.Context(), targetModel, gReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upstream Error: %v", err), http.StatusBadGateway)
		return
	}

	cResp, err := claude.MapResponse(gResp, cReq.Model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Response Mapping Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cResp)
}
