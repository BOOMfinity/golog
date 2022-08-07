package golog

import (
	"io"
	"os"
)

type logger struct {
	writer      io.Writer
	modules     []string
	hooks       *map[string]HookHandler
	level       Level
	logHandlers *map[string]LogHandler
}

func (l *logger) Panic() FatalMessage {
	return getMessage(l, LevelPanic).Stack().fatal()
}

func (l *logger) Fatal() FatalMessage {
	return l.Panic().ExitCode(1).fatal()
}

func (l *logger) Empty() {
	_, _ = l.writer.Write([]byte("\r\n"))
}

func (l *logger) ClearHooks() Logger {
	l.hooks = new(map[string]HookHandler)
	return l
}

func (l *logger) ClearAll() Logger {
	l.ClearHooks()
	l.ClearHandlers()
	return l
}

func (l *logger) ClearHandlers() Logger {
	l.logHandlers = new(map[string]LogHandler)
	return l
}

func (l *logger) Level() Level {
	return l.level
}

func (l *logger) Writer() io.Writer {
	return l.writer
}

func (l *logger) Hook(name string) HookHandler {
	return (*l.hooks)[name]
}

func (l *logger) Recover() {
	if err := recover(); err != nil {
		l.Error().Any("!! PANIC !!").Stack().Send("%v", err)
	}
}

func (l *logger) Info() Message {
	return getMessage(l, LevelInfo)
}

func (l *logger) Warn() Message {
	return getMessage(l, LevelWarn)
}

func (l *logger) Error() Message {
	return getMessage(l, LevelError)
}

func (l *logger) Debug() Message {
	return getMessage(l, LevelDebug)
}

func (l *logger) Module(name string) Logger {
	return &logger{
		level:       l.level,
		writer:      l.writer,
		hooks:       l.hooks,
		logHandlers: l.logHandlers,
		modules:     append(l.modules, name),
	}
}

func (l *logger) Modules() []string {
	return l.modules
}

func (l *logger) SetLevel(level Level) Logger {
	l.level = level
	return l
}

func (l *logger) SetWriter(wr io.Writer) Logger {
	l.writer = wr
	return l
}

func (l *logger) OnLog(msg Message) {
	for i := range *l.logHandlers {
		(*l.logHandlers)[i](msg)
	}
}

func (l *logger) AddOnLog(id string, fn LogHandler) Logger {
	(*l.logHandlers)[id] = fn
	return l
}

func (l *logger) CreateHook(name string, fn HookHandler) {
	(*l.hooks)[name] = fn
}

func New(name string) Logger {
	return NewWithLevel(name, LevelInfo)
}

func NewWithLevel(name string, level Level) Logger {
	log := &logger{
		writer:      os.Stdout,
		hooks:       new(map[string]HookHandler),
		logHandlers: new(map[string]LogHandler),
		modules:     []string{name},
		level:       level,
	}
	*log.hooks = map[string]HookHandler{}
	*log.logHandlers = map[string]LogHandler{}
	return log
}

func NewDefault() Logger {
	return New("app")
}
