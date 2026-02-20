// Package password provides password-related widgets for Fyne applications.
//
// This package includes:
//   - Password strength calculator with detailed analysis
//   - Password strength meter widget with visual feedback
//
// # Password Strength Calculation
//
// The strength calculator analyzes passwords based on multiple criteria:
//   - Length (minimum 8 characters recommended)
//   - Character diversity (uppercase, lowercase, numbers, symbols)
//   - Common pattern detection (repeated characters, sequences)
//   - Dictionary word detection
//
// Strength levels:
//   - Weak: Basic passwords that are easily guessable
//   - Medium: Moderate passwords with some diversity
//   - Strong: Well-formed passwords with good diversity
//   - VeryStrong: Excellent passwords meeting all security criteria
//
// # Basic Usage
//
//	// Calculate strength
//	strength := password.CalculateStrength("MyP@ssw0rd123")
//	fmt.Printf("Strength: %s\n", strength) // Output: Strong
//
//	// Use in a widget
//	passwordInput := widget.NewPasswordEntry()
//	strengthMeter := password.NewStrengthMeter()
//
//	passwordInput.OnChanged = func(text string) {
//		strengthMeter.SetPassword(text)
//	}
//
// # Thread Safety
//
// The strength calculator is stateless and thread-safe.
// The strength meter widget must be accessed from the UI goroutine only.
package password
