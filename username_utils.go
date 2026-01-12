package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// GenerateDefaultUsername creates a "Guest XY" username where:
// X = random number between 1-999
// Y = first letter of email (uppercase)
// Example: "Guest42J" for john@example.com
func GenerateDefaultUsername(email string) string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(999) + 1

	firstLetter := "A"
	if len(email) > 0 {
		firstLetter = strings.ToUpper(string(email[0]))
	}

	return fmt.Sprintf("Guest%d%s", randomNum, firstLetter)
}

// ProfanityWords contains a basic list of inappropriate words
// In production, consider using a comprehensive profanity filter library
var ProfanityWords = []string{
	"shit", "fuck", "damn", "ass", "bitch", "bastard",
	"cock", "dick", "pussy", "nigger", "nigga", "fag",
	"cunt", "whore", "slut", "nazi", "hitler",
	// Add more words or use a library like github.com/TwiN/go-away
}

// IsInappropriateName checks if username contains profanity or inappropriate content
func IsInappropriateName(username string) bool {
	lowerName := strings.ToLower(username)

	// Check for exact matches and substrings
	for _, word := range ProfanityWords {
		if strings.Contains(lowerName, strings.ToLower(word)) {
			return true
		}
	}

	return false
}

// ValidateUsername checks if username meets all requirements:
// - Length between 3 and 20 characters
// - Only alphanumeric characters and spaces
// - No inappropriate content
func ValidateUsername(username string) error {
	// Trim whitespace
	username = strings.TrimSpace(username)

	// Length check
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(username) > 20 {
		return fmt.Errorf("username must be 20 characters or less")
	}

	// Alphanumeric + spaces check
	matched, err := regexp.MatchString("^[a-zA-Z0-9 ]+$", username)
	if err != nil {
		return fmt.Errorf("failed to validate username format: %w", err)
	}
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, and spaces")
	}

	// Check for consecutive spaces
	if strings.Contains(username, "  ") {
		return fmt.Errorf("username cannot contain consecutive spaces")
	}

	// Profanity check
	if IsInappropriateName(username) {
		return fmt.Errorf("username contains inappropriate content")
	}

	return nil
}
