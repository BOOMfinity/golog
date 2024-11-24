package golog

import (
	"io"
	"testing"
)

func BenchmarkJsonSimple(b *testing.B) {
	b.Run("WithoutModule", func(b *testing.B) {
		log := NewCustom("test", NewJSONEngine(io.Discard))
		b.ResetTimer()
		for range b.N {
			printSimple(log)
		}
	})

	b.Run("WithModule", func(b *testing.B) {
		log := NewCustom("test", NewJSONEngine(io.Discard)).Module("another-module")
		b.ResetTimer()
		for range b.N {
			printSimple(log)
		}
	})
}

func BenchmarkJsonSimpleThreaded(b *testing.B) {
	b.Run("WithoutModule", func(b *testing.B) {
		log := NewCustom("test", NewJSONEngine(io.Discard))
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				printSimple(log)
			}
		})
	})
	b.Run("WithModule", func(b *testing.B) {
		log := NewCustom("test", NewJSONEngine(io.Discard)).Module("another-module")
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				printSimple(log)
			}
		})
	})
}

func BenchmarkJsonParams2(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printParams2(log)
	}
}

func BenchmarkJsonParams2Threaded(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printParams2(log)
		}
	})
}

func BenchmarkJsonParams10(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printParams10(log)
	}
}

func BenchmarkJsonParams10Threaded(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printParams10(log)
		}
	})
}

func BenchmarkJsonWithFormatString(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printFmtStringArg(log)
	}
}

func BenchmarkJsonWithFormatStringThreaded(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printFmtStringArg(log)
		}
	})
}

func BenchmarkJsonWithFormatInt(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printFmtIntArg(log)
	}
}

func BenchmarkJsonWithFormatIntThreaded(b *testing.B) {
	log := NewCustom("test", NewJSONEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printFmtIntArg(log)
		}
	})
}
