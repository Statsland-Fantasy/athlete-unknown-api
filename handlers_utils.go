package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isValidYear checks if a year string is valid (not empty and not a summary row label)
func isValidYear(year string) bool {
	if year == "" {
		return false
	}
	invalidValues := []string{"Season", "Career", "Yr", "Avg", "Average"}
	for _, invalid := range invalidValues {
		if year == invalid || strings.Contains(year, invalid) {
			return false
		}
	}
	return true
}

// formatYearsAsRanges converts a slice of year strings into consolidated ranges
// Example: ["2010", "2011", "2013", "2014"] -> "2010-2011, 2013-2014"
// If the last year matches the current season for the sport, it displays "Present"
// For basketball, handles season format like "2024-25" and extracts start year
func formatYearsAsRanges(years []string, sport string) string {
	if len(years) == 0 {
		return ""
	}

	// Convert strings to integers and sort
	var yearInts []int
	for _, y := range years {
		var yearInt int

		// For basketball, handle season format like "2024-25"
		if sport == "basketball" && strings.Contains(y, "-") {
			// Extract the first year from the season range (e.g., "2024" from "2024-25")
			parts := strings.Split(y, "-")
			if len(parts) > 0 {
				_, err := fmt.Sscanf(parts[0], "%d", &yearInt)
				if err == nil {
					yearInts = append(yearInts, yearInt)
				}
			}
		} else {
			// Regular year format
			_, err := fmt.Sscanf(y, "%d", &yearInt)
			if err == nil {
				yearInts = append(yearInts, yearInt)
			}
		}
	}

	if len(yearInts) == 0 {
		return ""
	}

	sort.Ints(yearInts)
	currentSeasonYear := GetCurrentSeasonYear(sport)

	// Build ranges
	var ranges []string
	rangeStart := yearInts[0]
	rangeEnd := yearInts[0]

	for i := 1; i < len(yearInts); i++ {
		if yearInts[i] == rangeEnd+1 {
			// Consecutive year, extend the range
			rangeEnd = yearInts[i]
		} else {
			// Gap found, save the current range and start a new one
			if rangeStart == rangeEnd {
				if rangeStart == currentSeasonYear {
					ranges = append(ranges, "Present")
				} else {
					// For basketball single seasons, add 1 to show the ending year
					if sport == "basketball" {
						ranges = append(ranges, fmt.Sprintf("%d-%d", rangeStart, rangeStart+1))
					} else {
						ranges = append(ranges, fmt.Sprintf("%d", rangeStart))
					}
				}
			} else {
				// For basketball, add 1 to the end year to show the actual ending year
				displayEndYear := rangeEnd
				if sport == "basketball" {
					displayEndYear = rangeEnd + 1
				}

				endStr := fmt.Sprintf("%d", displayEndYear)
				if rangeEnd == currentSeasonYear {
					endStr = "Present"
				}
				ranges = append(ranges, fmt.Sprintf("%d-%s", rangeStart, endStr))
			}
			rangeStart = yearInts[i]
			rangeEnd = yearInts[i]
		}
	}

	// Add the last range
	if rangeStart == rangeEnd {
		if rangeStart == currentSeasonYear {
			ranges = append(ranges, "Present")
		} else {
			// For basketball single seasons, add 1 to show the ending year
			if sport == "basketball" {
				ranges = append(ranges, fmt.Sprintf("%d-%d", rangeStart, rangeStart+1))
			} else {
				ranges = append(ranges, fmt.Sprintf("%d", rangeStart))
			}
		}
	} else {
		// For basketball, add 1 to the end year to show the actual ending year
		displayEndYear := rangeEnd
		if sport == "basketball" {
			displayEndYear = rangeEnd + 1
		}

		endStr := fmt.Sprintf("%d", displayEndYear)
		if rangeEnd == currentSeasonYear {
			endStr = "Present"
		}
		ranges = append(ranges, fmt.Sprintf("%d-%s", rangeStart, endStr))
	}

	return strings.Join(ranges, ", ")
}

