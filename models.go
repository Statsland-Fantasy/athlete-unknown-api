package main

import "time"

// Round represents a complete game round
type Round struct {
	RoundID     string     `json:"roundId" dynamodbav:"roundId"`
	Sport       string     `json:"sport" dynamodbav:"sport"`
	PlayDate    string     `json:"playDate" dynamodbav:"playDate"`
	Created     time.Time  `json:"created" dynamodbav:"created"`
	LastUpdated time.Time  `json:"lastUpdated" dynamodbav:"lastUpdated"`
	Title       string     `json:"title" dynamodbav:"title"`
	Player      Player     `json:"player" dynamodbav:"player"`
	Stats       RoundStats `json:"stats" dynamodbav:"stats"`
}

// RoundSummary contains minimal round information for list views
type RoundSummary struct {
	RoundID  string `json:"roundId" dynamodbav:"roundId"`
	Sport    string `json:"sport" dynamodbav:"sport"`
	PlayDate string `json:"playDate" dynamodbav:"playDate"`
}

// Player represents a player entity with comprehensive details
type Player struct {
	Sport                string `json:"sport" dynamodbav:"sport"`
	SportsReferenceURL   string `json:"sportsReferenceURL" dynamodbav:"sportsReferenceURL"`
	Name                 string `json:"name" dynamodbav:"name"`
	Bio                  string `json:"bio" dynamodbav:"bio"`
	PlayerInformation    string `json:"playerInformation" dynamodbav:"playerInformation"`
	DraftInformation     string `json:"draftInformation" dynamodbav:"draftInformation"`
	YearsActive          string `json:"yearsActive" dynamodbav:"yearsActive"`
	TeamsPlayedOn        string `json:"teamsPlayedOn" dynamodbav:"teamsPlayedOn"`
	JerseyNumbers        string `json:"jerseyNumbers" dynamodbav:"jerseyNumbers"`
	CareerStats          string `json:"careerStats" dynamodbav:"careerStats"`
	PersonalAchievements string `json:"personalAchievements" dynamodbav:"personalAchievements"`
	Photo                string `json:"photo" dynamodbav:"photo"`
	Initials             string `json:"initials" dynamodbav:"initials"`
	Nicknames            string `json:"nicknames" dynamodbav:"nicknames"`
}

// Stats represents common statistics fields shared across different stat types
type Stats struct {
	TotalPlays                 int             `json:"totalPlays" dynamodbav:"totalPlays"`
	PercentageCorrect          float64         `json:"percentageCorrect" dynamodbav:"percentageCorrect"`
	HighestScore               int             `json:"highestScore" dynamodbav:"highestScore"`
	AverageCorrectScore        float64         `json:"averageCorrectScore" dynamodbav:"averageCorrectScore"`
	AverageIncorrectGuesses    float64         `json:"averageIncorrectGuesses" dynamodbav:"averageIncorrectGuesses"`
	AverageNumberOfTileFlips   float64         `json:"averageNumberOfTileFlips" dynamodbav:"averageNumberOfTileFlips"`
	MostCommonFirstTileFlipped string          `json:"mostCommonFirstTileFlipped" dynamodbav:"mostCommonFirstTileFlipped"`
	MostCommonLastTileFlipped  string          `json:"mostCommonLastTileFlipped" dynamodbav:"mostCommonLastTileFlipped"`
	MostCommonTileFlipped      string          `json:"mostCommonTileFlipped" dynamodbav:"mostCommonTileFlipped"`
	LeastCommonTileFlipped     string          `json:"leastCommonTileFlipped" dynamodbav:"leastCommonTileFlipped"`
	FirstTileFlippedTracker    TileFlipTracker `json:"firstTileFlippedTracker" dynamodbav:"firstTileFlippedTracker"`
	LastTileFlippedTracker     TileFlipTracker `json:"lastTileFlippedTracker" dynamodbav:"lastTileFlippedTracker"`
	MostTileFlippedTracker     TileFlipTracker `json:"mostTileFlippedTracker" dynamodbav:"mostTileFlippedTracker"`
}

