package password

import (
	"testing"
)

func TestPasswordStrength_String(t *testing.T) {
	tests := []struct {
		strength PasswordStrength
		expected string
	}{
		{StrengthVeryWeak, "Very Weak"},
		{StrengthWeak, "Weak"},
		{StrengthFair, "Fair"},
		{StrengthGood, "Good"},
		{StrengthStrong, "Strong"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.strength.String(); got != tt.expected {
				t.Errorf("PasswordStrength.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCalculateStrength_EmptyPassword(t *testing.T) {
	calc := NewPasswordStrengthCalculator()
	strength, score := calc.CalculateStrength("")

	if strength != StrengthVeryWeak {
		t.Errorf("Empty password strength = %v, want %v", strength, StrengthVeryWeak)
	}
	if score != 0 {
		t.Errorf("Empty password score = %v, want 0", score)
	}
}

func TestCalculateStrength_VeryWeak(t *testing.T) {
	calc := NewPasswordStrengthCalculator()
	// Very weak passwords are those with extremely low scores (< 15)
	// Empty password is the only guaranteed very weak case
	veryWeakPasswords := []string{
		"", // Empty
	}

	for _, password := range veryWeakPasswords {
		t.Run(password, func(t *testing.T) {
			strength, _ := calc.CalculateStrength(password)
			if strength != StrengthVeryWeak {
				t.Errorf("Password %q strength = %v, want %v", password, strength, StrengthVeryWeak)
			}
		})
	}

	// Very short passwords should be at least Weak or VeryWeak
	shortPasswords := []string{
		"1111", // Repeated chars
		"a",    // Single char
		"12",   // Very short
	}

	for _, password := range shortPasswords {
		t.Run(password, func(t *testing.T) {
			strength, _ := calc.CalculateStrength(password)
			if strength > StrengthWeak {
				t.Errorf("Password %q strength = %v, want <= Weak", password, strength)
			}
		})
	}
}

func TestCalculateStrength_Weak(t *testing.T) {
	calc := NewPasswordStrengthCalculator()
	weakPasswords := []string{
		"password",
		"12345678",
		"abcdefgh",
	}

	for _, password := range weakPasswords {
		t.Run(password, func(t *testing.T) {
			strength, _ := calc.CalculateStrength(password)
			if strength > StrengthWeak {
				t.Errorf("Password %q strength = %v, want <= %v", password, strength, StrengthWeak)
			}
		})
	}
}

func TestCalculateStrength_Fair(t *testing.T) {
	calc := NewPasswordStrengthCalculator()
	fairPasswords := []string{
		"Password1",
		"MyPass123",
		"Test1234",
	}

	for _, password := range fairPasswords {
		t.Run(password, func(t *testing.T) {
			strength, _ := calc.CalculateStrength(password)
			if strength < StrengthFair || strength > StrengthGood {
				t.Errorf("Password %q strength = %v, want Fair or Good", password, strength)
			}
		})
	}
}

func TestCalculateStrength_Good(t *testing.T) {
	calc := NewPasswordStrengthCalculator()
	goodPasswords := []string{
		"MyP@ssw0rd123",
		"SecureP@ss1",
		"G00dP@ssword",
	}

	for _, password := range goodPasswords {
		t.Run(password, func(t *testing.T) {
			strength, _ := calc.CalculateStrength(password)
			if strength < StrengthGood {
				t.Errorf("Password %q strength = %v, want >= Good", password, strength)
			}
		})
	}
}

func TestCalculateStrength_Strong(t *testing.T) {
	calc := NewPasswordStrengthCalculator()
	strongPasswords := []string{
		"V3ry$tr0ng!P@ssw0rd",
		"C0mpl3x!P@ss#2024",
		"Sup3r$ecur3!Pass#123",
		"ExTr3m3ly!$tr0ng#Pwd",
	}

	for _, password := range strongPasswords {
		t.Run(password, func(t *testing.T) {
			strength, score := calc.CalculateStrength(password)
			if strength != StrengthStrong {
				t.Errorf("Password %q strength = %v (score: %d), want Strong", password, strength, score)
			}
		})
	}
}

func TestCalculateStrength_LengthBonus(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	shortPass := "A1b@"
	longPass := "A1b@c2D$e3F%g4H^"

	_, shortScore := calc.CalculateStrength(shortPass)
	_, longScore := calc.CalculateStrength(longPass)

	if longScore <= shortScore {
		t.Errorf("Longer password should have higher score: short=%d, long=%d", shortScore, longScore)
	}
}

func TestCalculateStrength_CharacterVariety(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	tests := []struct {
		name     string
		password string
		minScore int
	}{
		{"lowercase only", "abcdefgh", 0},
		{"upper+lower", "AbCdEfGh", 20},
		{"upper+lower+number", "AbCd1234", 30},
		{"all types", "AbCd12@#", 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, score := calc.CalculateStrength(tt.password)
			if score < tt.minScore {
				t.Errorf("Password %q score = %d, want >= %d", tt.password, score, tt.minScore)
			}
		})
	}
}

func TestCalculateStrength_CommonPatterns(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	withoutPattern := "Secur3!Pass"
	withPattern := "Password123"

	_, scoreWithout := calc.CalculateStrength(withoutPattern)
	_, scoreWith := calc.CalculateStrength(withPattern)

	if scoreWith >= scoreWithout {
		t.Errorf("Password with common pattern should have lower score: without=%d, with=%d", scoreWithout, scoreWith)
	}
}

func TestCalculateStrength_SequentialChars(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	withoutSequential := "P@ssw0rd!2024"
	withSequential := "P@ss123word"

	_, scoreWithout := calc.CalculateStrength(withoutSequential)
	_, scoreWith := calc.CalculateStrength(withSequential)

	if scoreWith >= scoreWithout {
		t.Errorf("Password with sequential chars should have lower score: without=%d, with=%d", scoreWithout, scoreWith)
	}
}

func TestCalculateStrength_RepeatedChars(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	withoutRepeated := "P@ssw0rd!2024"
	withRepeated := "P@ssw000rd"

	_, scoreWithout := calc.CalculateStrength(withoutRepeated)
	_, scoreWith := calc.CalculateStrength(withRepeated)

	if scoreWith >= scoreWithout {
		t.Errorf("Password with repeated chars should have lower score: without=%d, with=%d", scoreWithout, scoreWith)
	}
}

func TestCalculateStrength_UniqueCharacters(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	allUnique := "AbCd!@#$1234"
	someRepeated := "AAABBBcccddd"

	_, scoreUnique := calc.CalculateStrength(allUnique)
	_, scoreRepeated := calc.CalculateStrength(someRepeated)

	if scoreRepeated >= scoreUnique {
		t.Errorf("Password with more unique chars should have higher score: unique=%d, repeated=%d", scoreUnique, scoreRepeated)
	}
}

func TestCalculateStrength_ScoreBounds(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	// Test various passwords to ensure score stays in bounds
	passwords := []string{
		"",
		"a",
		"password",
		"P@ssw0rd!123",
		"V3ry!$tr0ng!P@ssw0rd#2024",
	}

	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			_, score := calc.CalculateStrength(password)
			if score < 0 || score > 100 {
				t.Errorf("Password %q score = %d, want between 0 and 100", password, score)
			}
		})
	}
}

