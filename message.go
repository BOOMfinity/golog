package golog

import (
	"bytes"
	"fmt"
	"github.com/VenomPCPL/golog/internal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var colorsEnabled = true

var (
	debugCode = []byte("\u001B[3m\u001B[35m")
	infoCode  = []byte("\u001b[36m")
	warnCode  = []byte("\u001b[1m\u001b[33m")
	errorCode = []byte("\u001b[1m\u001b[31m")
	resetCode = []byte("\u001B[0m")
)

func internalRecover() {
	if err := recover(); err != nil {
		println("=================================")
		println()
		println("GOLOG - INTERNAL PANIC")
		println()
		println(fmt.Sprintf("%s", err))
		println()
		println("=================================")
	}
}

const messageBufferSize = 1024
const stackBufferSize = 2048
const timestampFormat = "02.01.2006 15:04:05"

var (
	_messageBuffPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, messageBufferSize))
		},
	}
	_stackBuffPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, stackBufferSize))
		},
	}
	_messagePool = sync.Pool{
		New: func() any {
			return &message{
				buff:      make([]byte, 0, messageBufferSize),
				stack:     make([]byte, stackBufferSize),
				offset:    new(atomic.Uint64),
				arguments: make([]string, 0, 16),
			}
		},
	}
)

func getMessage(log Logger, level Level) *message {
	msg := _messagePool.Get().(*message)
	msg.empty = level > log.Level()
	msg.instance = log
	msg.level = level
	msg.buff = msg.buff[:0]
	msg.sendStack = false
	msg.stack = msg.stack[:cap(msg.stack)]
	msg.arguments = msg.arguments[:0]
	msg.offset.Store(0)
	msg.userMessage = ""
	msg.time = time.Now()
	msg.err = nil
	return msg
}

func freeMessage(msg *message) {
	if cap(msg.buff) > messageBufferSize {
		msg.buff = nil
		msg.buff = make([]byte, 0, messageBufferSize)
	}
	if cap(msg.stack) > stackBufferSize {
		msg.stack = nil
		msg.stack = make([]byte, stackBufferSize)
	}
	_messagePool.Put(msg)
}

type message struct {
	empty       bool
	instance    Logger
	level       Level
	arguments   []string
	userMessage string
	time        time.Time
	stack       []byte
	buff        []byte
	offset      *atomic.Uint64
	sendStack   bool
	err         error
}

func (m *message) Error() error {
	return m.err
}

func (m *message) SendError(err error) {
	m.err = err
	m.Stack().Send("%s", err)
}

func (m *message) Stack() Message {
	m.stack = m.GetStack()
	m.sendStack = true
	return m
}

func (m *message) GetStack() []byte {
	m.stack = m.stack[:cap(m.stack)]
	written := runtime.Stack(m.stack, false)
	m.stack = m.stack[:written]
	return m.stack
}

func (m *message) FileWithLine() Message {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		m.Add("%v:%v", file, line)
	}
	return m
}

func (m *message) Instance() Logger {
	return m.instance
}

func (m *message) Level() Level {
	return m.level
}

func (m *message) Arguments() []string {
	return m.arguments
}

func (m *message) UserMessage() string {
	return m.userMessage
}

func (m *message) Time() time.Time {
	return m.time
}

func (m *message) Use(hook string, arg any) Message {
	if m.empty {
		return m
	}
	defer internalRecover()
	if fn := m.Instance().Hook(hook); fn != nil {
		fn(m, arg)
	}
	return m
}

func (m *message) insSep() {
	m.buff = append(m.buff, ' ', '|', ' ')
}

func (m *message) Any(args ...any) Message {
	if m.empty {
		return m
	}
	for i := range args {
		m.arguments = append(m.arguments, fmt.Sprint(args[i]))
	}
	return m
}

func (m *message) Add(format string, args ...any) Message {
	if m.empty {
		return m
	}
	m.arguments = append(m.arguments, fmt.Sprintf(format, args...))
	return m
}

func (m *message) Send(format string, args ...any) {
	defer internalRecover()
	defer freeMessage(m)
	if m.empty {
		return
	}
	if colorsEnabled {
		switch m.level {
		case LevelInfo:
			m.buff = append(m.buff, infoCode...)
		case LevelDebug:
			m.buff = append(m.buff, debugCode...)
		case LevelError:
			m.buff = append(m.buff, errorCode...)
		case LevelWarn:
			m.buff = append(m.buff, warnCode...)
		}
	}
	m.buff = m.time.AppendFormat(m.buff, timestampFormat)
	m.insSep()
	m.buff = append(m.buff, internal.GetBytes(m.level.String())...)
	m.insSep()
	for i := range m.Instance().Modules() {
		if i != 0 {
			m.buff = append(m.buff, ' ')
		}
		m.buff = append(m.buff, internal.GetBytes(m.Instance().Modules()[i])...)
	}
	if len(m.arguments) != 0 {
		m.insSep()
		for i := range m.arguments {
			if i != 0 {
				m.insSep()
			}
			m.buff = append(m.buff, internal.GetBytes(m.arguments[i])...)
		}
	}
	m.buff = append(m.buff, ' ', '-', '>', ' ')
	offset := len(m.buff)
	m.buff = fmt.Appendf(m.buff, format, args...)
	m.userMessage = internal.ToString(m.buff[offset:len(m.buff)])
	m.buff = append(m.buff, resetCode...)
	if m.Instance().Writer() != nil {
		m.buff = append(m.buff, '\r', '\n')
		_, _ = m.Instance().Writer().Write(m.buff)
		if m.sendStack {
			if cap(m.stack)+2 > stackBufferSize {
				m.stack = append(m.stack, '\r', '\n')
			} else {
				m.stack[len(m.stack)] = '\r'
				m.stack[len(m.stack)+1] = '\n'
			}
			_, _ = m.Instance().Writer().Write(m.stack)
		}
	}
	m.Instance().OnLog(m)
}
