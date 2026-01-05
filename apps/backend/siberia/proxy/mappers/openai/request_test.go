package openai

import (
	"encoding/json"
	"testing"
)

func TestMapRequest_Basic(t *testing.T) {
	jsonReq := `
	{
		"model": "gpt-4",
		"messages": [
			{"role": "system", "content": "You are a helper."},
			{"role": "user", "content": "Hello"}
		],
		"temperature": 0.7
	}
	`
	var req ChatCompletionRequest
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		t.Fatalf("Failed to unmarshal test req: %v", err)
	}

	geminiReq, err := MapRequest(&req)
	if err != nil {
		t.Fatalf("MapRequest failed: %v", err)
	}

	// Checks
	if geminiReq.SystemInstruction == nil {
		t.Error("Expected SystemInstruction, got nil")
	} else if len(geminiReq.SystemInstruction.Parts) != 1 || geminiReq.SystemInstruction.Parts[0].Text != "You are a helper." {
		t.Errorf("System Instruction mismatch")
	}

	if len(geminiReq.Contents) != 1 {
		t.Fatalf("Expected 1 content message, got %d", len(geminiReq.Contents))
	}
	if geminiReq.Contents[0].Role != "user" {
		t.Errorf("Expected role 'user', got %s", geminiReq.Contents[0].Role)
	}
	if geminiReq.Contents[0].Parts[0].Text != "Hello" {
		t.Errorf("Expected text 'Hello', got %s", geminiReq.Contents[0].Parts[0].Text)
	}
	if geminiReq.GenerationConfig.Temperature != 0.7 {
		t.Errorf("Expected Temp 0.7, got %f", geminiReq.GenerationConfig.Temperature)
	}
}

func TestMapRequest_Tools(t *testing.T) {
	jsonReq := `
	{
		"model": "gpt-4",
		"messages": [{"role": "user", "content": "Code"}],
		"tools": [
			{
				"type": "function",
				"function": {
					"name": "run_code",
					"parameters": {"type": "object"}
				}
			}
		]
	}
	`
	var req ChatCompletionRequest
	if err := json.Unmarshal([]byte(jsonReq), &req); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	geminiReq, err := MapRequest(&req)
	if err != nil {
		t.Fatalf("MapRequest failed: %v", err)
	}

	if len(geminiReq.Tools) != 1 {
		t.Fatalf("Expected 1 tool definition")
	}
	if len(geminiReq.Tools[0].FunctionDeclarations) != 1 {
		t.Fatalf("Expected 1 function declaration")
	}
	if geminiReq.Tools[0].FunctionDeclarations[0].Name != "run_code" {
		t.Errorf("Tool name mismatch")
	}
}
