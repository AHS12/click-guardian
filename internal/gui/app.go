package gui

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"go-double-click-fix/internal/config"
	"go-double-click-fix/internal/hooks"
	"go-double-click-fix/internal/logger"
)

// Application represents the main GUI application
type Application struct {
	app       fyne.App
	window    fyne.Window
	hook      hooks.MouseHook
	logger    *logger.Logger
	config    *config.Config
	isRunning bool

	// UI components
	delayInput   *widget.Entry
	statusLabel  *widget.Label
	counterLabel *widget.Label
	startButton  *widget.Button
	stopButton   *widget.Button
	logText      *widget.RichText
	logContainer *container.Scroll
}

// NewApplication creates a new GUI application
func NewApplication() *Application {
	a := app.New()
	a.SetIcon(nil) // You can add an icon resource here if you have one

	cfg := config.DefaultConfig()
	w := a.NewWindow("Double Click Fix v1.0")

	// Create log display
	logText := widget.NewRichTextFromMarkdown("")
	logText.Wrapping = fyne.TextWrapWord
	logContainer := container.NewScroll(logText)
	logContainer.SetMinSize(fyne.NewSize(400, 200))

	logger := logger.NewLogger(logText, logContainer, cfg.MaxLogLines)

	return &Application{
		app:          a,
		window:       w,
		hook:         hooks.NewMouseHook(),
		logger:       logger,
		config:       cfg,
		logText:      logText,
		logContainer: logContainer,
	}
}

// Run starts the application
func (app *Application) Run() {
	app.setupUI()
	app.logger.Start()

	// Initialize log
	app.logger.Log("Double-Click Fix application started")
	if !app.hook.IsSupported() {
		app.logger.Log("❌ Mouse hooking not supported on this platform")
	} else {
		app.logger.Log("Enter a delay value and click 'Start Protection' to begin")
	}

	app.window.ShowAndRun()

	// Cleanup when window closes
	app.cleanup()
}

func (app *Application) setupUI() {
	// Delay input with default value
	app.delayInput = widget.NewEntry()
	app.delayInput.SetText(fmt.Sprintf("%d", app.config.DelayMs))
	app.delayInput.SetPlaceHolder("Delay in milliseconds (e.g., 50)")

	// Status label
	app.statusLabel = widget.NewLabel("Status: Stopped")
	app.statusLabel.Importance = widget.MediumImportance

	// Blocked clicks counter
	app.counterLabel = widget.NewLabel("Blocked Clicks: 0")
	app.counterLabel.Importance = widget.LowImportance

	// Clear log button
	clearButton := widget.NewButton("Clear Log", func() {
		app.logger.Clear()
	})

	// Control buttons
	app.startButton = widget.NewButton("Start Protection", app.startProtection)
	app.stopButton = widget.NewButton("Stop Protection", app.stopProtection)

	// Style buttons
	app.startButton.Importance = widget.HighImportance
	app.stopButton.Importance = widget.MediumImportance

	// Layout
	inputForm := container.NewBorder(
		widget.NewLabel("Double-Click Delay (ms):"), nil, nil, nil,
		app.delayInput,
	)

	buttonContainer := container.NewHBox(
		app.startButton,
		app.stopButton,
		widget.NewSeparator(),
		clearButton,
	)

	logSection := container.NewBorder(
		widget.NewLabel("Activity Log:"), nil, nil, nil,
		app.logContainer,
	)

	content := container.NewVBox(
		app.statusLabel,
		widget.NewSeparator(),
		inputForm,
		buttonContainer,
		widget.NewSeparator(),
		logSection,
	)

	app.window.SetContent(content)
	app.window.Resize(fyne.NewSize(float32(app.config.WindowWidth), float32(app.config.WindowHeight)))
	app.window.CenterOnScreen()
}

func (app *Application) startProtection() {
	if app.isRunning {
		return
	}

	if !app.hook.IsSupported() {
		app.logger.Log("❌ Mouse hooking not supported on this platform")
		return
	}

	delayText := strings.TrimSpace(app.delayInput.Text)
	delayMs, err := config.ParseDelay(delayText)
	if err != nil {
		app.logger.Log("ERROR: %v", err)
		return
	}

	// Update input field with validated value
	if delayText == "" {
		app.delayInput.SetText(fmt.Sprintf("%d", delayMs))
	}

	app.isRunning = true
	app.statusLabel.SetText("Status: Running")
	app.startButton.SetText("Running...")
	app.startButton.Disable()

	app.logger.Log("Starting double-click protection with %d ms delay", delayMs)

	err = app.hook.Start(time.Duration(delayMs)*time.Millisecond, app.logger.GetChannel())
	if err != nil {
		app.logger.Log("❌ Failed to start protection: %v", err)
		app.resetUI()
	}
}

func (app *Application) stopProtection() {
	if !app.isRunning {
		return
	}

	app.isRunning = false
	app.logger.Log("Stopping double-click protection")
	app.hook.Stop()
	app.resetUI()
}

func (app *Application) resetUI() {
	app.isRunning = false
	app.statusLabel.SetText("Status: Stopped")
	app.startButton.SetText("Start Protection")
	app.startButton.Enable()
}

func (app *Application) cleanup() {
	if app.isRunning {
		app.hook.Stop()
	}
	app.logger.Stop()
}
