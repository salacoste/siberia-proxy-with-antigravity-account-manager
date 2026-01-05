package claude

import (
	"encoding/json"
	"testing"
)

func TestMapRequest_ThinkingAndSystem(t *testing.T) {
	jsonReq := `
	{
		"model": "claude-3-opus",
		"system": "Be precise.",
		"messages": [
			{
				"role": "user", 
				"content": [
					{"type": "text", "text": "Solve this."}
				]
			},
			{
				"role": "assistant",
				"content": [
					{"type": "thinking", "thinking": "I need to calculate...", "signature": "sig123"},
					{"type": "text", "text": "Here is the answer."}
				]
			}
		],
		"max_tokens": 1000
	}
	`
	var req MessageRequest
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	geminiReq, err := MapRequest(&req)
	if err != nil {
		t.Fatalf("MapRequest failed: %v", err)
	}

	// scalable system prompt check
	if geminiReq.SystemInstruction == nil || geminiReq.SystemInstruction.Parts[0].Text != "Be precise." {
		t.Errorf("System prompt mismatch")
	}

	// Check thinking conversion (Should map to Text in Gemini currently)
	modelMsg := geminiReq.Contents[1]
	if len(modelMsg.Parts) != 2 {
		t.Fatalf("Expected 2 parts (Thinking->Text, Text), got %d", len(modelMsg.Parts))
	}
	if modelMsg.Parts[0].Text != "I need to calculate..." {
		t.Errorf("Thinking block not converted to Text")
	}
}

func TestMapRequest_CacheControlStripped(t *testing.T) {
	jsonReq := `
	{
		"model": "claude-3-5-sonnet",
		"messages": [
			{
				"role": "user",
				"content": [
					{
						"type": "text", 
						"text": "Hello", 
						"cache_control": {"type": "ephemeral"}
					}
				]
			}
		]
	}
	`
	var req MessageRequest
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	geminiReq, err := MapRequest(&req)
	if err != nil {
		t.Fatalf("MapRequest failed: %v", err)
	}

	// Gemini struct has no cache_control field, so implicit verification that it's ignored/not mapped to anything else
	if geminiReq.Contents[0].Parts[0].Text != "Hello" {
		t.Errorf("Text mismatch")
	}
}
