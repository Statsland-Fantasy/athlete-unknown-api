package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		expected     string
	}{
		{
			name:         "environment variable not set",
			key:          "TEST_VAR_NOT_SET",
			defaultValue: "default",
			setEnv:       false,
			expected:     "default",
		},
		{
			name:         "environment variable set",
			key:          "TEST_VAR_SET",
			defaultValue: "default",
			envValue:     "custom",
			setEnv:       true,
			expected:     "custom",
		},
		{
			name:         "environment variable set to empty string",
			key:          "TEST_VAR_EMPTY",
			defaultValue: "default",
			envValue:     "",
			setEnv:       true,
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment variable after test
			defer os.Unsetenv(tt.key)

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.expected {
				t.Errorf("getEnv() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		setupEnv       func()
		cleanupEnv     func()
		expectedConfig *Config
	}{
		{
			name: "default configuration",
			setupEnv: func() {
				// Clear all config-related environment variables
				os.Unsetenv("DYNAMODB_ENDPOINT")
				os.Unsetenv("ROUNDS_TABLE_NAME")
				os.Unsetenv("USER_STATS_TABLE_NAME")
				os.Unsetenv("AWS_REGION")
			},
			cleanupEnv: func() {},
			expectedConfig: &Config{
				DynamoDBEndpoint:   "",
				RoundsTableName:    "AthleteUnknownRoundsDev",
				UserStatsTableName: "AthleteUnknownUserStatsDev",
				AWSRegion:          "us-west-2",
			},
		},
		{
			name: "custom configuration from environment",
			setupEnv: func() {
				os.Setenv("DYNAMODB_ENDPOINT", "http://custom:9000")
				os.Setenv("ROUNDS_TABLE_NAME", "CustomRoundsTable")
				os.Setenv("USER_STATS_TABLE_NAME", "CustomUserStatsTable")
				os.Setenv("AWS_REGION", "us-east-1")
			},
			cleanupEnv: func() {
				os.Unsetenv("DYNAMODB_ENDPOINT")
				os.Unsetenv("ROUNDS_TABLE_NAME")
				os.Unsetenv("USER_STATS_TABLE_NAME")
				os.Unsetenv("AWS_REGION")
			},
			expectedConfig: &Config{
				DynamoDBEndpoint:   "http://custom:9000",
				RoundsTableName:    "CustomRoundsTable",
				UserStatsTableName: "CustomUserStatsTable",
				AWSRegion:          "us-east-1",
			},
		},
		{
			name: "partial custom configuration",
			setupEnv: func() {
				os.Unsetenv("DYNAMODB_ENDPOINT")
				os.Setenv("ROUNDS_TABLE_NAME", "CustomRoundsOnly")
				os.Unsetenv("USER_STATS_TABLE_NAME")
				os.Setenv("AWS_REGION", "eu-west-1")
			},
			cleanupEnv: func() {
				os.Unsetenv("ROUNDS_TABLE_NAME")
				os.Unsetenv("AWS_REGION")
			},
			expectedConfig: &Config{
				DynamoDBEndpoint:   "",
				RoundsTableName:    "CustomRoundsOnly",
				UserStatsTableName: "AthleteUnknownUserStatsDev",
				AWSRegion:          "eu-west-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			tt.setupEnv()
			defer tt.cleanupEnv()

			// Load config
			cfg := LoadConfig()

			// Verify config fields
			if cfg.DynamoDBEndpoint != tt.expectedConfig.DynamoDBEndpoint {
				t.Errorf("DynamoDBEndpoint = %v, want %v", cfg.DynamoDBEndpoint, tt.expectedConfig.DynamoDBEndpoint)
			}
			if cfg.RoundsTableName != tt.expectedConfig.RoundsTableName {
				t.Errorf("RoundsTableName = %v, want %v", cfg.RoundsTableName, tt.expectedConfig.RoundsTableName)
			}
			if cfg.UserStatsTableName != tt.expectedConfig.UserStatsTableName {
				t.Errorf("UserStatsTableName = %v, want %v", cfg.UserStatsTableName, tt.expectedConfig.UserStatsTableName)
			}
			if cfg.AWSRegion != tt.expectedConfig.AWSRegion {
				t.Errorf("AWSRegion = %v, want %v", cfg.AWSRegion, tt.expectedConfig.AWSRegion)
			}
		})
	}
}

func TestConfigStruct(t *testing.T) {
	// Test that Config struct can be created and fields accessed
	cfg := &Config{
		DynamoDBEndpoint:   "http://test:8000",
		RoundsTableName:    "TestRoundsTable",
		UserStatsTableName: "TestUserStatsTable",
		AWSRegion:          "test-region",
	}

	if cfg.DynamoDBEndpoint != "http://test:8000" {
		t.Errorf("DynamoDBEndpoint = %v, want %v", cfg.DynamoDBEndpoint, "http://test:8000")
	}
	if cfg.RoundsTableName != "TestRoundsTable" {
		t.Errorf("RoundsTableName = %v, want %v", cfg.RoundsTableName, "TestRoundsTable")
	}
	if cfg.UserStatsTableName != "TestUserStatsTable" {
		t.Errorf("UserStatsTableName = %v, want %v", cfg.UserStatsTableName, "TestUserStatsTable")
	}
	if cfg.AWSRegion != "test-region" {
		t.Errorf("AWSRegion = %v, want %v", cfg.AWSRegion, "test-region")
	}
}

func TestGetSportsReferenceHostname(t *testing.T) {
	tests := []struct {
		name     string
		sport    string
		expected string
	}{
		{
			name:     "baseball hostname",
			sport:    "baseball",
			expected: "baseball-reference.com",
		},
		{
			name:     "basketball hostname",
			sport:    "basketball",
			expected: "basketball-reference.com",
		},
		{
			name:     "football hostname",
			sport:    "football",
			expected: "pro-football-reference.com",
		},
		{
			name:     "unknown sport",
			sport:    "soccer",
			expected: "",
		},
		{
			name:     "empty sport",
			sport:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSportsReferenceHostname(tt.sport)
			if got != tt.expected {
				t.Errorf("GetSportsReferenceHostname(%q) = %v, want %v", tt.sport, got, tt.expected)
			}
		})
	}
}

func TestGetCurrentSeasonYear(t *testing.T) {
	tests := []struct {
		name     string
		sport    string
		expected int
	}{
		{
			name:     "baseball season year",
			sport:    "baseball",
			expected: 2025,
		},
		{
			name:     "basketball season year",
			sport:    "basketball",
			expected: 2025,
		},
		{
			name:     "football season year",
			sport:    "football",
			expected: 2025,
		},
		{
			name:     "unknown sport",
			sport:    "soccer",
			expected: 0,
		},
		{
			name:     "empty sport",
			sport:    "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCurrentSeasonYear(tt.sport)
			if got != tt.expected {
				t.Errorf("GetCurrentSeasonYear(%q) = %v, want %v", tt.sport, got, tt.expected)
			}
		})
	}
}

