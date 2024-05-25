package util

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,8}$`)

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if the password is strong and meets the criteria:
// - At least 8 characters long
// - Contains at least one digit
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one specialL character
func IsValidPassword(password string) (string, bool) {
	var (
		hasLower   = false
		hasUpper   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) < 8 {
		return "Password must contain at least 8 characters", false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain at least an uppercase character", false
	}
	if !hasLower {
		return "Password must contain at least a lowercase character", false
	}
	if !hasNumber {
		return "Password must contain at least a numeric character", false
	}
	if !hasSpecial {
		return "Password must contain at least a special character", false
	}

	return "", true
}
