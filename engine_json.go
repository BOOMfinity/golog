package golog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

func NewJSONEngine(writers ...io.Writer) WriteEngine {
	if len(writers) == 0 {
		writers = append(writers, stdout)
	}

	writer := io.MultiWriter(writers...)

	return func(msg Message) {
		buff := buffPool.Get()
		defer buffPool.Put(buff)
		buff.WriteString(`{"time:"`)
		buff.Write(time.Now().AppendFormat(buff.Bytes()[8:], time.RFC3339))
		buff.WriteString(`","level":"`)
		buff.WriteString(msg.Internal().Level.String())
		buff.WriteString(`","modules":[`)
		{
			loggers := loggerPool.Get()
			logger := msg.Internal().Logger
			i := 0
			for logger != nil {
				*loggers = append(*loggers, logger)
				logger = logger.Internal().Parent
			}
			slices.Reverse(*loggers)
			for i, logger = range *loggers {
				buff.WriteByte('"')
				buff.WriteString(logger.Internal().Module)
				if logger.Internal().Scope != "" {
					buff.WriteByte('@')
					buff.WriteString(logger.Internal().Scope)
				}
				buff.WriteByte('"')
				if i < len(*loggers)-1 {
					buff.WriteByte(',')
				}
			}
			*loggers = (*loggers)[:0]
			loggerPool.Put(loggers)
		}
		buff.WriteString(`]`)
		if len(msg.Internal().Params) > 0 || len(msg.Internal().Logger.Internal().Params) > 0 {
			sum := len(msg.Internal().Params) + len(msg.Internal().Logger.Internal().Params)
			num := 0
			buff.WriteString(`","params":{`)
			for _, val := range msg.Internal().Logger.Internal().Params {
				k, v := val[0], val[1]
				num++
				buff.WriteByte('"')
				_, _ = fmt.Fprint(buff, k)
				buff.WriteByte(':')
				_, _ = fmt.Fprint(buff, v)
				buff.WriteByte('"')
				if num < sum {
					buff.WriteByte(',')
				}
			}
			for _, val := range msg.Internal().Params {
				k, v := val[0], val[1]
				num++
				buff.WriteByte('"')
				_, _ = fmt.Fprint(buff, k)
				buff.WriteByte(':')
				_, _ = fmt.Fprint(buff, v)
				buff.WriteByte('"')
				if num < sum {
					buff.WriteByte(',')
				}
			}
			buff.WriteByte('}')
		}
		{
			if msg.Internal().Duration > 0 {
				buff.WriteString(`","duration":`)
				buff.WriteString(fmt.Sprint(msg.Internal().Duration.Milliseconds()))
			}
		}
		if !isEmpty(msg.Internal().Stack) {
			buff.WriteString(`,"stack":"`)
			buff.Write(bytes.ReplaceAll(bytes.ReplaceAll(msg.Internal().Stack, []byte("\n"), []byte("\\n")), []byte(`"`), []byte(`\"`)))
			buff.WriteByte('"')
		}
		buff.WriteString(`,"message":"`)
		buff.Write(msg.Internal().Message)
		buff.WriteByte('"')
		buff.WriteByte('}')
		buff.WriteByte('\n')
		globalLocker.Lock()
		_, _ = buff.WriteTo(writer)
		if msg.Internal().ExitCode > 0 {
			os.Exit(msg.Internal().ExitCode)
		}
		globalLocker.Unlock()
	}
}

func isEmpty(buff []byte) bool {
	for _, b := range buff {
		if b != 0 {
			return false
		}
	}
	return true
}
