package main

// incrementTileTracker increments the appropriate field in a TileFlipTracker based on the tile name
func incrementTileTracker(tracker *TileFlipTracker, tileName string) {
	if tracker == nil || tileName == "" {
		return
	}

	switch tileName {
	case "bio":
		tracker.Bio++
	case "playerInformation":
		tracker.PlayerInformation++
	case "draftInformation":
		tracker.DraftInformation++
	case "teamsPlayedOn":
		tracker.TeamsPlayedOn++
	case "jerseyNumbers":
		tracker.JerseyNumbers++
	case "careerStats":
		tracker.CareerStats++
	case "personalAchievements":
		tracker.PersonalAchievements++
	case "photo":
		tracker.Photo++
	case "yearsActive":
		tracker.YearsActive++
	}
}

// findMostCommonTile returns the tile name with the highest count in the tracker
func findMostCommonTile(tracker *TileFlipTracker) string {
	if tracker == nil {
		return ""
	}

	maxCount := 0
	mostCommon := ""

	tiles := map[string]int{
		"bio":                  tracker.Bio,
		"playerInformation":    tracker.PlayerInformation,
		"draftInformation":     tracker.DraftInformation,
		"teamsPlayedOn":        tracker.TeamsPlayedOn,
		"jerseyNumbers":        tracker.JerseyNumbers,
		"careerStats":          tracker.CareerStats,
		"personalAchievements": tracker.PersonalAchievements,
		"photo":                tracker.Photo,
		"yearsActive":          tracker.YearsActive,
	}

	for tileName, count := range tiles {
		if count > maxCount {
			maxCount = count
			mostCommon = tileName
		}
	}

	return mostCommon
}

// findLeastCommonTile returns the tile name with the lowest non-zero count in the tracker
func findLeastCommonTile(tracker *TileFlipTracker) string {
	if tracker == nil {
		return ""
	}

	minCount := -1
	leastCommon := ""

	tiles := map[string]int{
		"bio":                  tracker.Bio,
		"playerInformation":    tracker.PlayerInformation,
		"draftInformation":     tracker.DraftInformation,
		"teamsPlayedOn":        tracker.TeamsPlayedOn,
		"jerseyNumbers":        tracker.JerseyNumbers,
		"careerStats":          tracker.CareerStats,
		"personalAchievements": tracker.PersonalAchievements,
		"photo":                tracker.Photo,
		"yearsActive":          tracker.YearsActive,
	}

	for tileName, count := range tiles {
		if count > 0 && (minCount == -1 || count < minCount) {
			minCount = count
			leastCommon = tileName
		}
	}

	return leastCommon
}
