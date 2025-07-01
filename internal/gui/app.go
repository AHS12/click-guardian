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
	"click-guardian/internal/gui/components"
	"click-guardian/internal/gui/dialogs"
	"click-guardian/internal/gui/resources"
	"click-guardian/internal/hooks"
	"click-guardian/internal/logger"
	"click-guardian/internal/version"
	"click-guardian/pkg/platform"
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
	autoStartCheck      *widget.Check
	updateChan     chan int
	updateChanOnce sync.Once

	// System tray
	trayRestore *systray.MenuItem
	trayQuit    *systray.MenuItem

	// Cleanup control
	shutdownChan chan struct{}
	shutdownOnce sync.Once

	// Store checkbox states to avoid UI thread issues during shutdown
	minimizeToTrayEnabled bool
}

// NewApplication creates a new GUI application
func NewApplication() *Application {
	a := app.New()
	a.SetIcon(resources.GetAppIcon()) // Use our modern shield icon

	cfg := config.LoadConfig() // Load saved config instead of default

	// Set window title with version
	windowTitle := fmt.Sprintf("Click Guardian v%s", version.GetVersionString())
	w := a.NewWindow(windowTitle)

	// Create log display
	logText := widget.NewRichTextFromMarkdown("")
	logText.Wrapping = fyne.TextWrapWord
	logContainer := container.NewScroll(logText)
	logContainer.SetMinSize(fyne.NewSize(400, 200))

	logger := logger.NewLogger(logText, logContainer, cfg.MaxLogLines)

	return &Application{
		app:                   a,
		window:                w,
		hook:                  hooks.NewMouseHook(),
		logger:                logger,
		config:                cfg,
		logText:               logText,
		logContainer:          logContainer,
		updateChan:            make(chan int, 10),
		shutdownChan:          make(chan struct{}),
		minimizeToTrayEnabled: cfg.MinimizeToTray, // Use saved preference
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
	app.logger.Log("✅ Settings loaded: %dms delay, minimize to tray: %v", app.config.DelayMs, app.config.MinimizeToTray)

	// Start a goroutine to update the blocked clicks counter
	go app.updateBlockedClicksCounter()

	// Start a goroutine to handle UI updates safely
	go app.handleUIUpdates()

	// Set window close behavior based on stored state (avoid UI thread issues)
	app.window.SetCloseIntercept(app.handleWindowClose)

	app.window.ShowAndRun()

	// Cleanup when application finally quits
	app.cleanup()
}

// RunMinimized starts the application minimized to system tray
func (app *Application) RunMinimized() {
	app.setupUI()
	app.setupSystemTray()
	app.logger.Start()

	// Initialize log
	app.logger.Log("Click Guardian application started (minimized)")
	if !app.hook.IsSupported() {
		app.logger.Log("❌ Mouse hooking not supported on this platform")
	} else {
		app.logger.Log("Application started minimized to system tray")

		// Auto-start protection when launched minimized (from Windows startup)
		go func() {
			// Small delay to ensure everything is initialized
			time.Sleep(1 * time.Second)
			app.startProtection()
			app.logger.Log("🚀 Protection auto-started on Windows startup")
		}()
	}
	app.logger.Log("✅ Settings loaded: %dms delay, minimize to tray: %v", app.config.DelayMs, app.config.MinimizeToTray)

	// Start a goroutine to update the blocked clicks counter
	go app.updateBlockedClicksCounter()

	// Start a goroutine to handle UI updates safely
	go app.handleUIUpdates()

	// Set window close behavior based on stored state (avoid UI thread issues)
	app.window.SetCloseIntercept(app.handleWindowClose)

	// Start minimized to tray
	app.isHidden = true

	// Show window briefly then hide it (required for Fyne to initialize properly)
	app.window.Resize(fyne.NewSize(float32(app.config.WindowWidth), float32(app.config.WindowHeight)))
	app.window.SetFixedSize(true)
	app.window.CenterOnScreen()

	// Use a goroutine to hide the window after showing
	go func() {
		fyne.Do(func() {
			app.window.Show()
		})
		// Small delay to ensure window is properly initialized
		time.Sleep(100 * time.Millisecond)
		fyne.Do(func() {
			app.window.Hide()
		})
	}()

	// Run the application
	app.window.ShowAndRun()

	// Cleanup when application finally quits
	app.cleanup()
}

// RunWithAutoProtect starts the application and automatically enables protection
func (app *Application) RunWithAutoProtect() {
	app.setupUI()
	app.setupSystemTray()
	app.logger.Start()

	// Initialize log
	app.logger.Log("Click Guardian application started with auto-protect")
	if !app.hook.IsSupported() {
		app.logger.Log("❌ Mouse hooking not supported on this platform")
	} else {
		app.logger.Log("Enter a delay value or protection will start automatically")

		// Auto-start protection
		go func() {
			// Small delay to ensure everything is initialized
			time.Sleep(1 * time.Second)
			app.startProtection()
			app.logger.Log("🚀 Protection auto-started")
		}()
	}
	app.logger.Log("✅ Settings loaded: %dms delay, minimize to tray: %v", app.config.DelayMs, app.config.MinimizeToTray)

	// Start a goroutine to update the blocked clicks counter
	go app.updateBlockedClicksCounter()

	// Start a goroutine to handle UI updates safely
	go app.handleUIUpdates()

	// Set window close behavior based on stored state (avoid UI thread issues)
	app.window.SetCloseIntercept(app.handleWindowClose)

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
		// Save the new delay value to config
		app.config.DelayMs = int(value)
		if err := app.config.Save(); err != nil {
			app.logger.Log("⚠️ Failed to save delay setting: %v", err)
		}
	}

	// Minimize to tray checkbox
	app.minimizeToTrayCheck = widget.NewCheck("Minimize to system tray when closing", func(checked bool) {
		app.minimizeToTrayEnabled = checked
		// Save the minimize to tray preference
		app.config.MinimizeToTray = checked
		if err := app.config.Save(); err != nil {
			app.logger.Log("⚠️ Failed to save minimize to tray setting: %v", err)
		}
	})
	app.minimizeToTrayCheck.SetChecked(app.config.MinimizeToTray)

	// Auto-start checkbox
	app.autoStartCheck = widget.NewCheck("Start with Windows and auto-enable protection", app.onAutoStartChanged)
	app.autoStartCheck.SetChecked(false) // Default unchecked

	// Remove the separate auto-protect checkbox - it's now integrated

	// Check current auto-start status and update checkbox
	app.updateAutoStartStatus()

	// Clear log button
	clearButton := widget.NewButton("Clear Log", func() {
		app.logger.Clear()
	})

	// About button
	aboutButton := widget.NewButton("About", func() {
		dialogs.ShowAboutDialog(app.window)
	})

	// --- Layout ---

	// Main control section with toggle button
	controlSection := container.NewVBox(
		container.NewCenter(app.toggleButton),
	)

	// Status and statistics section with the new indicator
	hoverIndicator := components.NewHoverAware(statusIndicator, func() string {
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
		app.autoStartCheck,
	)

	configSection := widget.NewCard("", "", container.NewVBox(
		configHeader, // Left-aligned
		configContent,
	))

	// Log section with clear button integrated
	logTitle := canvas.NewText("Activity Log", color.White)
	logTitle.TextStyle = fyne.TextStyle{Bold: true}
	logTitle.TextSize = 16
	logHeader := container.NewHBox(widget.NewIcon(theme.ListIcon()), logTitle, layout.NewSpacer(), clearButton, aboutButton)

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

	// Update UI elements on main thread
	fyne.Do(func() {
		// Disable slider when protection is active
		app.delaySlider.Disable()

		app.statusIcon.FillColor = color.RGBA{R: 40, G: 167, B: 69, A: 255} // Green for active
		app.statusIcon.Refresh()
		// app.statusLabel.SetText("Protection Active")
		app.toggleButton.SetText("Stop Protection")
		app.toggleButton.Importance = widget.DangerImportance
	})

	app.isRunning = true
	app.logger.Log("Starting double-click protection with %d ms delay", delayMs)

	err := app.hook.Start(time.Duration(delayMs)*time.Millisecond, app.logger.GetChannel())
	if err != nil {
		app.logger.Log("❌ Failed to start protection: %v", err)
		fyne.Do(func() {
			app.statusIcon.FillColor = color.RGBA{R: 255, G: 193, B: 7, A: 255} // Yellow for failed
			app.statusIcon.Refresh()
		})
		// app.statusLabel.SetText("Protection Failed")
		app.resetUI()
	} else {
		// Update tray tooltip when protection starts successfully
		app.updateTrayTooltip()
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

	fyne.Do(func() {
		app.statusIcon.FillColor = color.RGBA{R: 220, G: 53, B: 69, A: 255} // Red for stopped
		app.statusIcon.Refresh()
		// app.statusLabel.SetText("Protection Stopped")
		app.toggleButton.SetText("Start Protection")
		app.toggleButton.Importance = widget.HighImportance

		// Re-enable slider when protection is stopped
		app.delaySlider.Enable()
	})

	// Update tray tooltip when protection stops
	app.updateTrayTooltip()
}

func (app *Application) cleanup() {
	app.shutdownOnce.Do(func() {
		fmt.Println("Cleanup called from main thread")
		// This is now mainly for the normal app termination path
		// The quitApplication function handles immediate shutdown
	})
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
	systray.SetIcon(resources.GetTrayIcon().Content())
	systray.SetTitle("Click Guardian")

	// Set initial tooltip
	app.updateTrayTooltip()

	app.trayRestore = systray.AddMenuItem("Show Click Guardian", "Restore the application window")
	systray.AddSeparator()
	app.trayQuit = systray.AddMenuItem("Quit Application", "Completely quit the application")

	// Handle tray menu clicks
	go func() {
		defer func() {
			// Cleanup when this goroutine exits
			if r := recover(); r != nil {
				fmt.Printf("Tray menu handler recovered from panic: %v\n", r)
			}
		}()

		for {
			select {
			case <-app.trayRestore.ClickedCh:
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("Recovery during show from tray: %v\n", r)
						}
					}()
					app.showFromTray()
				}()
			case <-app.trayQuit.ClickedCh:
				func() {
					defer func() {
						if r := recover(); r != nil {
							fmt.Printf("Recovery during quit from tray: %v\n", r)
						}
					}()
					app.quitApplication()
				}()
				return
			case <-app.shutdownChan:
				// Graceful shutdown
				return
			}
		}
	}()
}

