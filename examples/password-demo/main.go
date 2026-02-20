package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/digital-michael/fyne-components/pkg/password"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Password Strength Demo")
	myWindow.Resize(fyne.NewSize(400, 200))

	// Create password entry
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter your password...")

	// Create strength meter
	strengthMeter := password.NewPasswordStrengthMeter()

	// Update strength meter when password changes
	passwordEntry.OnChanged = func(text string) {
		strengthMeter.UpdatePassword(text)
	}

	// Create instructions label
	instructions := widget.NewLabel("Enter a password to see its strength:")

	// Create layout
	content := container.NewVBox(
		instructions,
		passwordEntry,
		strengthMeter,
		widget.NewLabel(""),
		widget.NewLabel("Strength criteria:"),
		widget.NewLabel("• Very Weak: < 8 characters or very simple"),
		widget.NewLabel("• Weak: Basic password (8+ chars, single type)"),
		widget.NewLabel("• Fair: 8+ chars, mixed case or numbers"),
		widget.NewLabel("• Good: 8+ chars, mixed case, numbers, and symbols"),
		widget.NewLabel("• Strong: 12+ chars, diverse, no patterns"),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
