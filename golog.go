package golog

import (
	"os"
	"sync"

	"github.com/BOOMfinity/go-utils/gpool"
)

var dateTimeBuffer = gpool.New[[]byte](gpool.OnInit[[]byte](func(data *[]byte) {
	*data = make([]byte, 0, 128)
}), gpool.OnPut[[]byte](func(data *[]byte) {
	*data = (*data)[:0:128]
}))

func Wrap(eng ...WriteEngine) WriteEngine {
	return func(log *Logger, data *MessageData) {
		for _, v := range eng {
			v(log, data)
		}
	}
}

type WriteEngine func(log *Logger, data *MessageData)

type globalOptions struct {
	mut                  sync.RWMutex
	stackTraceBufferSize int
	//loggerParametersSliceAllocation  int
	//loggerModulesSliceAllocation     int
	messageParametersSliceAllocation int
	messageBufferSize                int
	engineBufferSize                 int
	overrideMinimumMessageLevel      Level
	disableTerminalColors            bool
	marshalDetails                   bool
	detailsBufferSize                int
	includeStackOnError              bool
}

func (o *globalOptions) IncludeStackOnError() bool {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.includeStackOnError
}

func (o *globalOptions) StackTraceBufferSize() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.stackTraceBufferSize
}

/*func (o *globalOptions) LoggerParametersSliceAllocation() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.loggerParametersSliceAllocation
}

func (o *globalOptions) LoggerModulesSliceAllocation() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.loggerModulesSliceAllocation
}*/

func (o *globalOptions) MessageParametersSliceAllocation() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.messageParametersSliceAllocation
}

func (o *globalOptions) MessageBufferSize() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.messageBufferSize
}

func (o *globalOptions) EngineBufferSize() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.engineBufferSize
}

func (o *globalOptions) OverrideMinimumMessageLevel() Level {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.overrideMinimumMessageLevel
}

func (o *globalOptions) DisableTerminalColors() bool {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.disableTerminalColors
}

func (o *globalOptions) MarshalDetails() bool {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.marshalDetails
}

func (o *globalOptions) DetailsBufferSize() int {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.detailsBufferSize
}

func (o *globalOptions) SetStackTraceBufferSize(size int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.stackTraceBufferSize = size
}

/*func (o *globalOptions) SetLoggerParametersSliceAllocation(size int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.loggerParametersSliceAllocation = size
}

func (o *globalOptions) SetLoggerModulesSliceAllocation(size int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.loggerModulesSliceAllocation = size
}*/

func (o *globalOptions) SetMessageParametersSliceAllocation(size int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.messageParametersSliceAllocation = size
}

func (o *globalOptions) SetMessageBufferSize(size int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.messageBufferSize = size
}

func (o *globalOptions) SetEngineBufferSize(size int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.engineBufferSize = size
}

func (o *globalOptions) SetOverrideMinimumMessageLevel(level Level) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.overrideMinimumMessageLevel = level
}

func (o *globalOptions) SetDisableTerminalColors(disable bool) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.disableTerminalColors = disable
}

func (o *globalOptions) SetMarshalDetails(enable bool) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.marshalDetails = enable
}

func (o *globalOptions) SetDetailsBufferSize(size int) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.detailsBufferSize = size
}

func (o *globalOptions) SetIncludeStackOnError(include bool) {
	o.mut.Lock()
	defer o.mut.Unlock()
	o.includeStackOnError = include
}

var Config = globalOptions{
	stackTraceBufferSize:             512,
	messageParametersSliceAllocation: 25,
	messageBufferSize:                512,
	engineBufferSize:                 2048,
	overrideMinimumMessageLevel:      0,
	disableTerminalColors:            false,
	marshalDetails:                   true,
	detailsBufferSize:                1024,
	includeStackOnError:              false,
}

func init() {
	if os.Getenv("GOLOG_DISABLE_COLORS") != "" {
		Config.SetDisableTerminalColors(true)
	}
	if strl := os.Getenv("GOLOG_MINIMUM_MESSAGE_LEVEL"); strl != "" {
		level := levelFromString(strl)
		if level != 0 {
			Config.SetOverrideMinimumMessageLevel(level)
		}
	}
}
