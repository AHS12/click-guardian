package main

import (
	"fmt"
	"os"

	"click-guardian/internal/gui"
	"click-guardian/internal/version"
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
		case "--version", "-v":
			fmt.Println(version.GetFullVersionString())
			return
		case "--help", "-h":
			showHelp()
			return
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

// showHelp displays command-line usage information
func showHelp() {
	info := version.GetAppInfo()
	fmt.Printf("%s v%s\n", info.Name, info.Version)
	fmt.Printf("%s\n\n", info.Description)

	fmt.Println("Usage:")
	fmt.Printf("  %s [options]\n\n", os.Args[0])

	fmt.Println("Options:")
	fmt.Println("  --minimized      Start minimized to system tray")
	fmt.Println("  --auto-protect   Start with protection automatically enabled")
	fmt.Println("  --version, -v    Show version information")
	fmt.Println("  --help, -h       Show this help message")
	fmt.Println()

	fmt.Println("Examples:")
	fmt.Printf("  %s                    # Start normally\n", os.Args[0])
	fmt.Printf("  %s --minimized        # Start minimized to tray\n", os.Args[0])
	fmt.Printf("  %s --auto-protect     # Start with protection enabled\n", os.Args[0])
	fmt.Println()

	fmt.Printf("%s\n", info.Copyright)
}
