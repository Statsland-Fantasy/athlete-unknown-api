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