func TestDetectPatterns(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	tests := []struct {
		name     string
		password string
		wantZero bool
	}{
		{"no pattern", "Secur3!Pass", true},
		{"has password", "Password123", false},
		{"has 12345", "Pass12345word", false},
		{"has qwerty", "Qwerty789", false},
		{"has admin", "Admin123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			penalty := calc.detectPatterns(tt.password)
			if tt.wantZero && penalty != 0 {
				t.Errorf("detectPatterns(%q) = %d, want 0", tt.password, penalty)
			}
			if !tt.wantZero && penalty == 0 {
				t.Errorf("detectPatterns(%q) = 0, want > 0", tt.password)
			}
		})
	}
}

func TestHasSequentialChars(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"no sequential", "P@ssw0rd", false},
		{"has abc", "Pabc123", true},
		{"has 123", "Pass123", true},
		{"has 456", "My456Pass", true},
		{"too short", "ab", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.hasSequentialChars(tt.password, 3)
			if got != tt.want {
				t.Errorf("hasSequentialChars(%q, 3) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestHasRepeatedChars(t *testing.T) {
	calc := NewPasswordStrengthCalculator()

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"no repeated", "P@ssw0rd", false},
		{"has aaa", "Paaa123", true},
		{"has 111", "Pass111", true},
		{"has 000", "My000Pass", true},
		{"too short", "aa", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.hasRepeatedChars(tt.password, 3)
			if got != tt.want {
				t.Errorf("hasRepeatedChars(%q, 3) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}
