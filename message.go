package golog

import (
	"fmt"
	"os"
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
		return 0
	}
}

type MessageData struct {
	Details       any           `json:"details,omitempty"`
	Level         Level         `json:"level,omitempty"`
	Stack         []byte        `json:"stack,omitempty"`
	StackIncluded bool          `json:"-"`
	Error         error         `json:"-"`
	ExitCode      int           `json:"exit_code,omitempty"`
	Duration      time.Duration `json:"duration,omitempty"`
	Message       []byte        `json:"message,omitempty"`
	Params        []Parameter   `json:"params,omitempty"`
}

type Message struct {
	data   *MessageData
	parent *Logger
}

func (m Message) Stack() Message {
	l := runtime.Stack(m.data.Stack[:cap(m.data.Stack)], false)
	m.data.Stack = m.data.Stack[:l]
	m.data.StackIncluded = true
	return m
}

func (m Message) Throw(err error) {
	m.data.Error = err
	m.Stack().Send(err.Error())
}

func (m Message) Details(v any) Message {
	m.data.Details = v
	return m
}

func (m Message) Param(name string, value any) Message {
	m.data.Params = append(m.data.Params, Parameter{
		Name:  name,
		Value: value,
	})
	return m
}

func (m Message) Duration(d time.Duration) Message {
	m.data.Duration = d
	return m
}

func (m Message) Send(format string, args ...any) {
	defer dataPool.Put(m.data)
	if m.data.Level > m.parent.Level() {
		return
	}
	m.data.Message = fmt.Appendf(m.data.Message, format, args...)
	m.parent.engine(m.parent, m.data)
	if m.data.ExitCode != 0 {
		os.Exit(m.data.ExitCode)
	}
}

var dataPool = gpool.New[MessageData](gpool.OnInit[MessageData](func(m *MessageData) {
	m.Stack = make([]byte, Config.StackTraceBufferSize())
	m.Params = make([]Parameter, 0, Config.MessageParametersSliceAllocation())
	m.Message = make([]byte, 0, Config.MessageBufferSize())
}), gpool.OnPut[MessageData](func(m *MessageData) {
	clear(m.Stack)
	clear(m.Params)
	clear(m.Message)
	m.Stack = m.Stack[:Config.StackTraceBufferSize():Config.StackTraceBufferSize()]
	m.Params = m.Params[:0:Config.MessageParametersSliceAllocation()]
	m.Message = m.Message[:0:Config.MessageBufferSize()]
	m.Details = nil
	m.ExitCode = 0
	m.Duration = 0
	m.StackIncluded = false
}))

func newMessage(log *Logger, level Level) Message {
	data := dataPool.Get()
	data.Level = level
	return Message{
		data:   data,
		parent: log,
	}
}

func newErrorMessage(log *Logger, exitCode ...int) Message {
	message := newMessage(log, inlineif.IfElse(len(exitCode) > 0, LevelPanic, LevelError))
	if len(exitCode) > 0 {
		message.data.ExitCode = exitCode[0]
	}
	if message.data.Level == LevelPanic || Config.IncludeStackOnError() {
		message.Stack()
	}
	return message
}
