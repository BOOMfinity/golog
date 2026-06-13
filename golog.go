// Package golog is the entry point for the library.
//
// It provides a global logger, shortcut functions for each level, and an Init that defaults to colorfmt.
package golog

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/BOOMfinity/golog/v3/formats/colorfmt"
	"github.com/BOOMfinity/golog/v3/gcore"
)

// Init creates a new [gcore.Logger] with colorfmt as the default formatter.
//
// Pass a custom [gcore.Formatter] as the second argument to override.
func Init(name string, f ...gcore.Formatter) *gcore.Logger {
	formatter := colorfmt.Init()
	if len(f) > 0 {
		formatter = f[0]
	}
	return gcore.Init(name, formatter)
}

var def = new(atomic.Pointer[gcore.Logger])

// Default returns the global logger.
//
// The default instance is created on the first call with name "app" and colorfmt as a formatter.
func Default() *gcore.Logger {
	log := def.Load()
	if log == nil {
		log = Init("app")
		def.Store(log)
	}
	return log
}

// SetDefault replaces the global logger.
func SetDefault(log *gcore.Logger) {
	def.Store(log)
}

// Info logs at INFO level.
//
// Shortcut for Default().Info(c...).
func Info(c ...context.Context) *gcore.Message {
	return Default().Info(c...)
}

// Debug logs at DEBUG level.
//
// Shortcut for Default().Debug(c...)
func Debug(c ...context.Context) *gcore.Message {
	return Default().Debug(c...)
}

// Trace logs at TRACE level.
//
// Shortcut for Default().Trace(c...)
func Trace(c ...context.Context) *gcore.Message {
	return Default().Trace(c...)
}

// Error logs at ERROR level.
//
// Shortcut for Default().Error(c...)
func Error(c ...context.Context) *gcore.Message {
	return Default().Error(c...)
}

// Fatal logs at FATAL level.
//
// Shortcut for Default().Fatal(code, c...)
func Fatal(code int, c ...context.Context) *gcore.Message {
	return Default().Fatal(code, c...)
}

// Warn logs at WARNING level.
//
// Shortcut for Default().Warn(c...)
func Warn(c ...context.Context) *gcore.Message {
	return Default().Warn(c...)
}

// Recover catches a panic, logs it, and resumes execution.
func Recover(c ...context.Context) {
	if r := recover(); r != nil {
		msg := Default().Error(c...)
		if err, ok := r.(error); ok {
			msg.MsgErr(fmt.Errorf("recovered: %w", err))
		} else {
			msg.Msgf("recovered: %v", r)
		}
	}
}
