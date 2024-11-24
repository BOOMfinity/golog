package golog

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/BOOMfinity/go-utils/gpool"
	"github.com/BOOMfinity/go-utils/inlineif"
)

type Level uint

const (
	LevelPanic Level = iota + 1
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace
)

func (l Level) String() string {
	switch l {
	case LevelPanic:
		return "PANIC"
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
	default:
		panic("undefined logging level")
	}
}

func (l Level) Color() string {
	switch l {
	case LevelPanic:
		return errorColorCode
	case LevelError:
		return errorColorCode
	case LevelWarning:
		return warningColorCode
	case LevelInfo:
		return infoColorCode
	case LevelDebug:
		return debugColorCode
	case LevelTrace:
		return traceColorCode
	default:
		panic("undefined logging level")
	}
}

func levelFromString(str string) Level {
	switch strings.ToUpper(str) {
	case "PANIC":
		return LevelPanic
	case "ERROR":
		return LevelError
	case "WARNING":
		return LevelWarning
	case "INFO":
		return LevelInfo
	case "DEBUG":
		return LevelDebug
	case "TRACE":
		return LevelTrace
	default:
		return LevelInfo
	}
}

type Message interface {
	Internal() MessageInfo
	Details(v any) Message
	Param(name, value any) Message
	Duration(d time.Duration) Message
	Stack() Message

	Send(format string, args ...any)
	Throw(err error)
}

type MessageInfo struct {
	Details  any
	Logger   Logger
	Level    Level
	Stack    []byte
	ExitCode int
	Duration time.Duration
	Message  []byte
	Params   Params
}

type message struct {
	info MessageInfo
}

func (m *message) Stack() Message {
	if !isEmpty(m.info.Stack) {
		return m
	}
	runtime.Stack(m.info.Stack, false)
	return m
}

func (m *message) Throw(err error) {
	m.Stack().Send(err.Error())
}

func (m *message) Internal() MessageInfo {
	return m.info
}

func (m *message) Details(v any) Message {
	m.info.Details = v
	return m
}

func (m *message) Param(name, value any) Message {
	arr := paramsPool.Get()
	(*arr)[0] = name
	(*arr)[1] = value
	m.info.Params = append(m.info.Params, arr)
	return m
}

func (m *message) Duration(d time.Duration) Message {
	m.info.Duration = d
	return m
}

func (m *message) Send(format string, args ...any) {
	defer messagePool.Put(m)

	if m.Internal().Level > m.Internal().Logger.Internal().Level {
		return
	}
	m.info.Message = fmt.Appendf(m.info.Message, format, args...)
	m.info.Logger.Internal().Engine(m)
}

var messagePool = gpool.New[message](gpool.OnInit[message](func(m *message) {
	m.info.Stack = make([]byte, 1024*4)
	m.info.Params = make(Params, 0, 25)
	m.info.Message = make([]byte, 0, 1024*4)
}), gpool.OnPut[message](func(m *message) {
	m.info.Message = m.info.Message[:0]
	for _, param := range m.info.Params {
		paramsPool.Put(param)
	}
	m.info.Params = m.info.Params[:0]
	m.info.Details = nil
	m.info.ExitCode = 0
	clear(m.info.Stack)
	m.info.Duration = 0
}))

func newMessage(log Logger, level Level) Message {
	msg := messagePool.Get()
	msg.info.Level = level
	msg.info.Logger = log
	return msg
}

func newErrorMessage(log Logger, exitCode ...int) Message {
	msg := messagePool.Get()
	msg.info.Level = inlineif.IfElse(len(exitCode) > 0, LevelPanic, LevelError)
	if len(exitCode) > 0 {
		msg.info.ExitCode = exitCode[0]
	}
	msg.info.Logger = log
	return msg
}
