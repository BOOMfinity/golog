package golog

import "os"

type WriteEngine func(msg Message)

var (
	overrideLoggingLevel = os.Getenv("GOLOG_LEVEL")
	disableColors        = os.Getenv("GOLOG_DISABLE_COLORS") == "true"
)
