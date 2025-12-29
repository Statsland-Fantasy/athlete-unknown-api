package main

import (
	"testing"
	"time"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "item exists in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "item does not exist in slice",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "orange",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			item:     "apple",
			expected: false,
		},
		{
			name:     "empty item",
			slice:    []string{"apple", "banana"},
			item:     "",
			expected: false,
		},
		{
			name:     "empty item exists in slice",
			slice:    []string{"apple", "", "banana"},
			item:     "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(tt.slice, tt.item)
			if got != tt.expected {
				t.Errorf("contains(%v, %q) = %v, want %v", tt.slice, tt.item, got, tt.expected)
			}
		})
	}
}

func TestIsValidYear(t *testing.T) {
	tests := []struct {
		name     string
		year     string
		expected bool
	}{
		{
			name:     "valid year 2020",
			year:     "2020",
			expected: true,
		},
		{
			name:     "valid year 1999",
			year:     "1999",
			expected: true,
		},
		{
			name:     "empty string",
			year:     "",
			expected: false,
		},
		{
			name:     "Season label",
			year:     "Season",
			expected: false,
		},
		{
			name:     "Career label",
			year:     "Career",
			expected: false,
		},
		{
			name:     "Yr label",
			year:     "Yr",
			expected: false,
		},
		{
			name:     "Avg label",
			year:     "Avg",
			expected: false,
		},
		{
			name:     "Average label",
			year:     "Average",
			expected: false,
		},
		{
			name:     "contains Season",
			year:     "2020 Season",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidYear(tt.year)
			if got != tt.expected {
				t.Errorf("isValidYear(%q) = %v, want %v", tt.year, got, tt.expected)
			}
		})
	}
}

func TestGetOrdinalSuffix(t *testing.T) {
	tests := []struct {
		name     string
		numStr   string
		expected string
	}{
		{
			name:     "1st",
			numStr:   "1",
			expected: "st",
		},
		{
			name:     "2nd",
			numStr:   "2",
			expected: "nd",
		},
		{
			name:     "3rd",
			numStr:   "3",
			expected: "rd",
		},
		{
			name:     "4th",
			numStr:   "4",
			expected: "th",
		},
		{
			name:     "11th",
			numStr:   "11",
			expected: "th",
		},
		{
			name:     "12th",
			numStr:   "12",
			expected: "th",
		},
		{
			name:     "13th",
			numStr:   "13",
			expected: "th",
		},
		{
			name:     "21st",
			numStr:   "21",
			expected: "st",
		},
		{
			name:     "22nd",
			numStr:   "22",
			expected: "nd",
		},
		{
			name:     "23rd",
			numStr:   "23",
			expected: "rd",
		},
		{
			name:     "100th",
			numStr:   "100",
			expected: "th",
		},
		{
			name:     "101st",
			numStr:   "101",
			expected: "st",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOrdinalSuffix(tt.numStr)
			if got != tt.expected {
				t.Errorf("getOrdinalSuffix(%q) = %v, want %v", tt.numStr, got, tt.expected)
			}
		})
	}
}

func TestFormatYearsAsRanges(t *testing.T) {
	tests := []struct {
		name     string
		years    []string
		sport    string
		expected string
	}{
		{
			name:     "empty years",
			years:    []string{},
			sport:    "baseball",
			expected: "",
		},
		{
			name:     "single baseball year",
			years:    []string{"2020"},
			sport:    "baseball",
			expected: "2020",
		},
		{
			name:     "consecutive baseball years",
			years:    []string{"2018", "2019", "2020"},
			sport:    "baseball",
			expected: "2018-2020",
		},
		{
			name:     "non-consecutive baseball years",
			years:    []string{"2015", "2017", "2019"},
			sport:    "baseball",
			expected: "2015, 2017, 2019",
		},
		{
			name:     "mixed consecutive and non-consecutive baseball",
			years:    []string{"2015", "2016", "2017", "2019", "2020"},
			sport:    "baseball",
			expected: "2015-2017, 2019-2020",
		},
		{
			name:     "single basketball year",
			years:    []string{"2020"},
			sport:    "basketball",
			expected: "2020-2021",
		},
		{
			name:     "consecutive basketball years",
			years:    []string{"2018", "2019", "2020"},
			sport:    "basketball",
			expected: "2018-2021",
		},
		{
			name:     "basketball season format",
			years:    []string{"2020-21", "2021-22"},
			sport:    "basketball",
			expected: "2020-2022",
		},
		{
			name:     "current season baseball",
			years:    []string{"2023", "2024", "2025"},
			sport:    "baseball",
			expected: "2023-Present",
		},
		{
			name:     "current season basketball",
			years:    []string{"2023", "2024", "2025"},
			sport:    "basketball",
			expected: "2023-Present",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatYearsAsRanges(tt.years, tt.sport)
			if got != tt.expected {
				t.Errorf("formatYearsAsRanges(%v, %q) = %v, want %v", tt.years, tt.sport, got, tt.expected)
			}
		})
	}
}

