package claude

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
					"parts": [{"text": "Hello Claude user"}]
				}
			}
		],
		"usage_metadata": {
			"prompt_token_count": 50,
			"candidates_token_count": 10
		}
	}
	`
	var gResp mappers.GeminiResponse
	if err := json.Unmarshal([]byte(geminiJson), &gResp); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	cResp, err := MapResponse(&gResp, "claude-3-opus")
	if err != nil {
		t.Fatalf("MapResponse failed: %v", err)
	}

	if cResp.Content[0].Text != "Hello Claude user" {
		t.Errorf("Content mismatch")
	}
	if cResp.StopReason != "end_turn" {
		t.Errorf("StopReason mismatch, got %s", cResp.StopReason)
	}
	if cResp.Usage.InputTokens != 50 {
		t.Errorf("Usage mismatch")
	}
}

func TestMapResponse_ToolUse(t *testing.T) {
	geminiJson := `
	{
		"candidates": [
			{
				"content": {
					"parts": [
						{
							"function_call": {
								"name": "calc",
								"args": {"x": 1, "y": 2}
							}
						}
					]
				},
				"finish_reason": "STOP"
			}
		]
	}
	`
	var gResp mappers.GeminiResponse
	if err := json.Unmarshal([]byte(geminiJson), &gResp); err != nil {
		t.Fatalf("Unmarshal failed")
	}

	cResp, err := MapResponse(&gResp, "claude-3-5-sonnet")
	if err != nil {
		t.Fatalf("MapResponse failed")
	}

	block := cResp.Content[0]
	if block.Type != "tool_use" {
		t.Fatalf("Expected tool_use")
	}
	if block.Name != "calc" {
		t.Errorf("Name mismatch")
	}
	inputs := block.Input.(map[string]interface{})
	if inputs["x"].(float64) != 1 {
		t.Errorf("Args mismatch")
	}
}
