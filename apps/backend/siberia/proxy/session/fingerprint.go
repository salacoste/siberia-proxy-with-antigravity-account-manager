package session

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

type contextKey string

const SessionIDKey contextKey = "siberia_session_id"

// GenerateFingerprint creates a deterministic hash based on the model and the content of the first user message.
// This allows "Sticky Sessions" functionality where requests with the same prompt are routed to the same upstream account.
func GenerateFingerprint(model string, content string) string {
	// 1. Normalize Inputs
	model = strings.TrimSpace(strings.ToLower(model))
	content = strings.TrimSpace(content)

	if content == "" {
		return "" // No fingerprint possible
	}

	// 2. Create Input String
	input := fmt.Sprintf("%s:%s", model, content)

	// 3. Hash (SHA256 for collision resistance, though CRC32/FNV is faster, collisions here mean wrong account which is bad for caching)
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

// ExtractFirstUserMessage is a helper to pull the content from a GeminiRequest safely
func ExtractFirstUserMessage(req *mappers.GeminiRequest) string {
	if req == nil || len(req.Contents) == 0 {
		return ""
	}

	for _, c := range req.Contents {
		if c.Role == "user" && len(c.Parts) > 0 {
			// Concatenate all text parts of the first message
			var builder strings.Builder
			for _, p := range c.Parts {
				if p.Text != "" {
					builder.WriteString(p.Text)
				}
			}
			return builder.String()
		}
	}
	return ""
}
