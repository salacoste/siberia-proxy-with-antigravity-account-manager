package claude

import (
	"github.com/google/uuid"
	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// MapResponse converts Gemini Response to Claude MessageResponse
func MapResponse(geminiResp *mappers.GeminiResponse, model string) (*MessageResponse, error) {
	resp := &MessageResponse{
		ID:      "msg_" + uuid.New().String(),
		Type:    "message",
		Role:    "assistant",
		Model:   model,
		Content: []ContentBlock{},
	}

	if geminiResp == nil || len(geminiResp.Candidates) == 0 {
		return resp, nil
	}

	// Claude only supports 1 candidate
	cand := geminiResp.Candidates[0]
	resp.StopReason = mapStopReason(cand.FinishReason)

	for _, part := range cand.Content.Parts {
		if part.Text != "" {
			resp.Content = append(resp.Content, ContentBlock{
				Type: "text",
				Text: part.Text,
			})
		}
		if part.FunctionCall != nil {
			// Map to tool_use
			// Claude expects "input" as map
			resp.Content = append(resp.Content, ContentBlock{
				Type:  "tool_use",
				ID:    "toolu_" + uuid.New().String()[:8], // Generate ID if upstream doesn't provide
				Name:  part.FunctionCall.Name,
				Input: part.FunctionCall.Args,
			})
		}
	}

	// Usage
	if geminiResp.UsageMetadata != nil {
		resp.Usage = &Usage{
			InputTokens:  geminiResp.UsageMetadata.PromptTokenCount,
			OutputTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
		}
	}

	return resp, nil
}

func mapStopReason(finishReason string) string {
	switch finishReason {
	case "STOP":
		return "end_turn"
	case "MAX_TOKENS":
		return "max_tokens"
	case "SAFETY":
		return "stop_sequence" // Approximate
	default:
		return "end_turn"
	}
}
