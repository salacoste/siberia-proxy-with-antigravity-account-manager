package claude

import (
	"encoding/json"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// MapRequest converts Claude MessageRequest to GeminiRequest
func MapRequest(req *MessageRequest) (*mappers.GeminiRequest, error) {
	geminiReq := &mappers.GeminiRequest{
		Contents:         []mappers.GeminiContent{},
		GenerationConfig: &mappers.GeminiGenConfig{},
	}

	// 1. Config
	geminiReq.GenerationConfig.MaxOutputTokens = req.MaxTokens
	geminiReq.GenerationConfig.Temperature = req.Temperature
	geminiReq.GenerationConfig.TopP = req.TopP
	geminiReq.GenerationConfig.TopK = req.TopK
	geminiReq.GenerationConfig.StopSequences = req.StopSequences

	// 2. System Prompt
	// Claude passes system prompt as top-level field, Gemini uses separate struct
	if req.System != nil {
		geminiReq.SystemInstruction = &mappers.GeminiContent{Parts: []mappers.GeminiPart{}}
		switch v := req.System.(type) {
		case string:
			geminiReq.SystemInstruction.Parts = append(geminiReq.SystemInstruction.Parts, mappers.GeminiPart{Text: v})
		case []interface{}: // Array of blocks
			// Simple text extraction for now
			for _, b := range v {
				if bm, ok := b.(map[string]interface{}); ok {
					if t, ok := bm["text"].(string); ok {
						geminiReq.SystemInstruction.Parts = append(geminiReq.SystemInstruction.Parts, mappers.GeminiPart{Text: t})
					}
				}
			}
		}
	}

	// 3. Messages
	for _, msg := range req.Messages {
		gContent := mappers.GeminiContent{
			Role: mapRole(msg.Role),
		}

		var blocks []ContentBlock
		switch v := msg.Content.(type) {
		case string:
			blocks = []ContentBlock{{Type: "text", Text: v}}
		case []interface{}:
			// Need to re-marshal or manually map map[string]interface
			// This is painful with 'any'. Better to use helper or try convert.
			// Ideally we use json.RawMessage or custom Unmarshal.
			// For simplicity in this step, let's assume if it came from JSON decode into 'any', it's []interface{}.
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					cb := ContentBlock{}
					// Minimal mapping for now
					if t, ok := m["type"].(string); ok {
						cb.Type = t
					}
					if t, ok := m["text"].(string); ok {
						cb.Text = t
					}
					// ... map other fields if needed ...
					// Fast implementation:
					jsonDat, _ := json.Marshal(m)
					json.Unmarshal(jsonDat, &cb)
					blocks = append(blocks, cb)
				}
			}
		case []ContentBlock:
			blocks = v
		}

		for _, block := range blocks {
			// CRITICAL: Strip Cache Control (Implicitly done by not mapping it to anything in Gemini)

			// Handle different block types
			switch block.Type {
			case "text":
				gContent.Parts = append(gContent.Parts, mappers.GeminiPart{Text: block.Text})

			case "image":
				if block.Source != nil {
					// Gemini InlineData
					gContent.Parts = append(gContent.Parts, mappers.GeminiPart{
						InlineData: &mappers.GeminiBlob{
							MimeType: block.Source.MediaType,
							Data:     block.Source.Data,
						},
					})
				}

			case "tool_use":
				// Map to FunctionCall
				// Claude Input is usually a map, Gemini needs separate Name/Args
				args, _ := block.Input.(map[string]interface{})
				gContent.Parts = append(gContent.Parts, mappers.GeminiPart{
					FunctionCall: &mappers.GeminiFuncCall{
						Name: block.Name,
						Args: args,
					},
				})

			case "tool_result":
				// Map to FunctionResponse
				respMap := map[string]interface{}{
					"content": block.Content,
				}
				gContent.Parts = append(gContent.Parts, mappers.GeminiPart{
					FunctionResponse: &mappers.GeminiFuncResp{
						Name:     block.ToolUseID,
						Response: respMap,
					},
				})

			case "thinking":
				if block.Thinking != "" {
					gContent.Parts = append(gContent.Parts, mappers.GeminiPart{Text: block.Thinking})
				}
			}
		}
		geminiReq.Contents = append(geminiReq.Contents, gContent)
	}

	// 4. Tools
	if len(req.Tools) > 0 {
		gTool := mappers.GeminiTool{}
		for _, t := range req.Tools {
			gTool.FunctionDeclarations = append(gTool.FunctionDeclarations, mappers.GeminiFuncDecl{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.InputSchema,
			})
		}
		geminiReq.Tools = append(geminiReq.Tools, gTool)
	}

	return geminiReq, nil
}

func mapRole(role string) string {
	if role == "assistant" {
		return "model"
	}
	return "user"
}
