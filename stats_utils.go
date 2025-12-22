package main

// updateStatsWithResult updates statistics with a submitted result
// Works with both RoundStats and SportStats since they both embed Stats
func updateStatsWithResult(stats *Stats, result *Result) {
	// Update total plays
	stats.TotalPlays++

	// Update percentage correct and average correct score
	if result.IsCorrect {
		correctCount := int(stats.PercentageCorrect * float64(stats.TotalPlays-1) / 100)
		correctCount++
		stats.PercentageCorrect = float64(correctCount) * 100 / float64(stats.TotalPlays)

		// Update average correct score
		totalCorrectScore := stats.AverageCorrectScore * float64(correctCount-1)
		totalCorrectScore += float64(result.Score)
		stats.AverageCorrectScore = totalCorrectScore / float64(correctCount)
	}

	// Update highest score
	if result.Score > stats.HighestScore {
		stats.HighestScore = result.Score
	}

	// Update average number of tile flips
	totalTileFlips := stats.AverageNumberOfTileFlips * float64(stats.TotalPlays-1)
	totalTileFlips += float64(len(result.TilesFlipped))
	stats.AverageNumberOfTileFlips = totalTileFlips / float64(stats.TotalPlays)

	// Track tile flips
	if len(result.TilesFlipped) > 0 {
		// Track first tile flipped
		incrementTileTracker(&stats.FirstTileFlippedTracker, result.TilesFlipped[0])

		// Track last tile flipped
		incrementTileTracker(&stats.LastTileFlippedTracker, result.TilesFlipped[len(result.TilesFlipped)-1])

		// Track all tiles flipped
		for _, tile := range result.TilesFlipped {
			incrementTileTracker(&stats.MostTileFlippedTracker, tile)
		}

		// Recalculate most/least common tiles
		stats.MostCommonFirstTileFlipped = findMostCommonTile(&stats.FirstTileFlippedTracker)
		stats.MostCommonLastTileFlipped = findMostCommonTile(&stats.LastTileFlippedTracker)
		stats.MostCommonTileFlipped = findMostCommonTile(&stats.MostTileFlippedTracker)
		stats.LeastCommonTileFlipped = findLeastCommonTile(&stats.MostTileFlippedTracker)
	}
}
