# Password Component

Password strength validation and visual feedback components for Fyne applications.

## Features

- ✅ **Strength Calculation**: Comprehensive password strength analysis
- ✅ **Visual Meter Widget**: Real-time strength indicator with color-coded feedback
- ✅ **Multiple Criteria**: Length, character variety, patterns, common passwords
- ✅ **Entropy-Based**: Uses actual entropy calculation for accurate strength
- ✅ **Customizable**: Configurable minimum requirements
- ✅ **Lightweight**: Zero external dependencies beyond Fyne

## Installation

```bash
go get github.com/digital-michael/fyne-components
```

## Quick Start

### Basic Strength Checking

```go
package main

import (
    "fmt"
    "github.com/digital-michael/fyne-components/pkg/password"
)

func main() {
    // Check password strength
    strength := password.CalculateStrength("MyP@ssw0rd123")
    
    fmt.Printf("Score: %d/100\n", strength.Score)
    fmt.Printf("Level: %s\n", strength.Level)
    fmt.Printf("Entropy: %.2f bits\n", strength.Entropy)
    
    // Check if meets requirements
    if strength.Level == password.LevelStrong || 
       strength.Level == password.LevelVeryStrong {
        fmt.Println("✓ Password is acceptable")
    }
}
```

### Visual Strength Meter

```go
package main

import (
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "github.com/digital-michael/fyne-components/pkg/password"
)

func main() {
    myApp := app.New()
    window := myApp.NewWindow("Password Demo")

    // Create password entry
    passwordEntry := widget.NewPasswordEntry()
    passwordEntry.SetPlaceHolder("Enter password...")

    // Create strength meter
    meter := password.NewStrengthMeter()

    // Update meter as user types
    passwordEntry.OnChanged = func(text string) {
        strength := password.CalculateStrength(text)
        meter.SetStrength(strength)
    }

    // Layout
    content := container.NewVBox(
        widget.NewLabel("Password:"),
        passwordEntry,
        meter,
    )

    window.SetContent(content)
    window.ShowAndRun()
}
```

## Strength Levels

The strength calculator returns one of five levels:

| Level | Score Range | Description |
|-------|-------------|-------------|
| `LevelVeryWeak` | 0-20 | Easily cracked, unacceptable |
| `LevelWeak` | 21-40 | Poor password, not recommended |
| `LevelFair` | 41-60 | Acceptable but could be better |
| `LevelStrong` | 61-80 | Good password |
| `LevelVeryStrong` | 81-100 | Excellent password |

## Strength Calculation

### Scoring Criteria

The strength calculation considers multiple factors:

1. **Length** (25 points)
   - Minimum 8 characters required
   - Bonus for 12+ characters
   - Penalty for fewer than 8

