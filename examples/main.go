package main

import "github.com/VenomPCPL/golog"

func main() {
	log := golog.NewDefaultLogger().SetLevel(golog.LevelDebug)
	println()
	log.Debug().Send("debug msg")
	log.Info().Send("info msg")
	log.Warn().Send("warn msg")
	log.Error().Send("error msg")
	println()
	log.Debug().Str("custom str").Send("debug msg")
	log.Info().Str("custom str").Send("info msg")
	log.Warn().Str("custom str").Send("warn msg")
	log.Error().Str("custom str").Send("error msg")
	println()
	users := log.Module("users")
	users.Debug().Send("users debug msg")
	users.Info().Send("users info msg")
	users.Warn().Send("users warn msg")
	users.Error().Send("users error msg")
	println()
}