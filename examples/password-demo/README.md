# Password Strength Meter Demo

A simple demonstration of the password strength meter widget from fyne-components.

## Features

- Real-time password strength calculation
- Visual strength indicator with color coding
- Text label showing strength level
- Criteria display for each strength level

## Running the Demo

```bash
go run main.go
```

Or build and run:

```bash
go build
./password-demo
```

## Strength Levels

The password strength meter evaluates passwords based on multiple criteria:

- **Very Weak**: Less than 8 characters or extremely simple patterns
- **Weak**: Basic passwords (8+ characters, single character type)
- **Fair**: 8+ characters with mixed case or numbers
- **Good**: 8+ characters with mixed case, numbers, and symbols
- **Strong**: 12+ characters with diverse character types and no common patterns

## Implementation

```go
// Create password entry
passwordEntry := widget.NewPasswordEntry()

// Create strength meter
strengthMeter := password.NewPasswordStrengthMeter()

// Update strength on password change
passwordEntry.OnChanged = func(text string) {
    strengthMeter.UpdatePassword(text)
}
```

## Password Strength Calculator

The underlying password strength calculator is also available for direct use:

```go
calc := password.NewPasswordStrengthCalculator()
strength, score := calc.CalculateStrength("MyP@ssw0rd123")
fmt.Printf("Strength: %s, Score: %d\n", strength, score)
```