// onTrayExit is called when the system tray exits
func (app *Application) onTrayExit() {
	// Ensure cleanup is called when tray exits
	app.cleanup()
}

// minimizeToTray hides the window and shows a notification
func (app *Application) minimizeToTray() {
	fyne.Do(func() {
		app.window.Hide()
	})
	app.isHidden = true
	// Don't log during potential shutdown to avoid deadlocks
	if app.logger != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignore logger errors during shutdown
				}
			}()
			app.logger.Log("Application minimized to system tray")
		}()
	}
}

// showFromTray shows the window from the system tray
func (app *Application) showFromTray() {
	fyne.Do(func() {
		app.window.Show()
	})
	app.isHidden = false
	// Don't log during potential shutdown to avoid deadlocks
	if app.logger != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Ignore logger errors during shutdown
				}
			}()
			app.logger.Log("Application restored from system tray")
		}()
	}
}

// quitApplication properly closes the application
func (app *Application) quitApplication() {
	fmt.Println("User requested application quit")

	// Stop everything immediately and directly
	if app.isRunning {
		fmt.Println("Stopping mouse hook protection...")
		app.hook.Stop()
		app.isRunning = false
	}

	// Signal all goroutines to stop
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovery during shutdown signal: %v\n", r)
			}
		}()
		close(app.shutdownChan)
	}()

	// Stop logger
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovery during logger stop: %v\n", r)
			}
		}()
		app.logger.Stop()
	}()

	// Close update channel
	app.updateChanOnce.Do(func() {
		close(app.updateChan)
	})

	// Quit system tray
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovery during systray quit: %v\n", r)
			}
		}()
		systray.Quit()
	}()

	fmt.Println("Cleanup completed, quitting app...")

	// Quit the Fyne app directly on the main thread
	fyne.Do(func() {
		app.app.Quit()
	})
}

