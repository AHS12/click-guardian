package gui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"

	"click-guardian/internal/config"
	"click-guardian/internal/hooks"
	"click-guardian/internal/logger"
)

// Application represents the main GUI application
type Application struct {
	app       fyne.App
	window    fyne.Window
	hook      hooks.MouseHook
	logger    *logger.Logger
	config    *config.Config
	isRunning bool
	isHidden  bool

	// UI components
	delayInput          *widget.Entry
	statusLabel         *widget.Label
	statusIcon          *canvas.Circle
	counterLabel        *widget.Label
	toggleButton        *widget.Button
	logText             *widget.RichText
	logContainer        *container.Scroll
	minimizeToTrayCheck *widget.Check
	appIconWidget       *widget.Icon
	updateChan          chan int

	// System tray
	trayRestore *systray.MenuItem
	trayQuit    *systray.MenuItem
}

// NewApplication creates a new GUI application
func NewApplication() *Application {
	a := app.New()
	a.SetIcon(GetAppIcon()) // Use our modern shield icon

	cfg := config.DefaultConfig()
	w := a.NewWindow("Click Guardian v1.0")

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
		updateChan:   make(chan int, 10),
	}
}

// Run starts the application
func (app *Application) Run() {
	app.setupUI()
	app.setupSystemTray()
	app.logger.Start()

	// Initialize log
	app.logger.Log("Click Guardian application started")
	if !app.hook.IsSupported() {
		app.logger.Log("❌ Mouse hooking not supported on this platform")
	} else {
		app.logger.Log("Enter a delay value and click 'Start Protection' to begin")
	}

	// Start a goroutine to update the blocked clicks counter
	go app.updateBlockedClicksCounter()

	// Start a goroutine to handle UI updates safely
	go app.handleUIUpdates()

	// Set window close behavior based on checkbox
	app.window.SetCloseIntercept(func() {
		if app.minimizeToTrayCheck.Checked {
			app.minimizeToTray()
		} else {
			app.quitApplication()
		}
	})

	app.window.ShowAndRun()

	// Cleanup when application finally quits
	app.cleanup()
}

func (app *Application) setupUI() {
	// Large app icon at the top
	app.appIconWidget = widget.NewIcon(GetAppIcon())
	app.appIconWidget.Resize(fyne.NewSize(80, 80))

	// App title - large and prominent
	appTitle := widget.NewLabel("CLICK GUARDIAN")
	appTitle.Alignment = fyne.TextAlignCenter
	appTitle.Importance = widget.HighImportance
	appTitle.TextStyle = fyne.TextStyle{Bold: true}

	// Status indicator - larger circle for better visibility
	app.statusIcon = canvas.NewCircle(color.RGBA{R: 220, G: 53, B: 69, A: 255}) // Red for stopped
	app.statusIcon.Resize(fyne.NewSize(32, 32))

	// Large toggle button styled like a modern toggle switch
	app.toggleButton = widget.NewButton("Start Protection", app.toggleProtection)
	app.toggleButton.Importance = widget.HighImportance
	app.toggleButton.Resize(fyne.NewSize(200, 60))

	// Status text below toggle
	app.statusLabel = widget.NewLabel("Protection Stopped")
	app.statusLabel.Alignment = fyne.TextAlignCenter
	app.statusLabel.Importance = widget.MediumImportance

	// Protection counter - prominently displayed
	app.counterLabel = widget.NewLabel("Blocked Clicks: 0")
	app.counterLabel.Alignment = fyne.TextAlignCenter
	app.counterLabel.Importance = widget.MediumImportance

	// Delay input with default value
	app.delayInput = widget.NewEntry()
	app.delayInput.SetText(fmt.Sprintf("%d", app.config.DelayMs))
	app.delayInput.SetPlaceHolder("Delay in milliseconds")

	// Minimize to tray checkbox
	app.minimizeToTrayCheck = widget.NewCheck("Minimize to system tray when closing", nil)
	app.minimizeToTrayCheck.SetChecked(true)

	// Clear log button
	clearButton := widget.NewButton("Clear Log", func() {
		app.logger.Clear()
	})

	// Layout
	// Header with very large app icon centered
	headerContainer := container.NewVBox(
		container.NewCenter(app.appIconWidget),
		container.NewCenter(appTitle),
	)

	// Main control section with toggle button prominently placed
	controlSection := container.NewVBox(
		container.NewCenter(app.toggleButton),
		container.NewCenter(app.statusLabel),
	)

	// Status and statistics section - prominently displayed
	statusStatsContainer := container.NewVBox(
		container.NewCenter(
			container.NewHBox(
				app.statusIcon,
				widget.NewLabel("  "), // spacing
			),
		),
		container.NewCenter(app.counterLabel),
	)

	// Settings and input grouped together
	configSection := container.NewVBox(
		widget.NewLabel("Configuration:"),
		container.NewBorder(
			widget.NewLabel("Delay (ms):"), nil, nil, nil,
			app.delayInput,
		),
		app.minimizeToTrayCheck,
	)

	// Log section with clear button integrated
	logHeaderContainer := container.NewBorder(
		nil, nil, widget.NewLabel("Activity Log:"), clearButton,
	)

	logSection := container.NewBorder(
		logHeaderContainer, nil, nil, nil,
		app.logContainer,
	)

	content := container.NewVBox(
		headerContainer,
		widget.NewSeparator(),
		controlSection,
		widget.NewSeparator(),
		statusStatsContainer,
		widget.NewSeparator(),
		configSection,
		widget.NewSeparator(),
		logSection,
	)

	app.window.SetContent(content)
	app.window.Resize(fyne.NewSize(float32(app.config.WindowWidth), float32(app.config.WindowHeight)))
	app.window.SetFixedSize(true) // Disable resizing and maximize button
	app.window.CenterOnScreen()
}

