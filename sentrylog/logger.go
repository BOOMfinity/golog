package sentrylog

import (
	"errors"
	"github.com/VenomPCPL/golog"
	"github.com/getsentry/sentry-go"
	"strings"
	"time"
)

func handleSentry(level golog.Level) golog.LogHandler {
	return func(msg golog.Message) {
		if msg.Level() > level {
			return
		}
		if msg.Level() == golog.LevelError || msg.Level() == golog.LevelPanic {
			if msg.Error() != nil {
				sentry.CaptureException(msg.Error())
			} else {
				sentry.CaptureException(errors.New(msg.UserMessage()))
			}
			if msg.GetExitCode() != -1 {
				sentry.Flush(10 * time.Second)
			}
		} else {
			ev := sentry.NewEvent()
			switch msg.Level() {
			case golog.LevelError:
				ev.Level = sentry.LevelError
			case golog.LevelPanic:
				ev.Level = sentry.LevelFatal
			case golog.LevelInfo:
				ev.Level = sentry.LevelInfo
			case golog.LevelWarn:
				ev.Level = sentry.LevelWarning
			case golog.LevelDebug:
				ev.Level = sentry.LevelDebug
			}
			ev.Message = msg.UserMessage()
			ev.Extra["logger_path"] = strings.Join(msg.Instance().Modules()[:], " -> ")
			ev.Extra["arguments"] = strings.Join(msg.Arguments()[:], " : ")
			sentry.CaptureEvent(ev)
		}
	}
}

func UpgradeWithLevel(log golog.Logger, level golog.Level) golog.Logger {
	l := &logger{Logger: log, lvl: level}
	l.AddOnLog("$S_E_N_T_R_Y$", handleSentry(level))
	return l
}

func Upgrade(log golog.Logger) golog.Logger {
	return UpgradeWithLevel(log, golog.LevelError)
}

type logger struct {
	golog.Logger
	lvl golog.Level
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
	return UpgradeWithLevel(l.Logger.Module(name), l.lvl)
}
