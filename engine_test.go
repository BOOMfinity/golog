package golog

import (
	"fmt"
	"strings"
	"testing"
)

func runUserMessage(log *Logger, num int) func(b *testing.B) {
	log = log.Copy()
	return func(b *testing.B) {
		b.ReportAllocs()
		message := strings.Repeat("x", num)
		b.ResetTimer()
		for range b.N {
			log.Info().Send(message)
		}
	}
}

func runWithModules(log *Logger, num int) func(b *testing.B) {
	log = log.Copy()
	return func(b *testing.B) {
		b.ReportAllocs()
		for range num {
			log = log.Module(fmt.Sprintf("mod-%d", num))
		}
		b.ResetTimer()
		for range b.N {
			log.Info().Send("test")
		}
	}
}

func runJustMessage(log *Logger) func(b *testing.B) {
	log = log.Copy()
	return func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			log.Info().Send("test")
		}
	}
}

func runWithDetails(log *Logger) func(b *testing.B) {
	log = log.Copy()
	return func(b *testing.B) {
		b.Run("string", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				log.Info().Details("test").Send("test")
			}
		})
		b.Run("int", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				log.Info().Details(5).Send("test")
			}
		})
		b.Run("float", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				log.Info().Details(3.14).Send("test")
			}
		})
		b.Run("slice", func(b *testing.B) {
			b.ReportAllocs()
			arr := []int{1, 2, 3, 4, 5}
			b.ResetTimer()
			for range b.N {
				log.Info().Details(arr).Send("test")
			}
		})
	}
}
