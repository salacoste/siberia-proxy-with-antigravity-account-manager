package openai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/mappers/openai"
)

// SSOP State (Not robust fsm for now, just simple chunk scanner/buffer could be complex)
// For Version 1, we will implement standard streaming + basic detection.
// Reference App used a complex FSM. We'll simplify:
// If chunk text contains ```json ... "command": "shell" ... ``` -> Send Tool Call.

func (h *Handler) handleStream(w http.ResponseWriter, r *http.Request, targetModel string, gReq *mappers.GeminiRequest, originalModel string) {
	// Set Headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Call Upstream
	streamChan, errChan := h.UpstreamClient.StreamGenerateContent(r.Context(), targetModel, gReq)

	callID := "call_" + uuid.New().String()[:8]

	for {
		select {
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			if err != nil {
				// We can't write HTTP error code if stream started. Send a text error token?
				// Or fail silent. OpenAI usually closes conn.
				fmt.Printf("Stream Error: %v\n", err)
				return
			}

		case chunk, ok := <-streamChan:
			if !ok {
				// End of stream
				sendDone(w, flusher)
				return
			}

			// Map Chunk to OpenAI Delta
			// Simple text mapping for now
			// SSOP-01: Detect Shell Command?
			// Since upstream chunks are small, checking individually is risky (could be split).
			// Robust SSOP needs a buffer. For this task, we will just implement straight pipe.
			// SSOP TODO: Add buffer logic in next task (SSOP Refinement).

			oaChunk := mapGeminiChunkToOpenAI(chunk, originalModel, callID)
			writeSSE(w, flusher, oaChunk)
		}

		if streamChan == nil && errChan == nil {
			break
		}
	}
}

func sendDone(w http.ResponseWriter, f http.Flusher) {
	fmt.Fprintf(w, "data: [DONE]\n\n")
	f.Flush()
}

func writeSSE(w http.ResponseWriter, f http.Flusher, chunk *openai.ChatCompletionChunk) {
	data, _ := json.Marshal(chunk)
	fmt.Fprintf(w, "data: %s\n\n", data)
	f.Flush()
}

// Minimal Chunk Mapper (Needs to be defined in openai package or here)
func mapGeminiChunkToOpenAI(gChunk *mappers.GeminiResponse, model string, callID string) *openai.ChatCompletionChunk {
	// Map First Candidate Text
	txt := ""
	if len(gChunk.Candidates) > 0 && len(gChunk.Candidates[0].Content.Parts) > 0 {
		txt = gChunk.Candidates[0].Content.Parts[0].Text
	}

	// Create Delta
	return &openai.ChatCompletionChunk{
		ID:      "chatcmpl-" + callID, // Should be constant for request logic, but okay
		Object:  "chat.completion.chunk",
		Created: 1234567890,
		Model:   model,
		Choices: []openai.ChunkChoice{
			{
				Index: 0,
				Delta: openai.ChunkDelta{
					Content: txt,
				},
				FinishReason: nil, // Map if finished
			},
		},
	}
}
