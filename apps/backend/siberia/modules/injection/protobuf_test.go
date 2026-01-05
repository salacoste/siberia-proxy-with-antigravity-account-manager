package injection

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestReplaceField6(t *testing.T) {
	// Construct a fake protobuf message:
	// Field 1 (Varint): 150
	// Field 6 (String): "OLD_TOKEN"
	// Field 2 (String): "OtherData"

	tokenOld := []byte("OLD_TOKEN")
	tokenNew := []byte("NEW_AWESOME_TOKEN")
	otherData := []byte("OtherData")

	// Helper to build
	var buf bytes.Buffer
	// Field 1: Tag (1<<3)|0 = 8. Value 150 = 0x96 0x01
	buf.Write([]byte{0x08, 0x96, 0x01})
	// Field 6: Tag (6<<3)|2 = 0x32. Length 9. Value "OLD_TOKEN"
	buf.WriteByte(0x32)
	buf.WriteByte(byte(len(tokenOld)))
	buf.Write(tokenOld)
	// Field 2: Tag (2<<3)|2 = 0x12. Length 9. Value "OtherData"
	buf.WriteByte(0x12)
	buf.WriteByte(byte(len(otherData)))
	buf.Write(otherData)

	original := buf.Bytes()
	t.Logf("Original Hex: %s", hex.EncodeToString(original))

	// Execute replacement
	modified, err := ReplaceField6(original, tokenNew)
	if err != nil {
		t.Fatalf("ReplaceField6 failed: %v", err)
	}
	t.Logf("Modified Hex: %s", hex.EncodeToString(modified))

	// Verify Structure manually or by attempting to construct expected bytes
	// Expected: Field 1, Field 2, Field 6 (New) -- Order changed because we append 6 at end
	// Note: Our implementation strips 6 from middle and appends to end.

	var expected bytes.Buffer
	// Field 1
	expected.Write([]byte{0x08, 0x96, 0x01})
	// Field 2
	expected.WriteByte(0x12)
	expected.WriteByte(byte(len(otherData)))
	expected.Write(otherData)
	// Field 6 (New)
	expected.WriteByte(0x32)
	expected.WriteByte(byte(len(tokenNew)))
	expected.Write(tokenNew)

	if !bytes.Equal(modified, expected.Bytes()) {
		t.Errorf("Mismatch.\nExpected: %s\nGot:      %s", hex.EncodeToString(expected.Bytes()), hex.EncodeToString(modified))
	}
}

func TestReplaceField6_NoExisting(t *testing.T) {
	// Case where Field 6 doesn't exist initially
	// Field 1 only
	var buf bytes.Buffer
	buf.Write([]byte{0x08, 0x96, 0x01})
	original := buf.Bytes()

	tokenNew := []byte("NEW")

	modified, err := ReplaceField6(original, tokenNew)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	var expected bytes.Buffer
	expected.Write([]byte{0x08, 0x96, 0x01}) // Field 1 preserved
	// Field 6 appended
	expected.WriteByte(0x32)
	expected.WriteByte(byte(len(tokenNew)))
	expected.Write(tokenNew)

	if !bytes.Equal(modified, expected.Bytes()) {
		t.Errorf("Mismatch on append.\nExpected: %s\nGot:      %s", hex.EncodeToString(expected.Bytes()), hex.EncodeToString(modified))
	}
}
