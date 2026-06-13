// Package colorfmt sends logs in human-readable format with colors.
package colorfmt

import (
	"fmt"
	"os"
	"time"

	"github.com/BOOMfinity/golog/v3/gcore"
)

var (
	resetCode = "\u001B[0m"
	colors    = map[gcore.Level]string{
		gcore.LevelFatal:   "\u001B[1;31m",
		gcore.LevelError:   "\u001B[1;31m",
		gcore.LevelWarning: "\u001B[33m",
		gcore.LevelInfo:    "\u001B[36m",
		gcore.LevelDebug:   "\u001B[35m",
		gcore.LevelTrace:   "\u001B[3;39m",
	}
)

// Init returns a "pretty" formatter.
func Init(options ...Config) gcore.Formatter {
	var config Config
	if len(options) > 0 {
		config = options[0]
	}
	if config.Base.TimeFormat == "" {
		config.Base.TimeFormat = "02.01.2006 15:04:05"
	}
	return func(entry gcore.Context) {
		if config.Base.Severity != 0 && config.Base.Severity > entry.Level() {
			return
		}

		pooled := gcore.AcquireBuffer()
		defer gcore.FreeBuffer(pooled)

		buff := *pooled

		if !config.DisableColors && !DisableColors {
			buff = append(buff, colors[entry.Level()]...)
		}

		buff = time.Now().AppendFormat(buff, config.Base.TimeFormat)
		buff = append(buff, " | "...)
		buff = append(buff, entry.Level().String()...)
		buff = append(buff, " | "...)

		modules, scopes := entry.Modules()

		for i := range modules {
			buff = append(buff, modules[i]...)
			if scopes[i] != "" {
				buff = append(buff, '@')
				buff = append(buff, scopes[i]...)
			}
			buff = append(buff, ' ')
		}

		buff = buff[:len(buff)-1]

		global, local := entry.Attributes()

		if len(global) > 0 || len(local) > 0 {
			buff = append(buff, " | "...)
		}

		for i := range global {
			buff = append(buff, global[i].Name...)
			buff = append(buff, '(')
			buff = appendAttribute(buff, global[i].Value)
			buff = append(buff, ')')
			buff = append(buff, ' ')
		}

		for i := range local {
			buff = append(buff, local[i].Name...)
			buff = append(buff, '(')
			buff = appendAttribute(buff, local[i].Value)
			buff = append(buff, ')')
			buff = append(buff, ' ')
		}

		if len(global) > 0 || len(local) > 0 {
			buff = buff[:len(buff)-1]
		}

		buff = append(buff, " -> "...)

		switch {
		case entry.HasError() && entry.HasMessage():
			buff = append(buff, entry.Message()...)
			buff = append(buff, ": "...)
			buff = append(buff, entry.Error().Error()...)
		case entry.HasError() && !entry.HasMessage():
			buff = append(buff, entry.Error().Error()...)
		default:
			buff = append(buff, entry.Message()...)
		}

		if !config.DisableColors && !DisableColors {
			buff = append(buff, resetCode...)
		}

		if entry.HasStackTrace() {
			buff = append(buff, '\n')
			buff = append(buff, entry.StackTrace()...)
		}

		buff = append(buff, '\n')

		if _, err := config.Base.Writer(entry.Level()).Write(buff); err != nil {
			fmt.Fprintf(os.Stderr, "[golog/colorfmt] write failed: %v\n", err)
		}
	}
}
