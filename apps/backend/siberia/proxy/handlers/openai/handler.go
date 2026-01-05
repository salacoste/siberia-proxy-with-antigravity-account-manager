package openai

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers/openai"
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
	// Only support POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse Request
	var oaReq openai.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&oaReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Map to Gemini
	geminiReq, err := openai.MapRequest(&oaReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Mapping error: %v", err), http.StatusBadRequest)
		return
	}

	// 3. Routing/Upstream Call
	// Model ID: OpenAI client sends "gpt-4", we usually map this to a Gemini model or pass through.
	// For now, let's assume we pass "gemini-1.5-pro" or rotate.
	// The mapper doesn't decide the route, the handler or upstream does.
	// Ref doc says: "Maps gpt-4 -> gemini-1.5-pro".
	// Simple mapping for this story:
	targetModel := "models/gemini-1.5-pro-latest" // Default
	if strings.Contains(oaReq.Model, "flash") {
		targetModel = "models/gemini-1.5-flash-latest"
	}

	// 4. Streaming vs Unary
	if oaReq.Stream {
		h.handleStream(w, r, targetModel, geminiReq, oaReq.Model)
		return
	}

	// Unary Handler
	gResp, err := h.UpstreamClient.GenerateContent(r.Context(), targetModel, geminiReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upstream Error: %v", err), http.StatusBadGateway)
		return
	}

	// 5. Map Response
	oaResp, err := openai.MapResponse(gResp, oaReq.Model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Response Mapping Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oaResp)
}
