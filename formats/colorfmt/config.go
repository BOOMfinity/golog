package colorfmt

import (
	"os"

	"github.com/BOOMfinity/golog/v3/gcore"
)

var (
	// DisableColors disables colors globally when set to true.
	//
	// It can be configured at runtime or through GOLOG_DISABLE_COLORS environment variable.
	//
	// You can also disable colors per-handler by setting Config.DisableColors variable to true.
	DisableColors = false
)

// Config holds the colorfmt formatter settings.
type Config struct {
	Base          gcore.Config
	DisableColors bool
}

func init() {
	if os.Getenv("GOLOG_DISABLE_COLORS") != "" {
		DisableColors = true
	}
}
