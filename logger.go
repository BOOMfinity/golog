package golog

import (
	"strings"
	"sync"

	"github.com/BOOMfinity/go-utils/inlineif"
)

type LoggerInternal interface {
	Params() Params
	Level() Level
	Module() string
	Scope() string
	Parent() Logger
	Engine() WriteEngine
}

type Logger interface {
	Internal() LoggerInfo
	Param(name string, value any) Logger
	Level(l Level) Logger
	// Module creates Copy of Logger instance with given module name.
	//
	// It also resets scope and params.
	Module(name string) Logger
	// Scope helps you find where in code the log message was sent.
	//
	// It does not create Copy of Logger instance, so you can reuse the same log instance or create Copy manually.
	Scope(name string) Logger
	Copy() Logger
	ResetScope() Logger
	ResetParams() Logger
	// Recover should be used with defer. It prints error message when code panics and prevents program from exiting.
	//
	// You can use RecoverParams to customize method's behavior.
	Recover(params ...RecoverParams)

	Info() Message
	Debug() Message
	Warn() Message
	Error() Message
	Trace() Message
	Fatal(exitCode int) Message
}

type LoggerInfo struct {
	Level  Level
	Module string
	Scope  string
	Parent Logger
	Params Params
	Engine WriteEngine
}

type loggerImpl struct {
	eng  WriteEngine
	info LoggerInfo
}

func (l *loggerImpl) ResetScope() Logger {
	l.info.Scope = ""
	return l
}

func (l *loggerImpl) ResetParams() Logger {
	for i := range l.info.Params {
		l.info.Params[i] = nil
	}
	l.info.Params = l.info.Params[:0]
	return l
}

func (l *loggerImpl) Internal() LoggerInfo {
	return l.info
}

func (l *loggerImpl) Param(name string, value any) Logger {
	arr := &[2]any{}
	(*arr)[0] = name
	(*arr)[1] = value
	l.info.Params = append(l.info.Params, arr)
	return l
}

func (l *loggerImpl) Level(lvl Level) Logger {
	l.info.Level = inlineif.IfElse(overrideLoggingLevel != "", levelFromString(overrideLoggingLevel), lvl)
	return l
}

func (l *loggerImpl) Module(name string) Logger {
	cpy := l.doCopy()
	cpy.info.Module = name
	cpy.info.Scope = ""
	cpy.info.Parent = l
	cpy.info.Params = make(Params, 0, 25)
	return cpy
}

func (l *loggerImpl) Scope(name string) Logger {
	l.info.Scope = name
	return l
}

func (l *loggerImpl) doCopy() *loggerImpl {
	dst := make(Params, len(l.info.Params), 25)
	copy(dst, l.info.Params)
	return &loggerImpl{
		eng: l.eng,
		info: LoggerInfo{
			Level:  l.info.Level,
			Module: strings.Clone(l.info.Module),
			Scope:  strings.Clone(l.info.Scope),
			Parent: l.info.Parent,
			Params: dst,
			Engine: l.info.Engine,
		},
	}
}

func (l *loggerImpl) Copy() Logger {
	return l.doCopy()
}

func (l *loggerImpl) Info() Message {
	return newMessage(l, LevelInfo)
}

func (l *loggerImpl) Debug() Message {
	return newMessage(l, LevelDebug)
}

func (l *loggerImpl) Trace() Message {
	return newMessage(l, LevelTrace)
}

func (l *loggerImpl) Warn() Message {
	return newMessage(l, LevelWarning)
}

func (l *loggerImpl) Error() Message {
	return newErrorMessage(l)
}

func (l *loggerImpl) Fatal(exitCode int) Message {
	return newErrorMessage(l, exitCode)
}

type Pool interface {
	Get() Logger
	Put(l Logger)
}

type poolImpl struct {
	loggers *sync.Pool
}

func (s poolImpl) Get() Logger {
	return s.loggers.Get().(Logger)
}

func (s poolImpl) Put(l Logger) {
	l.ResetScope()
	l.ResetParams()
	s.loggers.Put(l)
}

func New(name string) Logger {
	return newLogger(name, NewColorEngine())
}

func newLogger(name string, eng WriteEngine) *loggerImpl {
	log := new(loggerImpl)
	log.info = LoggerInfo{}
	log.info.Level = inlineif.IfElse(overrideLoggingLevel != "", levelFromString(overrideLoggingLevel), LevelInfo)
	log.info.Module = name
	log.info.Parent = nil
	log.info.Params = make(Params, 0, 25)
	log.info.Engine = eng
	return log
}

func NewCustom(name string, eng WriteEngine) Logger {
	return newLogger(name, eng)
}

func NewPool(base Logger) Pool {
	return poolImpl{
		loggers: &sync.Pool{
			New: func() any {
				return base.Copy()
			},
		},
	}
}
