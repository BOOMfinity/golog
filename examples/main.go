package main

import "github.com/VenomPCPL/golog"

func main() {
	log := golog.NewDefault().SetLevel(golog.LevelDebug)
	log.Empty()
	log.Debug().Send("debug msg")
	log.Info().Send("info msg")
	log.Warn().Send("warn msg")
	log.Error().Send("error msg")
	log.Empty()
	log.Debug().Any("custom str").Send("debug msg")
	log.Info().Any("custom str").Send("info msg")
	log.Warn().Any("custom str").Send("warn msg")
	log.Error().Any("custom str").Send("error msg")
	log.Empty()
	users := log.Module("users")
	users.Debug().Send("users debug msg")
	users.Info().Send("users info msg")
	users.Warn().Send("users warn msg")
	users.Error().Send("users error msg")
	log.Empty()
}
