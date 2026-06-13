// Package jsonfmt sends logs in JSON.
package jsonfmt

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/BOOMfinity/golog/v3/gcore"
)

type handler struct {
	config Config
}

func (h *handler) handle(ctx gcore.Context) {
	if h.config.Base.Severity != 0 && h.config.Base.Severity > ctx.Level() {
		return
	}

	pooled := gcore.AcquireBuffer()
	defer gcore.FreeBuffer(pooled)

	buff := *pooled

	buff = append(buff, '{')
	buff = h.appendTime(h.appendKey(buff, "time"), time.Now(), h.config.Base.TimeFormat)
	buff = h.appendString(h.appendKey(buff, "level"), ctx.Level().String())

	modules, scopes := ctx.Modules()
	buff = h.appendKey(buff, "module")
	buff = append(buff, '"')
	for i := 0; i < len(modules); i++ {
		buff = append(buff, modules[i]...)
		if scopes[i] != "" {
			buff = append(buff, '@')
			buff = append(buff, scopes[i]...)
		}
		buff = append(buff, ' ')
	}
	buff = buff[:len(buff)-1]
	buff = append(buff, '"')

	if ctx.HasMessage() {
		if h.config.DisableMessageEscaping {
			buff = h.appendString(h.appendKey(buff, "message"), ctx.Message())
		} else {
			buff = h.appendEscapedString(h.appendKey(buff, "message"), ctx.Message())
		}
	}

	if ctx.HasError() {
		buff = h.appendEscapedString(h.appendKey(buff, "error"), ctx.Error().Error())
	}

	global, local := ctx.Attributes()

	if len(global) > 0 || len(local) > 0 {
		for i := range global {
			buff = h.appendAny(h.appendKey(buff, global[i].Name), global[i].Value)
		}
		for i := range local {
			buff = h.appendAny(h.appendKey(buff, local[i].Name), local[i].Value)
		}
	}

	if ctx.HasStackTrace() {
		buff = h.appendEscapedString(h.appendKey(buff, "stack"), ctx.StackTrace())
	}

	buff = append(buff, '}')
	buff = append(buff, '\n')

	if _, err := h.config.Base.Writer(ctx.Level()).Write(buff); err != nil {
		fmt.Fprintf(os.Stderr, "[golog/jsonfmt] write failed: %v\n", err)
	}
}

// Init returns a JSON formatter.
func Init(options ...Config) gcore.Formatter {
	var config Config
	if len(options) > 0 {
		config = options[0]
	}
	if config.Base.TimeFormat == "" {
		config.Base.TimeFormat = time.RFC3339
	}
	if config.Marshaler == nil {
		config.Marshaler = json.Marshal
	}
	return new(handler{config: config}).handle
}
