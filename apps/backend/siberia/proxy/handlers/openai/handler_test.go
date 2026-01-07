package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/mappers/openai"
)

// Mock Client
type MockClient struct {
	Resp   *mappers.GeminiResponse
	Err    error
	Chunks []*mappers.GeminiResponse
}

func (m *MockClient) GenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (*mappers.GeminiResponse, string, error) {
	return m.Resp, "test-identity", m.Err
}

func (m *MockClient) StreamGenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (<-chan *mappers.GeminiResponse, <-chan error) {
	ch := make(chan *mappers.GeminiResponse, len(m.Chunks))
	errCh := make(chan error, 1) // buffered

	go func() {
		for _, c := range m.Chunks {
			ch <- c
			time.Sleep(10 * time.Millisecond) // sim delay
		}
		close(ch)
		close(errCh)
	}()
	return ch, errCh
}

func (m *MockClient) GenerateImage(ctx context.Context, req *mappers.ImageRequest) (*mappers.ImageResponse, string, error) {
	return nil, "test-identity", nil
}

func TestHandler_Unary(t *testing.T) {
	mockC := &MockClient{
		Resp: &mappers.GeminiResponse{
			Candidates: []mappers.GeminiCandidate{
				{
					Content: mappers.GeminiContent{
						Parts: []mappers.GeminiPart{{Text: "Hello"}},
					},
					FinishReason: "STOP",
				},
			},
		},
	}
	h := NewHandler(mockC)

	reqBody := `{"model":"gpt-4", "messages":[{"role":"user", "content":"Hi"}]}`
	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var oaResp openai.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&oaResp); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(oaResp.Choices) != 1 || oaResp.Choices[0].Message.Content != "Hello" {
		t.Errorf("Response mismatch")
	}
}

func TestHandler_Stream(t *testing.T) {
	mockC := &MockClient{
		Chunks: []*mappers.GeminiResponse{
			{
				Candidates: []mappers.GeminiCandidate{
					{Content: mappers.GeminiContent{Parts: []mappers.GeminiPart{{Text: "He"}}}},
				},
			},
			{
				Candidates: []mappers.GeminiCandidate{
					{Content: mappers.GeminiContent{Parts: []mappers.GeminiPart{{Text: "llo"}}}},
				},
			},
		},
	}
	h := NewHandler(mockC)

	reqBody := `{"model":"gpt-4", "messages":[{"role":"user", "content":"Hi"}], "stream": true}`
	req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.Header.Get("Content-Type") != "text/event-stream" {
		t.Errorf("Expected SSE content type")
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !contains(bodyStr, "He") || !contains(bodyStr, "llo") {
		t.Errorf("Stream missing content: %s", bodyStr)
	}
	if !contains(bodyStr, "[DONE]") {
		t.Errorf("Stream not closed with DONE")
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
