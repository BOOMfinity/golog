package gcore

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var stackPool = sync.Pool{
	New: func() any {
		return make([]uintptr, 128)
	},
}

func appendStacktrace(buff []byte) []byte {
	pcs := stackPool.Get()
	defer stackPool.Put(pcs)
	n := runtime.Callers(0, pcs.([]uintptr))
	frames := runtime.CallersFrames(pcs.([]uintptr)[:n])
	for {
		f, more := frames.Next()
		if strings.Contains(f.Function, "runtime.") || strings.Contains(f.Function, "golog") {
			if !more {
				break
			}
			continue
		}
		buff = append(buff, '\n')
		buff = append(buff, f.Function...)
		buff = append(buff, "()"...)
		buff = append(buff, "\n\t"...)
		buff = append(buff, f.File...)
		buff = append(buff, ':')
		buff = strconv.AppendInt(buff, int64(f.Line), 10)
		if !more {
			break
		}
	}
	return buff
}
