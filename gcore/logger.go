package gcore

import (
	"context"
	"fmt"
	"sync"
)

var (
	PooledBufferSize                    = 1024 * 5
	PreAllocatedStackBufferSize         = 1024
	PreAllocatedMessageBufferSize       = 1024
	PreAllocatedSliceSize               = 32
	OverrideLoggingLevel          Level = 0
)

// Attribute represents a key-value pair.
type Attribute struct {
	Name  string
	Value any
}

// Hook is executed on every log entry, before the [Formatter].
// The minimum logging level is ignored.
//
// Hooks are executed synchronously, so they should be as lightweight as possible.
type Hook func(entry Context)

// Formatter is a [Hook] that prints the log entries.
type Formatter Hook

// Multiple combines multiple [Formatter](s) into one.
func Multiple(formats ...Formatter) Formatter {
	return func(entry Context) {
		for _, f := range formats {
			f(entry)
		}
	}
}

// Logger holds basic information and configuration for all logs.
//
// It is designed to be immutable, every change in the configuration returns a deeply cloned instance.
// This makes it safe for concurrent use.
type Logger struct {
	modules      []string
	scopes       []string
	attrs        []Attribute
	hooks        []Hook
	ctx          context.Context
	level        Level
	formatter    Formatter
	stackOnError bool
}

func (l *Logger) cmr(fn func(cpy *Logger)) *Logger {
	cpy := l.Clone()
	fn(cpy)
	return cpy
}

// WithFormat returns clone of the [Logger] with a changed [Formatter].
func (l *Logger) WithFormat(f Formatter) *Logger {
	return l.cmr(func(cpy *Logger) {
		cpy.formatter = f
	})
}

// WithHook returns clone of the [Logger] with [Hook](s) added.
// Hooks are always executed before the [Formatter] and ignore the minimum severity [Level].
func (l *Logger) WithHook(hook ...Hook) *Logger {
	return l.cmr(func(cpy *Logger) {
		cpy.hooks = append(cpy.hooks, hook...)
	})
}

// WithAutoStack returns clone of the [Logger] with automatic stack traces configured.
// If enabled, it will attach stack traces to the [LevelError] and [LevelFatal] severity levels automatically.
func (l *Logger) WithAutoStack(enabled bool) *Logger {
	return l.cmr(func(cpy *Logger) {
		cpy.stackOnError = enabled
	})
}

// WithLevel returns clone of the [Logger] with a modified minimum severity [Level].
// Only logs with equal or higher severity will be processed.
func (l *Logger) WithLevel(level Level) *Logger {
	return l.cmr(func(cpy *Logger) {
		cpy.level = level
	})
}

// WithContext returns clone of the [Logger] with the context attached.
// This context will be available in [Hook](s) and the [Formatter].
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return l.cmr(func(cpy *Logger) {
		cpy.ctx = ctx
	})
}

// WithAttribute returns clone of the [Logger] with a new [Attribute].
func (l *Logger) WithAttribute(name string, value any) *Logger {
	return l.cmr(func(cpy *Logger) {
		var attr Attribute
		attr.Name = name
		attr.Value = value
		cpy.attrs = append(cpy.attrs, attr)
	})
}

// WithScope returns clone of the [Logger] with a scope assigned to the current module.
func (l *Logger) WithScope(scope string) *Logger {
	return l.cmr(func(cpy *Logger) {
		cpy.scopes[len(cpy.scopes)-1] = scope
	})
}

// Recover is a helper that catches panics and logs them automatically.
//
// It MUST be used with `defer` to work correctly.
func (l *Logger) Recover(c ...context.Context) {
	if r := recover(); r != nil {
		msg := l.Error(c...)
		if err, ok := r.(error); ok {
			msg.MsgErr(fmt.Errorf("recovered: %w", err))
		} else {
			msg.Msgf("recovered: %v", r)
		}
	}
}

// Fatal creates [Message] with [LevelFatal] severity.
func (l *Logger) Fatal(code int, c ...context.Context) *Message {
	msg := acquireMessage(l, LevelFatal, c...)
	if l.stackOnError {
		msg.Stack()
	}
	msg.exitCode = code
	return msg
}

// Error creates [Message] with [LevelError] severity.
func (l *Logger) Error(c ...context.Context) *Message {
	msg := acquireMessage(l, LevelError, c...)
	if l.stackOnError {
		msg.Stack()
	}
	return msg
}

