package main

import (
	"os"

	"click-guardian/internal/gui"
)

func main() {
	// Check for command line arguments
	startMinimized := false
	autoProtect := false

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--minimized":
			startMinimized = true
		case "--auto-protect":
			autoProtect = true
		}
	}

	app := gui.NewApplication()
	if startMinimized {
		app.RunMinimized()
	} else {
		if autoProtect {
			app.RunWithAutoProtect()
		} else {
			app.Run()
		}
	}
}
