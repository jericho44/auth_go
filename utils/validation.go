package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail performs email validation
func IsValidEmail(email string) bool {
	// Basic validation - contains @ and .
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	// More comprehensive email regex validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidUsername validates username format
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}

	// Username can contain letters, numbers, and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// IsValidPassword validates password strength
func IsValidPassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	// Add more password strength requirements if needed
	// For now, just check minimum length
	return true
}