// handleWindowClose handles the window close event safely
func (app *Application) handleWindowClose() {
	if app.minimizeToTrayEnabled {
		app.minimizeToTray()
	} else {
		app.quitApplication()
	}
}

// updateBlockedClicksCounter continuously updates the blocked clicks display
func (app *Application) updateBlockedClicksCounter() {
	ticker := time.NewTicker(500 * time.Millisecond) // Update every 500ms
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if app.hook != nil {
				count := app.hook.GetBlockedCount()
				// Send the count to the UI update channel
				select {
				case app.updateChan <- count:
				default:
					// Don't block if channel is full
				}
			}
		case <-app.shutdownChan:
			// Graceful shutdown
			return
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

				// Update tray tooltip when count changes (only if protection is active)
				if app.isRunning {
					blockedCount := app.hook.GetBlockedCount()
					tooltip := fmt.Sprintf("Click Guardian - Active\nBlocked clicks: %d", blockedCount)
					systray.SetTooltip(tooltip)
				}
			}

			app.lastBlockedCount = count
			app.counterText.Text = fmt.Sprintf("%d", count)
			app.counterText.Refresh()
		})
	}
}

// updateTrayTooltip updates the system tray tooltip with current status and blocked count
func (app *Application) updateTrayTooltip() {
	// Ensure this runs on the main thread
	fyne.Do(func() {
		if app.isRunning {
			blockedCount := app.hook.GetBlockedCount()
			tooltip := fmt.Sprintf("Click Guardian - Active\nBlocked clicks: %d", blockedCount)
			systray.SetTooltip(tooltip)
		} else {
			systray.SetTooltip("Click Guardian - Inactive")
		}
	})
}

// onAutoStartChanged handles the auto-start checkbox state change
func (app *Application) onAutoStartChanged(checked bool) {
	if checked {
		err := platform.EnableAutoStart()
		if err != nil {
			app.logger.Log("❌ Failed to enable auto-start: %v", err)
			// Revert checkbox state if failed
			app.autoStartCheck.SetChecked(false)
		} else {
			app.logger.Log("✅ Auto-start with Windows enabled")
		}
	} else {
		err := platform.DisableAutoStart()
		if err != nil {
			app.logger.Log("❌ Failed to disable auto-start: %v", err)
			// Revert checkbox state if failed
			app.autoStartCheck.SetChecked(true)
		} else {
			app.logger.Log("✅ Auto-start with Windows disabled")
		}
	}
}

// updateAutoStartStatus checks if auto-start is currently enabled and updates the checkbox
func (app *Application) updateAutoStartStatus() {
	isEnabled := platform.IsAutoStartEnabled()
	app.autoStartCheck.SetChecked(isEnabled)
}
