package main

import (
	"github.com/BOOMfinity/golog"
)

func main() {
	log := golog.New("main").SetLevel(golog.LevelDebug)
	defer log.Recover()
	log.Info().Send("Hello World!")
	log = golog.New("web-server")
	userIP := "127.0.0.1"
	reqPath := "/users/test-1"
	reqMethod := "GET"
	log.Info().Add("%v %v", reqMethod, reqPath).Any(userIP).Send("new request!")
	log = golog.New("cmd").SetLevel(golog.LevelDebug)
	log.Info().Send("Ready!")
	apiLog := log.Module("api")
	apiLog.Warn().Send("Debug mode is enabled")
	usersLog := apiLog.Module("users")
	usersLog.Debug().Send("Nickname changed")
	log.Error().FileWithLine().Stack().Send("error message with stack and file")
	x()
}

func x() {
	panic("fatal error D:")
}
