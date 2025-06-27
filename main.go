package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	a.SetIcon(nil) // You can add an icon resource here if you have one
	w := a.NewWindow("Double Click Fix v1.0")

	logChannel := make(chan string, 100)
	var isRunning bool

	// Delay input with default value
	delayInput := widget.NewEntry()
	delayInput.SetText("50") // Default 50ms
	delayInput.SetPlaceHolder("Delay in milliseconds (e.g., 50)")

	// Status label
	statusLabel := widget.NewLabel("Status: Stopped")
	statusLabel.Importance = widget.MediumImportance

	// Blocked clicks counter
	counterLabel := widget.NewLabel("Blocked Clicks: 0")
	counterLabel.Importance = widget.LowImportance

	// Log display with rich text
	logText := widget.NewRichTextFromMarkdown("")
	logText.Wrapping = fyne.TextWrapWord
	logContainer := container.NewScroll(logText)
	logContainer.SetMinSize(fyne.NewSize(400, 200))

	// Clear log button
	clearButton := widget.NewButton("Clear Log", func() {
		logText.ParseMarkdown("**Log cleared**\n\n")
	})

	// Declare buttons first
	var startButton, stopButton *widget.Button

	startButton = widget.NewButton("Start Protection", func() {
		if isRunning {
			return
		}

		delayText := strings.TrimSpace(delayInput.Text)
		if delayText == "" {
			delayText = "50"
			delayInput.SetText("50")
		}

		delayMs, err := strconv.Atoi(delayText)
		if err != nil || delayMs < 1 || delayMs > 5000 {
			logChannel <- "ERROR: Invalid delay. Please enter a value between 1 and 5000 milliseconds."
			return
		}

		isRunning = true
		statusLabel.SetText("Status: Running")
		startButton.SetText("Running...")
		startButton.Disable()

		logChannel <- fmt.Sprintf("Starting double-click protection with %d ms delay", delayMs)
		startHook(time.Duration(delayMs)*time.Millisecond, logChannel)
	})

	stopButton = widget.NewButton("Stop Protection", func() {
		if !isRunning {
			return
		}

		isRunning = false
		statusLabel.SetText("Status: Stopped")
		startButton.SetText("Start Protection")
		startButton.Enable()

		logChannel <- "Stopping double-click protection"
		stopHook()
	})

	// Style buttons
	startButton.Importance = widget.HighImportance
	stopButton.Importance = widget.MediumImportance

	// Handle log messages
	go func() {
		var logHistory []string
		for msg := range logChannel {
			timestamp := time.Now().Format("15:04:05")
			logEntry := fmt.Sprintf("**[%s]** %s", timestamp, msg)

			logHistory = append(logHistory, logEntry)

			// Keep only the last 100 entries to prevent memory issues
			if len(logHistory) > 100 {
				logHistory = logHistory[1:]
			}

			// Update the log display on the main UI thread
			allLogs := strings.Join(logHistory, "\n\n")
			fyne.Do(func() {
				logText.ParseMarkdown(allLogs)
				logContainer.ScrollToBottom()
			})
		}
	}()

	// Layout
	inputForm := container.NewBorder(
		widget.NewLabel("Double-Click Delay (ms):"), nil, nil, nil,
		delayInput,
	)

	buttonContainer := container.NewHBox(
		startButton,
		stopButton,
		widget.NewSeparator(),
		clearButton,
	)

	logSection := container.NewBorder(
		widget.NewLabel("Activity Log:"), nil, nil, nil,
		logContainer,
	)

	content := container.NewVBox(
		statusLabel,
		widget.NewSeparator(),
		inputForm,
		buttonContainer,
		widget.NewSeparator(),
		logSection,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(500, 450))
	w.CenterOnScreen()

	// Initialize log
	logChannel <- "Double-Click Fix application started"
	logChannel <- "Enter a delay value and click 'Start Protection' to begin"

	w.ShowAndRun()

	// Cleanup when window closes
	close(logChannel)
}
