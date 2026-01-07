package claude

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/mappers/claude"
)

// Mock Client (Same as OpenAI mock essentially, should move to shared util if reused more)
type MockClient struct {
	Resp      *mappers.GeminiResponse
	Err       error
	Chunks    []*mappers.GeminiResponse
	LastModel string
}

func (m *MockClient) GenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (*mappers.GeminiResponse, string, error) {
	m.LastModel = model
	return m.Resp, "test-identity", m.Err
}

func (m *MockClient) StreamGenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (<-chan *mappers.GeminiResponse, <-chan error) {
	m.LastModel = model
	ch := make(chan *mappers.GeminiResponse, len(m.Chunks))
	errCh := make(chan error, 1)

	go func() {
		for _, c := range m.Chunks {
			ch <- c
			time.Sleep(5 * time.Millisecond)
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
					Content:      mappers.GeminiContent{Parts: []mappers.GeminiPart{{Text: "Claude reply"}}},
					FinishReason: "STOP",
				},
			},
		},
	}
	h := NewHandler(mockC)

	reqBody := `{"model":"claude-3-opus", "messages":[{"role":"user", "content":"Hi"}]}`
	req := httptest.NewRequest("POST", "/v1/messages", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var cResp claude.MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	if cResp.Content[0].Text != "Claude reply" {
		t.Errorf("Content mismatch")
	}
}

func TestHandler_Downgrade(t *testing.T) {
	mockC := &MockClient{Resp: &mappers.GeminiResponse{}}
	h := NewHandler(mockC)

	// Haiku should map to Flash
	reqBody := `{"model":"claude-3-haiku", "messages":[{"role":"user", "content":"Fast"}]}`
	req := httptest.NewRequest("POST", "/v1/messages", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if !strings.Contains(mockC.LastModel, "flash") {
		t.Errorf("Expected model downgrade to flash, got %s", mockC.LastModel)
	}
}

func TestHandler_StreamEvents(t *testing.T) {
	mockC := &MockClient{
		Chunks: []*mappers.GeminiResponse{
			{Candidates: []mappers.GeminiCandidate{{Content: mappers.GeminiContent{Parts: []mappers.GeminiPart{{Text: "Hel"}}}}}},
			{Candidates: []mappers.GeminiCandidate{{Content: mappers.GeminiContent{Parts: []mappers.GeminiPart{{Text: "lo"}}}}}},
		},
	}
	h := NewHandler(mockC)

	reqBody := `{"model":"claude-3-opus", "messages":[{"role":"user", "content":"Stream"}], "stream": true}`
	req := httptest.NewRequest("POST", "/v1/messages", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	// Check for Claude Event Names
	if !strings.Contains(bodyStr, "event: message_start") {
		t.Errorf("Missing message_start")
	}
	if !strings.Contains(bodyStr, "event: content_block_delta") {
		t.Errorf("Missing content_block_delta")
	}
	if !strings.Contains(bodyStr, "event: message_stop") {
		t.Errorf("Missing message_stop")
	}
}
