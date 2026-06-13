[![Go Reference](https://pkg.go.dev/badge/github.com/BOOMfinity/golog/v3.svg)](https://pkg.go.dev/github.com/BOOMfinity/golog/v3)

## 💫 Features

**Fast**: ~300 ns/op, zero allocations on the hot path

**Configurable**: independent formatter, hooks, level per instance

**Batteries included**: pretty and JSON output, Sentry & Discord plugins

**Hierarchical**: modules and scopes for structured output

## 🛫 Quick start

```
go get github.com/BOOMfinity/golog/v3
```

```go
golog.Info().Msg("Hello World!")
golog.Error().Cause(err).Msg("something broke")

db := golog.Default().Module("db", "queries")
db.Info().Attr("duration", "12ms").Msg("SELECT users")
```

## 🔌 Plugins

Plugins are separate modules, install them individually.

### Sentry [![Go Reference](https://pkg.go.dev/badge/github.com/BOOMfinity/golog/v3/plugins/sentrylog.svg)](https://pkg.go.dev/github.com/BOOMfinity/golog/v3/plugins/sentrylog)

Captures exceptions and forwards logs to [sentry.io](https://sentry.io).

```
go get github.com/BOOMfinity/golog/v3/plugins/sentrylog
```

```go
log = log.WithHook(sentrylog.Init(true, true)) // exceptions + logs
```

### Discord [![Go Reference](https://pkg.go.dev/badge/github.com/BOOMfinity/golog/v3/plugins/dislog.svg)](https://pkg.go.dev/github.com/BOOMfinity/golog/v3/plugins/dislog)

Sends logs to a Discord channel via webhooks.

```
go get github.com/BOOMfinity/golog/v3/plugins/dislog
```

```go
log = log.WithHook(dislog.Init(webhookClient))
```
