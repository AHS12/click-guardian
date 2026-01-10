package dialogs

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"click-guardian/internal/config"
	"click-guardian/internal/hooks"
)

type SettingsCallbacks struct {
	OnDelayChanged            func(int)
	OnMinimizeToTrayChanged   func(bool)
	OnAutoStartChanged        func(bool)
	OnProtectedButtonsChanged func([]string)
	OnDragFixChanged          func(bool)
	OnDragFixThresholdChanged func(int)
	OnPauseDurationChanged    func(int)
}

// ShowSettingsDialog shows a comprehensive settings dialog
func ShowSettingsDialog(parent fyne.Window, cfg *config.Config, initialAutoStart bool, callbacks SettingsCallbacks) {
	// --- Tab 1: General ---

	// Delay Slider
	delayLabel := widget.NewLabel(fmt.Sprintf("%d ms", cfg.DelayMs))
	delaySlider := widget.NewSlider(5, 500)
	delaySlider.Value = float64(cfg.DelayMs)
	delaySlider.Step = 5
	delaySlider.OnChanged = func(v float64) {
		delayLabel.SetText(fmt.Sprintf("%.0f ms", v))
		callbacks.OnDelayChanged(int(v))
	}

	delayContainer := container.NewVBox(
		widget.NewLabel("Double-Click Threshold"),
		widget.NewRichTextFromMarkdown("_Adjust the time window for detecting double-clicks._"),
		container.NewBorder(nil, nil, nil, delayLabel, delaySlider),
	)

	// Application Behavior
	minTrayCheck := widget.NewCheck("Minimize to system tray when closing", func(checked bool) {
		callbacks.OnMinimizeToTrayChanged(checked)
	})
	minTrayCheck.Checked = cfg.MinimizeToTray

	autoStartCheck := widget.NewCheck("Start with Windows", func(checked bool) {
		callbacks.OnAutoStartChanged(checked)
	})
	autoStartCheck.Checked = initialAutoStart

	behaviorContainer := container.NewVBox(
		widget.NewLabel("Application Behavior"),
		minTrayCheck,
		autoStartCheck,
	)

	generalContent := container.NewVBox(
		delayContainer,
		widget.NewSeparator(),
		behaviorContainer,
	)

	// --- Tab 2: Protected Buttons ---

	availableButtons := hooks.GetMouseButtons()

	// Create checkboxes for each mouse button
	checkBoxes := make([]*widget.Check, len(availableButtons))
	buttonOptions := make([]fyne.CanvasObject, len(availableButtons))

	saveSelection := func() {
		var selected []string
		for j, cb := range checkBoxes {
			if cb.Checked {
				selected = append(selected, availableButtons[j].ID)
			}
		}

		// Ensure at least one button is selected
		if len(selected) == 0 {
			selected = []string{"left"}
			for j, button := range availableButtons {
				if button.ID == "left" {
					checkBoxes[j].SetChecked(true)
					break
				}
			}
		}

		callbacks.OnProtectedButtonsChanged(selected)
	}

	for i, button := range availableButtons {
		option, check := createButtonOption(button.Name, button.ID, cfg)
		checkBoxes[i] = check
		buttonOptions[i] = option

		check.OnChanged = func(checked bool) {
			saveSelection()
		}
	}

	buttonsContent := container.NewVBox(
		widget.NewLabel("Select buttons to protect:"),
		container.NewVBox(buttonOptions...),
	)

	// --- Tab 3: Advanced (Drag Fix) ---

	dragFixCheck := widget.NewCheck("Enable Drag Fix (Experimental)", func(checked bool) {
		callbacks.OnDragFixChanged(checked)
	})
	// Drag Fix Threshold Slider
	dragThresholdLabel := widget.NewLabel(fmt.Sprintf("%d ms", cfg.DragFixThreshold))
	dragThresholdSlider := widget.NewSlider(5, 100)
	dragThresholdSlider.Value = float64(cfg.DragFixThreshold)
	dragThresholdSlider.Step = 5
	dragThresholdSlider.OnChanged = func(v float64) {
		dragThresholdLabel.SetText(fmt.Sprintf("%.0f ms", v))
		callbacks.OnDragFixThresholdChanged(int(v))
	}

	// Pause Duration Slider (1-5s)
	pauseDurationLabel := widget.NewLabel(fmt.Sprintf("%d s", cfg.PauseDuration))
	pauseDurationSlider := widget.NewSlider(1, 5)
	pauseDurationSlider.Value = float64(cfg.PauseDuration)
	pauseDurationSlider.Step = 1
	pauseDurationSlider.OnChanged = func(v float64) {
		pauseDurationLabel.SetText(fmt.Sprintf("%.0f s", v))
		callbacks.OnPauseDurationChanged(int(v))
	}

	// Helper to update slider state
	updateSliderState := func(enabled bool) {
		if enabled {
			dragThresholdSlider.Enable()
			pauseDurationSlider.Enable()
		} else {
			dragThresholdSlider.Disable()
			pauseDurationSlider.Disable()
		}
	}

	dragFixCheck = widget.NewCheck("Enable Drag Fix (Experimental)", func(checked bool) {
		callbacks.OnDragFixChanged(checked)
		updateSliderState(checked)
	})
	dragFixCheck.Checked = cfg.DragFix

	// Initialize state
	updateSliderState(cfg.DragFix)

	dragFixContainer := container.NewVBox(
		dragFixCheck,
		widget.NewSeparator(),
		widget.NewLabel("Drag Fix Threshold"),
		container.NewBorder(nil, nil, nil, dragThresholdLabel, dragThresholdSlider),
		widget.NewSeparator(),
		widget.NewLabel("High Privilege Pause Duration"),
		container.NewBorder(nil, nil, nil, pauseDurationLabel, pauseDurationSlider),
	)

	dragFixInfo := widget.NewRichTextFromMarkdown(
		"**Drag Fix (Anti-Bounce):**\n" +
			"Prevents accidental drops by ignoring momentary release signals (bouncing) during drags.\n\n" +
			"**High Privilege Pause:**\n" +
			"Pauses protection when interacting with Admin apps (e.g. VMware) to prevent stuck drags.")
	dragFixInfo.Wrapping = fyne.TextWrapWord

	advancedContent := container.NewVBox(
		container.NewPadded(dragFixContainer),
		widget.NewCard("Information", "", container.NewPadded(dragFixInfo)),
	)

	// --- Assemble Dialog ---

	tabs := container.NewAppTabs(
		container.NewTabItem("General", container.NewPadded(generalContent)),
		container.NewTabItem("Buttons", container.NewPadded(buttonsContent)),
		container.NewTabItem("Advanced", container.NewPadded(advancedContent)),
	)

	d := dialog.NewCustom("Settings", "Close", container.NewPadded(tabs), parent)
	d.Resize(fyne.NewSize(400, 500))
	d.Show()
}

func createButtonOption(buttonName, buttonID string, currentConfig *config.Config) (fyne.CanvasObject, *widget.Check) {
	// Check if this button is currently protected
	isProtected := false
	for _, protectedButton := range currentConfig.ProtectedButtons {
		if protectedButton == buttonID {
			isProtected = true
			break
		}
	}

	// Create a checkbox for this button
	toggle := widget.NewCheck("", nil)
	toggle.SetChecked(isProtected)

	// Create button name label
	nameLabel := widget.NewLabel(buttonName)
	nameLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create the main option container
	optionContent := container.NewBorder(
		nil, nil, nameLabel, toggle,
		container.NewPadded(widget.NewLabel("")),
	)

	// Create a card-like appearance with padding
	card := container.NewPadded(optionContent)

	return card, toggle
}
