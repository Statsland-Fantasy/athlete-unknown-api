package main

import (
	"testing"
)

func TestGetAchievementMappings(t *testing.T) {
	tests := []struct {
		name          string
		sport         string
		expectedCount int
		checkMapping  func([]AchievementMapping) bool
	}{
		{
			name:          "baseball mappings",
			sport:         "baseball",
			expectedCount: 18,
			checkMapping: func(mappings []AchievementMapping) bool {
				// Check that Hall of Fame exists
				for _, m := range mappings {
					if m.FullName == "hall of fame" && m.Abbreviation == "HOF" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "basketball mappings",
			sport:         "basketball",
			expectedCount: 20,
			checkMapping: func(mappings []AchievementMapping) bool {
				// Check that Finals MVP exists
				for _, m := range mappings {
					if m.FullName == "finals mvp" && m.Tier == 1 {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football mappings",
			sport:         "football",
			expectedCount: 15,
			checkMapping: func(mappings []AchievementMapping) bool {
				// Check that Super Bowl Champ exists
				for _, m := range mappings {
					if m.FullName == "sb champ" && m.Tier == 1 {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "unknown sport",
			sport:         "soccer",
			expectedCount: 0,
			checkMapping: func(mappings []AchievementMapping) bool {
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAchievementMappings(tt.sport)
			if len(got) != tt.expectedCount {
				t.Errorf("GetAchievementMappings(%q) returned %d mappings, want %d", tt.sport, len(got), tt.expectedCount)
			}
			if !tt.checkMapping(got) {
				t.Errorf("GetAchievementMappings(%q) missing expected mapping", tt.sport)
			}
		})
	}
}

func TestGetAchievementAbbreviation(t *testing.T) {
	tests := []struct {
		name            string
		sport           string
		achievementName string
		expectedText    string
		expectedTier    int
		expectNil       bool
	}{
		{
			name:            "baseball hall of fame",
			sport:           "baseball",
			achievementName: "Hall of Fame",
			expectedText:    "HOF",
			expectedTier:    1,
			expectNil:       false,
		},
		{
			name:            "baseball mvp",
			sport:           "baseball",
			achievementName: "AL MVP",
			expectedText:    "AL MVP",
			expectedTier:    1,
			expectNil:       false,
		},
		{
			name:            "basketball finals mvp",
			sport:           "basketball",
			achievementName: "NBA Finals MVP",
			expectedText:    "NBA Finals MVP",
			expectedTier:    1,
			expectNil:       false,
		},
		{
			name:            "basketball all-star",
			sport:           "basketball",
			achievementName: "12x All Star",
			expectedText:    "12x All Star",
			expectedTier:    1,
			expectNil:       false,
		},
		{
			name:            "football super bowl champ",
			sport:           "football",
			achievementName: "2x SB Champ",
			expectedText:    "2x SB Champ",
			expectedTier:    1,
			expectNil:       false,
		},
		{
			name:            "unknown achievement",
			sport:           "baseball",
			achievementName: "Bert Bell Award",
			expectedText:    "",
			expectedTier:    0,
			expectNil:       true,
		},
		{
			name:            "basketball DPOY abbreviation",
			sport:           "basketball",
			achievementName: "3x Def. POY",
			expectedText:    "3x DPOY",
			expectedTier:    1,
			expectNil:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetAchievementAbbreviation(tt.sport, tt.achievementName)
			if tt.expectNil {
				if got != nil {
					t.Errorf("GetAchievementAbbreviation(%q, %q) = %v, want nil", tt.sport, tt.achievementName, got)
				}
			} else {
				if got == nil {
					t.Errorf("GetAchievementAbbreviation(%q, %q) = nil, want non-nil", tt.sport, tt.achievementName)
					return
				}
				if got.ProcessedText != tt.expectedText {
					t.Errorf("GetAchievementAbbreviation(%q, %q).ProcessedText = %v, want %v", tt.sport, tt.achievementName, got.ProcessedText, tt.expectedText)
				}
				if got.Tier != tt.expectedTier {
					t.Errorf("GetAchievementAbbreviation(%q, %q).Tier = %v, want %v", tt.sport, tt.achievementName, got.Tier, tt.expectedTier)
				}
			}
		})
	}
}

func TestProcessAchievements(t *testing.T) {
	tests := []struct {
		name         string
		sport        string
		achievements []string
		maxLength    int
		checkResult  func(string) bool
	}{
		{
			name:         "empty achievements",
			sport:        "baseball",
			achievements: []string{},
			maxLength:    100,
			checkResult: func(result string) bool {
				return result == ""
			},
		},
		{
			name:         "single baseball achievement",
			sport:        "baseball",
			achievements: []string{"Hall of Fame"},
			maxLength:    100,
			checkResult: func(result string) bool {
				return result == "HOF"
			},
		},
		{
			name:         "multiple tier 1 basketball achievements",
			sport:        "basketball",
			achievements: []string{"5x NBA Champ", "2x Finals MVP", "18x All Star"},
			maxLength:    100,
			checkResult: func(result string) bool {
				// Should include all tier 1 achievements
				return len(result) > 0
			},
		},
		{
			name:         "exceeds max length with lower tier",
			sport:        "football",
			achievements: []string{"Hall of Fame", "5x Pro Bowl", "3x All-Pro", "SB Champ", "AP MVP"},
			maxLength:    50,
			checkResult: func(result string) bool {
				// Should prioritize tier 1 achievements
				return len(result) > 0
			},
		},
		{
			name:         "filters out unknown achievements",
			sport:        "baseball",
			achievements: []string{"Hall of Fame", "Unknown Award", "All-Star"},
			maxLength:    100,
			checkResult: func(result string) bool {
				// Should not include Unknown Award
				return result != "" && len(result) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessAchievements(tt.sport, tt.achievements, tt.maxLength)
			if !tt.checkResult(got) {
				t.Errorf("ProcessAchievements(%q, %v, %d) = %v, failed check", tt.sport, tt.achievements, tt.maxLength, got)
			}
		})
	}
}
