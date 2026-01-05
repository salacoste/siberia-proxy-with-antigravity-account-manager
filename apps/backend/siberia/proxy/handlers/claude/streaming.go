package claude

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// Events used in Claude Streaming
type Event struct {
	Type         string      `json:"type"`
	Message      interface{} `json:"message,omitempty"` // For message_start
	Index        int         `json:"index,omitempty"`
	ContentBlock interface{} `json:"content_block,omitempty"` // For content_block_start
	Delta        interface{} `json:"delta,omitempty"`         // For content_block_delta
	Usage        interface{} `json:"usage,omitempty"`         // For message_delta
	StopReason   *string     `json:"stop_reason,omitempty"`
}

func (h *Handler) handleStream(w http.ResponseWriter, r *http.Request, targetModel string, gReq *mappers.GeminiRequest, originalModel string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	streamChan, errChan := h.UpstreamClient.StreamGenerateContent(r.Context(), targetModel, gReq)

	msgID := "msg_" + uuid.New().String()[:8]
	hasStarted := false

	for {
		select {
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			if err != nil {
				fmt.Printf("Claude Stream Error: %v\n", err)
				return // Abort stream
			}
		case chunk, ok := <-streamChan:
			if !ok {
				// End of Stream
				sendStopEvents(w, flusher)
				return
			}

			if !hasStarted {
				sendStartEvents(w, flusher, msgID, originalModel)
				hasStarted = true
			}

			// extract text
			txt := ""
			if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
				txt = chunk.Candidates[0].Content.Parts[0].Text
			}

			if txt != "" {
				sendDelta(w, flusher, txt)
			}
		}

		if streamChan == nil && errChan == nil {
			break
		}
	}
}

func sendEvent(w http.ResponseWriter, f http.Flusher, eventName string, data interface{}) {
	jsonData, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventName, jsonData)
	f.Flush()
}

func sendStartEvents(w http.ResponseWriter, f http.Flusher, id, model string) {
	// 1. message_start
	msgStart := Event{
		Type: "message_start",
		Message: map[string]interface{}{
			"id":            id,
			"type":          "message",
			"role":          "assistant",
			"content":       []interface{}{},
			"model":         model,
			"stop_reason":   nil,
			"stop_sequence": nil,
			"usage":         map[string]int{"input_tokens": 10}, // Placeholder
		},
	}
	sendEvent(w, f, "message_start", msgStart)

	// 2. content_block_start (Assume one text block for now)
	blockStart := Event{
		Type:         "content_block_start",
		Index:        0,
		ContentBlock: map[string]string{"type": "text", "text": ""},
	}
	sendEvent(w, f, "content_block_start", blockStart)
}

func sendDelta(w http.ResponseWriter, f http.Flusher, text string) {
	delta := Event{
		Type:  "content_block_delta",
		Index: 0,
		Delta: map[string]string{"type": "text_delta", "text": text},
	}
	sendEvent(w, f, "content_block_delta", delta)
}

func sendStopEvents(w http.ResponseWriter, f http.Flusher) {
	// 1. content_block_stop
	stopBlock := Event{
		Type:  "content_block_stop",
		Index: 0,
	}
	sendEvent(w, f, "content_block_stop", stopBlock)

	// 2. message_delta (stop reason)
	s := "end_turn"
	msgDelta := Event{
		Type:  "message_delta",
		Delta: map[string]interface{}{"stop_reason": s, "stop_sequence": nil},
		Usage: map[string]int{"output_tokens": 10}, // Placeholder
	}
	sendEvent(w, f, "message_delta", msgDelta)

	// 3. message_stop
	msgStop := Event{Type: "message_stop"}
	sendEvent(w, f, "message_stop", msgStop)
}
