package golog

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var buffPool = &sync.Pool{
	New: func() interface{} {
		d := make([]byte, 2048)
		return &d
	},
}

var messagePool = &sync.Pool{
	New: func() interface{} {
		buf := make(WritableBuffer, 0, 512)
		rawBuf := make(WritableBuffer, 0, 1024)
		return &Message{buf: &buf, rawBuf: &rawBuf}
	},
}

func putMessage(x *Message) {
	if x.buf.Len() > 1024 {
		return
	}
	buffPool.Put(x.stack)
	x.stack = nil
	messagePool.Put(x)
}

func newBuffer() *[]byte {
	return buffPool.Get().(*[]byte)
}

func newMessage(logger *Logger, level Level) *Message {
	msg := messagePool.Get().(*Message)
	msg.empty = level < logger.level
	msg.sent = false
	msg.stack = nil
	msg.buf.Reset()
	msg.logger = logger
	msg.level = level
	msg.rawBuf.Reset()
	*msg.rawBuf = appendColors(*msg.rawBuf, msg.level)
	for _, hook := range logger.hooks {
		func() {
			defer func() {
				if err := recover(); err != nil {
					println()
					println("=======================")
					println()
					println("Recovered panic: Hooks")
					println()
					println(fmt.Print(err))
					println()
					println()
				}
			}()
			hook(msg, nil)
		}()
	}
	return msg
}

type Message struct {
	empty        bool
	logger       *Logger
	buf          *WritableBuffer
	rawBuf       *WritableBuffer
	stack        *[]byte
	level        Level
	sent         bool
	sendStack    bool
	noWriteHooks bool
}

// FileWithLine adds file and line where THIS function was used
func (m *Message) FileWithLine() *Message {
	_, file, line, ok := runtime.Caller(4)
	if ok {
		m.Fmt("%v:%v", file, line)
	}
	return m
}

func (m *Message) GetStack() []byte {
	if len(*m.stack) == 0 {
		m.allocStack()
	}
	return *m.stack
}

func (m *Message) allocStack() {
	buf := newBuffer()
	m.stack = buf
	runtime.Stack(*buf, false)
}

// Stack prints stack trace of current goroutine
//
// Printed stack will be sent as uncolored text under log message
func (m *Message) Stack() *Message {
	m.allocStack()
	m.sendStack = true
	return m
}

// NoWriteHooks disables write hooks for this Message
func (m *Message) NoWriteHooks() *Message {
	if m.empty {
		return m
	}
	m.noWriteHooks = true
	return m
}

// Use executes named hook declared with Logger.NamedHook
func (m *Message) Use(name string, arg interface{}) *Message {
	if m.empty {
		return m
	}
	if _hook, ok := m.logger.namedHooks.Load(name); ok {
		defer func() {
			if err := recover(); err != nil {
				println()
				println("=======================")
				println()
				println("Recovered panic: Hooks")
				println()
				println(fmt.Print(err))
				println()
				println()
			}
		}()
		hook := _hook.(HookExecutor)
		hook(m, arg)
	}
	return m
}

// Str adds string to the output
func (m *Message) Str(str string) *Message {
	if m.empty {
		return m
	}
	*m.buf = appendType(m.buf, " | ")
	*m.buf = append(*m.buf, unsafeBytes(str)...)
	return m
}

// Level returns level of the log message
func (m *Message) Level() Level {
	return m.level
}

// Any accepts any type and adds it to the output
func (m *Message) Any(v interface{}) *Message {
	if m.empty {
		return m
	}
	return m.Fmt("%v", v)
}

// Fmt adds formatted string to the output
func (m *Message) Fmt(format string, values ...interface{}) *Message {
	if m.empty {
		return m
	}
	*m.buf = appendType(m.buf, " | ")
	_, _ = fmt.Fprintf(m.buf, format, values...)
	return m
}

// Send writes output to the Logger io.Writer
func (m *Message) Send(format string, values ...interface{}) {
	if m.empty {
		return
	}
	if m.sent {
		panic("You cannot use the same message type many times")
	}
	defer putMessage(m)
	*m.rawBuf = appendTime(m.rawBuf, time.Now(), "2006-01-02 15:04:05")
	*m.rawBuf = appendType(m.rawBuf, " | ")
	*m.rawBuf = appendLevel(*m.rawBuf, m.level)
	*m.rawBuf = appendType(m.rawBuf, " | ")
	for i := range m.logger.modules {
		*m.rawBuf = append(*m.rawBuf, unsafeBytes(m.logger.modules[i])...)
		if i < len(m.logger.modules)-1 {
			*m.rawBuf = appendType(m.rawBuf, " ")
		}
	}
	*m.rawBuf = append(*m.rawBuf, (*m.buf)...)
	*m.rawBuf = appendType(m.rawBuf, " -> ")
	fmt.Fprintf(m.rawBuf, format, values...)
	if !m.noWriteHooks && len(m.logger.writeHooks) > 0 {
		strBuf := make(WritableBuffer, 0, 1024)
		fmt.Fprintf(&strBuf, format, values...)
		func() {
			for _, hook := range m.logger.writeHooks {
				hook(m, *m.rawBuf, strBuf)
			}
		}()
	}
	*m.rawBuf = appendReset(*m.rawBuf)
	*m.rawBuf = append(*m.rawBuf, []byte("\r\n")...)
	if m.logger.writer != nil {
		m.logger.writer.Write(*m.rawBuf)
		if m.sendStack && len(*m.stack) > 0 {
			m.logger.writer.Write(append(*m.stack, []byte("\r\n")...))
		}
	}
	m.sent = true
}
