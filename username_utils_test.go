package main

import (
	"strings"
	"testing"
)

func TestGenerateDefaultUsername(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		wantPrefix    string
		wantLastChar  string
		wantMinLength int
		wantMaxLength int
	}{
		{
			name:          "Standard email",
			email:         "john@example.com",
			wantPrefix:    "Guest",
			wantLastChar:  "J",
			wantMinLength: 7, // Guest1J
			wantMaxLength: 9, // Guest999J
		},
		{
			name:          "Email starting with lowercase",
			email:         "alice@test.com",
			wantPrefix:    "Guest",
			wantLastChar:  "A",
			wantMinLength: 7,
			wantMaxLength: 9,
		},
		{
			name:          "Empty email",
			email:         "",
			wantPrefix:    "Guest",
			wantLastChar:  "A",
			wantMinLength: 7,
			wantMaxLength: 9,
		},
		{
			name:          "Email with special char",
			email:         "$pecial@test.com",
			wantPrefix:    "Guest",
			wantLastChar:  "$",
			wantMinLength: 7,
			wantMaxLength: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateDefaultUsername(tt.email)

			// Check prefix
			if !strings.HasPrefix(got, tt.wantPrefix) {
				t.Errorf("GenerateDefaultUsername() = %v, want prefix %v", got, tt.wantPrefix)
			}

			// Check last character
			if !strings.HasSuffix(got, tt.wantLastChar) {
				t.Errorf("GenerateDefaultUsername() = %v, want suffix %v", got, tt.wantLastChar)
			}

			// Check length range
			if len(got) < tt.wantMinLength || len(got) > tt.wantMaxLength {
				t.Errorf("GenerateDefaultUsername() length = %v, want between %v and %v", len(got), tt.wantMinLength, tt.wantMaxLength)
			}
		})
	}
}

func TestIsInappropriateName(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{
			name:     "Clean username",
			username: "JohnDoe123",
			want:     false,
		},
		{
			name:     "Clean username with spaces",
			username: "John Doe",
			want:     false,
		},
		{
			name:     "Contains profanity - exact match",
			username: "fuck",
			want:     true,
		},
		{
			name:     "Contains profanity - substring",
			username: "fuckthis",
			want:     true,
		},
		{
			name:     "Contains profanity - uppercase",
			username: "DAMN",
			want:     true,
		},
		{
			name:     "Contains profanity - mixed case",
			username: "BiTcH123",
			want:     true,
		},
		{
			name:     "Innocent word that contains bad substring",
			username: "Assassin",
			want:     true, // Will flag because contains "ass"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsInappropriateName(tt.username)
			if got != tt.want {
				t.Errorf("IsInappropriateName(%v) = %v, want %v", tt.username, got, tt.want)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid username - alphanumeric",
			username: "JohnDoe123",
			wantErr:  false,
		},
		{
			name:     "Valid username - with spaces",
			username: "John Doe",
			wantErr:  false,
		},
		{
			name:     "Valid username - minimum length",
			username: "Joe",
			wantErr:  false,
		},
		{
			name:     "Valid username - maximum length",
			username: "TwentyCharactersName",
			wantErr:  false,
		},
		{
			name:     "Too short",
			username: "AB",
			wantErr:  true,
			errMsg:   "at least 3 characters",
		},
		{
			name:     "Too long",
			username: "ThisUsernameIsWayTooLongToBeValid",
			wantErr:  true,
			errMsg:   "20 characters or less",
		},
		{
			name:     "Contains special characters",
			username: "John@Doe",
			wantErr:  true,
			errMsg:   "can only contain letters, numbers, and spaces",
		},
		{
			name:     "Contains underscore",
			username: "John_Doe",
			wantErr:  true,
			errMsg:   "can only contain letters, numbers, and spaces",
		},
		{
			name:     "Contains consecutive spaces",
			username: "John  Doe",
			wantErr:  true,
			errMsg:   "cannot contain consecutive spaces",
		},
		{
			name:     "Contains profanity",
			username: "shit123",
			wantErr:  true,
			errMsg:   "inappropriate content",
		},
		{
			name:     "Only spaces (trimmed to empty)",
			username: "   ",
			wantErr:  true,
			errMsg:   "at least 3 characters",
		},
		{
			name:     "Valid with leading/trailing spaces (should trim)",
			username: "  JohnDoe  ",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername(%v) error = %v, wantErr %v", tt.username, err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateUsername(%v) error = %v, want error containing %v", tt.username, err.Error(), tt.errMsg)
				}
			}
		})
	}
}