// toggleProtection handles both start and stop protection
func (app *Application) toggleProtection() {
	if app.isRunning {
		app.stopProtection()
	} else {
		app.startProtection()
	}
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
	app.statusIcon.FillColor = color.RGBA{R: 40, G: 167, B: 69, A: 255} // Green for active
	app.statusIcon.Refresh()
	app.statusLabel.SetText("Protection Active")
	app.toggleButton.SetText("Stop Protection")
	app.toggleButton.Importance = widget.DangerImportance

	app.logger.Log("Starting double-click protection with %d ms delay", delayMs)

	err = app.hook.Start(time.Duration(delayMs)*time.Millisecond, app.logger.GetChannel())
	if err != nil {
		app.logger.Log("❌ Failed to start protection: %v", err)
		app.statusIcon.FillColor = color.RGBA{R: 255, G: 193, B: 7, A: 255} // Yellow for failed
		app.statusIcon.Refresh()
		app.statusLabel.SetText("Protection Failed")
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
	app.statusIcon.FillColor = color.RGBA{R: 220, G: 53, B: 69, A: 255} // Red for stopped
	app.statusIcon.Refresh()
	app.statusLabel.SetText("Protection Stopped")
	app.toggleButton.SetText("Start Protection")
	app.toggleButton.Importance = widget.HighImportance
}

func (app *Application) cleanup() {
	if app.isRunning {
		app.hook.Stop()
	}
	app.logger.Stop()
	close(app.updateChan)
	systray.Quit()
}

// setupSystemTray initializes the system tray
func (app *Application) setupSystemTray() {
	go func() {
		systray.Run(app.onTrayReady, app.onTrayExit)
	}()
}

// onTrayReady is called when the system tray is ready
func (app *Application) onTrayReady() {
	// Set the system tray icon
	systray.SetIcon(getTrayIcon())
	systray.SetTitle("Click Guardian")
	systray.SetTooltip("Click Guardian - Double-click Protection")

	app.trayRestore = systray.AddMenuItem("Show Click Guardian", "Restore the application window")
	systray.AddSeparator()
	app.trayQuit = systray.AddMenuItem("Quit Application", "Completely quit the application")

	// Handle tray menu clicks
	go func() {
		for {
			select {
			case <-app.trayRestore.ClickedCh:
				app.showFromTray()
			case <-app.trayQuit.ClickedCh:
				app.quitApplication()
				return
			}
		}
	}()
}

// onTrayExit is called when the system tray exits
func (app *Application) onTrayExit() {
	// Cleanup if needed
}

// minimizeToTray hides the window and shows a notification
func (app *Application) minimizeToTray() {
	app.window.Hide()
	app.isHidden = true
	app.logger.Log("Application minimized to system tray")
}

// showFromTray shows the window from the system tray
func (app *Application) showFromTray() {
	app.window.Show()
	app.isHidden = false
	app.logger.Log("Application restored from system tray")
}

// quitApplication properly closes the application
func (app *Application) quitApplication() {
	app.cleanup()
	systray.Quit()
	app.app.Quit()
}

// updateBlockedClicksCounter continuously updates the blocked clicks display
func (app *Application) updateBlockedClicksCounter() {
	ticker := time.NewTicker(500 * time.Millisecond) // Update every 500ms
	defer ticker.Stop()

	for range ticker.C {
		if app.hook != nil {
			count := app.hook.GetBlockedCount()
			// Send the count to the UI update channel
			select {
			case app.updateChan <- count:
			default:
				// Don't block if channel is full
			}
		}
	}
}

// handleUIUpdates safely handles UI updates from the main thread
func (app *Application) handleUIUpdates() {
	for count := range app.updateChan {
		fyne.Do(func() {
			app.counterLabel.SetText(fmt.Sprintf("Blocked Clicks: %d", count))
		})
		
	}
}