// Scoped creates a clone of the logger with a new scope,
// then runs fn with that sub-logger.
//
// Note: The sub-logger is returned to the pool after fn returns.
func (l *Logger) Scoped(name string, fn func(log *Logger)) {
	if name == "" {
		panic("scope must not be empty")
	}
	scoped := loggers.Get().(*Logger)
	defer freeLogger(scoped)
	scoped.attrs = append(scoped.attrs, l.attrs...)
	scoped.modules = append(scoped.modules, l.modules...)
	scoped.scopes = append(scoped.scopes, l.scopes...)
	scoped.ctx = l.ctx
	scoped.level = l.level
	scoped.stackOnError = l.stackOnError
	scoped.formatter = l.formatter
	scoped.hooks = append(scoped.hooks, l.hooks...)
	scoped.scopes[len(scoped.scopes)-1] = name
	fn(scoped)
}

// Warn creates [Message] with [LevelWarning] severity.
func (l *Logger) Warn(c ...context.Context) *Message {
	return acquireMessage(l, LevelWarning, c...)
}

// Info creates [Message] with [LevelInfo] severity.
func (l *Logger) Info(c ...context.Context) *Message {
	return acquireMessage(l, LevelInfo, c...)
}

// Debug creates [Message] with [LevelDebug] severity.
func (l *Logger) Debug(c ...context.Context) *Message {
	return acquireMessage(l, LevelDebug, c...)
}

// Trace creates [Message] with [LevelTrace] severity.
func (l *Logger) Trace(c ...context.Context) *Message {
	return acquireMessage(l, LevelTrace, c...)
}

// Module creates a sub-logger by appending a module and scope to the current [Logger].
func (l *Logger) Module(name string, scope ...string) *Logger {
	if name == "" {
		panic("module name must not be empty")
	}
	cpy := l.Clone()
	cpy.modules = make([]string, len(cpy.modules)+1)
	copy(cpy.modules, l.modules)
	cpy.modules[len(cpy.modules)-1] = name
	cpy.scopes = make([]string, len(cpy.scopes)+1)
	copy(cpy.scopes, l.scopes)
	if len(scope) > 0 {
		cpy = cpy.WithScope(scope[0])
	}
	return cpy
}

// Clone returns a deeply cloned [Logger] instance.
//
// If scope is provided, it is set on the cloned instance - equivalent to Clone().WithScope(name).
func (l *Logger) Clone(scope ...string) *Logger {
	log := new(Logger)
	log.ctx = l.ctx
	log.level = l.level
	log.formatter = l.formatter
	log.stackOnError = l.stackOnError
	log.hooks = make([]Hook, len(l.hooks))
	copy(log.hooks, l.hooks)
	log.attrs = make([]Attribute, len(l.attrs))
	copy(log.attrs, l.attrs)
	log.modules = make([]string, len(l.modules))
	copy(log.modules, l.modules)
	log.scopes = make([]string, len(l.scopes))
	copy(log.scopes, l.scopes)
	if len(scope) > 0 {
		log = log.WithScope(scope[0])
	}
	return log
}

// Init returns a new [Logger] with specific name and [Formatter].
//
// All parameters for this function are required.
//
// By default, it sets the minimum logging level to [LevelInfo].
func Init(name string, f Formatter) *Logger {
	if f == nil {
		panic("f (Formatter) must not be nil")
	}
	if name == "" {
		name = "app"
	}
	log := &Logger{}
	log.modules = make([]string, 1)
	log.scopes = make([]string, 1)
	log.modules[0] = name
	log.level = LevelInfo
	log.stackOnError = true
	log.formatter = f
	return log
}

func freeLogger(log *Logger) {
	log.ctx = nil
	log.formatter = nil
	log.level = 0
	log.stackOnError = true
	clear(log.attrs)
	log.attrs = log.attrs[:0]
	clear(log.modules)
	log.modules = log.modules[:0]
	clear(log.scopes)
	log.scopes = log.scopes[:0]
	clear(log.hooks)
	log.hooks = log.hooks[:0]
	if cap(log.attrs) > PreAllocatedSliceSize || cap(log.modules) > PreAllocatedSliceSize || cap(log.scopes) > PreAllocatedSliceSize || cap(log.hooks) > PreAllocatedSliceSize {
		return
	}
	loggers.Put(log)
}

var loggers = sync.Pool{
	New: func() any {
		return new(Logger{
			modules: make([]string, 0, PreAllocatedSliceSize),
			scopes:  make([]string, 0, PreAllocatedSliceSize),
			attrs:   make([]Attribute, 0, PreAllocatedSliceSize),
			hooks:   make([]Hook, 0, PreAllocatedSliceSize),
		})
	},
}
