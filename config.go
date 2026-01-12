package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file before any package variables are initialized
	_ = godotenv.Load()

	// Initialize date variables after .env is loaded
	FIRST_ROUND_DATE_STRING = getEnv("FIRST_ROUND_DATE", "2026-02-08")
	FIRST_ROUND_DATE = mustParseDate(FIRST_ROUND_DATE_STRING)
}

// Config holds application configuration
type Config struct {
	DynamoDBEndpoint   string
	RoundsTableName    string
	UserStatsTableName string
	AWSRegion          string
	AIUpscalerEnabled  bool
	AIUpscalerAPIKey   string
	AIUpscalerAPIURL   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		DynamoDBEndpoint:   getEnv("DYNAMODB_ENDPOINT", ""),
		RoundsTableName:    getEnv("ROUNDS_TABLE_NAME", "AthleteUnknownRoundsDev"),
		UserStatsTableName: getEnv("USER_STATS_TABLE_NAME", "AthleteUnknownUserStatsDev"),
		AWSRegion:          getEnv("AWS_REGION", "us-west-2"),
		AIUpscalerEnabled:  getEnvBool("AI_UPSCALER_ENABLED", false),
		AIUpscalerAPIKey:   getEnv("AI_UPSCALER_API_KEY", ""),
		AIUpscalerAPIURL:   getEnv("AI_UPSCALER_API_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}

// GetSportsReferenceHostname returns the hostname for the given sport
func GetSportsReferenceHostname(sport string) string {
	switch sport {
	case SportBaseball:
		return "baseball-reference.com"
	case SportBasketball:
		return "basketball-reference.com"
	case SportFootball:
		return "pro-football-reference.com"
	default:
		return ""
	}
}

// GetCurrentSeasonYear returns the current season year for the given sport
func GetCurrentSeasonYear(sport string) int {
	switch sport {
	case SportBaseball:
		return 2025 // MLB season begins in March/April & ends in October
	case SportBasketball:
		return 2025 // NBA season begins in October & ends in June
	case SportFootball:
		return 2025 // NFL season begins in September & ends in February
	default:
		return 0
	}
}

// FIRST_ROUND_DATE_STRING is the string representation (YYYY-MM-DD) from env or default
var FIRST_ROUND_DATE_STRING string

// FIRST_ROUND_DATE is the parsed time.Time version for calculations
var FIRST_ROUND_DATE time.Time

// mustParseDate parses a date string or panics (should only be called at startup)
func mustParseDate(dateStr string) time.Time {
	date, err := time.Parse(DateFormatYYYYMMDD, dateStr)
	if err != nil {
		panic(fmt.Sprintf("invalid FIRST_ROUND_DATE format '%s': %v", dateStr, err))
	}
	return date
}

// GenerateRoundID generates a round ID by concatenating the sport and the number of days since FIRST_ROUND_DATE
func GenerateRoundID(sport string, playDate string) (string, error) {
	// Parse the playDate
	date, err := time.Parse(DateFormatYYYYMMDD, playDate)
	if err != nil {
		return "", fmt.Errorf("invalid playDate format: %w", err)
	}

	// Calculate the number of days since FIRST_ROUND_DATE
	daysSince := int(date.Sub(FIRST_ROUND_DATE).Hours() / 24)

	// Generate the round ID. Split sport and round number by "#"
	roundID := fmt.Sprintf("%s#%d", sport, daysSince)
	return roundID, nil
}

// GetAllowedCORSOrigins returns the list of allowed CORS origins from environment or defaults
func GetAllowedCORSOrigins() []string {
	// Get from environment variable (comma-separated list)
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv != "" {
		origins := strings.Split(originsEnv, ",")
		// Trim whitespace from each origin
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
		return origins
	}

	// Default origins for development
	// In production, ALLOWED_ORIGINS environment variable should be set
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		// In production with no ALLOWED_ORIGINS set, return empty slice
		// This will effectively block all CORS requests, which is safer than wildcard
		return []string{}
	}

	// Development defaults
	return []string{
		"http://localhost:3000", // React default
		"http://localhost:5173", // Vite default
		"http://localhost:4200", // Angular default
		"http://localhost:8080", // Various frameworks
		"http://127.0.0.1:3000", // Localhost IP variant
		"http://127.0.0.1:5173", // Localhost IP variant
	}
}
