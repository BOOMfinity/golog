package internal

import "os"

func IsDebugEnabled() bool {
	env := os.Getenv("GOLOG_DEBUG")
	if env == "" {
		env = os.Getenv("GDEBUG")
	}
	return env == "true" || env == "on"
}
