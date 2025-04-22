package golog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
	"unsafe"

	"github.com/BOOMfinity/go-utils/gpool"
)

type encoder struct {
	json *json.Encoder
	buff *bytes.Buffer
}

type structure struct {
	Time     string      `json:"time"`
	Level    string      `json:"level"`
	Context  []Parameter `json:"context"`
	Module   []string    `json:"module"`
	Params   []Parameter `json:"params"`
	Stack    string      `json:"stack,omitempty"`
	Duration int64       `json:"duration,omitempty"`
	Details  any         `json:"details,omitempty"`
	Message  string      `json:"message"`
}

var structures = gpool.New[structure](gpool.OnPut(func(s *structure) {
	clear(s.Module)
	clear(s.Params)
	s.Time = ""
	s.Level = ""
	s.Module = nil
	s.Context = nil
	s.Stack = ""
	s.Duration = 0
	s.Details = nil
	s.Message = ""
}))

var encoders = gpool.New[encoder](gpool.OnInit[encoder](func(e *encoder) {
	e.buff = bytes.NewBuffer(make([]byte, 0, Config.EngineBufferSize()))
	e.json = json.NewEncoder(e.buff)
}))

func JSONEngine(writers ...io.Writer) WriteEngine {
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	writer := io.MultiWriter(writers...)

	return func(log *Logger, data *MessageData) {
		struc := structures.Get()
		defer structures.Put(struc)

		dateTimeBuff := dateTimeBuffer.Get()
		defer dateTimeBuffer.Put(dateTimeBuff)
		*dateTimeBuff = time.Now().AppendFormat(*dateTimeBuff, log.DateTimeFormat(time.RFC3339Nano))
		struc.Time = unsafe.String(unsafe.SliceData(*dateTimeBuff), len(*dateTimeBuff))
		struc.Message = unsafe.String(unsafe.SliceData(data.Message), len(data.Message))
		struc.Level = data.Level.String()
		struc.Details = data.Details
		if data.StackIncluded {
			struc.Stack = unsafe.String(unsafe.SliceData(data.Stack), len(data.Stack))
		}
		struc.Params = data.Params
		if data.Duration > 0 {
			struc.Duration = data.Duration.Milliseconds()
		}
		struc.Module = log.Modules()
		struc.Context = log.Params()
		enc := encoders.Get()
		defer encoders.Put(enc)
		enc.buff.Reset()
		if err := enc.json.Encode(struc); err != nil {
			panic(fmt.Errorf("failed to encode structure: %w", err))
		}
		_, _ = writer.Write(enc.buff.Bytes())
	}
}