func TestGenerateRoundID(t *testing.T) {
	tests := []struct {
		name      string
		sport     string
		playDate  string
		expected  string
		expectErr bool
	}{
		{
			name:      "baseball on first round date",
			sport:     "baseball",
			playDate:  "2025-01-01",
			expected:  "baseball#0",
			expectErr: false,
		},
		{
			name:      "basketball one day after first round date",
			sport:     "basketball",
			playDate:  "2025-01-02",
			expected:  "basketball#1",
			expectErr: false,
		},
		{
			name:      "football 30 days after first round date",
			sport:     "football",
			playDate:  "2025-01-31",
			expected:  "football#30",
			expectErr: false,
		},
		{
			name:      "baseball 365 days after first round date",
			sport:     "baseball",
			playDate:  "2026-01-01",
			expected:  "baseball#365",
			expectErr: false,
		},
		{
			name:      "date before first round date",
			sport:     "basketball",
			playDate:  "2024-12-31",
			expected:  "basketball#-1",
			expectErr: false,
		},
		{
			name:      "invalid date format",
			sport:     "baseball",
			playDate:  "2025/01/01",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "invalid date",
			sport:     "football",
			playDate:  "not-a-date",
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRoundID(tt.sport, tt.playDate)
			if tt.expectErr {
				if err == nil {
					t.Errorf("GenerateRoundID(%q, %q) expected error, got nil", tt.sport, tt.playDate)
				}
			} else {
				if err != nil {
					t.Errorf("GenerateRoundID(%q, %q) unexpected error: %v", tt.sport, tt.playDate, err)
				}
				if got != tt.expected {
					t.Errorf("GenerateRoundID(%q, %q) = %v, want %v", tt.sport, tt.playDate, got, tt.expected)
				}
			}
		})
	}
}
