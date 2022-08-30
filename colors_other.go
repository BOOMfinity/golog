//go:build !windows
// +build !windows

package golog

import (
	"os"
)

func init() {
	if os.Getenv("GOLOG_COLORS_DISABLED") == "off" {
		colorsEnabled = false
		return
	}
}
