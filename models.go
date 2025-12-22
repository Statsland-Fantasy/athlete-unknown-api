package main

import "time"

// Player represents a player entity with comprehensive details
type Player struct {
	Sport                 string `json:"sport" dynamodbav:"sport"`
	SportsReferenceURL    string `json:"sportsReferenceURL" dynamodbav:"sportsReferenceURL"`
	Name                  string `json:"name" dynamodbav:"name"`
	Bio                   string `json:"bio" dynamodbav:"bio"`
	PlayerInformation     string `json:"playerInformation" dynamodbav:"playerInformation"`
	DraftInformation      string `json:"draftInformation" dynamodbav:"draftInformation"`
	YearsActive           string `json:"yearsActive" dynamodbav:"yearsActive"`
	TeamsPlayedOn         string `json:"teamsPlayedOn" dynamodbav:"teamsPlayedOn"`
	JerseyNumbers         string `json:"jerseyNumbers" dynamodbav:"jerseyNumbers"`
	CareerStats           string `json:"careerStats" dynamodbav:"careerStats"`
	PersonalAchievements  string `json:"personalAchievements" dynamodbav:"personalAchievements"`
	Photo                 string `json:"photo" dynamodbav:"photo"`
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
}

// Stats represents common statistics fields shared across different stat types
type Stats struct {
	TotalPlays                  int             `json:"totalPlays" dynamodbav:"totalPlays"`
	PercentageCorrect           float64         `json:"percentageCorrect" dynamodbav:"percentageCorrect"`
	HighestScore                int             `json:"highestScore" dynamodbav:"highestScore"`
	AverageCorrectScore         float64         `json:"averageCorrectScore" dynamodbav:"averageCorrectScore"`
	AverageNumberOfTileFlips    float64         `json:"averageNumberOfTileFlips" dynamodbav:"averageNumberOfTileFlips"`
	MostCommonFirstTileFlipped  string          `json:"mostCommonFirstTileFlipped" dynamodbav:"mostCommonFirstTileFlipped"`
	MostCommonLastTileFlipped   string          `json:"mostCommonLastTileFlipped" dynamodbav:"mostCommonLastTileFlipped"`
	MostCommonTileFlipped       string          `json:"mostCommonTileFlipped" dynamodbav:"mostCommonTileFlipped"`
	LeastCommonTileFlipped      string          `json:"leastCommonTileFlipped" dynamodbav:"leastCommonTileFlipped"`
	FirstTileFlippedTracker     TileFlipTracker `json:"firstTileFlippedTracker" dynamodbav:"firstTileFlippedTracker"`
	LastTileFlippedTracker      TileFlipTracker `json:"lastTileFlippedTracker" dynamodbav:"lastTileFlippedTracker"`
	MostTileFlippedTracker      TileFlipTracker `json:"mostTileFlippedTracker" dynamodbav:"mostTileFlippedTracker"`
}

// RoundStats represents statistics for a specific round
type RoundStats struct {
	PlayDate string `json:"playDate" dynamodbav:"playDate"`
	Name     string `json:"name" dynamodbav:"name"`
	Sport    string `json:"sport" dynamodbav:"sport"`
	Stats           `json:",inline" dynamodbav:",inline"`
}

// Round represents a complete game round
type Round struct {
	RoundID               string     `json:"roundId" dynamodbav:"roundId"`
	Sport                 string     `json:"sport" dynamodbav:"sport"`
	PlayDate              string     `json:"playDate" dynamodbav:"playDate"`
	Created               time.Time  `json:"created" dynamodbav:"created"`
	LastUpdated           time.Time  `json:"lastUpdated" dynamodbav:"lastUpdated"`
	Theme				  string      `json:"theme" dynamodbav:"theme"`
	Player                Player     `json:"player" dynamodbav:"player"`
	Stats                 RoundStats `json:"stats" dynamodbav:"stats"`
}

// Result represents a game result submission
type Result struct {
	Score        int      `json:"score"`
	IsCorrect    bool     `json:"isCorrect"`
	TilesFlipped []string `json:"tilesFlipped"`
}

// SportStats represents statistics for a specific sport for all users
type SportStats struct {
	Sport              string `json:"sport" dynamodbav:"sport"`
	Stats              `json:",inline" dynamodbav:",inline"`
}

// UserStats represents comprehensive statistics for a user
type UserStats struct {
	UserId                      string          `json:"userId" dynamodbav:"userId"`
	UserName                    string          `json:"userName" dynamodbav:"userName"`
	UserCreated                 time.Time       `json:"userCreated" dynamodbav:"userCreated"`
	CurrentDailyStreak 			int    			`json:"currentDailyStreak" dynamodbav:"currentDailyStreak"`
	LastDayPlayed               string          `json:"lastDayPlayed" dynamodbav:"lastDayPlayed"`
	Sports                      []SportStats    `json:"sports" dynamodbav:"sports"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Code      string                 `json:"code"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}
