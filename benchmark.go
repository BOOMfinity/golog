package golog

import (
	"github.com/BOOMfinity/go-utils/gpool"
	"time"
)

type BenchmarkContext struct {
	last  time.Time
	start time.Time
}

func (ctx *BenchmarkContext) Total() time.Duration {
	return time.Now().Sub(ctx.start)
}

func (ctx *BenchmarkContext) Update() {
	ctx.last = time.Now()
}

func (ctx *BenchmarkContext) Elapsed() time.Duration {
	dur := time.Now().Sub(ctx.last)
	ctx.Update()
	return dur
}

func CreateBenchmarkContext() *BenchmarkContext {
	return &BenchmarkContext{last: time.Now(), start: time.Now()}
}

var benchPool = gpool.New[BenchmarkContext]()

func AcquireBenchmarkContext() *BenchmarkContext {
	bench := benchPool.Get()
	bench.start = time.Now()
	bench.last = time.Now()
	return bench
}

func ReleaseBenchmarkContext(ctx *BenchmarkContext) {
	benchPool.Put(ctx)
}
