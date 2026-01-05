package openai

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// OpenAI Chat Completion Response
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Message `json:"delta,omitempty"` // For streaming
	FinishReason *string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// MapResponse converts a Gemini Response to OpenAI Response
func MapResponse(geminiResp *mappers.GeminiResponse, model string) (*ChatCompletionResponse, error) {
	resp := &ChatCompletionResponse{
		ID:      "chatcmpl-" + uuid.New().String(),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []Choice{},
	}

	if geminiResp == nil || len(geminiResp.Candidates) == 0 {
		return resp, nil
	}

	for i, cand := range geminiResp.Candidates {
		choice := Choice{
			Index: i,
		}

		// Message Content
		choice.Message = &Message{
			Role: "assistant",
		}

		var contentTxt string
		for _, part := range cand.Content.Parts {
			if part.Text != "" {
				contentTxt += part.Text
			}
			if part.FunctionCall != nil {
				// Map FunctionCall -> ToolCall
				argsBytes, _ := json.Marshal(part.FunctionCall.Args)
				choice.Message.ToolCalls = append(choice.Message.ToolCalls, ToolCall{
					ID:   "call_" + uuid.New().String()[:8], // Gemini doesn't ensure ID uniqueness per call often
					Type: "function",
					Function: ToolCallFunction{
						Name:      part.FunctionCall.Name,
						Arguments: string(argsBytes),
					},
				})
			}
		}
		choice.Message.Content = contentTxt

		// Finish Reason
		reason := mapFinishReason(cand.FinishReason)
		choice.FinishReason = &reason

		resp.Choices = append(resp.Choices, choice)
	}

	// Usage (Approximation if not provided by Gemini v1internal correctly)
	if geminiResp.UsageMetadata != nil {
		resp.Usage = &Usage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		}
	}

	return resp, nil
}

func mapFinishReason(geminiReason string) string {
	// Gemini: STOP, MAX_TOKENS, SAFETY, RECITATION, OTHER
	switch geminiReason {
	case "STOP":
		return "stop"
	case "MAX_TOKENS":
		return "length"
	case "SAFETY":
		return "content_filter"
	case "RECITATION":
		return "content_filter"
	default:
		return "stop" // Default fallback
	}
}
