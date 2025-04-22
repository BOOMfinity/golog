package golog

import (
	"fmt"
	"sync"
)

type Logger struct {
	mut            sync.RWMutex
	level          Level
	modules        []string
	params         []Parameter
	engine         WriteEngine
	dateTimeFormat string
}

func (l *Logger) SetDateTimeFormat(format string) {
	l.mut.Lock()
	l.dateTimeFormat = format
	l.mut.Unlock()
}

func (l *Logger) DateTimeFormat(fallback ...string) string {
	l.mut.RLock()
	defer l.mut.RUnlock()
	if l.dateTimeFormat == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return "02.01.2006 15:04:05"
	}
	format := l.dateTimeFormat
	return format
}

func (l *Logger) Modules() []string {
	return l.modules
}

func (l *Logger) Params() []Parameter {
	l.mut.RLock()
	params := l.params
	l.mut.RUnlock()
	return params
}

func (l *Logger) ClearParams() *Logger {
	l.mut.Lock()
	l.params = []Parameter{}
	l.mut.Unlock()
	return l
}

func (l *Logger) Param(name string, value any) *Logger {
	l.mut.Lock()
	params := make([]Parameter, len(l.params), len(l.params)+1)
	copy(params, l.params)
	params = append(params, Parameter{
		Name:  name,
		Value: value,
	})
	l.params = params
	l.mut.Unlock()
	return l
}

func (l *Logger) Level() Level {
	if overrideLevel := Config.OverrideMinimumMessageLevel(); overrideLevel != 0 {
		return overrideLevel
	}
	l.mut.RLock()
	lvl := l.level
	l.mut.RUnlock()
	return lvl
}

func (l *Logger) SetLevel(lvl Level) *Logger {
	l.mut.Lock()
	l.level = lvl
	l.mut.Unlock()
	return l
}

func (l *Logger) Module(name string, scope ...string) *Logger {
	cpy := l.doCopy()
	n := name
	if len(scope) > 0 {
		n = fmt.Sprintf("%s@%s", n, scope[0])
	}
	cpy.modules = append(cpy.modules, n)
	return cpy
}

func (l *Logger) doCopy(scope ...string) *Logger {
	l.mut.RLock()
	params := make([]Parameter, len(l.params))
	copy(params, l.params)
	modules := make([]string, len(l.modules))
	copy(modules, l.modules)
	if len(scope) > 0 && len(modules) > 0 {
		modules[len(l.modules)-1] = fmt.Sprintf("%s@%s", modules[len(l.modules)-1], scope[0])
	}
	cpy := &Logger{
		level:   l.level,
		modules: modules,
		params:  params,
		engine:  l.engine,
	}
	l.mut.RUnlock()
	return cpy
}

func (l *Logger) Copy(scope ...string) *Logger {
	return l.doCopy(scope...)
}

func (l *Logger) Info() Message {
	return newMessage(l, LevelInfo)
}

func (l *Logger) Debug() Message {
	return newMessage(l, LevelDebug)
}

func (l *Logger) Trace() Message {
	return newMessage(l, LevelTrace)
}

func (l *Logger) Warn() Message {
	return newMessage(l, LevelWarning)
}

func (l *Logger) Error() Message {
	return newErrorMessage(l)
}

func (l *Logger) Fatal(exitCode int) Message {
	return newErrorMessage(l, exitCode)
}

func New(name string, engine ...WriteEngine) *Logger {
	if len(engine) > 0 {
		return newLogger(name, engine[0])
	}
	return newLogger(name, ColorEngine())
}

func newLogger(name string, eng WriteEngine) *Logger {
	log := new(Logger)
	log.level = LevelInfo
	log.modules = append(log.modules, name)
	log.engine = eng
	return log
}
