package legacy

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// MockUpstreamClient implements upstream.Client
type MockUpstreamClient struct {
	CapturedModel   string
	CapturedRequest *mappers.GeminiRequest
	ReturnResponse  *mappers.GeminiResponse
}

func (m *MockUpstreamClient) GenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (*mappers.GeminiResponse, string, error) {
	m.CapturedModel = model
	m.CapturedRequest = req
	return m.ReturnResponse, "mock-identity", nil
}

func (m *MockUpstreamClient) StreamGenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (<-chan *mappers.GeminiResponse, <-chan error) {
	return nil, nil // Not tested here
}

func (m *MockUpstreamClient) GenerateImage(ctx context.Context, req *mappers.ImageRequest) (*mappers.ImageResponse, string, error) {
	return nil, "", nil // Not tested here
}

func (m *MockUpstreamClient) IsAvailable() bool { return true }

func TestHandler_ServeHTTP(t *testing.T) {
	// Setup Mock
	mockClient := &MockUpstreamClient{
		ReturnResponse: &mappers.GeminiResponse{
			Candidates: []mappers.GeminiCandidate{
				{
					Content: mappers.GeminiContent{
						Role: "model",
						Parts: []mappers.GeminiPart{
							{Text: "print('Hello World')"},
						},
					},
					FinishReason: "STOP",
				},
			},
		},
	}

	handler := NewHandler(mockClient)

	// Create Request
	payload := LegacyCompletionRequest{
		Model:        "copilot-codex",
		Input:        []string{"def hello():", "    "},
		Instructions: "You are a Python expert.",
		MaxTokens:    50,
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/v1/completions", bytes.NewBuffer(bodyBytes))
	w := httptest.NewRecorder()

	// Execute
	handler.ServeHTTP(w, req)

	// Verify Response Code
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	// Verify Response Body
	var legResp LegacyCompletionResponse
	json.NewDecoder(resp.Body).Decode(&legResp)

	if len(legResp.Choices) == 0 {
		t.Fatalf("Expected choices, got 0")
	}

	if legResp.Choices[0].Text != "print('Hello World')" {
		t.Errorf("Expected 'print('Hello World')', got '%s'", legResp.Choices[0].Text)
	}

	if legResp.Object != "text_completion" {
		t.Errorf("Expected object 'text_completion', got '%s'", legResp.Object)
	}

	// Verify Transformation Logic (Mock Capture)
	// Input should be joined by newline
	// Gemini req mapping is complex (Chat -> Gemini).
	// We at least verify the call was made.
	if mockClient.CapturedModel == "" {
		t.Error("Mock client was not called")
	}
	// Default model logic check
	if mockClient.CapturedModel != "models/gemini-1.5-pro-latest" {
		t.Errorf("Expected default gemini model, got %s", mockClient.CapturedModel)
	}
}
