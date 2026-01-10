package main

// incrementTileTracker increments the appropriate field in a TileFlipTracker based on the tile name
func incrementTileTracker(tracker *TileFlipTracker, tileName string) {
	if tracker == nil || tileName == "" {
		return
	}

	switch tileName {
	case TileBio:
		tracker.Bio++
	case TilePlayerInformation:
		tracker.PlayerInformation++
	case TileDraftInformation:
		tracker.DraftInformation++
	case TileTeamsPlayedOn:
		tracker.TeamsPlayedOn++
	case TileJerseyNumbers:
		tracker.JerseyNumbers++
	case TileCareerStats:
		tracker.CareerStats++
	case TilePersonalAchievements:
		tracker.PersonalAchievements++
	case TilePhoto:
		tracker.Photo++
	case TileYearsActive:
		tracker.YearsActive++
	case TileInitials:
		tracker.Initials++
	case TileNicknames:
		tracker.Nicknames++
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
		TileBio:                  tracker.Bio,
		TilePlayerInformation:    tracker.PlayerInformation,
		TileDraftInformation:     tracker.DraftInformation,
		TileTeamsPlayedOn:        tracker.TeamsPlayedOn,
		TileJerseyNumbers:        tracker.JerseyNumbers,
		TileCareerStats:          tracker.CareerStats,
		TilePersonalAchievements: tracker.PersonalAchievements,
		TilePhoto:                tracker.Photo,
		TileYearsActive:          tracker.YearsActive,
		TileInitials:             tracker.Initials,
		TileNicknames:            tracker.Nicknames,
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
		TileBio:                  tracker.Bio,
		TilePlayerInformation:    tracker.PlayerInformation,
		TileDraftInformation:     tracker.DraftInformation,
		TileTeamsPlayedOn:        tracker.TeamsPlayedOn,
		TileJerseyNumbers:        tracker.JerseyNumbers,
		TileCareerStats:          tracker.CareerStats,
		TilePersonalAchievements: tracker.PersonalAchievements,
		TilePhoto:                tracker.Photo,
		TileYearsActive:          tracker.YearsActive,
		TileInitials:             tracker.Initials,
		TileNicknames:            tracker.Nicknames,
	}

	for tileName, count := range tiles {
		if count >= 0 && (minCount == -1 || count < minCount) {
			minCount = count
			leastCommon = tileName
		}
	}

	return leastCommon
}
