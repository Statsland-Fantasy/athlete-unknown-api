package main

import (
	"regexp"
	"strings"
)

// getPlayerInitials extracts the initials from a player's full name
// Returns the first letter of each word in uppercase with periods
// Preserves suffixes like Jr., Sr., and Roman numerals (e.g., III, IV)
// Empty names return empty string
func getPlayerInitials(name string) string {
	// Trim whitespace and check for empty string
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	// Split name into words
	words := strings.Fields(name)
	if len(words) == 0 {
		return ""
	}

	// Regex to match suffixes: Jr., Sr., or Roman numerals
	suffixRegex := regexp.MustCompile(`^(Jr\.?|Sr\.?|[IVX]+)$`)

	// Check if the last word is a suffix
	var suffix string
	if len(words) > 0 && suffixRegex.MatchString(words[len(words)-1]) {
		suffix = words[len(words)-1]
		words = words[:len(words)-1] // Remove suffix from words to process
	}

	// Extract first letter of each word, handling hyphens
	var initials strings.Builder
	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		// Split by hyphen to handle hyphenated names
		parts := strings.Split(word, "-")
		for _, part := range parts {
			if len(part) > 0 {
				// Add separator before if not the first initial
				if initials.Len() > 0 {
					// Check if the last character is a hyphen
					str := initials.String()
					if str[len(str)-1] != '-' {
						initials.WriteString(".")
					}
				}

				// Convert to uppercase and take first character
				initials.WriteString(strings.ToUpper(string(part[0])))
			}
		}
	}

	// Add final period
	if initials.Len() > 0 {
		initials.WriteString(".")
	}

	// Add suffix if present
	if suffix != "" {
		if initials.Len() > 0 {
			initials.WriteString(" ")
		}
		initials.WriteString(strings.ToUpper(suffix))
	}

	return initials.String()
}

// updateStatsWithResult updates statistics with a submitted result
// Works with both RoundStats and SportStats since they both embed Stats
func updateStatsWithResult(stats *Stats, result *Result) {
	// Update total plays
	stats.TotalPlays++

	// Calculate current correct count from previous percentage
	correctCount := int(stats.PercentageCorrect * float64(stats.TotalPlays-1) / 100)

	// Update correct count and average correct score if this result is correct
	if result.IsCorrect {
		correctCount++

		// Update average correct score
		totalCorrectScore := stats.AverageCorrectScore * float64(correctCount-1)
		totalCorrectScore += float64(result.Score)
		stats.AverageCorrectScore = totalCorrectScore / float64(correctCount)
	}

	// Update percentage correct (whether result was correct or incorrect)
	stats.PercentageCorrect = float64(correctCount) * 100 / float64(stats.TotalPlays)

	// Update highest score
	if result.Score > stats.HighestScore {
		stats.HighestScore = result.Score
	}

	// Calculate current incorrect guesses from previous average
	incorrectGuesses := int(stats.AverageIncorrectGuesses * float64(stats.TotalPlays-1))
	incorrectGuesses += result.IncorrectGuesses
	stats.AverageIncorrectGuesses = float64(incorrectGuesses) / float64(stats.TotalPlays)

	// Update average number of tile flips
	totalTileFlips := stats.AverageNumberOfTileFlips * float64(stats.TotalPlays-1)
	totalTileFlips += float64(len(result.FlippedTiles))
	stats.AverageNumberOfTileFlips = totalTileFlips / float64(stats.TotalPlays)

	// Track tile flips
	if len(result.FlippedTiles) > 0 {
		// Track first tile flipped
		incrementTileTracker(&stats.FirstTileFlippedTracker, result.FlippedTiles[0])

		// Track last tile flipped
		incrementTileTracker(&stats.LastTileFlippedTracker, result.FlippedTiles[len(result.FlippedTiles)-1])

		// Track all tiles flipped
		for _, tile := range result.FlippedTiles {
			incrementTileTracker(&stats.MostTileFlippedTracker, tile)
		}

		// Recalculate most/least common tiles
		stats.MostCommonFirstTileFlipped = findMostCommonTile(&stats.FirstTileFlippedTracker)
		stats.MostCommonLastTileFlipped = findMostCommonTile(&stats.LastTileFlippedTracker)
		stats.MostCommonTileFlipped = findMostCommonTile(&stats.MostTileFlippedTracker)
		stats.LeastCommonTileFlipped = findLeastCommonTile(&stats.MostTileFlippedTracker)
	}
}
