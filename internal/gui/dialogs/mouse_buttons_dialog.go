package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"click-guardian/internal/config"
	"click-guardian/internal/hooks"
)

// ShowMouseButtonsDialog shows a modern dialog for selecting which mouse buttons to protect
func ShowMouseButtonsDialog(window fyne.Window, currentConfig *config.Config, onSave func([]string)) {
	// Get available mouse buttons
	availableButtons := hooks.GetMouseButtons()
	
	// Create instruction text
	instruction := widget.NewLabel("Select which mouse buttons should be protected from double-clicks:")
	instruction.Wrapping = fyne.TextWrapWord
	instruction.Alignment = fyne.TextAlignCenter
	
	// Create checkboxes for each mouse button
	checkBoxes := make([]*widget.Check, len(availableButtons))
	buttonOptions := make([]fyne.CanvasObject, len(availableButtons))
	
	// Create a function to handle saving
	saveSelection := func() {
		// Collect selected buttons
		var selected []string
		for j, cb := range checkBoxes {
			if cb.Checked {
				selected = append(selected, availableButtons[j].ID)
			}
		}
		
		// If no buttons selected, default to left button
		if len(selected) == 0 {
			selected = []string{"left"}
			// Ensure left button checkbox is checked
			for j, button := range availableButtons {
				if button.ID == "left" {
					checkBoxes[j].SetChecked(true)
					break
				}
			}
		}
		
		onSave(selected)
	}
	
	// Set initial states based on current config and create button options
	for i, button := range availableButtons {
		// Create a modern option for button selection
		option, check := createButtonOption(button.Name, button.ID, currentConfig)
		checkBoxes[i] = check
		buttonOptions[i] = option
		
		// Add immediate save functionality to each checkbox
		check.OnChanged = func(checked bool) {
			saveSelection()
		}
	}
	
	// Create options section
	optionsContainer := container.NewVBox(buttonOptions...)
	
	// Create a subtle footer with information
	footerText := widget.NewRichTextFromMarkdown(
		"**Tips:**\n" +
		"• Protected buttons will have double-click filtering applied\n" +
		"• Settings are saved automatically when changed\n" +
		"• At least one button must be protected")
	footerText.Wrapping = fyne.TextWrapWord
	
	// Assemble the complete dialog content
	content := container.NewVBox(
		instruction,
		layout.NewSpacer(),
		optionsContainer,
		layout.NewSpacer(),
		widget.NewSeparator(),
		footerText,
	)
	
	// Create and show the dialog
	d := dialog.NewCustom("Mouse Button Settings", "Close", content, window)
	d.Resize(fyne.NewSize(400, 400))
	d.Show()
}

// createButtonOption creates a modern option for button selection
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