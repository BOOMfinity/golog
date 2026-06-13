// Package dislog is a plugin that send logs to the Discord channel through webhook.
package dislog

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/BOOMfinity/golog/v3"
	"github.com/BOOMfinity/golog/v3/gcore"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
)

var (
	mu    sync.Mutex
	hooks []*hook
)

func init() {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		for range ticker.C {
			FlushAll()
		}
	}()
}

// FlushAll sends all remaining logs to the Discord.
func FlushAll() {
	mu.Lock()
	toFlush := make([]*hook, len(hooks))
	copy(toFlush, hooks)
	mu.Unlock()

	for _, h := range toFlush {
		h._flush(true)
	}
}

var _internal = golog.Default().Clone("dislog")

// Config holds the dislog hook settings.
type Config struct {
	// Minimum level required for a log to be sent to Discord.
	// Defaults to [gcore.LevelWarning].
	MinimumLevel gcore.Level
}

const (
	LayoutComponentsLimit = 5
)

type hook struct {
	wh      *webhook.Client
	entries []discord.LayoutComponent
	mu      sync.Mutex
	config  Config
}

func (h *hook) _flush(lock bool) {
	if lock {
		h.mu.Lock()
		defer h.mu.Unlock()
	}

	if len(h.entries) == 0 {
		return
	}

	_internal.Trace().Attr("lock", lock).Msgf("Flushing %d entries", h.entries)

	if _, err := h.wh.CreateMessage(discord.WebhookMessageCreate{
		Components: h.entries,
		Flags:      discord.MessageFlagIsComponentsV2,
	}, rest.CreateWebhookMessageParams{
		WithComponents: true,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "[golog/dislog] webhook failed: %v\n", err)
	}

	clear(h.entries)
	h.entries = h.entries[:0]
}

func (h *hook) createTextDisplay(ctx gcore.Context) (d discord.TextDisplayComponent) {
	if ctx.HasMessage() {
		d.Content += strings.Clone(ctx.Message())
	} else {
		d.Content += ctx.Error().Error()
	}

	global, local := ctx.Attributes()

	if len(global) > 0 || len(local) > 0 {
		d.Content += "\n\n`"
	}

	for i, p := range global {
		d.Content += fmt.Sprintf("%s(%v)", p.Name, p.Value)
		if i < len(global)-1 || len(local) > 0 {
			d.Content += " "
		}
	}

	for i, p := range local {
		d.Content += fmt.Sprintf("%s(%v)", p.Name, p.Value)
		if i < len(local)-1 {
			d.Content += " "
		}
	}

	if len(global) > 0 || len(local) > 0 {
		d.Content += "`"
	}

	if ctx.HasStackTrace() {
		available := 1800 - len(d.Content)
		stack := strings.Clone(ctx.StackTrace())
		if len(stack) > available {
			stack = stack[:available] + "..."
		}
		d.Content += "\n\n```" + stack + "```"
	}

	d.Content += "\n\n-# "

	if ctx.HasExitCode() {
		d.Content += fmt.Sprintf("exit code %d \\| ", ctx.ExitCode())
	}

	modules, scopes := ctx.Modules()

	for i := range modules {
		d.Content += modules[i]
		if scopes[i] != "" {
			d.Content += "@"
			d.Content += scopes[i]
		}
		if i+1 < len(modules) {
			d.Content += " -> "
		}
	}

	return
}

func (h *hook) color(ctx gcore.Context) int {
	switch ctx.Level() {
	case gcore.LevelFatal:
		return 0xFF0000
	case gcore.LevelError:
		return 0xE74C3C
	case gcore.LevelWarning:
		return 0xF1C40F
	case gcore.LevelInfo:
		return 0x3498DB
	case gcore.LevelDebug:
		return 0x9B59B6
	case gcore.LevelTrace:
		return 0x95A5A6
	}
	panic("invalid logging level")
}

// Init returns a hook that sends logs to Discord via a webhook.
//
// Logs are flushed to Discord every 15 seconds, or when the per-message
// component limit is reached.
//
// Use [FlushAll] to send all remaining logs before the program exits.
func Init(wh *webhook.Client, cfg ...Config) gcore.Hook {
	var config Config
	if len(cfg) > 0 {
		config = cfg[0]
	}
	if config.MinimumLevel == 0 {
		config.MinimumLevel = gcore.LevelWarning
	}
	h := &hook{
		wh:      wh,
		entries: make([]discord.LayoutComponent, 0, LayoutComponentsLimit),
		config:  config,
	}
	mu.Lock()
	hooks = append(hooks, h)
	mu.Unlock()
	return func(entry gcore.Context) {
		if entry.Level() >= h.config.MinimumLevel {
			h.mu.Lock()
			defer h.mu.Unlock()

			td := h.createTextDisplay(entry)

			if len(h.entries) >= LayoutComponentsLimit {
				h._flush(false)
			}

			color := h.color(entry)

			var container *discord.ContainerComponent

			if len(h.entries) > 0 {
				container = h.entries[len(h.entries)-1].(*discord.ContainerComponent)
			}

			if container == nil || container.AccentColor != color || len(container.Components) >= 8 {
				container = &discord.ContainerComponent{
					AccentColor: color,
				}
				h.entries = append(h.entries, container)
			}

			if len(container.Components) == 0 {
				container.Components = []discord.ContainerSubComponent{td}
			} else {
				container.Components = append(container.Components, discord.SeparatorComponent{}, td)
			}

			if entry.Level() == gcore.LevelFatal {
				h._flush(false)
			}
		}
	}
}
