package main

import (
	"strings"
)

// AchievementMapping maps full award names to their abbreviations
type AchievementMapping struct {
	FullName     string
	Abbreviation string
	Tier         int
}

// ProcessedAchievement holds an achievement with its tier and processed text
type ProcessedAchievement struct {
	OriginalText  string
	ProcessedText string
	Tier          int
}

// GetAchievementMappings returns all achievement mappings
func GetAchievementMappings(sport string) []AchievementMapping {
	if sport == SportBaseball {
		return []AchievementMapping{
			// baseball - more specific matches first
			{FullName: "ws mvp", Abbreviation: "", Tier: 1},
			{FullName: "alcs mvp", Abbreviation: "", Tier: 2},
			{FullName: "nlcs mvp", Abbreviation: "", Tier: 2},
			{FullName: "as mvp", Abbreviation: "", Tier: 3},
			{FullName: "hall of fame", Abbreviation: "HOF", Tier: 1},
			{FullName: "world series", Abbreviation: "WS Champ", Tier: 1},
			{FullName: "mvp", Abbreviation: "", Tier: 1},
			{FullName: "cy young", Abbreviation: "", Tier: 1},
			{FullName: "rookie of the year", Abbreviation: "ROY", Tier: 1},
			{FullName: "all-star", Abbreviation: "", Tier: 1},
			{FullName: "gold glove", Abbreviation: "GG", Tier: 2},
			{FullName: "silver slugger", Abbreviation: "SS", Tier: 2},
			{FullName: "platinum glove", Abbreviation: "", Tier: 2},
			{FullName: "batting title", Abbreviation: "", Tier: 3},
			{FullName: "era title", Abbreviation: "", Tier: 3},
			{FullName: "triple crown", Abbreviation: "", Tier: 3},
			{FullName: "hr derby champ", Abbreviation: "", Tier: 4},
			{FullName: "clemente", Abbreviation: "", Tier: 4},
			// exclude: TSN Major League Player of the Year, Wilson Overall Def Player
		}
	} else if sport == SportBasketball {
		return []AchievementMapping{
			// basketball - more specific matches first
			{FullName: "finals mvp", Abbreviation: "", Tier: 1},
			{FullName: "as mvp", Abbreviation: "", Tier: 3},
			{FullName: "ist mvp", Abbreviation: "", Tier: 3},
			{FullName: "hall of fame", Abbreviation: "HOF", Tier: 1},
			{FullName: "nba champ", Abbreviation: "", Tier: 1},
			{FullName: "mvp", Abbreviation: "", Tier: 1},
			{FullName: "roy", Abbreviation: "", Tier: 1},
			{FullName: "def. poy", Abbreviation: "DPOY", Tier: 1},
			{FullName: "sixth man", Abbreviation: "", Tier: 1},
			{FullName: "most improved", Abbreviation: "MIPOY", Tier: 1},
			{FullName: "all star", Abbreviation: "", Tier: 1},
			{FullName: "all-nba", Abbreviation: "", Tier: 1},
			{FullName: "all-defensive", Abbreviation: "All-Def", Tier: 1},
			{FullName: "all-rookie", Abbreviation: "", Tier: 2},
			{FullName: "scoring champ", Abbreviation: "", Tier: 2},
			{FullName: "nba 75th anniv. team", Abbreviation: "75th Anniv.", Tier: 3},
			{FullName: "trb champ", Abbreviation: "REB Champ", Tier: 3},
			{FullName: "ast champ", Abbreviation: "", Tier: 3},
			{FullName: "stl champ", Abbreviation: "", Tier: 3},
			{FullName: "blk champ", Abbreviation: "", Tier: 3},
		}
	} else if sport == SportFootball {
		return []AchievementMapping{
			// football
			{FullName: "hall of fame", Abbreviation: "HOF", Tier: 1},
			{FullName: "sb champ", Abbreviation: "", Tier: 1},
			{FullName: "nfl champ", Abbreviation: "", Tier: 1},
			{FullName: "pro bowl", Abbreviation: "", Tier: 1},
			{FullName: "all-pro", Abbreviation: "", Tier: 1},
			{FullName: "ap mvp", Abbreviation: "MVP", Tier: 1},
			{FullName: "sb", Abbreviation: "", Tier: 1}, // should be SB ____ MVP
			{FullName: "ap off. poy", Abbreviation: "OPOY", Tier: 1},
			{FullName: "ap def. poy", Abbreviation: "DPOY", Tier: 1},
			{FullName: "ap off. roy", Abbreviation: "OROY", Tier: 1},
			{FullName: "ap def. roy", Abbreviation: "DROY", Tier: 1},
			{FullName: "ap comeback player", Abbreviation: "CPOY", Tier: 1},
			{FullName: "walter payton moty", Abbreviation: "", Tier: 2},
			{FullName: "hof all-", Abbreviation: "", Tier: 2}, // should be HOF All-2000s/1920s, etc Team
			{FullName: "alan page award", Abbreviation: "", Tier: 2},
		}
	}

	return []AchievementMapping{}
}

