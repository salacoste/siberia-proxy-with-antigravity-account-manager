package legacy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers/openai"
	"github.com/salacoste/siberia/siberia/proxy/middleware"
	"github.com/salacoste/siberia/siberia/proxy/upstream"
)

// LegacyCompletionRequest represents the Codex-style payload
type LegacyCompletionRequest struct {
	Model        string   `json:"model"`
	Input        []string `json:"input"`        // Code lines
	Instructions string   `json:"instructions"` // System prompt
	MaxTokens    int      `json:"max_tokens,omitempty"`
	Temperature  float64  `json:"temperature,omitempty"`
	Stop         []string `json:"stop,omitempty"`
}

// LegacyCompletionResponse represents the OpenAI v1/completions response
type LegacyCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	LogProbs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

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

	// 1. Decode Legacy Request
	var legReq LegacyCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&legReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Transform to Chat Format
	// Input (lines) -> User Message
	// Instructions -> System Message
	inputContent := strings.Join(legReq.Input, "\n")

	chatReq := openai.ChatCompletionRequest{
		Model: legReq.Model,
		Messages: []openai.Message{
			{Role: "system", Content: legReq.Instructions},
			{Role: "user", Content: inputContent},
		},
		MaxTokens:   legReq.MaxTokens,
		Temperature: legReq.Temperature,
		Stop:        legReq.Stop, // openai.ChatCompletionRequest Stop is `any`, we pass []string
	}

	// 3. Map to Gemini Request
	geminiReq, err := openai.MapRequest(&chatReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Mapping error: %v", err), http.StatusBadRequest)
		return
	}

	// 4. Determine Target Model
	targetModel := "models/gemini-1.5-pro-latest" // Default
	if strings.Contains(legReq.Model, "flash") {
		targetModel = "models/gemini-1.5-flash-latest"
	} else if strings.Contains(legReq.Model, "gpt-3.5") {
		targetModel = "models/gemini-1.5-flash-latest" // GPT-3.5 equivalent
	}

	// 5. Call Upstream
	// Note: We only support unary for now as per story scope implies text completion
	gResp, identity, err := h.UpstreamClient.GenerateContent(r.Context(), targetModel, geminiReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upstream Error: %v", err), http.StatusBadGateway)
		return
	}

	// 6. Map Back to Chat Response (Intermediate)
	chatResp, err := openai.MapResponse(gResp, legReq.Model)
	if err != nil {
		http.Error(w, fmt.Sprintf("Response Mapping Error: %v", err), http.StatusInternalServerError)
		return
	}

	// 7. Transform to Legacy Response
	legacyResp := LegacyCompletionResponse{
		ID:      chatResp.ID,
		Object:  "text_completion",
		Created: chatResp.Created,
		Model:   chatResp.Model,
	}

	if chatResp.Usage != nil {
		legacyResp.Usage = Usage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		}
	}

	// Extract content
	if len(chatResp.Choices) > 0 {
		content := ""
		// Chat choices have Message.Content
		if chatResp.Choices[0].Message.Content != nil {
			content = fmt.Sprintf("%v", chatResp.Choices[0].Message.Content)
		}

		legacyResp.Choices = []Choice{
			{
				Text:         content,
				Index:        0,
				LogProbs:     nil,
				FinishReason: safeDeref(chatResp.Choices[0].FinishReason),
			},
		}
	} else {
		legacyResp.Choices = []Choice{}
	}

	// 8. Send Response
	maskedIdentity := "unknown"
	if len(identity) > 3 {
		maskedIdentity = identity[:3] + "***"
	} else if identity != "" {
		maskedIdentity = "***"
	}
	middleware.SetAttribution(w, "openai-legacy-shim", targetModel, maskedIdentity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(legacyResp)
}

func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