// TileFlipTracker tracks tile flip counts
type TileFlipTracker struct {
	Bio                  int `json:"bio" dynamodbav:"bio"`
	PlayerInformation    int `json:"playerInformation" dynamodbav:"playerInformation"`
	DraftInformation     int `json:"draftInformation" dynamodbav:"draftInformation"`
	TeamsPlayedOn        int `json:"teamsPlayedOn" dynamodbav:"teamsPlayedOn"`
	JerseyNumbers        int `json:"jerseyNumbers" dynamodbav:"jerseyNumbers"`
	CareerStats          int `json:"careerStats" dynamodbav:"careerStats"`
	PersonalAchievements int `json:"personalAchievements" dynamodbav:"personalAchievements"`
	Photo                int `json:"photo" dynamodbav:"photo"`
	YearsActive          int `json:"yearsActive" dynamodbav:"yearsActive"`
	Initials             int `json:"initials" dynamodbav:"initials"`
	Nicknames            int `json:"nicknames" dynamodbav:"nicknames"`
}

// RoundStats represents statistics for a specific round
type RoundStats struct {
	PlayDate string `json:"playDate" dynamodbav:"playDate"`
	Name     string `json:"name" dynamodbav:"name"`
	Sport    string `json:"sport" dynamodbav:"sport"`
	Stats    `json:",inline" dynamodbav:",inline"`
}

// Result represents a game result submission
type Result struct {
	Score            int      `json:"score" dynamodbav:"score"`
	IsCorrect        bool     `json:"isCorrect" dynamodbav:"isCorrect"`
	FlippedTiles     []string `json:"flippedTiles" dynamodbav:"flippedTiles"`
	IncorrectGuesses int      `json:"incorrectGuesses" dynamodbav:"incorrectGuesses"`
}

// User represents comprehensive statistics for a user
type User struct {
	UserId             string           `json:"userId" dynamodbav:"userId"`
	UserName           string           `json:"userName" dynamodbav:"userName"`
	UserCreated        time.Time        `json:"userCreated" dynamodbav:"userCreated"`
	CurrentDailyStreak int              `json:"currentDailyStreak" dynamodbav:"currentDailyStreak"`
	LastDayPlayed      string           `json:"lastDayPlayed" dynamodbav:"lastDayPlayed"`
	Sports             []UserSportStats `json:"sports" dynamodbav:"sports"`
	StoryMissions      []StoryMission   `json:"storyMissions" dynamodbav:"storyMissions"`
}

// SportStats represents statistics for a specific sport for all users
type UserSportStats struct {
	Sport   string         `json:"sport" dynamodbav:"sport"`
	Stats   Stats          `json:"stats" dynamodbav:"stats"`
	History []RoundHistory `json:"history" dynamodbav:"history"`
}

// RoundHistory represents the results of past rounds played
type RoundHistory struct {
	PlayDate string `json:"playDate" dynamodbav:"playDate"`
	Result
}

// StoryTemplate represents the configured story saved in database and then resolved in API to be dislayed in the FE
type StoryTemplate struct {
	StoryId    string `json:"storyId" dynamodbav:"storyId"`
	Title      string `json:"title" dynamodbav:"title"`
	ModalType  string `json:"modalType" dynamodbav:"modalType"`
	DateOffset string `json:"dateOffset" dynamodbav:"dateOffset"`
	BodyText   string `json:"bodyText" dynamodbav:"bodyText"`
}

// StoryMission represents the status of each story mission. Stored in User object
type StoryMission struct {
	Criteria     string `json:"criteria"` // ex: Solve 1 Case
	Title        string `json:"title"`    // ex: A Smashing Debut
	DateAchieved string `json:"date"`
	PlayerName   string `json:"playerName"`
	StoryId      string `json:"storyId"`
}

// Story represents the response to FE with resolved text, date, names, etc for simple pass through display
type Story struct {
	StoryId   string `json:"storyId"`
	Title     string `json:"title"`
	ModalType string `json:"modalType"`
	Date      string `json:"date"`
	BodyText  string `json:"bodyText"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Code      string                 `json:"code"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}
