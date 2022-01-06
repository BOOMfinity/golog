package golog

import (
	"fmt"
	"sync"
	"time"
)

type nullMessage struct{}

func (n *nullMessage) Str(str string) Message {
	return n
}

func (n *nullMessage) Any(v interface{}) Message {
	return n
}

func (n *nullMessage) Level() Level {
	return 0
}

func (n *nullMessage) Fmt(format string, values ...interface{}) Message {
	return n
}

func (n *nullMessage) Use(name string, arg interface{}) Message {
	return n
}

func (n *nullMessage) NoWriteHooks() Message {
	return n
}

func (n *nullMessage) Send(format string, values ...interface{}) {}

type Message interface {
	// Str adds string to the output
	Str(str string) Message
	// Any accepts any type and adds it to the output
	Any(v interface{}) Message
	// Fmt adds formatted string to the output
	Fmt(format string, values ...interface{}) Message
	// Use executes named hook declared with Logger.NamedHook
	Use(name string, arg interface{}) Message
	// Send writes output to the Logger io.Writer
	Send(format string, values ...interface{})
	// Level returns level of the log message
	Level() Level
	// NoWriteHooks disables write hooks for this Message
	NoWriteHooks() Message
}

var globalNullMessage = &nullMessage{}

var messagePool = &sync.Pool{
	New: func() interface{} {
		buf := make(WritableBuffer, 0, 512)
		rawBuf := make(WritableBuffer, 0, 1024)
		return &message{buf: &buf, rawBuf: &rawBuf}
	},
}

func putMessage(x *message) {
	if x.buf.Len() > 1024 {
		return
	}
	messagePool.Put(x)
}

func newMessage(logger *logger, level Level) *message {
	msg := messagePool.Get().(*message)
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

type message struct {
	logger       *logger
	buf          *WritableBuffer
	rawBuf       *WritableBuffer
	level        Level
	sent         bool
	noWriteHooks bool
}

func (m *message) NoWriteHooks() Message {
	m.noWriteHooks = true
	return m
}

func (m *message) Use(name string, arg interface{}) Message {
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

func (m *message) Str(str string) Message {
	*m.buf = appendType(m.buf, " | ")
	*m.buf = append(*m.buf, unsafeBytes(str)...)
	return m
}

func (m *message) Level() Level {
	return m.level
}

func (m *message) Any(v interface{}) Message {
	return m.Fmt("%v", v)
}

func (m *message) Fmt(format string, values ...interface{}) Message {
	*m.buf = appendType(m.buf, " | ")
	_, _ = fmt.Fprintf(m.buf, format, values...)
	return m
}

func (m *message) Send(format string, values ...interface{}) {
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