// GetAchievementAbbreviation returns the abbreviated form of an achievement name
// If no mapping exists, returns the original name and tier 999 (lowest priority)
func GetAchievementAbbreviation(sport, achievementName string) *ProcessedAchievement {
	// Normalize the input to lowercase for comparison
	normalizedName := strings.ToLower(achievementName)

	mappings := GetAchievementMappings(sport)

	// Sort mappings by priority: longer/more specific strings should be checked first
	// This ensures "ist mvp" matches before "mvp", "finals mvp" before "mvp", etc.
	// We process them in order, and mappings are already ordered by specificity in GetAchievementMappings()

	for _, mapping := range mappings {
		if strings.Contains(normalizedName, mapping.FullName) {
			if mapping.Abbreviation == "" {
				return &ProcessedAchievement{
					OriginalText:  achievementName,
					ProcessedText: achievementName,
					Tier:          mapping.Tier,
				}
			}
			return &ProcessedAchievement{
				OriginalText:  achievementName,
				ProcessedText: strings.ReplaceAll(normalizedName, mapping.FullName, mapping.Abbreviation),
				Tier:          mapping.Tier,
			}
		}
	}

	// If no mapping found, do not return achievement (ex: Bert Bell Award)
	return nil
}

// ProcessAchievements takes a slice of achievement names and returns abbreviated versions
// If the combined string exceeds maxLength, it filters out lower-priority achievements
func ProcessAchievements(sport string, achievements []string, maxLength int) string {
	var processed []ProcessedAchievement
	for _, achievement := range achievements {
		abbreviatedAchievement := GetAchievementAbbreviation(sport, achievement)
		if abbreviatedAchievement != nil {
			processed = append(processed, *abbreviatedAchievement)
		}
	}

	// Try to fit achievements within maxLength by progressively adding the lowest tier achievements first (Tier 1), and so on until the character limit is reached
	var filteredAchievements []string
	for maxTier := 1; maxTier <= 9; maxTier++ {

		newFilteredAchievements := filteredAchievements
		for _, p := range processed {
			if p.Tier == maxTier { // only test one tier value at a time
				newFilteredAchievements = append(newFilteredAchievements, p.ProcessedText)
			}
		}

		testLengthString := strings.Join(newFilteredAchievements, ", ")
		if len(testLengthString) > maxLength { // achievements string is greater than allowed string length
			if maxTier == 1 {
				return strings.Join(newFilteredAchievements, ", ") // always return all Tier 1 stats even if longer than maxLength
			}
			return strings.Join(filteredAchievements, ", ")
		}

		// for next loop, update achievements list
		filteredAchievements = newFilteredAchievements
	}

	return strings.Join(filteredAchievements, ", ")
}
