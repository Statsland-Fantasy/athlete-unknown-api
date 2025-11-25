package main

import "time"

// Player represents a player entity with comprehensive details
type Player struct {
	Sport                 string `json:"sport"`
	SportsReferenceURL    string `json:"sportsReferenceURL"`
	Name                  string `json:"name"`
	Bio                   string `json:"bio"`
	PlayerInformation     string `json:"playerInformation"`
	DraftInformation      string `json:"draftInformation"`
	YearsActive           string `json:"yearsActive"`
	TeamsPlayedOn         string `json:"teamsPlayedOn"`
	JerseyNumbers         string `json:"jerseyNumbers"`
	CareerStats           string `json:"careerStats"`
	PersonalAchievements  string `json:"personalAchievements"`
	Photo                 string `json:"photo"`
}

// TileFlipTracker tracks tile flip counts
type TileFlipTracker struct {
	Bio                  int `json:"bio"`
	PlayerInformation    int `json:"playerInformation"`
	DraftInformation     int `json:"draftInformation"`
	TeamsPlayedOn        int `json:"teamsPlayedOn"`
	JerseyNumbers        int `json:"jerseyNumbers"`
	CareerStats          int `json:"careerStats"`
	PersonalAchievements int `json:"personalAchievements"`
	Photo                int `json:"photo"`
	YearsActive          int `json:"yearsActive"`
}

// RoundStats represents statistics for a specific round
type RoundStats struct {
	PlayDate                    string          `json:"playDate"`
	Name                        string          `json:"name"`
	Sport                       string          `json:"sport"`
	TotalPlays                  int             `json:"totalPlays"`
	PercentageCorrect           float64         `json:"percentageCorrect"`
	HighestScore                int             `json:"highestScore"`
	AverageCorrectScore         float64         `json:"averageCorrectScore"`
	MostCommonFirstTileFlipped  string          `json:"mostCommonFirstTileFlipped"`
	MostCommonLastTileFlipped   string          `json:"mostCommonLastTileFlipped"`
	MostCommonTileFlipped       string          `json:"mostCommonTileFlipped"`
	LeastCommonTileFlipped      string          `json:"leastCommonTileFlipped"`
	FirstTileFlippedTracker     TileFlipTracker `json:"firstTileFlippedTracker,omitempty"`
	LastTileFlippedTracker      TileFlipTracker `json:"lastTileFlippedTracker,omitempty"`
	MostTileFlippedTracker      TileFlipTracker `json:"mostTileFlippedTracker,omitempty"`
}

// Round represents a complete game round
type Round struct {
	RoundID               string     `json:"roundId"`
	Sport                 string     `json:"sport"`
	PlayDate              string     `json:"playDate"`
	Created               time.Time  `json:"created"`
	LastUpdated           time.Time  `json:"lastUpdated"`
	PreviouslyPlayedDates []string   `json:"previouslyPlayedDates"`
	Player                Player     `json:"player"`
	Stats                 RoundStats `json:"stats"`
}

// Result represents a game result submission
type Result struct {
	Score        int      `json:"score"`
	IsCorrect    bool     `json:"isCorrect"`
	TilesFlipped []string `json:"tilesFlipped"`
}

// SportStats represents statistics for a specific sport for a user
type SportStats struct {
	Sport                       string          `json:"sport"`
	CurrentDailyStreak          int             `json:"currentDailyStreak"`
	TotalPlays                  int             `json:"totalPlays"`
	PercentageCorrect           float64         `json:"percentageCorrect"`
	HighestScore                int             `json:"highestScore"`
	AverageCorrectScore         float64         `json:"averageCorrectScore"`
	MostCommonFirstTileFlipped  string          `json:"mostCommonFirstTileFlipped"`
	MostCommonLastTileFlipped   string          `json:"mostCommonLastTileFlipped"`
	MostCommonTileFlipped       string          `json:"mostCommonTileFlipped"`
	LeastCommonTileFlipped      string          `json:"leastCommonTileFlipped"`
	FirstTileFlippedTracker     TileFlipTracker `json:"firstTileFlippedTracker,omitempty"`
	LastTileFlippedTracker      TileFlipTracker `json:"lastTileFlippedTracker,omitempty"`
	MostTileFlippedTracker      TileFlipTracker `json:"mostTileFlippedTracker,omitempty"`
}

// UserStats represents comprehensive statistics for a user
type UserStats struct {
	UserID                      string          `json:"userId"`
	UserCreated                 time.Time       `json:"userCreated"`
	Sports                      []SportStats    `json:"sports"`
	TotalPlays                  int             `json:"totalPlays"`
	PercentageCorrect           float64         `json:"percentageCorrect"`
	HighestScore                int             `json:"highestScore"`
	AverageCorrectScore         float64         `json:"averageCorrectScore"`
	MostCommonFirstTileFlipped  string          `json:"mostCommonFirstTileFlipped"`
	MostCommonLastTileFlipped   string          `json:"mostCommonLastTileFlipped"`
	MostCommonTileFlipped       string          `json:"mostCommonTileFlipped"`
	LeastCommonTileFlipped      string          `json:"leastCommonTileFlipped"`
	FirstTileFlippedTracker     TileFlipTracker `json:"firstTileFlippedTracker,omitempty"`
	LastTileFlippedTracker      TileFlipTracker `json:"lastTileFlippedTracker,omitempty"`
	MostTileFlippedTracker      TileFlipTracker `json:"mostTileFlippedTracker,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Code      string                 `json:"code"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}
