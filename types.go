package golog

import (
	"io"
	"time"
)

type Level uint8

const (
	LevelPanic Level = iota + 1
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

func (l Level) String() string {
	switch l {
	case LevelPanic:
		return "PANIC"
	case LevelError:
		return "ERROR"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

type Logger interface {
	Recover()

	Info() Message
	Warn() Message
	Error() Message
	Debug() Message
	// Panic is a special log message type with extra decorations.
	//
	// By default, it doesn't exit process.
	Panic() FatalMessage
	// Fatal is wrapper for Logger.Panic that exits process after sending log message.
	//
	// Default exit code is 1.
	Fatal() FatalMessage
	Empty()

	Module(name string) Logger
	Modules() []string
	SetLevel(level Level) Logger
	Level() Level
	SetWriter(wr io.Writer) Logger
	Writer() io.Writer

	AddOnLog(id string, fn LogHandler) Logger
	OnLog(msg Message)

	ClearHooks() Logger
	ClearAll() Logger
	ClearHandlers() Logger

	CreateHook(name string, fn HookHandler)
	Hook(name string) HookHandler
}

type HookHandler func(msg Message, arg any)
type LogHandler func(msg Message)

type Message interface {
	Instance() Logger
	Level() Level
	Arguments() []string
	UserMessage() string
	Time() time.Time
	Error() error
	GetExitCode() int

	Stack() Message
	GetStack() []byte
	FileWithLine() Message
	SendError(err error)

	Use(hook string, arg any) Message
	Any(arg ...any) Message
	Add(format string, args ...any) Message
	Send(format string, args ...any)

	fatal() FatalMessage
}

type FatalMessage interface {
	Message
	ExitCode(code int) Message
}
