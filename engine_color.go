package golog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/BOOMfinity/go-utils/gpool"
	"github.com/gookit/color"
)

var buffPool = gpool.New[bytes.Buffer](gpool.OnInit[bytes.Buffer](func(b *bytes.Buffer) {
	b.Grow(Config.EngineBufferSize())
}), gpool.OnPut(func(b *bytes.Buffer) {
	b.Reset()
}))

func ColorEngine(writers ...io.Writer) WriteEngine {
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	writer := io.MultiWriter(writers...)

	return func(log *Logger, data *MessageData) {
		buff := buffPool.Get()
		defer buffPool.Put(buff)
		if !Config.DisableTerminalColors() {
			buff.WriteString(data.Level.Color())
		}
		{
			dateBuff := dateTimeBuffer.Get()
			*dateBuff = time.Now().AppendFormat(*dateBuff, log.DateTimeFormat("02.01.2006 15:04:05"))
			buff.Write(*dateBuff)
			dateTimeBuffer.Put(dateBuff)
		}
		buff.WriteString(" | ")
		buff.WriteString(data.Level.String())
		buff.WriteString(" | ")
		for _, module := range log.Modules() {
			buff.WriteString(module)
			buff.WriteByte(' ')
		}
		{
			params := log.Params()
			if len(params) > 0 {
				buff.WriteString("| ")
				for _, parameter := range params {
					buff.WriteString(parameter.Name)
					buff.WriteByte('(')
					_, _ = fmt.Fprint(buff, parameter.Value)
					buff.WriteByte(')')
					buff.WriteByte(' ')
				}
			}
		}
		{
			if len(data.Params) > 0 {
				buff.WriteString("| ")
				for _, p := range data.Params {
					buff.WriteString(p.Name)
					buff.WriteByte('(')
					_, _ = fmt.Fprint(buff, p.Value)
					buff.WriteByte(')')
					buff.WriteByte(' ')
				}
			}
		}
		{
			if data.Duration > 0 {
				buff.WriteString("| ")
				buff.WriteString(data.Duration.Round(time.Millisecond).String())
				buff.WriteByte(' ')
			}
		}
		buff.WriteString("-> ")
		buff.Write(data.Message)
		{
			if data.Details != nil {
				buff.WriteByte('\n')
				if !Config.MarshalDetails() {
					_, _ = fmt.Fprint(buff, data.Details)
				} else {
					enc := encoders.Get()
					defer encoders.Put(enc)
					enc.buff.Reset()
					_ = enc.json.Encode(data.Details)
					buff.Write(enc.buff.Bytes())
					buff.Truncate(buff.Len() - 1)
				}
			}
		}
		if !Config.DisableTerminalColors() {
			buff.WriteString(color.ResetSet)
		}
		buff.WriteByte('\n')
		if data.StackIncluded {
			buff.Write(data.Stack)
		}
		_, _ = writer.Write(buff.Bytes())
		if data.ExitCode != 0 {
			os.Exit(data.ExitCode)
		}
	}
}
