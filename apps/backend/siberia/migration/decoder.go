package migration

import (
	"fmt"

	"google.golang.org/protobuf/encoding/protowire"
)

// ExtractRefreshToken walks the binary blob to find Field 6 -> Field 3 (Refresh Token).
// Logic based on observed structure:
// Outer Message (jetskiStateSync): Field 6 (WireType 2) -> Contains Inner Message
// Inner Message (agentManagerInitState): Field 3 (WireType 2) -> Contains the Refresh Token string.
func ExtractRefreshToken(blob []byte) (string, error) {
	// 1. Scan Outer Message for Field 6
	innerBlob, err := findFieldBytes(blob, 6)
	if err != nil {
		return "", fmt.Errorf("failed to find outer field 6: %w", err)
	}

	// 2. Scan Inner Message for Field 3
	tokenBytes, err := findFieldBytes(innerBlob, 3)
	if err != nil {
		return "", fmt.Errorf("failed to find inner field 3 (token): %w", err)
	}

	return string(tokenBytes), nil
}

// findFieldBytes scans a byte slice for a specific field number with WireType 2 (Length Delimited)
// and returns its content.
func findFieldBytes(data []byte, targetFieldNum int32) ([]byte, error) {
	for len(data) > 0 {
		num, typ, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, fmt.Errorf("invalid proto tag")
		}
		data = data[n:]

		if int32(num) == int32(targetFieldNum) {
			if typ != protowire.BytesType {
				return nil, fmt.Errorf("field %d has unexpected wire type %d (expected 2/Bytes)", num, typ)
			}
			// Consume value (length + bytes)
			val, m := protowire.ConsumeBytes(data)
			if m < 0 {
				return nil, fmt.Errorf("invalid bytes value for field %d", num)
			}
			return val, nil
		}

		// Skip other fields
		m := protowire.ConsumeFieldValue(num, typ, data)
		if m < 0 {
			return nil, fmt.Errorf("failed to skip field %d", num)
		}
		data = data[m:]
	}

	return nil, fmt.Errorf("field %d not found", targetFieldNum)
}

// Helper to manually construct simple proto msg (Length prefixed)
func manualEncode(num int, content []byte) []byte {
	// This is tricky without the library helpers fully exposing write.
	// But `protowire.AppendBytes` handles the length prefixing of the value.
	// `protowire.AppendTag` handles the Key.
	// But note: AppendTag prepends to the slice provided as first arg? No, it appends.

	// Wait, `AppendTag(b, num, type)` appends the Tag varint to b.
	// `AppendBytes(b, v)` appends length(v) + v to b.

	// So:
	b := protowire.AppendTag(nil, protowire.Number(int32(num)), protowire.BytesType)
	b = protowire.AppendBytes(b, content)
	return b
}

// Encode for testing purposes (public so tests can use it)
func EncodeProtoField(num int, content []byte) []byte {
	return manualEncode(num, content)
}