2. **Character Variety** (30 points)
   - Lowercase letters (a-z)
   - Uppercase letters (A-Z)
   - Numbers (0-9)
   - Special characters (!@#$%^&*...)

3. **Entropy** (35 points)
   - Measures actual randomness
   - Based on character set size and length
   - Higher entropy = harder to crack

4. **Pattern Detection** (10 points)
   - Deducts for common patterns
   - Deducts for sequential characters (abc, 123)
   - Deducts for repeated characters (aaa, 111)
   - Checks against common password list

### Strength Struct

```go
type Strength struct {
    Score      int         // 0-100
    Level      Level       // VeryWeak to VeryStrong
    Entropy    float64     // Bits of entropy
    Length     int         // Password length
    HasLower   bool        // Contains lowercase
    HasUpper   bool        // Contains uppercase
    HasNumber  bool        // Contains digits
    HasSpecial bool        // Contains special chars
    Feedback   []string    // Suggestions for improvement
}
```

## Examples

### Example 1: Registration Form Validation

```go
func validatePassword(password string) error {
    strength := password.CalculateStrength(password)
    
    if strength.Level == password.LevelVeryWeak {
        return fmt.Errorf("password too weak: %s", 
            strings.Join(strength.Feedback, ", "))
    }
    
    if strength.Score < 50 {
        return fmt.Errorf("password must score at least 50/100 (current: %d)", 
            strength.Score)
    }
    
    return nil
}
```

### Example 2: Real-time Feedback

```go
passwordEntry.OnChanged = func(text string) {
    strength := password.CalculateStrength(text)
    meter.SetStrength(strength)
    
    // Show feedback
    if len(strength.Feedback) > 0 {
        feedbackLabel.SetText(strings.Join(strength.Feedback, "\n"))
    } else {
        feedbackLabel.SetText("✓ Strong password")
    }
}
```

### Example 3: Custom Minimum Requirements

```go
func meetsCustomRequirements(pwd string) bool {
    strength := password.CalculateStrength(pwd)
    
    return strength.Length >= 10 &&
           strength.HasUpper &&
           strength.HasLower &&
           strength.HasNumber &&
           strength.HasSpecial &&
           strength.Score >= 60
}
```

### Example 4: Password Strength Report

```go
func printPasswordReport(pwd string) {
    strength := password.CalculateStrength(pwd)
    
    fmt.Printf("Password Analysis\n")
    fmt.Printf("================\n")
    fmt.Printf("Length:    %d characters\n", strength.Length)
    fmt.Printf("Score:     %d/100\n", strength.Score)
    fmt.Printf("Level:     %s\n", strength.Level)
    fmt.Printf("Entropy:   %.2f bits\n", strength.Entropy)
    fmt.Printf("\nCharacter Types:\n")
    fmt.Printf("  Lowercase: %v\n", strength.HasLower)
    fmt.Printf("  Uppercase: %v\n", strength.HasUpper)
    fmt.Printf("  Numbers:   %v\n", strength.HasNumber)
    fmt.Printf("  Special:   %v\n", strength.HasSpecial)
    
    if len(strength.Feedback) > 0 {
        fmt.Printf("\nSuggestions:\n")
        for _, fb := range strength.Feedback {
            fmt.Printf("  • %s\n", fb)
        }
    }
}
```

## Strength Meter Widget

The visual strength meter provides color-coded feedback:

```go
meter := password.NewStrengthMeter()
meter.SetStrength(strength)
```

### Visual Indicators

- **Very Weak** (0-20): Red bar, minimal fill
- **Weak** (21-40): Orange bar, low fill
- **Fair** (41-60): Yellow bar, medium fill
- **Strong** (61-80): Light green bar, high fill
- **Very Strong** (81-100): Dark green bar, full fill

### Customization

The meter is a standard Fyne widget and can be used anywhere:

```go
// In a border layout
content := container.NewBorder(
    header,
    meter,  // Bottom
    nil,
    nil,
    mainContent,
)

// In a form
form := &widget.Form{
    Items: []*widget.FormItem{
        {Text: "Password", Widget: passwordEntry},
        {Text: "Strength", Widget: meter},
    },
}
```

## API Reference

### Strength Calculation

```go
// Calculate password strength
func CalculateStrength(password string) Strength

// Strength levels
const (
    LevelVeryWeak   Level = "Very Weak"
    LevelWeak       Level = "Weak"
    LevelFair       Level = "Fair"
    LevelStrong     Level = "Strong"
    LevelVeryStrong Level = "Very Strong"
)
```

### Strength Meter Widget

```go
// Create new meter
func NewStrengthMeter() *StrengthMeter

// Update meter display
func (m *StrengthMeter) SetStrength(strength Strength)

// Get current strength
func (m *StrengthMeter) GetStrength() Strength

// Standard widget methods
func (m *StrengthMeter) CreateRenderer() fyne.WidgetRenderer
func (m *StrengthMeter) Refresh()
```

## Testing

The password component includes comprehensive tests:

```bash
cd pkg/password
go test -v
```

Test coverage includes:
- All strength levels (63 test cases)
- Character type detection
- Pattern detection
- Entropy calculation
- Common password checking
- Edge cases (empty, very long passwords)

Example test output:
```
=== RUN   TestCalculateStrength/Very_Weak_-_Empty
=== RUN   TestCalculateStrength/Very_Weak_-_Short
=== RUN   TestCalculateStrength/Weak_-_Only_lowercase
=== RUN   TestCalculateStrength/Fair_-_Mixed_case
=== RUN   TestCalculateStrength/Strong_-_Good_variety
=== RUN   TestCalculateStrength/Very_Strong_-_Excellent
--- PASS: TestCalculateStrength (0.00s)
PASS
ok      github.com/digital-michael/fyne-components/pkg/password
```

## Common Patterns

### Login Form with Validation

```go
type LoginForm struct {
    passwordEntry *widget.Entry
    meter         *password.StrengthMeter
    submitBtn     *widget.Button
}

func (f *LoginForm) Create() fyne.CanvasObject {
    f.passwordEntry = widget.NewPasswordEntry()
    f.meter = password.NewStrengthMeter()
    
    f.passwordEntry.OnChanged = func(text string) {
        strength := password.CalculateStrength(text)
        f.meter.SetStrength(strength)
        
        // Enable submit only if strong enough
        f.submitBtn.Enable()
        if strength.Score < 50 {
            f.submitBtn.Disable()
        }
    }
    
    f.submitBtn = widget.NewButton("Create Account", f.onSubmit)
    
    return container.NewVBox(
        widget.NewLabel("Choose password:"),
        f.passwordEntry,
        f.meter,
        f.submitBtn,
    )
}
```

### Password Change Dialog

```go
func ShowPasswordChangeDialog(parent fyne.Window) {
    currentEntry := widget.NewPasswordEntry()
    newEntry := widget.NewPasswordEntry()
    confirmEntry := widget.NewPasswordEntry()
    meter := password.NewStrengthMeter()
    
    newEntry.OnChanged = func(text string) {
        strength := password.CalculateStrength(text)
        meter.SetStrength(strength)
    }
    
    items := []*widget.FormItem{
        {Text: "Current:", Widget: currentEntry},
        {Text: "New:", Widget: newEntry},
        {Text: "Strength:", Widget: meter},
        {Text: "Confirm:", Widget: confirmEntry},
    }
    
    dialog.ShowForm("Change Password", "Save", "Cancel", items, 
        func(accepted bool) {
            if !accepted {
                return
            }
            // Validate and save...
        }, parent)
}
```

## Security Considerations

### What the Calculator Checks

✅ **Does Check:**
- Password length
- Character variety (lower, upper, digits, special)
- Entropy (randomness)
- Common password patterns
- Sequential characters
- Repeated characters

❌ **Does Not Check:**
- Password reuse across sites
- Personal information (name, birthday)
- Dictionary words in other languages
- Timing attacks
- Storage security

### Best Practices

1. **Require Strong Passwords**: Set minimum score of 50-60
2. **Provide Feedback**: Show the strength meter and suggestions
3. **Educate Users**: Explain what makes a strong password
4. **Never Store Plaintext**: Use proper password hashing (bcrypt, argon2)
5. **Consider Length Over Complexity**: A long passphrase is often stronger than short complex password

### Recommended Minimums

For different security requirements:

**Low Security** (basic accounts):
```go
score >= 40  // Fair level
length >= 8
```

**Medium Security** (standard applications):
```go
score >= 60  // Strong level
length >= 10
hasUpper && hasLower && hasNumber
```

**High Security** (sensitive data):
```go
score >= 75  // Very Strong level
length >= 12
hasUpper && hasLower && hasNumber && hasSpecial
entropy >= 40.0
```

## Demo Application

See the [password-demo](../../examples/password-demo/) for a complete working example with:
- Interactive password entry
- Real-time strength meter
- Detailed strength breakdown
- Feedback suggestions
- All strength levels demonstrated

Run the demo:
```bash
cd examples/password-demo
go run main.go
```

## License

See the [LICENSE](../../LICENSE) file for details.

## Related

- [Password Demo](../../examples/password-demo/) - Interactive demonstration
- [Table Widget](../table/) - Sortable, filterable table component
- [Design Principles](../../DESIGN-PRINCIPLES.md) - Architecture guidelines
