// Package sentrylog is a plugin that integrates golog with sentry.
package sentrylog

import (
	"context"
	"encoding/json"

	"github.com/BOOMfinity/golog/v3/gcore"
	"github.com/getsentry/sentry-go"
)

func appendAttrs(entry sentry.LogEntry, attrs []gcore.Attribute) {
	for _, p := range attrs {
		switch v := p.Value.(type) {
		case string:
			entry.String(p.Name, v)
		case float32:
			entry.Float64(p.Name, float64(v))
		case float64:
			entry.Float64(p.Name, v)
		case uint:
			entry.Int64(p.Name, int64(v))
		case uint8:
			entry.Int64(p.Name, int64(v))
		case uint16:
			entry.Int64(p.Name, int64(v))
		case uint32:
			entry.Int64(p.Name, int64(v))
		case uint64:
			entry.Int64(p.Name, int64(v))
		case int:
			entry.Int64(p.Name, int64(v))
		case int8:
			entry.Int64(p.Name, int64(v))
		case int16:
			entry.Int64(p.Name, int64(v))
		case int32:
			entry.Int64(p.Name, int64(v))
		case int64:
			entry.Int64(p.Name, v)
		case bool:
			entry.Bool(p.Name, v)
		default:
			if b, err := json.Marshal(v); err == nil {
				entry.String(p.Name, string(b))
			}
		}
	}
}

// Init returns a hook which works with the Sentry API.
//
// If exceptions is true, errors will be captured as Sentry exceptions.
//
// If logs is true, Sentry will be used as a log collector.
//
// To filter which logs are collected,
// configure the Sentry instance with [sentry.ClientOptions].
func Init(exceptions, logs bool) gcore.Hook {
	var log sentry.Logger
	if logs {
		log = sentry.NewLogger(context.Background())
	}
	return func(entry gcore.Context) {
		if exceptions && entry.HasError() {
			hub := sentry.GetHubFromContext(entry.Context())
			if hub == nil {
				hub = sentry.CurrentHub()
			}
			hub.CaptureException(entry.Error())
		}
		if log != nil {
			var le sentry.LogEntry
			switch entry.Level() {
			case gcore.LevelFatal:
				le = log.LFatal()
			case gcore.LevelError:
				le = log.Error()
			case gcore.LevelWarning:
				le = log.Warn()
			case gcore.LevelInfo:
				le = log.Info()
			case gcore.LevelDebug:
				le = log.Debug()
			case gcore.LevelTrace:
				le = log.Trace()
			}
			modules, scopes := entry.Modules()
			modSlice := make([]string, 0, len(modules))
			for i := range modules {
				modSlice = append(modSlice, modules[i])
				if scopes[i] != "" {
					modSlice[i] += "@" + scopes[i]
				}
			}
			le.StringSlice("module", modSlice)
			global, local := entry.Attributes()
			appendAttrs(le, global)
			appendAttrs(le, local)
			le.WithCtx(entry.Context()).Emit(entry.Message())
		}
	}
}
