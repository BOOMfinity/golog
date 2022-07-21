package discord

import (
	"context"
	"fmt"

	"github.com/BOOMfinity-Developers/bfcord/discord"
	"github.com/BOOMfinity-Developers/bfcord/discord/colors"
	"github.com/BOOMfinity-Developers/bfcord/webhook"
	"github.com/VenomPCPL/golog"
	"github.com/andersfylling/snowflake/v5"
)

type hookContext struct {
	footer string
}

type Option func(ctx context.Context)

func SetFooter(str string) Option {
	return func(ctx context.Context) {
		c, ok := ctx.Value("$hcontext").(*hookContext)
		if !ok || c == nil {
			return
		}
		c.footer = str
	}
}

func InjectDiscordHook(logger *golog.Logger, id snowflake.Snowflake, token string, level golog.Level, opts ...Option) {
	c := new(hookContext)
	for _, opt := range opts {
		opt(context.WithValue(context.Background(), "$hcontext", c))
	}
	wh := webhook.NewClient(id, token, nil)
	logger.WriteHook(func(m *golog.Message, _ []byte, ui []byte) {
		if m.Level() < level {
			return
		}
		var color colors.Color
		if m.Level() == golog.LevelDebug {
			color = colors.Purple
		} else if m.Level() == golog.LevelInfo {
			color = colors.Blue
		} else if m.Level() == golog.LevelWarn {
			color = colors.Orange
		} else if m.Level() == golog.LevelError {
			color = colors.Red
		}
		_, err := wh.Execute(discord.WebhookMessageCreateOptions{
			MessageCreateOptions: discord.MessageCreateOptions{
				Embeds: []discord.MessageEmbed{
					{
						Color:       color,
						Title:       m.Level().String(),
						Description: fmt.Sprintf("```%v```", string(ui)),
						Footer:      &discord.EmbedFooter{Text: c.footer},
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}
	})
}
