package claude

import (
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

		for _, block := range msg.Content {
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
				// Gemini v1internal expects {name, response: {content: ...}} strictly
				respMap := map[string]interface{}{
					"content": block.Content,
				}
				// Note: Gemini actually needs the function name corresponding to this result.
				// Claude provides `tool_use_id`. We might need to look up the name or pass ID if Gemini supports.
				// Internal API often implies strict {name: "func_name", response: {}} mapping.
				// If we don't have the name from request history easily, this is tricky.
				// However, for proxying, we assume 1:1 map.
				// IMPORTANT: Reference App handles this by tracking tool calls or assuming simplistic mapping.
				// We'll map to generic structure for now.
				gContent.Parts = append(gContent.Parts, mappers.GeminiPart{
					FunctionResponse: &mappers.GeminiFuncResp{
						Name: block.ToolUseID, // This might fail if Gemini expects Name not ID.
						// Ref app usually fixes this by context lookup or pass-thru.
						Response: respMap,
					},
				})

			case "thinking":
				// Upgrade/Downgrade logic handled here?
				// Gemini doesn't support "thinking" input directly usually, unless enabled on model.
				// We typically treat it as text or drop it if Gemini strictly doesn't support it.
				// For now, treat as Text to preserve context.
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
