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

func InjectDiscordHook(logger golog.Logger, id snowflake.ID, token string, level golog.Level, opts ...Option) {
	c := new(hookContext)
	for _, opt := range opts {
		opt(context.WithValue(context.Background(), "$hcontext", c))
	}
	wh := webhook.NewClient(id, token, nil)
	logger.OnWrite("discord", func(str string, l golog.Level) {
		if l <= level {
			var color colors.Color
			if l == golog.Debug {
				color = colors.Blue
			} else if l == golog.Info {
				color = colors.Green
			} else if l == golog.Warning {
				color = colors.Yellow
			} else if l == golog.Error || l == golog.Fatal {
				color = colors.Red
			}
			_, err := wh.Execute(discord.WebhookMessageCreateOptions{
				MessageCreateOptions: discord.MessageCreateOptions{
					Content: fmt.Sprintf("`%v`", str),
					Embed: &discord.MessageEmbed{
						Color:       color,
						Description: fmt.Sprintf(`%v`, str),
						Footer:      &discord.EmbedFooter{Text: c.footer},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		}
	})
}
