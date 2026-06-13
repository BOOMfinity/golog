package gcore

import "sync"

var buffers = sync.Pool{
	New: func() any {
		return new(make([]byte, 0, PooledBufferSize))
	},
}

// AcquireBuffer returns a pre-allocated buffer with capacity equal to [PooledBufferSize].
func AcquireBuffer() *[]byte {
	return buffers.Get().(*[]byte)
}

// FreeBuffer puts the buffer back into the pool.
// If the buffer's capacity exceeds [PooledBufferSize], it is discarded.
func FreeBuffer(buff *[]byte) {
	if buff == nil {
		return
	}
	clear(*buff)
	*buff = (*buff)[:0]
	if cap(*buff) > PooledBufferSize {
		return
	}
	buffers.Put(buff)
}