// formatDraftInformation transforms draft text from verbose format to concise format
// Example: "Draft: Buffalo Bills in the 1st round (4th overall) of the 2014 NFL Draft." with draftSchool="Syracuse" -> "2014: 1st Rd (4th Ovr) from Syracuse"
// Example: "Draft: Washington Wizards, 1st round (18th pick, 18th overall), 2025 NBA Draft" with draftSchool="Duke" -> "2025: 1st Rd (18th Ovr) from Duke"
// Example: "Draft: Drafted by the Los Angeles Angels of Anaheim in the 1st round (25th) of the 2009 MLB June Amateur Draft from Vanderbilt University" -> "2009: 1st Rd (25th Ovr) from Vanderbilt University"
// For multiple drafts, it selects the most recent year (the draft they actually signed with)
func formatDraftInformation(draftText, sport, draftSchool string) string {
	// Return "Undrafted" if the text indicates no draft
	if strings.Contains(strings.ToLower(draftText), "undrafted") {
		return "Undrafted"
	}

	if sport == "baseball" {
		// Baseball-specific pattern with multiple variations:
		// 1. "Xth round (Yth) of the YYYY MLB ... Draft from [School] (City, State)" - with pick number
		// 2. "Xth round of the YYYY MLB ... Draft from [School] (City, State)" - without pick number

		// Try pattern with pick number first (includes city/state in parentheses)
		pattern := `(\d+)(?:st|nd|rd|th)\s+round\s+\((\d+)(?:st|nd|rd|th)\)\s+of\s+the\s+(\d{4})\s+MLB[^)]*?Draft\s+from\s+([^)]+\s*\([^)]+\))(?:\.|$|\s+and\s+)`
		re := regexp.MustCompile(pattern)
		allMatches := re.FindAllStringSubmatch(draftText, -1)

		// If no matches, try pattern without pick number but with "from" clause (includes city/state)
		if len(allMatches) == 0 {
			pattern = `(\d+)(?:st|nd|rd|th)\s+round\s+of\s+the\s+(\d{4})\s+MLB[^)]*?Draft\s+from\s+([^)]+\s*\([^)]+\))(?:\.|$|\s+and\s+)`
			re = regexp.MustCompile(pattern)
			allMatches = re.FindAllStringSubmatch(draftText, -1)

			if len(allMatches) > 0 {
				// Find the draft with the latest year
				var bestMatch []string
				var latestYear int = 0

				for _, matches := range allMatches {
					yearStr := matches[2] // Year is at index 2 for this pattern
					var yearNum int
					fmt.Sscanf(yearStr, "%d", &yearNum)

					if yearNum > latestYear {
						latestYear = yearNum
						bestMatch = matches
					}
				}

				if bestMatch != nil {
					round := bestMatch[1]
					year := bestMatch[2]
					school := strings.TrimSpace(bestMatch[3])

					roundSuffix := getOrdinalSuffix(round)

					// No pick number available, so just show round and school
					return fmt.Sprintf("%s: %s%s Rd from %s", year, round, roundSuffix, school)
				}
			}
		}

		// Try simpler pattern with pick number but no "from" clause
		if len(allMatches) == 0 {
			pattern = `(\d+)(?:st|nd|rd|th)\s+round\s+\((\d+)(?:st|nd|rd|th)\)\s+of\s+the\s+(\d{4})\s+MLB`
			re = regexp.MustCompile(pattern)
			allMatches = re.FindAllStringSubmatch(draftText, -1)
		}

		if len(allMatches) > 0 {
			// Find the draft with the latest year
			var bestMatch []string
			var latestYear int = 0

			for _, matches := range allMatches {
				yearStr := matches[3]
				var yearNum int
				fmt.Sscanf(yearStr, "%d", &yearNum)

				if yearNum > latestYear {
					latestYear = yearNum
					bestMatch = matches
				}
			}

			if bestMatch != nil {
				round := bestMatch[1]
				overall := bestMatch[2]
				year := bestMatch[3]
				school := ""
				if len(bestMatch) > 4 && bestMatch[4] != "" {
					school = strings.TrimSpace(bestMatch[4])
				}

				roundSuffix := getOrdinalSuffix(round)
				overallSuffix := getOrdinalSuffix(overall)

				if school != "" {
					return fmt.Sprintf("%s: %s%s Rd (%s%s Ovr) from %s", year, round, roundSuffix, overall, overallSuffix, school)
				}
				return fmt.Sprintf("%s: %s%s Rd (%s%s Ovr)", year, round, roundSuffix, overall, overallSuffix)
			}
		}
	}

	// Non-baseball format (NFL/NBA)
	// Regex pattern to extract: year, round number, and overall pick
	// Handles multiple formats:
	// - "1st round (4th overall) of the 2014" - with "overall" keyword
	// - "1st round (18th pick, 18th overall), 2025" - with "pick" and "overall"
	// - "2014 NFL Draft ... 1st round (4th overall)" - year before round
	pattern := `(\d+)(?:st|nd|rd|th)\s+round\s+\((?:[^)]*?(\d+)(?:st|nd|rd|th)\s+overall|(\d+)(?:st|nd|rd|th))\)[^0-9]*?(\d{4})|(\d{4})[^0-9]*?(\d+)(?:st|nd|rd|th)\s+round\s+\((?:[^)]*?(\d+)(?:st|nd|rd|th)\s+overall|(\d+)(?:st|nd|rd|th))\)`

	re := regexp.MustCompile(pattern)
	allMatches := re.FindAllStringSubmatch(draftText, -1)

	if len(allMatches) == 0 {
		// If pattern doesn't match, return original text
		return draftText
	}

	// If multiple drafts found, select the one with the latest year (most recent)
	var bestMatch []string
	var latestYear int = 0

	for _, matches := range allMatches {
		var yearStr string

		// Extract year from the match
		if matches[4] != "" {
			yearStr = matches[4] // First pattern
		} else {
			yearStr = matches[5] // Second pattern
		}

		// Convert year to int for comparison
		var yearNum int
		fmt.Sscanf(yearStr, "%d", &yearNum)

		// Keep the match with the latest (highest) year
		if yearNum > latestYear {
			latestYear = yearNum
			bestMatch = matches
		}
	}

	if bestMatch == nil {
		return draftText
	}

	var year, round, overall string

	// Check which pattern matched (pattern has two alternatives)
	if bestMatch[1] != "" {
		// First pattern: round info before year
		round = bestMatch[1]
		// Check if we captured "overall" format or just the number
		if bestMatch[2] != "" {
			overall = bestMatch[2] // Has "overall" keyword
		} else {
			overall = bestMatch[3] // No "overall" keyword
		}
		year = bestMatch[4]
	} else {
		// Second pattern: year before round info
		year = bestMatch[5]
		round = bestMatch[6]
		// Check if we captured "overall" format or just the number
		if bestMatch[7] != "" {
			overall = bestMatch[7] // Has "overall" keyword
		} else {
			overall = bestMatch[8] // No "overall" keyword
		}
	}

	// Get the ordinal suffix for round
	roundSuffix := getOrdinalSuffix(round)

	// Get the ordinal suffix for overall
	overallSuffix := getOrdinalSuffix(overall)

	// Format: "2014: 1st Rd (4th Ovr) from School" for football/basketball
	if draftSchool != "" {
		return fmt.Sprintf("%s: %s%s Rd (%s%s Ovr) from %s", year, round, roundSuffix, overall, overallSuffix, draftSchool)
	}
	return fmt.Sprintf("%s: %s%s Rd (%s%s Ovr)", year, round, roundSuffix, overall, overallSuffix)
}

