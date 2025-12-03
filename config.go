package main

import (
	"os"
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
