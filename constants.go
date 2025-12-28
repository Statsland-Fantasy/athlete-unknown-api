package main

// Sport constants
const (
	SportBaseball   = "baseball"
	SportBasketball = "basketball"
	SportFootball   = "football"
)

// Permission constants
const (
	PermissionReadUserStats      = "read:athlete-unknown:user-stats"
	PermissionReadUpcomingRounds = "read:athlete-unknown:upcoming-rounds"
	PermissionReadRounds         = "read:athlete-unknown:rounds"
	PermissionReadRoundStats     = "read:athlete-unknown:round-stats"
	PermissionSubmitResults      = "submit:athlete-unknown:results"
)

// Context key constants
const (
	ConstantUserId      = "userId"
	ConstantPermissions = "permissions"
	ConstantRoles       = "roles"
	ConstantIsAdmin     = "isAdmin"
)

// Role constants
const (
	RolePlayer     = "Player"
	RolePlaytester = "Playtester"
	RoleAdmin      = "Admin"
)

// Auth0 custom claim namespace
const (
	Auth0ClaimRoles  = "https://statslandfantasy.com/roles"
	Auth0ClaimUserId = "https://statslandfantasy.com/user_id"
)

// Tile name constants
const (
	TileBio                  = "bio"
	TilePlayerInformation    = "playerInformation"
	TileDraftInformation     = "draftInformation"
	TileTeamsPlayedOn        = "teamsPlayedOn"
	TileJerseyNumbers        = "jerseyNumbers"
	TileCareerStats          = "careerStats"
	TilePersonalAchievements = "personalAchievements"
	TilePhoto                = "photo"
	TileYearsActive          = "yearsActive"
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

// AllSports returns a slice of all supported sports
func AllSports() []string {
	return []string{
		SportBaseball,
		SportBasketball,
		SportFootball,
	}
}

// Allowed domains for web scraping (whitelist)
var AllowedScrapingDomains = []string{
	"baseball-reference.com",
	"www.baseball-reference.com",
	"basketball-reference.com",
	"www.basketball-reference.com",
	"pro-football-reference.com",
	"www.pro-football-reference.com",
}

// IsValidSport checks if a sport is valid
func IsValidSport(sport string) bool {
	switch sport {
	case SportBaseball, SportBasketball, SportFootball:
		return true
	default:
		return false
	}
}
