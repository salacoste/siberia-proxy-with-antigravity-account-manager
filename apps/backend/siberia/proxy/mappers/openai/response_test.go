package openai

import (
	"encoding/json"
	"testing"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

func TestMapResponse_Basic(t *testing.T) {
	geminiJson := `
	{
		"candidates": [
			{
				"index": 0,
				"finish_reason": "STOP",
				"content": {
					"role": "model",
					"parts": [{"text": "Hello world"}]
				}
			}
		],
		"usage_metadata": {
			"prompt_token_count": 10,
			"candidates_token_count": 5,
			"total_token_count": 15
		}
	}
	`
	var gResp mappers.GeminiResponse
	if err := json.Unmarshal([]byte(geminiJson), &gResp); err != nil {
		t.Fatalf("Unmarshall failed: %v", err)
	}

	oaResp, err := MapResponse(&gResp, "gpt-4")
	if err != nil {
		t.Fatalf("MapResponse failed: %v", err)
	}

	if oaResp.Model != "gpt-4" {
		t.Errorf("Model mismatch")
	}
	if len(oaResp.Choices) != 1 {
		t.Fatalf("Expected 1 choice")
	}
	if oaResp.Choices[0].Message.Content != "Hello world" {
		t.Errorf("Content mismatch")
	}
	if *oaResp.Choices[0].FinishReason != "stop" {
		t.Errorf("Finish reason mismatch, got %s", *oaResp.Choices[0].FinishReason)
	}
	if oaResp.Usage.TotalTokens != 15 {
		t.Errorf("Usage mismatch")
	}
}

func TestMapResponse_FunctionCall(t *testing.T) {
	// Gemini v1internal style function call
	geminiJson := `
	{
		"candidates": [
			{
				"index": 0,
				"finish_reason": "STOP",
				"content": {
					"role": "model",
					"parts": [
						{
							"function_call": {
								"name": "get_weather",
								"args": {"city": "Paris"}
							}
						}
					]
				}
			}
		]
	}
	`
	var gResp mappers.GeminiResponse
	if err := json.Unmarshal([]byte(geminiJson), &gResp); err != nil {
		t.Fatalf("Unmarshall failed: %v", err)
	}

	oaResp, err := MapResponse(&gResp, "gpt-4")
	if err != nil {
		t.Fatalf("MapResponse failed: %v", err)
	}

	msg := oaResp.Choices[0].Message
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("Expected 1 tool call")
	}
	if msg.ToolCalls[0].Function.Name != "get_weather" {
		t.Errorf("Function name mismatch")
	}
	if msg.ToolCalls[0].Function.Arguments == "" {
		t.Errorf("Arguments empty")
	}
}
