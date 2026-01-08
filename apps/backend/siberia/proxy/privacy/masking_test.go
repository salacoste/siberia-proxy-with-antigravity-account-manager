package privacy

import (
	"testing"
)

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"john.doe@example.com", "j***@example.com"},
		{"a@b.com", "***@b.com"},
		{"ab@b.com", "a***@b.com"},
		{"", "<empty>"},
		{"invalid", "<invalid-email>"},
	}

	for _, tt := range tests {
		got := MaskEmail(tt.input)
		if got != tt.expected {
			t.Errorf("MaskEmail(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestAnonymizeID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1234567890abcdef", "1234...cdef"},
		{"short", "<redacted>"},
		{"12345678", "<redacted>"},
		{"", "<empty>"},
		{"  ", "<empty>"},
	}

	for _, tt := range tests {
		got := AnonymizeID(tt.input)
		if got != tt.expected {
			t.Errorf("AnonymizeID(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestStableHash(t *testing.T) {
	val := "test-value"
	hash1 := StableHash(val)
	hash2 := StableHash(val)

	if hash1 != hash2 {
		t.Errorf("StableHash is not deterministic")
	}

	if len(hash1) != 64 {
		t.Errorf("StableHash length = %d, want 64", len(hash1))
	}
}
