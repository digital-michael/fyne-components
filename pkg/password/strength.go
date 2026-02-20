package password

import (
	"strings"
	"unicode"
)

// PasswordStrength represents the strength level of a password
type PasswordStrength int

const (
	// StrengthVeryWeak indicates a very weak password
	StrengthVeryWeak PasswordStrength = iota
	// StrengthWeak indicates a weak password
	StrengthWeak
	// StrengthFair indicates a fair password
	StrengthFair
	// StrengthGood indicates a good password
	StrengthGood
	// StrengthStrong indicates a strong password
	StrengthStrong
)

// String returns the string representation of the strength level
func (s PasswordStrength) String() string {
	switch s {
	case StrengthVeryWeak:
		return "Very Weak"
	case StrengthWeak:
		return "Weak"
	case StrengthFair:
		return "Fair"
	case StrengthGood:
		return "Good"
	case StrengthStrong:
		return "Strong"
	default:
		return "Unknown"
	}
}

// PasswordStrengthCalculator calculates password strength based on various criteria
type PasswordStrengthCalculator struct{}

// NewPasswordStrengthCalculator creates a new password strength calculator
func NewPasswordStrengthCalculator() *PasswordStrengthCalculator {
	return &PasswordStrengthCalculator{}
}

// CalculateStrength calculates the strength of a password
// Returns a strength level (0-4) and a score (0-100)
func (c *PasswordStrengthCalculator) CalculateStrength(password string) (PasswordStrength, int) {
	if password == "" {
		return StrengthVeryWeak, 0
	}

	score := 0

	// Length scoring (0-30 points)
	length := len(password)
	switch {
	case length >= 16:
		score += 30
	case length >= 12:
		score += 25
	case length >= 10:
		score += 20
	case length >= 8:
		score += 15
	case length >= 6:
		score += 10
	default:
		score += 5
	}

	// Character variety scoring (0-40 points)
	var hasLower, hasUpper, hasNumber, hasSpecial bool
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsDigit(char) {
			hasNumber = true
		} else if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	if hasLower {
		score += 10
	}
	if hasUpper {
		score += 10
	}
	if hasNumber {
		score += 10
	}
	if hasSpecial {
		score += 10
	}

	// Diversity bonus (0-15 points)
	charTypes := 0
	if hasLower {
		charTypes++
	}
	if hasUpper {
		charTypes++
	}
	if hasNumber {
		charTypes++
	}
	if hasSpecial {
		charTypes++
	}

	switch charTypes {
	case 4:
		score += 15
	case 3:
		score += 10
	case 2:
		score += 5
	}

	// Unique character bonus (0-10 points)
	uniqueChars := make(map[rune]bool)
	for _, char := range password {
		uniqueChars[char] = true
	}
	uniqueRatio := float64(len(uniqueChars)) / float64(length)
	if uniqueRatio >= 0.8 {
		score += 10
	} else if uniqueRatio >= 0.6 {
		score += 7
	} else if uniqueRatio >= 0.4 {
		score += 5
	}

	// Pattern penalties (0 to -15 points)
	score -= c.detectPatterns(password)

	// Ensure score stays within bounds
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}

	// Convert score to strength level
	strength := c.scoreToStrength(score)

	return strength, score
}

// detectPatterns looks for common patterns and returns a penalty score
func (c *PasswordStrengthCalculator) detectPatterns(password string) int {
	penalty := 0
	lower := strings.ToLower(password)

	// Common patterns - check for exact or prominent matches
	commonPatterns := []string{
		"password", "12345678", "qwerty", "111111", "000000",
		"admin", "user", "login", "123456", "654321",
	}

	for _, pattern := range commonPatterns {
		if strings.Contains(lower, pattern) {
			penalty += 10
			break // Only penalize once for patterns
		}
	}

	// Sequential characters (abc, 123, etc.)
	if c.hasSequentialChars(password, 3) {
		penalty += 5
	}

	// Repeated characters (aaa, 111, etc.)
	if c.hasRepeatedChars(password, 3) {
		penalty += 5
	}

	return penalty
}

// hasSequentialChars checks if the password has sequential characters
func (c *PasswordStrengthCalculator) hasSequentialChars(password string, minLength int) bool {
	if len(password) < minLength {
		return false
	}

	for i := 0; i <= len(password)-minLength; i++ {
		isSequential := true
		for j := 1; j < minLength; j++ {
			if password[i+j] != password[i+j-1]+1 {
				isSequential = false
				break
			}
		}
		if isSequential {
			return true
		}
	}
	return false
}

// hasRepeatedChars checks if the password has repeated characters
func (c *PasswordStrengthCalculator) hasRepeatedChars(password string, minLength int) bool {
	if len(password) < minLength {
		return false
	}

	for i := 0; i <= len(password)-minLength; i++ {
		allSame := true
		for j := 1; j < minLength; j++ {
			if password[i+j] != password[i] {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}
	return false
}

// scoreToStrength converts a numeric score to a strength level
func (c *PasswordStrengthCalculator) scoreToStrength(score int) PasswordStrength {
	switch {
	case score >= 75:
		return StrengthStrong
	case score >= 55:
		return StrengthGood
	case score >= 35:
		return StrengthFair
	case score >= 15:
		return StrengthWeak
	default:
		return StrengthVeryWeak
	}
}
