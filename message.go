package golog

import (
	"fmt"
	"sync"
	"time"
)

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
	messagePool.Put(x)
}

func newMessage(logger *Logger, level Level) *Message {
	msg := messagePool.Get().(*Message)
	msg.empty = level < logger.level
	msg.sent = false
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
	level        Level
	sent         bool
	noWriteHooks bool
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
	}
	m.sent = true
}
