// Package golog is simple (but fast), zero allocation logger with a few extra things like hooks, colors
package golog

import (
	"io"
	"os"
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
	LevelFatal
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
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type HookExecutor func(m Message, arg interface{})
type WriteHookExecutor func(m Message, msg []byte, userMsg []byte)

type Logger interface {
	// GlobalHook will be added to each log message
	GlobalHook(fn HookExecutor) Logger
	// ClearHooks deletes all global hooks. Wow
	ClearHooks() Logger
	// NamedHook can be used with Message.Use
	NamedHook(name string, fn HookExecutor) Logger
	// RemoveNamedHook removes named hook by its name. Yes. Amazing, isn't it?
	RemoveNamedHook(name string) Logger
	// SetWriter allows you to specify where the logs will be delivered
	SetWriter(w io.Writer) Logger
	// SetLevel sets the minimum level of log messages to be sent to io.Writer
	SetLevel(lv Level) Logger
	// Module allows you to create new Logger instance, BUT adds module with the given name
	//
	// Format: <timestamp> | <level> | <module 1> <module 2> <module 3> ...
	Module(name string) Logger
	// WriteHook
	//
	// Don't use them if you want to be still fast as f.
	WriteHook(fn WriteHookExecutor) Logger
	// ClearWriteHooks
	//
	// Just read Logger.ClearHooks and replace "global hooks" with "write hooks"
	ClearWriteHooks() Logger

	// Info sends log message with "INFO" prefix
	Info() Message
	// Warn sends log message with "WARN" prefix
	Warn() Message
	// Debug sends log message with "DEBUG" prefix
	Debug() Message
	// Error sends log message with "ERROR" prefix
	Error() Message
	// Fatal sends log message with "FATAL" prefix and exits the program
	Fatal() Message
}

type logger struct {
	level      Level
	writer     io.Writer
	modules    []string
	hooks      []HookExecutor
	writeHooks []WriteHookExecutor
	namedHooks *sync.Map
}

func (l *logger) WriteHook(fn WriteHookExecutor) Logger {
	l.writeHooks = append(l.writeHooks, fn)
	return l
}

func (l *logger) ClearWriteHooks() Logger {
	l.writeHooks = l.writeHooks[:0]
	return l
}

func (l *logger) Info() Message {
	if l.level > LevelInfo {
		return globalNullMessage
	}
	return newMessage(l, LevelInfo)
}

func (l *logger) Warn() Message {
	if l.level > LevelWarn {
		return globalNullMessage
	}
	return newMessage(l, LevelWarn)
}

func (l *logger) Debug() Message {
	if l.level > LevelDebug && !forcedDebugMode() {
		return globalNullMessage
	}
	return newMessage(l, LevelDebug)
}

func (l *logger) Error() Message {
	if l.level > LevelError {
		return globalNullMessage
	}
	return newMessage(l, LevelError)
}

func (l *logger) Fatal() Message {
	defer os.Exit(1)
	if l.level > LevelFatal {
		return globalNullMessage
	}
	return newMessage(l, LevelFatal)
}

func (l *logger) ClearNamedHooks() Logger {
	l.namedHooks = new(sync.Map)
	return l
}

func (l *logger) NamedHook(name string, fn HookExecutor) Logger {
	l.namedHooks.Store(name, fn)
	return l
}

func (l *logger) RemoveNamedHook(name string) Logger {
	l.namedHooks.Delete(name)
	return l
}

func (l *logger) GlobalHook(fn HookExecutor) Logger {
	l.hooks = append(l.hooks, fn)
	return l
}

func (l *logger) ClearHooks() Logger {
	l.hooks = l.hooks[:0]
	return l
}

func (l *logger) SetWriter(w io.Writer) Logger {
	l.writer = w
	return l
}

func (l *logger) SetLevel(lv Level) Logger {
	l.level = lv
	return l
}

func (l *logger) Module(name string) Logger {
	nl := &logger{
		modules:    append(append([]string{}, l.modules...), []string{name}...),
		hooks:      append([]HookExecutor{}, l.hooks...),
		writeHooks: append([]WriteHookExecutor{}, l.writeHooks...),
		level:      l.level,
		writer:     l.writer,
		namedHooks: l.namedHooks,
	}
	return nl
}

// NewCustomLogger allows you to create FULLY CUSTOM L O G G E R, including name, level, AND WRITER
func NewCustomLogger(name string, level Level, writer io.Writer) Logger {
	l := new(logger)
	l.namedHooks = new(sync.Map)
	l.level = level
	l.writer = writer
	l.modules = []string{name}
	return l
}

// NewLoggerWithLevel same as NewLogger but with custom logging level
func NewLoggerWithLevel(name string, level Level) Logger {
	return NewCustomLogger(name, level, os.Stdout)
}

// NewLogger creates Logger instance with custom name but default (LevelInfo) logging level
func NewLogger(name string) Logger {
	return NewLoggerWithLevel(name, LevelInfo)
}

// NewDefaultLogger gives you the easiest way to get new logger, but all properties are default
func NewDefaultLogger() Logger {
	return NewLogger("main")
}
