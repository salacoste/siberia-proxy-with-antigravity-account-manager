package proxy

import (
	"testing"
)

func TestBufferPool(t *testing.T) {
	// Get a buffer
	bufPtr := GetBuffer()
	if bufPtr == nil {
		t.Fatal("GetBuffer returned nil")
	}
	if len(*bufPtr) != DefaultBufferSize {
		t.Errorf("Expected buffer length %d, got %d", DefaultBufferSize, len(*bufPtr))
	}

	// Modify buffer
	(*bufPtr)[0] = 0xFF

	// Put it back
	PutBuffer(bufPtr)

	// Get another one (might be the same)
	bufPtr2 := GetBuffer()
	if bufPtr2 == nil {
		t.Fatal("GetBuffer returned nil")
	}
	if len(*bufPtr2) != DefaultBufferSize {
		t.Errorf("Expected buffer length %d, got %d", DefaultBufferSize, len(*bufPtr2))
	}
}
