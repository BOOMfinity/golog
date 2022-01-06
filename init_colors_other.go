//go:build !windows
// +build !windows

package golog

import (
	"os"
)

func init() {
	if os.Getenv("GCOLORS") == "off" {
		colorsEnabled = false
		return
	}
}
