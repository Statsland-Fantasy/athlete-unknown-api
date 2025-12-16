package main

import (
	"fmt"
	"os"
	"time"
)

// Config holds application configuration
type Config struct {
	DynamoDBEndpoint   string
	RoundsTableName    string
	UserStatsTableName string
	AWSRegion          string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		DynamoDBEndpoint:   getEnv("DYNAMODB_ENDPOINT", "http://localhost:8000"),
		RoundsTableName:    getEnv("ROUNDS_TABLE_NAME", "AthleteUnknownRoundsDev"),
		UserStatsTableName: getEnv("USER_STATS_TABLE_NAME", "AthleteUnknownUserStatsDev"),
		AWSRegion:          getEnv("AWS_REGION", "us-west-2"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetSportsReferenceHostname returns the hostname for the given sport
func GetSportsReferenceHostname(sport string) string {
	switch sport {
	case "baseball":
		return "baseball-reference.com"
	case "basketball":
		return "basketball-reference.com"
	case "football":
		return "pro-football-reference.com"
	default:
		return ""
	}
}

// GetCurrentSeasonYear returns the current season year for the given sport
func GetCurrentSeasonYear(sport string) int {
	switch sport {
	case "baseball":
		return 2025 // MLB season begins in March/April & ends in October
	case "basketball":
		return 2025 // NBA season begins in October & ends in June
	case "football":
		return 2025 // NFL season begins in September & ends in February
	default:
		return 0
	}
}

// FIRST_ROUND_DATE is the reference date for calculating round IDs
var FIRST_ROUND_DATE = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

// GenerateRoundID generates a round ID by concatenating the sport and the number of days since FIRST_ROUND_DATE
func GenerateRoundID(sport string, playDate string) (string, error) {
	// Parse the playDate
	date, err := time.Parse("2006-01-02", playDate)
	if err != nil {
		return "", fmt.Errorf("invalid playDate format: %w", err)
	}

	// Calculate the number of days since FIRST_ROUND_DATE
	daysSince := int(date.Sub(FIRST_ROUND_DATE).Hours() / 24)

	// Generate the round ID
	roundID := fmt.Sprintf("%s%d", sport, daysSince)
	return roundID, nil
}
