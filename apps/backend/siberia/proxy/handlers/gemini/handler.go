package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/middleware"
	"github.com/salacoste/siberia/siberia/proxy/session"
	"github.com/salacoste/siberia/siberia/proxy/upstream"
)

type Handler struct {
	UpstreamClient upstream.Client
}

func NewHandler(client upstream.Client) *Handler {
	return &Handler{
		UpstreamClient: client,
	}
}

// ServeHTTP handles Native Gemini /v1beta/models/...:generateContent requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse Model from Path
	// Path should be like /v1beta/models/{model}:{action}
	// e.g. /v1beta/models/gemini-1.5-pro:generateContent
	parts := strings.Split(r.URL.Path, "/")
	// Minimal validation
	var model string
	for _, p := range parts {
		if strings.HasPrefix(p, "models/") || strings.Contains(p, "gemini") {
			// Extract pure model part before :
			if idx := strings.Index(p, ":"); idx != -1 {
				model = p[:idx]
			} else {
				model = p
			}
			break
		}
	}
	// Fallback/Cleanup
	// If path is .../models/gemini-1.5-pro:generateContent
	// p might be gemini-1.5-pro
	if model == "" {
		// Try to parse more strictly if needed, but for now let upstream handle defaults or errors?
		// Or assume client sends valid path.
		// Let's rely on decoding body mainly.
		// Actually upstream client NEEDS model string to select key or build path.
		// Let's take the segments.
		// /v1beta/models/gemini...
		// [ , v1beta, models, gemini-1.5-flash:generateContent ]
		if len(parts) > 3 && parts[2] == "models" {
			seg := parts[3]
			if idx := strings.Index(seg, ":"); idx != -1 {
				model = "models/" + seg[:idx]
			} else {
				model = "models/" + seg
			}
		} else {
			// Default?
			model = "models/gemini-1.5-pro-latest"
		}
	}

	// Ensure prefix
	if model != "" && !strings.HasPrefix(model, "models/") {
		model = "models/" + model
	}

	// 2. Decode Body
	var gReq mappers.GeminiRequest
	if err := json.NewDecoder(r.Body).Decode(&gReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 3. Routing/Streaming Logic
	isStream := strings.HasSuffix(r.URL.Path, ":streamGenerateContent")

	if isStream {
		h.handleStream(w, r, model, &gReq)
		return
	}

	// 4. Unary Call
	// 4a. Inject Session ID from Header if present (Sticky Sessions)
	ctx := r.Context()
	if sessionID := r.Header.Get("x-siberia-session-id"); sessionID != "" {
		ctx = context.WithValue(ctx, session.SessionIDKey, sessionID)
	}

	gResp, identity, err := h.UpstreamClient.GenerateContent(ctx, model, &gReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upstream Error: %v", err), http.StatusBadGateway)
		return
	}

	// 5. Encode Response
	middleware.SetAttribution(w, "google", model, identity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gResp)
}

func (h *Handler) handleStream(w http.ResponseWriter, r *http.Request, model string, req *mappers.GeminiRequest) {
	middleware.SetAttribution(w, "google", model, "unknown")
	w.Header().Set("Content-Type", "application/json") // Gemini stream is usually SSE or JSON array stream?
	// Gemini API stream is NOT SSE normally. It's a bracketed JSON list if using REST?
	// ACTUALLY: `streamGenerateContent` returns a stream of `GenerateContentResponse`.
	// "By default, the response is a stream of JSON objects".

	// Wait, standard Gemini REST:
	// returns: [{...}, {...}] partial JSONs?
	// OR: Server-Sent Events?
	// Documentation says: "The response is a standard HTTP response... writing a JSON list... The response will contain a series of GenerateContentResponse objects..."
	// Actually, most clients expect SSE "data: ..." lines if using browser?
	// But `google-generative-ai` SDKs often use raw JSON stream `[{},{},...]`.
	// BUT, `curl` tests show: [`{...},\r\n{...}`]

	// Let's assume we output pure JSON chunks.

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// We'll write start bracket for array?
	// Actually UpstreamClient.StreamGenerateContent returns chunks.
	// If we want to emulate Gemini REST exactly, we should probably output `[` then objects separated by `,` then `]`.

	w.Write([]byte("["))
	flusher.Flush()

	// Inject Session ID for stream too
	ctx := r.Context()
	if sessionID := r.Header.Get("x-siberia-session-id"); sessionID != "" {
		ctx = context.WithValue(ctx, session.SessionIDKey, sessionID)
	}

	ch, errCh := h.UpstreamClient.StreamGenerateContent(ctx, model, req)

	first := true
	for {
		select {
		case chunk, ok := <-ch:
			if !ok {
				// Done
				w.Write([]byte("]"))
				flusher.Flush()
				return
			}

			if !first {
				w.Write([]byte(",\n"))
			}
			first = false

			if err := json.NewEncoder(w).Encode(chunk); err != nil {
				return
			}
			flusher.Flush()

		case err := <-errCh:
			if err != nil {
				// Can't change status code now.
				// Maybe write error object?
				fmt.Printf("[GeminiHandler] Stream Error: %v\n", err)
			}
			return
		case <-r.Context().Done():
			return
		}
	}
}
