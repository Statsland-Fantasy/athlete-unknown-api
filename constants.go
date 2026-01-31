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
	PermissionMigrateUserStats   = "migrate:athlete-unknown:user-stats"
)

// Context key constants
const (
	ConstantUserId      = "userId"
	ConstantUsername    = "username"
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
	TileInitials             = "initials"
	TileNicknames            = "nicknames"
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
		TileInitials,
		TileNicknames,
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

// HTTP Status reason phrases
const (
	StatusBadRequest          = "Bad Request"
	StatusInternalServerError = "Internal Server Error"
	StatusNotFound            = "Not Found"
	StatusConflict            = "Conflict"
)

// Error codes
const (
	ErrorMissingRequiredParameter = "MISSING_REQUIRED_PARAMETER"
	ErrorInvalidParameter         = "INVALID_PARAMETER"
	ErrorDatabaseError            = "DATABASE_ERROR"
	ErrorRoundNotFound            = "ROUND_NOT_FOUND"
	ErrorInvalidRequestBody       = "INVALID_REQUEST_BODY"
	ErrorMissingRequiredField     = "MISSING_REQUIRED_FIELD"
	ErrorInvalidPlayDate          = "INVALID_PLAY_DATE"
	ErrorRoundAlreadyExists       = "ROUND_ALREADY_EXISTS"
	ErrorNoUpcomingRounds         = "NO_UPCOMING_ROUNDS"
	ErrorStatsNotFound            = "STATS_NOT_FOUND"
	ErrorUserStatsNotFound        = "USER_STATS_NOT_FOUND"
	ErrorUserAlreadyMigrated      = "USER_ALREADY_MIGRATED"
	ErrorConfigurationError       = "CONFIGURATION_ERROR"
	ErrorInvalidURL               = "INVALID_URL"
	ErrorScrapingError            = "SCRAPING_ERROR"
	ErrorNoPlayersFound           = "NO_PLAYERS_FOUND"
	ErrorMultiplePlayersFound     = "MULTIPLE_PLAYERS_FOUND"
	ErrorInvalidSearchResultURL   = "INVALID_SEARCH_RESULT_URL"
)

// Date format constants
const (
	DateFormatYYYYMMDD = "2006-01-02"
)

// Query parameter names
const (
	QueryParamSport              = "sport"
	QueryParamPlayDate           = "playDate"
	QueryParamStartDate          = "startDate"
	QueryParamEndDate            = "endDate"
	QueryParamUserId             = "userId"
	QueryParamName               = "name"
	QueryParamSportsReferenceURL = "sportsReferenceURL"
	QueryParamTheme              = "theme"
)

// JSON response field names
const (
	JSONFieldError     = "error"
	JSONFieldMessage   = "message"
	JSONFieldCode      = "code"
	JSONFieldTimestamp = "timestamp"
	JSONFieldDetails   = "details"
)

// IsValidSport checks if a sport is valid
func IsValidSport(sport string) bool {
	switch sport {
	case SportBaseball, SportBasketball, SportFootball:
		return true
	default:
		return false
	}
}
