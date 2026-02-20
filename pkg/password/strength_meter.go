package password

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// PasswordStrengthMeter is a widget that displays password strength visually
type PasswordStrengthMeter struct {
	widget.BaseWidget

	calculator  *PasswordStrengthCalculator
	password    string
	strength    PasswordStrength
	score       int
	strengthBar *canvas.Rectangle
	labelWidget *widget.Label
	container   *fyne.Container
}

// NewPasswordStrengthMeter creates a new password strength meter widget
func NewPasswordStrengthMeter() *PasswordStrengthMeter {
	meter := &PasswordStrengthMeter{
		calculator: NewPasswordStrengthCalculator(),
	}

	// Create visual elements
	meter.strengthBar = canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255})
	meter.strengthBar.SetMinSize(fyne.NewSize(200, 8))

	meter.labelWidget = widget.NewLabel("")
	meter.labelWidget.TextStyle = fyne.TextStyle{Bold: true}

	// Create container with bar and label
	meter.container = container.NewVBox(
		meter.strengthBar,
		meter.labelWidget,
	)

	meter.ExtendBaseWidget(meter)
	meter.UpdatePassword("") // Initialize with empty state

	return meter
}

// UpdatePassword updates the meter based on the new password
func (m *PasswordStrengthMeter) UpdatePassword(password string) {
	m.password = password
	m.strength, m.score = m.calculator.CalculateStrength(password)

	// Update visual appearance
	m.updateAppearance()
}

// updateAppearance updates the color and label based on current strength
func (m *PasswordStrengthMeter) updateAppearance() {
	var barColor color.Color
	var labelText string

	switch m.strength {
	case StrengthVeryWeak:
		barColor = color.RGBA{R: 220, G: 53, B: 69, A: 255} // Red
		labelText = "Password Strength: Very Weak"
	case StrengthWeak:
		barColor = color.RGBA{R: 255, G: 193, B: 7, A: 255} // Yellow/Orange
		labelText = "Password Strength: Weak"
	case StrengthFair:
		barColor = color.RGBA{R: 255, G: 193, B: 7, A: 255} // Yellow
		labelText = "Password Strength: Fair"
	case StrengthGood:
		barColor = color.RGBA{R: 40, G: 167, B: 69, A: 255} // Light Green
		labelText = "Password Strength: Good"
	case StrengthStrong:
		barColor = color.RGBA{R: 25, G: 135, B: 84, A: 255} // Dark Green
		labelText = "Password Strength: Strong"
	default:
		barColor = color.RGBA{R: 200, G: 200, B: 200, A: 255} // Gray
		labelText = "Password Strength: Unknown"
	}

	// Update bar color and size based on score
	m.strengthBar.FillColor = barColor

	// Calculate bar width based on score (0-100 maps to 0-200 pixels)
	barWidth := float32(m.score) * 2.0
	if barWidth < 20 {
		barWidth = 20 // Minimum visible width
	}
	m.strengthBar.SetMinSize(fyne.NewSize(barWidth, 8))

	// Update label
	m.labelWidget.SetText(labelText)

	m.Refresh()
}

// CreateRenderer creates the renderer for this widget
func (m *PasswordStrengthMeter) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.container)
}

// GetStrength returns the current strength level
func (m *PasswordStrengthMeter) GetStrength() PasswordStrength {
	return m.strength
}

// GetScore returns the current score (0-100)
func (m *PasswordStrengthMeter) GetScore() int {
	return m.score
}
