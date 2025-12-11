package main

// Tile name constants
const (
	TileBio                   = "bio"
	TilePlayerInformation     = "playerInformation"
	TileDraftInformation      = "draftInformation"
	TileTeamsPlayedOn         = "teamsPlayedOn"
	TileJerseyNumbers         = "jerseyNumbers"
	TileCareerStats           = "careerStats"
	TilePersonalAchievements  = "personalAchievements"
	TilePhoto                 = "photo"
	TileYearsActive           = "yearsActive"
)

// AllTiles returns a slice of all tile names
func AllTiles() []string {
	return []string{
		TileBio,
		TilePlayerInformation,
		TileDraftInformation,
		TileTeamsPlayedOn,
		TileJerseyNumbers,
		TileCareerStats,
		TilePersonalAchievements,
		TilePhoto,
		TileYearsActive,
	}
}
