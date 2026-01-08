package proxy

import "sync"

const DefaultBufferSize = 32 * 1024 // 32KB

var bufferPool = sync.Pool{
	New: func() interface{} {
		// Allocate a new byte slice with capacity
		b := make([]byte, DefaultBufferSize)
		return &b
	},
}

// GetBuffer returns a byte slice from the pool.
func GetBuffer() *[]byte {
	return bufferPool.Get().(*[]byte)
}

// PutBuffer returns a byte slice to the pool.
func PutBuffer(b *[]byte) {
	if b == nil || cap(*b) < DefaultBufferSize {
		return
	}
	// Reset not strictly necessary as we overwrite, but good practice
	// We don't shrink it here.
	bufferPool.Put(b)
}
