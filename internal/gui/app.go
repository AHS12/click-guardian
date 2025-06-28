package gui

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/systray"

	"click-guardian/internal/config"
	"click-guardian/internal/hooks"
	"click-guardian/internal/logger"
)

// Application represents the main GUI application
type Application struct {
	app              fyne.App
	window           fyne.Window
	hook             hooks.MouseHook
	logger           *logger.Logger
	config           *config.Config
	isRunning        bool
	isHidden         bool
	lastBlockedCount int

	// UI components
	delaySlider         *widget.Slider
	delayValueLabel     *widget.Label
	statusLabel         *widget.Label
	statusIcon          *canvas.Circle
	counterText         *canvas.Text
	blockedLabelText    *canvas.Text
	toggleButton        *widget.Button
	logText             *widget.RichText
	logContainer        *container.Scroll
	minimizeToTrayCheck *widget.Check
	// appIconWidget       *widget.Icon
	updateChan     chan int
	updateChanOnce sync.Once

	// System tray
	trayRestore *systray.MenuItem
	trayQuit    *systray.MenuItem
}

// NewApplication creates a new GUI application
func NewApplication() *Application {
	a := app.New()
	a.SetIcon(GetAppIcon()) // Use our modern shield icon

	cfg := config.DefaultConfig()
	w := a.NewWindow("Click Guardian")

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
	// Status indicator circle
	app.statusIcon = canvas.NewCircle(color.RGBA{R: 220, G: 53, B: 69, A: 255}) // Red for stopped

	// Large toggle button
	app.toggleButton = widget.NewButton("Start Protection", app.toggleProtection)
	app.toggleButton.Importance = widget.HighImportance

	// Status label (will be placed inside the circle)
	app.statusLabel = widget.NewLabel("Protection Stopped")
	app.statusLabel.Alignment = fyne.TextAlignCenter

	// Counter label (will be placed inside the circle)
	app.counterText = canvas.NewText("0", color.White)
	app.counterText.Alignment = fyne.TextAlignCenter
	app.counterText.TextStyle = fyne.TextStyle{Bold: true}
	app.counterText.TextSize = 48

	// Blocked label text
	app.blockedLabelText = canvas.NewText("Blocked Double Clicks", color.White)
	app.blockedLabelText.Alignment = fyne.TextAlignCenter
	app.blockedLabelText.TextSize = 14

	// Create the status indicator by stacking the circle and the labels
	statusIndicator := container.NewStack(
		app.statusIcon,
		container.NewCenter(container.NewVBox(
			app.blockedLabelText,
			app.counterText,
			layout.NewSpacer(),
		)),
	)

	// Delay slider with min 5ms, max 500ms
	app.delaySlider = widget.NewSlider(5, 500)
	app.delaySlider.SetValue(float64(app.config.DelayMs))
	app.delaySlider.Step = 5 // 5ms increments

	// Value label to show current slider value
	app.delayValueLabel = widget.NewLabel(fmt.Sprintf("%d ms", app.config.DelayMs))
	app.delayValueLabel.Alignment = fyne.TextAlignCenter

	// Update label when slider value changes
	app.delaySlider.OnChanged = func(value float64) {
		app.delayValueLabel.SetText(fmt.Sprintf("%.0f ms", value))
	}

	// Minimize to tray checkbox
	app.minimizeToTrayCheck = widget.NewCheck("Minimize to system tray when closing", nil)
	app.minimizeToTrayCheck.SetChecked(true)

	// Clear log button
	clearButton := widget.NewButton("Clear Log", func() {
		app.logger.Clear()
	})

	// --- Layout ---

	// Main control section with toggle button
	controlSection := container.NewVBox(
		container.NewCenter(app.toggleButton),
	)

	// Status and statistics section with the new indicator
	hoverIndicator := NewHoverAware(statusIndicator, func() string {
		if app.isRunning {
			return "Status: Active"
		} else {
			return "Status: Inactive"
		}
	})
	statusIndicatorContainer := container.NewGridWrap(fyne.NewSize(200, 200), hoverIndicator)

	statusStatsContainer := container.NewVBox(
		container.NewCenter(statusIndicatorContainer),
	)

	// Settings and input grouped together
	configTitle := canvas.NewText("Configuration", color.White)
	configTitle.TextStyle = fyne.TextStyle{Bold: true}
	configTitle.TextSize = 16 // Reduced font size
	configHeader := container.NewHBox(widget.NewIcon(theme.SettingsIcon()), configTitle)

	configContent := container.NewVBox(
		container.NewVBox(
			widget.NewLabel("Delay (ms):"),
			app.delaySlider,
			container.NewCenter(app.delayValueLabel),
		),
		app.minimizeToTrayCheck,
	)

	configSection := widget.NewCard("", "", container.NewVBox(
		configHeader, // Left-aligned
		configContent,
	))

	// Log section with clear button integrated
	logTitle := canvas.NewText("Activity Log", color.White)
	logTitle.TextStyle = fyne.TextStyle{Bold: true}
	logTitle.TextSize = 16
	logHeader := container.NewHBox(widget.NewIcon(theme.ListIcon()), logTitle, layout.NewSpacer(), clearButton)

	logSection := widget.NewCard("", "", container.NewVBox(
		logHeader,
		app.logContainer,
	))

	// Assemble the final content
	content := container.NewVBox(
		canvas.NewRectangle(color.Transparent), // Top spacer
		statusStatsContainer,
		canvas.NewRectangle(color.Transparent), // Bottom spacer
		controlSection,
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

	// Get delay value from slider
	delayMs := int(app.delaySlider.Value)

	// Disable slider when protection is active
	app.delaySlider.Disable()

	app.isRunning = true
	app.statusIcon.FillColor = color.RGBA{R: 40, G: 167, B: 69, A: 255} // Green for active
	app.statusIcon.Refresh()
	// app.statusLabel.SetText("Protection Active")
	app.toggleButton.SetText("Stop Protection")
	app.toggleButton.Importance = widget.DangerImportance

	app.logger.Log("Starting double-click protection with %d ms delay", delayMs)

	err := app.hook.Start(time.Duration(delayMs)*time.Millisecond, app.logger.GetChannel())
	if err != nil {
		app.logger.Log("❌ Failed to start protection: %v", err)
		app.statusIcon.FillColor = color.RGBA{R: 255, G: 193, B: 7, A: 255} // Yellow for failed
		app.statusIcon.Refresh()
		// app.statusLabel.SetText("Protection Failed")
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
	// app.statusLabel.SetText("Protection Stopped")
	app.toggleButton.SetText("Start Protection")
	app.toggleButton.Importance = widget.HighImportance

	// Re-enable slider when protection is stopped
	app.delaySlider.Enable()
}

func (app *Application) cleanup() {
	if app.isRunning {
		app.hook.Stop()
	}
	app.logger.Stop()
	app.updateChanOnce.Do(func() {
		close(app.updateChan)
	})
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
			// Animate if the count has increased
			if count > app.lastBlockedCount {
				originalSize := app.counterText.TextSize
				anim := canvas.NewSizeAnimation(
					fyne.NewSize(originalSize, originalSize),
					fyne.NewSize(originalSize+10, originalSize+10),
					time.Millisecond*100,
					func(s fyne.Size) {
						app.counterText.TextSize = s.Height
						app.counterText.Refresh()
					},
				)
				anim.AutoReverse = true
				anim.Start()
			}

			app.lastBlockedCount = count
			app.counterText.Text = fmt.Sprintf("%d", count)
			app.counterText.Refresh()
		})
	}
}
