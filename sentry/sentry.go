package sentry

import (
	"fmt"
	"slices"
	"time"

	"github.com/BOOMfinity/golog/v2"
	"github.com/getsentry/sentry-go"
)

func sentryLevel(level golog.Level) sentry.Level {
	switch level {
	case golog.LevelTrace:
		return sentry.LevelDebug
	case golog.LevelDebug:
		return sentry.LevelDebug
	case golog.LevelInfo:
		return sentry.LevelInfo
	case golog.LevelWarning:
		return sentry.LevelWarning
	case golog.LevelError:
		return sentry.LevelError
	case golog.LevelPanic:
		return sentry.LevelFatal
	}
	panic("invalid logging level")
}

const identifier = "boomfinity.golog"

func New(opts sentry.ClientOptions, levels ...golog.Level) (golog.WriteEngine, error) {
	c, err := sentry.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize sentry client: %v", err)
	}
	c.SetSDKIdentifier(identifier)
	hub := sentry.NewHub(c, sentry.NewScope())
	return NewWithHub(hub, levels...)
}

func NewWithHub(hub *sentry.Hub, levels ...golog.Level) (golog.WriteEngine, error) {
	if len(levels) == 0 {
		levels = append(levels, golog.LevelPanic, golog.LevelError)
	}
	return func(log *golog.Logger, data *golog.MessageData) {
		if !slices.Contains(levels, data.Level) {
			println("skip")
			return
		}
		ev := sentry.NewEvent()
		ev.Timestamp = time.Now()
		ev.Logger = "golog"
		ev.Level = sentryLevel(data.Level)
		if len(data.Message) > 0 {
			ev.Message = string(data.Message)
		}
		ctx := sentry.Context{
			"module": log.Modules(),
		}
		if len(data.Params) > 0 {
			ctx["message_params"] = data.Params
		}
		if len(log.Params()) > 0 {
			ctx["params"] = log.Params()
		}
		if data.ExitCode != 0 {
			ctx["exit_code"] = data.ExitCode
		}
		if data.Duration > 0 {
			ctx["duration"] = data.Duration.String()
		}
		if data.Details != nil {
			ctx["details"] = data.Details
		}
		if data.StackIncluded {
			ev.SetException(data.Error, -1)
		}
		ev.Contexts["golog"] = ctx
		hub.CaptureEvent(ev)
		hub.Flush(2 * time.Second)
	}, nil
}
