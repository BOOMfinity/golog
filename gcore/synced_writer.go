package gcore

import (
	"io"
	"sync"
)

// SyncedWriter is a thread-safe wrapper for [io.Writer].
// It ensures that you can write to it concurrently.
type SyncedWriter struct {
	w  io.Writer
	mu sync.Mutex
}

// Write implements [io.Writer].
func (sw *SyncedWriter) Write(p []byte) (n int, err error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.w.Write(p)
}

// NewSyncedWriter creates the [SyncedWriter] from [io.Writer].
func NewSyncedWriter(w io.Writer) *SyncedWriter {
	return &SyncedWriter{
		w: w,
	}
}
