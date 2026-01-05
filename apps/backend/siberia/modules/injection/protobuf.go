package injection

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Wire Types
const (
	WireVarint          = 0
	WireFixed64         = 1
	WireLengthDelimited = 2
	WireStartGroup      = 3 // Deprecated
	WireEndGroup        = 4 // Deprecated
	WireFixed32         = 5
)

// ReplaceField6 removes any existing Field 6 from the protobuf blob and appends
// a new Field 6 with the provided value (as Length Delimited).
// This is used for injecting tokens into VS Code's session state.
func ReplaceField6(data []byte, newValue []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(make([]byte, 0, len(data)+len(newValue)+10))
	reader := bytes.NewReader(data)

	// 1. Copy everything EXCEPT Field 6
	for {
		// Read Tag (Varint)
		tag, err := binary.ReadUvarint(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tag: %w", err)
		}

		fieldNum := tag >> 3
		wireType := tag & 7

		if fieldNum == 6 {
			// Skip this field
			if err := skipValue(reader, int(wireType)); err != nil {
				return nil, fmt.Errorf("failed to skip field 6: %w", err)
			}
			continue
		}

		// Copy Tag
		writeVarint(buffer, tag)

		// Copy Value
		if err := copyValue(buffer, reader, int(wireType)); err != nil {
			return nil, fmt.Errorf("failed to copy field %d: %w", fieldNum, err)
		}
	}

	// 2. Append new Field 6
	// Tag: (6 << 3) | WireLengthDelimited
	field6Tag := (uint64(6) << 3) | uint64(WireLengthDelimited)
	writeVarint(buffer, field6Tag)

	// Length
	writeVarint(buffer, uint64(len(newValue)))

	// Value
	buffer.Write(newValue)

	return buffer.Bytes(), nil
}

// skipValue advances the reader past the value for the given wire type.
func skipValue(r *bytes.Reader, wireType int) error {
	switch wireType {
	case WireVarint:
		_, err := binary.ReadUvarint(r)
		return err
	case WireFixed64:
		_, err := r.Seek(8, io.SeekCurrent)
		if err != nil {
			return err
		}
		return nil
	case WireLengthDelimited:
		length, err := binary.ReadUvarint(r)
		if err != nil {
			return err
		}
		_, err = r.Seek(int64(length), io.SeekCurrent)
		if err != nil {
			return err
		}
		return nil
	case WireFixed32:
		_, err := r.Seek(4, io.SeekCurrent)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown wire type: %d", wireType)
	}
}

// copyValue reads the value from r and writes it to w.
func copyValue(w *bytes.Buffer, r *bytes.Reader, wireType int) error {
	switch wireType {
	case WireVarint:
		val, err := binary.ReadUvarint(r)
		if err != nil {
			return err
		}
		writeVarint(w, val)
		return nil
	case WireFixed64:
		buf := make([]byte, 8)
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
		w.Write(buf)
		return nil
	case WireLengthDelimited:
		length, err := binary.ReadUvarint(r)
		if err != nil {
			return err
		}
		writeVarint(w, length)
		if length > 0 {
			buf := make([]byte, length)
			if _, err := io.ReadFull(r, buf); err != nil {
				return err
			}
			w.Write(buf)
		}
		return nil
	case WireFixed32:
		buf := make([]byte, 4)
		if _, err := io.ReadFull(r, buf); err != nil {
			return err
		}
		w.Write(buf)
		return nil
	default:
		return fmt.Errorf("unknown wire type: %d", wireType)
	}
}

func writeVarint(w *bytes.Buffer, val uint64) {
	buf := make([]byte, 10) // Max varint size
	n := binary.PutUvarint(buf, val)
	w.Write(buf[:n])
}
