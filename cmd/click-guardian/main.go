package main

import (
	"fmt"
	"os"
	"time"

	"click-guardian/internal/gui"
	"click-guardian/internal/version"
	"click-guardian/pkg/platform"

	"github.com/juju/mutex/v2"
)

// realClock implements mutex.Clock using the real system clock
type realClock struct{}

// After waits for the duration to elapse and then sends the current time on the returned channel
func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Now returns the current clock time
func (realClock) Now() time.Time {
	return time.Now()
}

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

	// Define a unique name for your app's mutex
	spec := mutex.Spec{
		Name:    "click-guardian-single-instance", // Must be unique per app (valid format)
		Clock:   realClock{},                      // Use real-time clock
		Delay:   500 * time.Millisecond,           // Polling interval
		Timeout: 1 * time.Second,                  // How long to wait for the mutex
	}

	// Try to acquire the mutex
	releaser, err := mutex.Acquire(spec)
	if err != nil {
		// If mutex acquisition fails, another instance is running
		showAlreadyRunningMessage()
		os.Exit(1)
	}
	defer releaser.Release() // Release mutex when the app exits

	// If we reach here, this is the only instance
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

// showAlreadyRunningMessage shows a dialog informing the user that another instance is running
func showAlreadyRunningMessage() {
	// Show a native Windows message box
	platform.ShowMessageBox("Click Guardian", "Click Guardian is already running.\n\nLook for the icon in your system tray or check your taskbar.")
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
