package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// Mock Client
type MockClient struct {
	Resp      *mappers.GeminiResponse
	Err       error
	Chunks    []*mappers.GeminiResponse
	LastModel string
}

func (m *MockClient) GenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (*mappers.GeminiResponse, error) {
	m.LastModel = model
	return m.Resp, m.Err
}

func (m *MockClient) StreamGenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (<-chan *mappers.GeminiResponse, <-chan error) {
	m.LastModel = model
	ch := make(chan *mappers.GeminiResponse, len(m.Chunks))
	errCh := make(chan error, 1)

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

func TestHandler_Unary(t *testing.T) {
	mockC := &MockClient{
		Resp: &mappers.GeminiResponse{
			Candidates: []mappers.GeminiCandidate{
				{
					Content: mappers.GeminiContent{
						Parts: []mappers.GeminiPart{{Text: "Hello Gemini"}},
					},
					FinishReason: "STOP",
				},
			},
		},
	}
	h := NewHandler(mockC)

	reqBody := `{"contents":[{"parts":[{"text":"Hi"}]}]}`
	req := httptest.NewRequest("POST", "/v1beta/models/gemini-1.5-flash:generateContent", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	if mockC.LastModel != "models/gemini-1.5-flash" {
		t.Errorf("Model parsing failed. Got %s", mockC.LastModel)
	}

	var gResp mappers.GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if len(gResp.Candidates) != 1 || gResp.Candidates[0].Content.Parts[0].Text != "Hello Gemini" {
		t.Errorf("Response mismatch")
	}
}

func TestHandler_Stream(t *testing.T) {
	mockC := &MockClient{
		Chunks: []*mappers.GeminiResponse{
			{Candidates: []mappers.GeminiCandidate{{Content: mappers.GeminiContent{Parts: []mappers.GeminiPart{{Text: "He"}}}}}},
			{Candidates: []mappers.GeminiCandidate{{Content: mappers.GeminiContent{Parts: []mappers.GeminiPart{{Text: "llo"}}}}}},
		},
	}
	h := NewHandler(mockC)

	reqBody := `{"contents":[{"parts":[{"text":"Hi"}]}]}`
	req := httptest.NewRequest("POST", "/v1beta/models/gemini-pro:streamGenerateContent", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected application/json content type")
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	// Check format [Chunk, Chunk]
	if !contains(bodyStr, "[") || !contains(bodyStr, "]") {
		t.Errorf("Expected JSON Array format. Got: %s", bodyStr)
	}
	if !contains(bodyStr, "He") || !contains(bodyStr, "llo") {
		t.Errorf("Stream missing content: %s", bodyStr)
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
