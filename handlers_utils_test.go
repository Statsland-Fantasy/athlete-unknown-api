package main

import (
	"testing"
	"time"
)

// TestUpdateDailyStreak tests the updateDailyStreak helper function
// This function uses engagement-based tracking: it tracks consecutive real-life calendar days
// the user plays ANY round, regardless of which round's playDate they choose to play
func TestUpdateDailyStreak(t *testing.T) {
	tests := []struct {
		name                  string
		existingUser          *User
		currentDate           string // Real-life date (e.g., today's date), not the round's playDate
		expectedStreak        int
		expectedLastDayPlayed string
	}{
		{
			name:                  "nil user stats - should do nothing",
			existingUser:          nil,
			currentDate:           "2024-01-15",
			expectedStreak:        0,
			expectedLastDayPlayed: "",
		},
		{
			name: "consecutive day - increment streak",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 5,
				LastDayPlayed:      "2024-01-14",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-01-15",
			expectedStreak:        6,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "same day - keep streak unchanged",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 3,
				LastDayPlayed:      "2024-01-15",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-01-15",
			expectedStreak:        3,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "missed one day - reset streak to 1",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 10,
				LastDayPlayed:      "2024-01-13",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-01-15",
			expectedStreak:        1,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "missed multiple days - reset streak to 1",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 7,
				LastDayPlayed:      "2024-01-10",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-01-20",
			expectedStreak:        1,
			expectedLastDayPlayed: "2024-01-20",
		},
		{
			name: "streak continues after weekend - friday to saturday",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 4,
				LastDayPlayed:      "2024-01-12",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-01-13",
			expectedStreak:        5,
			expectedLastDayPlayed: "2024-01-13",
		},
		{
			name: "existing user with empty lastDayPlayed",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 2,
				LastDayPlayed:      "",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-01-15",
			expectedStreak:        2,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "month transition - consecutive day",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 15,
				LastDayPlayed:      "2024-01-31",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-02-01",
			expectedStreak:        16,
			expectedLastDayPlayed: "2024-02-01",
		},
		{
			name: "year transition - consecutive day",
			existingUser: &User{
				UserId:             "user123",
				CurrentDailyStreak: 20,
				LastDayPlayed:      "2023-12-31",
				UserCreated:        time.Now(),
			},
			currentDate:           "2024-01-01",
			expectedStreak:        21,
			expectedLastDayPlayed: "2024-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the user stats if not nil
			var user *User
			if tt.existingUser != nil {
				user = &User{
					UserId:             tt.existingUser.UserId,
					CurrentDailyStreak: tt.existingUser.CurrentDailyStreak,
					LastDayPlayed:      tt.existingUser.LastDayPlayed,
					UserCreated:        tt.existingUser.UserCreated,
					Sports:             tt.existingUser.Sports,
				}
			}

			// Call the function
			updateDailyStreak(user, tt.currentDate)

			// Verify the results
			if user == nil {
				if tt.expectedStreak != 0 || tt.expectedLastDayPlayed != "" {
					t.Errorf("Expected nil user to remain nil")
				}
			} else {
				if user.CurrentDailyStreak != tt.expectedStreak {
					t.Errorf("Expected streak %d, got %d", tt.expectedStreak, user.CurrentDailyStreak)
				}

				if user.LastDayPlayed != tt.expectedLastDayPlayed {
					t.Errorf("Expected lastDayPlayed %s, got %s", tt.expectedLastDayPlayed, user.LastDayPlayed)
				}
			}
		})
	}
}

func TestValidateSportsReferenceURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		shouldErr bool
	}{
		{
			name:      "valid baseball reference URL",
			url:       "https://www.baseball-reference.com/players/t/troutmi01.shtml",
			shouldErr: false,
		},
		{
			name:      "valid basketball reference URL",
			url:       "https://www.basketball-reference.com/players/j/jamesle01.html",
			shouldErr: false,
		},
		{
			name:      "valid pro-football reference URL",
			url:       "https://www.pro-football-reference.com/players/M/MahoPa00.htm",
			shouldErr: false,
		},
		{
			name:      "valid URL without www",
			url:       "https://baseball-reference.com/players/t/troutmi01.shtml",
			shouldErr: false,
		},
		{
			name:      "invalid domain",
			url:       "https://malicious-site.com/players/test",
			shouldErr: true,
		},
		{
			name:      "missing scheme",
			url:       "baseball-reference.com/players/test",
			shouldErr: true,
		},
		{
			name:      "invalid scheme",
			url:       "ftp://baseball-reference.com/players/test",
			shouldErr: true,
		},
		{
			name:      "empty URL",
			url:       "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSportsReferenceURL(tt.url)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for URL %s, but got none", tt.url)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for URL %s, but got: %v", tt.url, err)
			}
		})
	}
}
