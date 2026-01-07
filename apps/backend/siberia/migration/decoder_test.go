package migration

import (
	"testing"
)

func TestExtractRefreshToken(t *testing.T) {
	// Construct a fake blob:
	// Outer: Field 6 -> Inner
	// Inner: Field 3 -> "my-refresh-token"

	expectedToken := "1//test-token-123"

	// Encode Inner
	innerMsg := EncodeProtoField(3, []byte(expectedToken))

	// Encode Outer
	outerMsg := EncodeProtoField(6, innerMsg)

	// Test
	got, err := ExtractRefreshToken(outerMsg)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if got != expectedToken {
		t.Errorf("Expected %s, got %s", expectedToken, got)
	}
}

func TestExtractRefreshToken_Invalid(t *testing.T) {
	// Test 1: Empty
	_, err := ExtractRefreshToken([]byte{})
	if err == nil {
		t.Error("Expected error for empty blob")
	}

	// Test 2: Field 6 Missing
	wrongMsg := EncodeProtoField(5, []byte("wrong"))
	_, err = ExtractRefreshToken(wrongMsg)
	if err == nil {
		t.Error("Expected error for missing field 6")
	}

	// Test 3: Field 6 Present but Inner Field 3 Missing
	innerWrong := EncodeProtoField(2, []byte("wrong-inner"))
	outerWrong := EncodeProtoField(6, innerWrong)
	_, err = ExtractRefreshToken(outerWrong)
	if err == nil {
		t.Error("Expected error for missing inner field 3")
	}
}
