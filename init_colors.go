package golog

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"runtime"
)

func init() {
	if os.Getenv("GCOLORS") == "off" {
		colorsEnabled = false
		return
	}

	if runtime.GOOS == "windows" {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("[golog] Recovered from panic while initializing windows VT processing. Colors will be disabled.")
				colorsEnabled = false
				return
			}
		}()
		stdout := windows.Handle(os.Stdout.Fd())
		var originalMode uint32

		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}
