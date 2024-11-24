package golog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/BOOMfinity/go-utils/gpool"
	"github.com/gookit/color"
)

var globalLocker sync.Mutex
var stdout = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")

var TimeFormat = "02.01.2006 15:04:05"

var buffPool = gpool.New[bytes.Buffer](gpool.OnInit[bytes.Buffer](func(b *bytes.Buffer) {
	b.Grow(1024 * 8)
}))

var loggerPool = gpool.New[[]Logger](gpool.OnInit[[]Logger](func(i *[]Logger) {
	*i = make([]Logger, 0, 25)
}))

func NewColorEngine(writers ...io.Writer) WriteEngine {
	if len(writers) == 0 {
		writers = append(writers, stdout)
	}

	writer := io.MultiWriter(writers...)

	return func(msg Message) {
		buff := buffPool.Get()
		defer buff.Reset()
		defer buffPool.Put(buff)
		if !disableColors {
			buff.WriteString(msg.Internal().Level.Color())
		}
		buff.Write(time.Now().AppendFormat(buff.Bytes(), TimeFormat))
		buff.WriteString(" | ")
		buff.WriteString(msg.Internal().Level.String())
		buff.WriteString(" | ")
		{
			loggers := loggerPool.Get()
			logger := msg.Internal().Logger
			for logger != nil {
				*loggers = append(*loggers, logger)
				logger = logger.Internal().Parent
			}
			slices.Reverse(*loggers)
			for _, logger = range *loggers {
				buff.WriteString(logger.Internal().Module)
				if logger.Internal().Scope != "" {
					buff.WriteByte('@')
					buff.WriteString(logger.Internal().Scope)
				}
				buff.WriteByte(' ')
			}
			*loggers = (*loggers)[:0]
			loggerPool.Put(loggers)
		}
		{
			if len(msg.Internal().Params) > 0 || len(msg.Internal().Logger.Internal().Params) > 0 {
				buff.WriteString("| ")
			}
			if len(msg.Internal().Logger.Internal().Params) > 0 {
				for _, val := range msg.Internal().Logger.Internal().Params {
					k, v := val[0], val[1]
					_, _ = fmt.Fprint(buff, k)
					buff.WriteByte('(')
					_, _ = fmt.Fprint(buff, v)
					buff.WriteByte(')')
					buff.WriteByte(' ')
				}
			}
			if len(msg.Internal().Params) > 0 {
				for _, val := range msg.Internal().Params {
					k, v := val[0], val[1]
					_, _ = fmt.Fprint(buff, k)
					buff.WriteByte('(')
					_, _ = fmt.Fprint(buff, v)
					buff.WriteByte(')')
					buff.WriteByte(' ')
				}
			}
		}
		{
			if msg.Internal().Duration > 0 {
				buff.WriteString("| ")
				buff.WriteString(msg.Internal().Duration.Round(time.Millisecond).String())
				buff.WriteByte(' ')
			}
		}
		buff.WriteString("-> ")
		buff.Write(msg.Internal().Message)
		{
			if msg.Internal().Details != nil {
				buff.WriteByte('\n')
				_, _ = fmt.Fprintf(buff, "%v", msg.Internal().Details)
			}
		}
		if !disableColors {
			buff.WriteString(color.ResetSet)
		}
		buff.WriteByte('\n')
		if !isEmpty(msg.Internal().Stack) {
			buff.Write(msg.Internal().Stack)
		}

		globalLocker.Lock()
		_, _ = writer.Write(buff.Bytes())
		if msg.Internal().ExitCode != 0 {
			os.Exit(msg.Internal().ExitCode)
		}
		globalLocker.Unlock()
	}
}
