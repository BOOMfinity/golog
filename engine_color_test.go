package golog

import (
	"io"
	"testing"
)

func BenchmarkColor(b *testing.B) {
	log := New("test", ColorEngine(io.Discard))
	b.Run("JustMessage", runJustMessage(log))
	b.Run("Modules", func(b *testing.B) {
		b.Run("1", runWithModules(log, 1))
		b.Run("3", runWithModules(log, 3))
		b.Run("5", runWithModules(log, 5))
		b.Run("10", runWithModules(log, 10))
		b.Run("20", runWithModules(log, 20))
	})
	b.Run("WithDetails", func(b *testing.B) {
		b.Run("ParseToJSON=true", runWithDetails(log))
		Config.SetMarshalDetails(false)
		b.Run("ParseToJSON=false", runWithDetails(log))
		Config.SetMarshalDetails(true)
	})
	b.Run("UserMessage", func(b *testing.B) {
		b.Run("10", runUserMessage(log, 10))
		b.Run("25", runUserMessage(log, 25))
		b.Run("50", runUserMessage(log, 50))
		b.Run("100", runUserMessage(log, 100))
		b.Run("200", runUserMessage(log, 200))
		b.Run("400", runUserMessage(log, 400))
		b.Run("800", runUserMessage(log, 800))
		b.Run("1600", runUserMessage(log, 1600))
	})
}
