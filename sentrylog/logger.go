package sentrylog

import (
	"errors"
	"github.com/VenomPCPL/golog"
	"github.com/getsentry/sentry-go"
	"time"
)

func handleSentry(level golog.Level) golog.LogHandler {
	return func(msg golog.Message) {
		if msg.Level() > level {
			return
		}
		if msg.Level() == golog.LevelError {
			if msg.Error() != nil {
				sentry.CaptureException(msg.Error())
			} else {
				sentry.CaptureException(errors.New(msg.UserMessage()))
			}
		} else {
			sentry.CaptureMessage(msg.UserMessage())
		}
	}
}

func UpgradeWithLevel(log golog.Logger, level golog.Level) golog.Logger {
	l := &logger{Logger: log}
	l.AddOnLog("$S_E_N_T_R_Y$", handleSentry(level))
	return l
}

func Upgrade(log golog.Logger) golog.Logger {
	return UpgradeWithLevel(log, golog.LevelError)
}

type logger struct {
	golog.Logger
}

func (l *logger) Recover() {
	if err := recover(); err != nil {
		msg := l.Error().Any("!! PANIC !!")
		if _err, ok := err.(error); ok {
			msg.SendError(_err)
		} else {
			msg.Stack().Send("%v", err)
		}
		sentry.Flush(2 * time.Second)
	}
}

func (l *logger) Module(name string) golog.Logger {
	return Upgrade(l.Logger.Module(name))
}
