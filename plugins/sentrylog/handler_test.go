package sentrylog

import (
	"github.com/BOOMfinity/golog/v3/formats/colorfmt"
	"github.com/BOOMfinity/golog/v3/gcore"
)

func ExampleInit() {
	log := gcore.Init("app", colorfmt.Init())
	log = log.WithHook(Init(true, true))
}
