// Package golog is simple (but fast), zero allocation logger with a few extra things like hooks, colors
package golog

import (
	"io"
	"os"
	"runtime"
	"sync"
)

var (
	forcedDebugMode = func() bool {
		return os.Getenv("GDEBUG") == "on"
	}
)

// Level determines which messages will be sent to the output
type Level uint8

const (
	LevelDebug Level = iota + 1
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type HookExecutor func(m *Message, arg interface{})
type WriteHookExecutor func(m *Message, msg []byte, userMsg []byte)

type Logger struct {
	level      Level
	writer     io.Writer
	modules    []string
	hooks      []HookExecutor
	writeHooks []WriteHookExecutor
	namedHooks *sync.Map
}

// WriteHook
//
// Don't use them if you want to be still fast as f.
func (l *Logger) WriteHook(fn WriteHookExecutor) *Logger {
	l.writeHooks = append(l.writeHooks, fn)
	return l
}

// ClearWriteHooks
//
// Just read Logger.ClearHooks and replace "global hooks" with "write hooks"
func (l *Logger) ClearWriteHooks() *Logger {
	l.writeHooks = l.writeHooks[:0]
	return l
}

// Info sends log message with "INFO" prefix
func (l *Logger) Info() *Message {
	return newMessage(l, LevelInfo)
}

// Warn sends log message with "WARN" prefix
func (l *Logger) Warn() *Message {
	return newMessage(l, LevelWarn)
}

// Debug sends log message with "DEBUG" prefix
func (l *Logger) Debug() *Message {
	return newMessage(l, LevelDebug)
}

// Error sends log message with "ERROR" prefix
func (l *Logger) Error() *Message {
	return newMessage(l, LevelError)
}

// NamedHook can be used with Message.Use
func (l *Logger) NamedHook(name string, fn HookExecutor) *Logger {
	l.namedHooks.Store(name, fn)
	return l
}

// RemoveNamedHook removes named hook by its name. Yes. Amazing, isn't it?
func (l *Logger) RemoveNamedHook(name string) *Logger {
	l.namedHooks.Delete(name)
	return l
}

// GlobalHook will be added to each log message
func (l *Logger) GlobalHook(fn HookExecutor) *Logger {
	l.hooks = append(l.hooks, fn)
	return l
}

// ClearHooks deletes all global hooks. Wow
func (l *Logger) ClearHooks() *Logger {
	l.hooks = l.hooks[:0]
	return l
}

var empty = []byte("\r\n")

// EmptyLine prints an empty line
func (l *Logger) EmptyLine() {
	_, _ = l.writer.Write(empty)
}

// SetWriter allows you to specify where the logs will be delivered
func (l *Logger) SetWriter(w io.Writer) *Logger {
	l.writer = w
	return l
}

// SetLevel sets the minimum level of log messages to be sent to io.Writer
func (l *Logger) SetLevel(lv Level) *Logger {
	l.level = lv
	return l
}

// Module allows you to create new Logger instance, BUT adds module with the given name
//
// Format: <timestamp> | <level> | <module 1> <module 2> <module 3> ...
func (l *Logger) Module(name string) *Logger {
	nl := &Logger{
		modules:    append(append([]string{}, l.modules...), []string{name}...),
		hooks:      append([]HookExecutor{}, l.hooks...),
		writeHooks: append([]WriteHookExecutor{}, l.writeHooks...),
		level:      l.level,
		writer:     l.writer,
		namedHooks: l.namedHooks,
	}
	return nl
}

// Stack returns stack of current goroutine
func (l *Logger) Stack() (data []byte, ok bool) {
	data = make([]byte, 1024*5)
	ok = runtime.Stack(data, false) != 0
	return
}

// Recover function can be used for default panic handling
//
// Uses predefined template that cannot be changed
func (l *Logger) Recover() {
	if v := recover(); v != nil {
		l.Error().FileWithLine().Stack().Send("%v", v)
	}
}

// NewCustomLogger allows you to create FULLY CUSTOM L O G G E R, including name, level, AND WRITER
func NewCustomLogger(name string, level Level, writer io.Writer) *Logger {
	l := new(Logger)
	l.namedHooks = new(sync.Map)
	l.level = level
	l.writer = writer
	l.modules = []string{name}
	return l
}

// NewLoggerWithLevel same as NewLogger but with custom logging level
func NewLoggerWithLevel(name string, level Level) *Logger {
	return NewCustomLogger(name, level, os.Stdout)
}

// NewLogger creates Logger instance with custom name but default (LevelInfo) logging level
func NewLogger(name string) *Logger {
	return NewLoggerWithLevel(name, LevelInfo)
}

// NewDefaultLogger gives you the easiest way to get new logger, but all properties are default
func NewDefaultLogger() *Logger {
	return NewLogger("main")
}
