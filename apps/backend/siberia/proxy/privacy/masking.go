package privacy

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// MaskEmail masks an email address for privacy logging.
// Format: "u***@domain.com" or "***" if invalid.
func MaskEmail(email string) string {
	if email == "" {
		return "<empty>"
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "<invalid-email>"
	}

	local, domain := parts[0], parts[1]

	var maskedLocal string
	if len(local) <= 1 {
		maskedLocal = "***"
	} else {
		// Take first char, mask the rest
		maskedLocal = string(local[0]) + "***"
	}

	return fmt.Sprintf("%s@%s", maskedLocal, domain)
}

// AnonymizeID truncates a long ID to prevent full leakage while allowing basic correlation.
// Format: "start...end"
func AnonymizeID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return "<empty>"
	}
	if len(id) <= 8 {
		return "<redacted>"
	}

	start := id[:4]
	end := id[len(id)-4:]
	return fmt.Sprintf("%s...%s", start, end)
}

// StableHash returns a SHA256 hex string of the value.
// Used for correlating logs without storing raw PII.
func StableHash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}
