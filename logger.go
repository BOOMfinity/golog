// Package golog is a basic logger with a few extra things like custom values / functions.
//
// Format:
// <time> | <level> | <name and modules> | <custom1> | <custom2> | <custom...> -> <message>
package golog

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
)

// Custom functions used by all loggers in the program.
var globalCustoms sync.Map

func AddGlobalCustom(name string, fn CustomHandler) {
	globalCustoms.Store(name, fn)
}

// Level defines the logging level.
//
// Logger will write the output to the console only if its Level is higher or equal to the current Level of Option.
type Level uint8

const (
	Fatal Level = iota
	Error
	Warning
	Info
	Debug
)

func (x Level) String() string {
	switch x {
	case Fatal:
		return "[**] FATAL [**]"
	case Error:
		return "** ERROR **"
	case Warning:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "* DEBUG *"
	}
	return ""
}

var (
	defaultTemplate, _ = template.New(`golog`).
		Funcs(template.FuncMap{
			"join": func(a []string, b string) string {
				return strings.Join(a, b)
			},
			"sprintf": func(format string, values []interface{}) string {
				return fmt.Sprintf(format, values...)
			},
			"modulesJoin": func(customs []string) (res string) {
				for i := range customs {
					res += "| " + customs[i]
				}
				return
			},
		}).
		Parse(`{{ .time.Format "2006-01-02 15:04:05 (Z07:00)" }} | {{ .type }} | {{ .name }}{{if ne (len .modules) 0}} {{ join .modules " " }}{{end}}{{ if ne (len .customs) 0 }} {{ modulesJoin .customs }}{{ end }} -> {{ sprintf .format .values }}
`)
)

type Logger interface {
	// AddCustom defines function that can be used by Option.Custom.
	AddCustom(name string, fn CustomHandler)
	// Module creates new Logger with the same Level and name but adds new module with name specified in first argument.
	//
	// Custom functions are also the same (and still updating).
	Module(name string) Logger
	// Fatal writes output to the console AND exits program with status 1.
	Fatal() Option
	Error() Option
	Warn() Option
	Info() Option
	Debug() Option
	// SetLevel changes the logging Level of current Logger.
	SetLevel(lvl Level) Logger
}

type CustomHandler func(arg interface{}) string

type logger struct {
	name      string
	modules   []string
	customs   *sync.Map
	showLevel Level
}

func (l *logger) SetLevel(lvl Level) Logger {
	l.showLevel = lvl
	return l
}

func (l *logger) AddCustom(name string, fn CustomHandler) {
	l.customs.Store(name, fn)
}

func (l *logger) Module(name string) Logger {
	n := NewLoggerWithLevel(l.name, l.showLevel).(*logger)
	n.modules = append(n.modules, name)
	n.customs = l.customs
	return n
}

func (l *logger) Info() Option {
	if l.showLevel < Info {
		return &nullOption{}
	}
	return newOption(l, Info)
}

func (l *logger) Warn() Option {
	if l.showLevel < Warning {
		return &nullOption{}
	}
	return newOption(l, Warning)
}

func (l *logger) Debug() Option {
	if l.showLevel < Debug {
		return &nullOption{}
	}
	return newOption(l, Debug)
}

func (l *logger) Error() Option {
	if l.showLevel < Error {
		return &nullOption{}
	}
	return newOption(l, Error)
}

func (l *logger) Fatal() Option {
	if l.showLevel < Fatal {
		return &nullOption{}
	}
	return newOption(l, Fatal)
}

// NewLoggerWithLevel allows creating fully custom Logger.
func NewLoggerWithLevel(name string, level Level) Logger {
	log := new(logger)
	log.name = name
	log.showLevel = level
	log.customs = new(sync.Map)
	return log
}

// NewLogger creates Logger with Info logging level and custom name.
func NewLogger(name string) Logger {
	return NewLoggerWithLevel(name, Info)
}

// NewDefaultLogger creates Logger with Info logging level and default name - "main".
func NewDefaultLogger() Logger {
	return NewLoggerWithLevel("main", Info)
}
