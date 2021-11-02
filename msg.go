package golog

import (
	"fmt"
	"os"
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

func (o option) Send(name string, values ...interface{}) {
	defaultTemplate.Execute(o.logger.writer, map[string]interface{}{
		"time":    time.Now(),
		"type":    o.level.String(),
		"format":  name,
		"values":  values,
		"modules": o.logger.modules,
		"customs": o.custom,
		"name":    o.logger.name,
	})
	if o.level == Fatal {
		os.Exit(1)
	}
}
