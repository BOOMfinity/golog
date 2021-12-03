package golog

import (
	"fmt"
	"strings"
	"time"
)

type Message interface {
	// Custom allows you to use previously added (Logger.AddCustom) custom functions.
	//
	// It also uses custom functions defined by AddGlobalCustom.
	Custom(name string, arg interface{}) Message
	// Add allows you to easily add custom value to the output.
	Add(value interface{}) Message
	// Format adds formatted string as custom value to the output.
	Format(format string, values ...interface{}) Message
	// Send writes output to the stdout.
	Send(name string, values ...interface{})
}

type nullMessage struct{}

func (n *nullMessage) Custom(_ string, _ interface{}) Message {
	return n
}

func (n *nullMessage) Add(_ interface{}) Message {
	return n
}

func (n *nullMessage) Format(_ string, _ ...interface{}) Message {
	return n
}

func (n *nullMessage) Send(_ string, _ ...interface{}) {
	return
}

func newMessage(l *logger, level Level) Message {
	return &option{
		logger: l,
		level:  level,
	}
}

type option struct {
	logger *logger
	custom []string
	level  Level
}

func (o *option) Custom(name string, arg interface{}) Message {
	var (
		custom interface{}
		ok     bool
	)
	if custom, ok = o.logger.customs.Load(name); !ok {
		custom, ok = globalCustoms.Load(name)
	}
	if !ok {
		return o
	}
	if fn, ok := custom.(CustomHandler); ok {
		o.custom = append(o.custom, fn(arg))
	}
	return o
}

func (o *option) Add(value interface{}) Message {
	o.custom = append(o.custom, fmt.Sprint(value))
	return o
}

func (o *option) Format(format string, values ...interface{}) Message {
	o.custom = append(o.custom, fmt.Sprintf(format, values...))
	return o
}

func (o *option) Send(name string, values ...interface{}) {
	var str string
	if len(o.logger.modules) > 1 {
		str = fmt.Sprintf(`%v | %v | %v | %v -> %v`,
			time.Now().Format("2006-01-02 15:04:05 (Z07:00)"), o.level.String(),
			strings.Join(o.logger.modules, " "), strings.Join(o.custom, " | "), fmt.Sprintf(name, values...))
	} else {
		str = fmt.Sprintf(`%v | %v | %v -> %v`,
			time.Now().Format("2006-01-02 15:04:05 (Z07:00)"), o.level.String(),
			strings.Join(o.logger.modules, " "), fmt.Sprintf(name, values...))
	}
	o.logger.execHandlers.Range(func(key, value interface{}) bool {
		value.(ExecHandler)(str, o.level)
		return true
	})
	if colors {
		switch o.level {
		case Info:
			str = infoStyle.Sprint(str)
		case Warning:
			str = warningStyle.Sprint(str)
		case Debug:
			str = debugStyle.Sprint(str)
		case Error, Fatal:
			str = errorStyle.Sprint(str)
		}
	}
	fmt.Fprintln(o.logger.writer, str)
}
