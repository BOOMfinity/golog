# Golog - Another logger written in Go

![image](https://media.discordapp.net/attachments/871820540762005565/1012454715759546439/unknown.png)

The main golog features are:
- **Zero**(1) allocations
- Color support (also on windows!)
- Debug mode controllable via environment variable
- Looks good(2)
- **Modules** thanks to which you know where the logs come from
- **Arguments** configurable for each log message independently
- **Hooks** which add specific information to every log message that uses a hook with the same name
- **Sentry** and **Discord** integration (via plugins)
- Built-in panic recovery function
- File path with line and stack trace support
- Potential for many integrations and plugins

We also plan to add PostgreSQL (and maybe MongoDB) integration with a dashboard.

_(1) - **IT IS** zero allocation logger if you don't use many arguments or specific formatting. Internally it does not make any allocation._

_(2) - We don't provide any way to change message style or colors_

Most methods and types are not documented as they are pretty straightforward - their names say everything.

Installation
---

Just run `go get github.com/BOOMfinity/golog` in your project directory

Plugins
---

[Sentry.io](https://sentry.io/) integration - [Repository](https://github.com/BOOMfinity/golog-sentry) - `go get github.com/BOOMfinity/golog-sentry`

~~Discord webhook integration - [Repository](https://github.com/BOOMfinity/golog-discord) - `go get github.com/BOOMfinity/golog-discord`~~ (Currently not available as our Discord library is private)

Examples
---

#### Simple "Hello World!" message:

```go
package main

import "github.com/BOOMfinity/golog"

func main() {
	log := golog.New("main")
	log.Info().Send("Hello World!")
}
```

#### Simple http server log message:

```go
package main

import "github.com/BOOMfinity/golog"

func main() {
	log := golog.New("web-server")
	userIP := "127.0.0.1"
	reqPath := "/users/test-1"
	reqMethod := "GET"
	log.Info().Add("%v %v", reqMethod, reqPath).Any(userIP).Send("new request!")
}
```

#### Modules example

```go
package main

import "github.com/BOOMfinity/golog"

func main() {
	log := golog.New("cmd").SetLevel(golog.LevelDebug)
	log.Info().Send("Ready!")
	apiLog := log.Module("api")
	apiLog.Warn().Send("Debug mode is enabled")
	usersLog := apiLog.Module("users")
	usersLog.Debug().Send("Nickname changed")
}
```

**Errors and panics**

```go
package main

import "github.com/BOOMfinity/golog"

func main() {
	log := golog.New("cmd")
	defer log.Recover()
	log.Error().FileWithLine().Stack().Send("error message with stack and file")
	x() // panic will be caught by log.Recover
}

func x() {
	panic("fatal error D:")
}
```

Internal environment variables
---

There are two special, global environment variables that can be used with golog.

- `GOLOG_DEBUG or GDEBUG` = `on | true` - If set to "true" it will globally enable "debug" level for all instances in application.


- `GOLOG_COLORS_DISABLED` = `on | yes | true` - If for some reason you want to disable colors, set this env variable to "true".

