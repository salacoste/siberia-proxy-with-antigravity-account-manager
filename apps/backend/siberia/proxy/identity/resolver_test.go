package identity

import (
	"strings"
	"testing"
)

func TestGenerateMockProjectID(t *testing.T) {
	id := GenerateMockProjectID()
	parts := strings.Split(id, "-")

	if len(parts) != 3 {
		t.Errorf("expected 3 parts, got %d for %s", len(parts), id)
	}

	if len(parts[2]) != 5 {
		t.Errorf("expected 5 char suffix, got %d for %s", len(parts[2]), id)
	}
}

// Note: TestFetchProjectID would require networking or a mock server.
// Skipping network test for CI safety, focusing on logic coverage.
