package golog

import (
	"fmt"
	"os"
	"time"
)

type Option interface {
	// Custom allows you to use previously added custom functions via Logger.AddCustom.
	//
	// It also uses custom functions defined by AddGlobalCustom.
	Custom(name string, arg interface{}) Option
	// Add allows you to easily add custom value to the output.
	Add(value interface{}) Option
	// Format adds formatted string as custom value to the output.
	Format(format string, values ...interface{}) Option
	// Send writes output to the stdout.
	Send(name string, values ...interface{})
}

type nullOption struct{}

func (n *nullOption) Custom(_ string, _ interface{}) Option {
	return n
}

func (n *nullOption) Add(_ interface{}) Option {
	return n
}

func (n *nullOption) Format(_ string, _ ...interface{}) Option {
	return n
}

func (n *nullOption) Send(_ string, _ ...interface{}) {
	return
}

func newOption(l *logger, level Level) Option {
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

func (o *option) Custom(name string, arg interface{}) Option {
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

func (o *option) Add(value interface{}) Option {
	o.custom = append(o.custom, fmt.Sprint(value))
	return o
}

func (o *option) Format(format string, values ...interface{}) Option {
	o.custom = append(o.custom, fmt.Sprintf(format, values...))
	return o
}

func (o option) Send(name string, values ...interface{}) {
	defaultTemplate.Execute(os.Stdout, map[string]interface{}{
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
