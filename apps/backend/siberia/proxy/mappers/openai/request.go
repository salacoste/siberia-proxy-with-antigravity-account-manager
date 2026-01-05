package openai

import (
	"encoding/json"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// MapRequest converts an OpenAI ChatCompletionRequest to a GeminiRequest
func MapRequest(req *ChatCompletionRequest) (*mappers.GeminiRequest, error) {
	geminiReq := &mappers.GeminiRequest{
		Contents:         []mappers.GeminiContent{},
		GenerationConfig: &mappers.GeminiGenConfig{},
	}

	// 1. Generation Config
	geminiReq.GenerationConfig.Temperature = req.Temperature
	geminiReq.GenerationConfig.TopP = req.TopP
	geminiReq.GenerationConfig.CandidateCount = req.N
	geminiReq.GenerationConfig.MaxOutputTokens = req.MaxTokens

	if req.Stop != nil {
		switch v := req.Stop.(type) {
		case string:
			geminiReq.GenerationConfig.StopSequences = []string{v}
		case []interface{}:
			var stops []string
			for _, s := range v {
				if str, ok := s.(string); ok {
					stops = append(stops, str)
				}
			}
			geminiReq.GenerationConfig.StopSequences = stops
		}
	}

	// 2. Messages & System Instruction
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			// Extract as System Instruction
			if geminiReq.SystemInstruction == nil {
				geminiReq.SystemInstruction = &mappers.GeminiContent{
					Parts: []mappers.GeminiPart{},
				}
			}
			text := extractTextContent(msg.Content)
			geminiReq.SystemInstruction.Parts = append(geminiReq.SystemInstruction.Parts, mappers.GeminiPart{Text: text})
			continue
		}

		// Regular Message
		geminiContent := mappers.GeminiContent{
			Role: mapRole(msg.Role),
		}

		// Content
		parts, err := mapContent(msg.Content)
		if err != nil {
			return nil, err
		}
		geminiContent.Parts = append(geminiContent.Parts, parts...)

		// Tool Calls
		if len(msg.ToolCalls) > 0 {
			for _, tc := range msg.ToolCalls {
				if tc.Type == "function" {
					var args map[string]interface{}
					if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
						// Fallback: treat as empty map or error? internal google api is strict
						args = make(map[string]interface{})
					}
					geminiContent.Parts = append(geminiContent.Parts, mappers.GeminiPart{
						FunctionCall: &mappers.GeminiFuncCall{
							Name: tc.Function.Name,
							Args: args,
						},
					})
				}
			}
		}

		geminiReq.Contents = append(geminiReq.Contents, geminiContent)
	}

	// 3. Tools Definition
	if len(req.Tools) > 0 {
		geminiTool := mappers.GeminiTool{}
		for _, t := range req.Tools {
			if t.Type == "function" {
				geminiTool.FunctionDeclarations = append(geminiTool.FunctionDeclarations, mappers.GeminiFuncDecl{
					Name:        t.Function.Name,
					Description: t.Function.Description,
					Parameters:  t.Function.Parameters,
				})
			}
		}
		geminiReq.Tools = append(geminiReq.Tools, geminiTool)
	}

	return geminiReq, nil
}

func mapRole(role string) string {
	if role == "assistant" {
		return "model"
	}
	return "user"
}

func extractTextContent(content any) string {
	switch v := content.(type) {
	case string:
		return v
	case []interface{}:
		var sb strings.Builder
		for _, p := range v {
			if partMap, ok := p.(map[string]interface{}); ok {
				if t, ok := partMap["type"].(string); ok && t == "text" {
					if txt, ok := partMap["text"].(string); ok {
						sb.WriteString(txt)
					}
				}
			}
		}
		return sb.String()
	}
	return ""
}

func mapContent(content any) ([]mappers.GeminiPart, error) {
	var parts []mappers.GeminiPart
	if content == nil {
		return parts, nil
	}

	switch v := content.(type) {
	case string:
		parts = append(parts, mappers.GeminiPart{Text: v})
	case []interface{}: // "content": [...]
		// Need to parse complex array (text, image_url)
		// Since we defined ContentPart struct but `json.Unmarshal` into `any` gives map[string]interface{},
		// let's iterate the map.
		for _, p := range v {
			partMap, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			if typ, ok := partMap["type"].(string); ok {
				if typ == "text" {
					if txt, ok := partMap["text"].(string); ok {
						parts = append(parts, mappers.GeminiPart{Text: txt})
					}
				} else if typ == "image_url" {
					// Handle Image (Placeholder for now, usually needs Base64 fetching or inline)
					// Verify if we have URL or Base64. Gemini needs inline data usually.
					// Implementation note: We might need to download URL or pass as is if supported.
					// For now, simple text fallback/warning.
					// TODO: Implement image download/inline logic from Reference App.
				}
			}
		}
	}
	return parts, nil
}
