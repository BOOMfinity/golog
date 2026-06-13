package gcore

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	// SyncedStdout is the [os.Stdout] protected with a [sync.Mutex].
	//
	// It is used by default for non-error logs.
	//
	// DO NOT edit this variable.
	SyncedStdout = NewSyncedWriter(os.Stdout)
	// SyncedStderr is the [os.Stderr] protected with a [sync.Mutex].
	//
	// It is used by default for errors.
	//
	// DO NOT edit this variable.
	SyncedStderr = NewSyncedWriter(os.Stderr)
)

// Config provides basic information for the [Formatter].
//
// All writers must be safe for concurrent writes.
// Use [NewSyncedWriter] to wrap your own writer with a [sync.Mutex].
type Config struct {
	StandardOutput io.Writer
	ErrorOutput    io.Writer
	TimeFormat     string
	Severity       Level
}

// Writer returns the [io.Writer] based on the log level.
//
// It falls back to [SyncedStdout] or [SyncedStderr] if no writers are set.
// If only standard output is set, it will be used for all logs.
func (c Config) Writer(l Level) io.Writer {
	var (
		w  io.Writer = SyncedStdout
		we io.Writer = SyncedStderr
	)
	if c.StandardOutput != nil {
		w = c.StandardOutput
		we = nil
	}
	if c.ErrorOutput != nil {
		we = c.ErrorOutput
	}
	if l >= LevelError && we != nil {
		return we
	}
	return w
}

func init() {
	if os.Getenv("GOLOG_LOGGING_LEVEL") != "" {
		level := os.Getenv("GOLOG_LOGGING_LEVEL")
		n, err := strconv.Atoi(level)
		if err != nil || n < 1 || n > 6 {
			println(fmt.Sprintf("[golog] invalid logging level: %s", level))
		} else {
			OverrideLoggingLevel = Level(n)
		}
	}
}