// getOrdinalSuffix returns the ordinal suffix (st, nd, rd, th) for a number string
func getOrdinalSuffix(numStr string) string {
	// Get the last digit
	lastChar := numStr[len(numStr)-1]
	lastDigit := int(lastChar - '0')

	// Handle special cases for 11, 12, 13
	if len(numStr) >= 2 {
		lastTwo := numStr[len(numStr)-2:]
		if lastTwo == "11" || lastTwo == "12" || lastTwo == "13" {
			return "th"
		}
	}

	// Standard rules
	switch lastDigit {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

// abbreviatePositions abbreviates position names in the player information string
// This function looks for "Position:" or "Positions:" and abbreviates the position names
// Example: "Position: Shooting Guard" -> "Position: SG"
// Example: "Positions: Shooting Guard and Point Guard" -> "Positions: SG and PG"
func abbreviatePositions(playerInfo string) string {
	// Position abbreviation map - add your specific abbreviations here
	// Order matters: longer position names should be replaced first to avoid partial replacements
	positionReplacements := []struct {
		full  string
		abbrev string
	}{
		// Basketball positions (ordered by length, longest first)
		{"Shooting Guard", "SG"},
		{"Point Guard", "PG"},
		{"Small Forward", "SF"},
		{"Power Forward", "PF"},
		{"Forward", "F"},
		{"Center", "C"},
		{"Guard", "G"},

		// Football positions
		// Already abbreviated

		// Baseball positions
		{"Designated Hitter", "DH"},
		{"Second Baseman", "2B"},
		{"First Baseman", "1B"},
		{"Third Baseman", "3B"},
		{"Centerfielder", "CF"},
		{"Rightfielder", "RF"},
		{"Leftfielder", "LF"},
		{"Outfielder", "OF"},
		{"Shortstop", "SS"},
		{"Catcher", "C"},
		{"Pitcher", "P"},
		
	}

	result := playerInfo

	// Replace each position with its abbreviation in order (longest first)
	for _, replacement := range positionReplacements {
		result = strings.ReplaceAll(result, replacement.full, replacement.abbrev)
	}	

	result = strings.ReplaceAll(result, " and ", ", ") // no "ands", all comma separated

	return result
}

// updateDailyStreak updates the user's daily streak and last day played based on the play date
// For new users (userStats is nil), it initializes the streak to 1 and sets the last day played
// For existing users, it:
// - Increments the streak if the play date is exactly 1 day after the last day played
// - Resets the streak to 1 if more than 1 day has passed since the last day played
// - Keeps the streak unchanged if playing on the same day
// Always updates lastDayPlayed to the current play date
func updateDailyStreak(userStats *UserStats, playDate string) {
	if userStats == nil {
		return
	}

	// Check if we have a previous play date to compare against
	if userStats.LastDayPlayed != "" {
		lastPlayed, err := time.Parse("2006-01-02", userStats.LastDayPlayed)
		currentPlay, err2 := time.Parse("2006-01-02", playDate)

		if err == nil && err2 == nil {
			daysDiff := int(currentPlay.Sub(lastPlayed).Hours() / 24)

			if daysDiff == 1 {
				// Consecutive day - increment streak
				userStats.CurrentDailyStreak++
			} else if daysDiff > 1 {
				// Missed a day - reset streak to 1
				userStats.CurrentDailyStreak = 1
			}
			// If daysDiff == 0, same day - don't change streak
		}
	}

	// Update last day played
	userStats.LastDayPlayed = playDate
}