func TestFormatDraftInformation(t *testing.T) {
	tests := []struct {
		name        string
		draftText   string
		sport       string
		draftSchool string
		expected    string
	}{
		{
			name:        "undrafted player",
			draftText:   "Player was undrafted",
			sport:       "football",
			draftSchool: "",
			expected:    "Undrafted",
		},
		{
			name:        "NFL draft with school",
			draftText:   "Draft: Buffalo Bills in the 1st round (4th overall) of the 2014 NFL Draft.",
			sport:       "football",
			draftSchool: "Syracuse",
			expected:    "2014: 1st Rd (4th Ovr) from Syracuse",
		},
		{
			name:        "NBA draft with school",
			draftText:   "Draft: Washington Wizards, 1st round (18th pick, 18th overall), 2025 NBA Draft",
			sport:       "basketball",
			draftSchool: "Duke",
			expected:    "2025: 1st Rd (18th Ovr) from Duke",
		},
		{
			name:        "NBA draft without school",
			draftText:   "Draft: Los Angeles Lakers, 2nd round (48th pick, 48th overall), 2020 NBA Draft",
			sport:       "basketball",
			draftSchool: "",
			expected:    "2020: 2nd Rd (48th Ovr)",
		},
		{
			name:        "MLB draft with school",
			draftText:   "Drafted by the Los Angeles Angels of Anaheim in the 1st round (25th) of the 2009 MLB June Amateur Draft from Vanderbilt University (Nashville, TN)",
			sport:       "baseball",
			draftSchool: "",
			expected:    "2009: 1st Rd (25th Ovr) from Vanderbilt University (Nashville, TN)",
		},
		{
			name:        "NFL 11th pick",
			draftText:   "Draft: Green Bay Packers in the 1st round (11th overall) of the 2019 NFL Draft.",
			sport:       "football",
			draftSchool: "Michigan",
			expected:    "2019: 1st Rd (11th Ovr) from Michigan",
		},
		{
			name:        "NFL 22nd pick",
			draftText:   "Draft: Tennessee Titans in the 1st round (22nd overall) of the 2018 NFL Draft.",
			sport:       "football",
			draftSchool: "Alabama",
			expected:    "2018: 1st Rd (22nd Ovr) from Alabama",
		},
		{
			name:        "NFL 3rd pick",
			draftText:   "Draft: San Francisco 49ers in the 1st round (3rd overall) of the 2021 NFL Draft.",
			sport:       "football",
			draftSchool: "North Dakota State",
			expected:    "2021: 1st Rd (3rd Ovr) from North Dakota State",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDraftInformation(tt.draftText, tt.sport, tt.draftSchool)
			if got != tt.expected {
				t.Errorf("formatDraftInformation() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestUpdateDailyStreak tests the updateDailyStreak helper function
func TestUpdateDailyStreak(t *testing.T) {
	tests := []struct {
		name                     string
		existingUserStats        *UserStats
		playDate                 string
		expectedStreak           int
		expectedLastDayPlayed    string
	}{
		{
			name:                  "nil user stats - should do nothing",
			existingUserStats:     nil,
			playDate:              "2024-01-15",
			expectedStreak:        0,
			expectedLastDayPlayed: "",
		},
		{
			name: "consecutive day - increment streak",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 5,
				LastDayPlayed:      "2024-01-14",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-01-15",
			expectedStreak:        6,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "same day - keep streak unchanged",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 3,
				LastDayPlayed:      "2024-01-15",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-01-15",
			expectedStreak:        3,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "missed one day - reset streak to 1",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 10,
				LastDayPlayed:      "2024-01-13",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-01-15",
			expectedStreak:        1,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "missed multiple days - reset streak to 1",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 7,
				LastDayPlayed:      "2024-01-10",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-01-20",
			expectedStreak:        1,
			expectedLastDayPlayed: "2024-01-20",
		},
		{
			name: "streak continues after weekend - friday to saturday",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 4,
				LastDayPlayed:      "2024-01-12",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-01-13",
			expectedStreak:        5,
			expectedLastDayPlayed: "2024-01-13",
		},
		{
			name: "existing user with empty lastDayPlayed",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 2,
				LastDayPlayed:      "",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-01-15",
			expectedStreak:        2,
			expectedLastDayPlayed: "2024-01-15",
		},
		{
			name: "month transition - consecutive day",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 15,
				LastDayPlayed:      "2024-01-31",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-02-01",
			expectedStreak:        16,
			expectedLastDayPlayed: "2024-02-01",
		},
		{
			name: "year transition - consecutive day",
			existingUserStats: &UserStats{
				UserId:             "user123",
				CurrentDailyStreak: 20,
				LastDayPlayed:      "2023-12-31",
				UserCreated:        time.Now(),
			},
			playDate:              "2024-01-01",
			expectedStreak:        21,
			expectedLastDayPlayed: "2024-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of the user stats if not nil
			var userStats *UserStats
			if tt.existingUserStats != nil {
				userStats = &UserStats{
					UserId:             tt.existingUserStats.UserId,
					CurrentDailyStreak: tt.existingUserStats.CurrentDailyStreak,
					LastDayPlayed:      tt.existingUserStats.LastDayPlayed,
					UserCreated:        tt.existingUserStats.UserCreated,
					Sports:             tt.existingUserStats.Sports,
				}
			}

			// Call the function
			updateDailyStreak(userStats, tt.playDate)

			// Verify the results
			if userStats == nil {
				if tt.expectedStreak != 0 || tt.expectedLastDayPlayed != "" {
					t.Errorf("Expected nil userStats to remain nil")
				}
			} else {
				if userStats.CurrentDailyStreak != tt.expectedStreak {
					t.Errorf("Expected streak %d, got %d", tt.expectedStreak, userStats.CurrentDailyStreak)
				}

				if userStats.LastDayPlayed != tt.expectedLastDayPlayed {
					t.Errorf("Expected lastDayPlayed %s, got %s", tt.expectedLastDayPlayed, userStats.LastDayPlayed)
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
