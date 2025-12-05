package main

import (
	"fmt"
	"sort"
	"strings"
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
