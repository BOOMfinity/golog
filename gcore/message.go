package gcore

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// Level is a severity level of log entry.
type Level uint

const (
	LevelFatal   Level = 6
	LevelError   Level = 5
	LevelWarning Level = 4
	LevelInfo    Level = 3
	LevelDebug   Level = 2
	LevelTrace   Level = 1
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarning:
		return "WARNING"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	}
	panic("invalid logging level")
}

func acquireMessage(parent *Logger, level Level, c ...context.Context) *Message {
	msg := messages.Get().(*Message)
	msg.level = level
	msg.parent = parent
	if len(c) > 0 {
		msg.ctx = c[0]
	}
	return msg
}

// Message represents a single log.
//
// Messages returned by the [Logger] are pre-allocated, and
// they are freed after calling Msg, Msgf or MsgErr.
type Message struct {
	parent   *Logger
	message  []byte
	level    Level
	ctx      context.Context
	attrs    []Attribute
	err      error
	stack    []byte
	exitCode int
}

// Stack captures the goroutine call stack and attaches it to the log.
func (msg *Message) Stack() *Message {
	msg.stack = appendStacktrace(msg.stack)
	return msg
}

// Attr adds an [Attribute] to the log.
func (msg *Message) Attr(name string, val any) *Message {
	msg.attrs = append(msg.attrs, Attribute{Name: name, Value: val})
	return msg
}

// Cause attaches an error to the log.
func (msg *Message) Cause(err error) *Message {
	msg.err = err
	return msg
}

// Msg sets a message and sends the log.
func (msg *Message) Msg(str string) {
	msg.message = append(msg.message, str...)
	msg.handle()
}

// Msgf sets a formatted message and sends the log.
func (msg *Message) Msgf(format string, args ...any) {
	msg.message = fmt.Appendf(msg.message, format, args...)
	msg.handle()
}

// MsgErr sends the log using error as a message.
func (msg *Message) MsgErr(err error) {
	msg.Cause(err).handle()
}

func (msg *Message) handle() {
	if msg.exitCode != 0 {
		defer os.Exit(msg.exitCode)
	}
	defer msg.resetAndPutToPool()

	if len(msg.message) == 0 && msg.err == nil {
		return
	}

	for i := range msg.parent.hooks {
		msg.parent.hooks[i]((*contextImpl)(msg))
	}

	required := msg.parent.level

	if OverrideLoggingLevel > 0 {
		required = OverrideLoggingLevel
	}

	if msg.level < required {
		return
	}

	msg.parent.formatter((*contextImpl)(msg))
}

func (msg *Message) init() {
	msg.attrs = make([]Attribute, 0, PreAllocatedSliceSize)
	msg.message = make([]byte, 0, PreAllocatedMessageBufferSize)
	msg.stack = make([]byte, 0, PreAllocatedStackBufferSize)
}

func (msg *Message) resetAndPutToPool() {
	canPutBack := cap(msg.message) == PreAllocatedMessageBufferSize && cap(msg.attrs) == PreAllocatedSliceSize && cap(msg.stack) == PreAllocatedStackBufferSize
	if len(msg.attrs) > 0 {
		clear(msg.attrs)
		msg.attrs = msg.attrs[:0]
	}
	if len(msg.message) > 0 {
		clear(msg.message)
		msg.message = msg.message[:0]
	}
	if len(msg.stack) > 0 {
		clear(msg.stack)
		msg.stack = msg.stack[:0]
	}

	msg.parent = nil
	msg.level = 0
	msg.ctx = nil
	msg.exitCode = 0
	msg.err = nil

	if canPutBack {
		messages.Put(msg)
	}
}

var messages = sync.Pool{
	New: func() any {
		msg := new(Message)
		msg.init()
		return msg
	},
}
