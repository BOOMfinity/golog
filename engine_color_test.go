package golog

import (
	"io"
	"testing"
)

func BenchmarkColorSimple(b *testing.B) {
	b.Run("WithoutModule", func(b *testing.B) {
		b.ReportAllocs()
		log := NewCustom("test", NewColorEngine(io.Discard))
		b.ResetTimer()
		for range b.N {
			printSimple(log)
		}
	})

	b.Run("WithModule", func(b *testing.B) {
		b.ReportAllocs()
		log := NewCustom("test", NewColorEngine(io.Discard)).Module("another-module")
		b.ResetTimer()
		for range b.N {
			printSimple(log)
		}
	})
}

func BenchmarkColorSimpleThreaded(b *testing.B) {
	b.Run("WithoutModule", func(b *testing.B) {
		b.ReportAllocs()
		log := NewCustom("test", NewColorEngine(io.Discard))
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				printSimple(log)
			}
		})
	})
	b.Run("WithModule", func(b *testing.B) {
		b.ReportAllocs()
		log := NewCustom("test", NewColorEngine(io.Discard)).Module("another-module")
		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				printSimple(log)
			}
		})
	})
}

func BenchmarkColorParams2(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printParams2(log)
	}
}

func BenchmarkColorParams2Threaded(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printParams2(log)
		}
	})
}

func BenchmarkColorParams10(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printParams10(log)
	}
}

func BenchmarkColorParams10Threaded(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printParams10(log)
		}
	})
}

func BenchmarkColorWithFormatString(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printFmtStringArg(log)
	}
}

func BenchmarkColorWithFormatStringThreaded(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printFmtStringArg(log)
		}
	})
}

func BenchmarkColorWithFormatInt(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	for range b.N {
		printFmtIntArg(log)
	}
}

func BenchmarkColorWithFormatIntThreaded(b *testing.B) {
	b.ReportAllocs()
	log := NewCustom("test", NewColorEngine(io.Discard))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			printFmtIntArg(log)
		}
	})
}
