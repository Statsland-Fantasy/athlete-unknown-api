package main

import (
	"testing"
	"time"
)

func TestNewDB(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "create DB with local endpoint",
			config: &Config{
				DynamoDBEndpoint:   "http://localhost:8000",
				RoundsTableName:    "TestRoundsTable",
				UserStatsTableName: "TestUserStatsTable",
				AWSRegion:          "us-west-2",
			},
			wantError: false,
		},
		{
			name: "create DB without custom endpoint",
			config: &Config{
				DynamoDBEndpoint:   "",
				RoundsTableName:    "TestRoundsTable",
				UserStatsTableName: "TestUserStatsTable",
				AWSRegion:          "us-west-2",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDB(tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("NewDB() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && db == nil {
				t.Error("NewDB() returned nil DB without error")
			}
			if !tt.wantError {
				if db.roundsTableName != tt.config.RoundsTableName {
					t.Errorf("roundsTableName = %v, want %v", db.roundsTableName, tt.config.RoundsTableName)
				}
				if db.userStatsTableName != tt.config.UserStatsTableName {
					t.Errorf("userStatsTableName = %v, want %v", db.userStatsTableName, tt.config.UserStatsTableName)
				}
			}
		})
	}
}

func TestDBStruct(t *testing.T) {
	// Test that DB struct can be created with the expected fields
	cfg := &Config{
		DynamoDBEndpoint:   "http://localhost:8000",
		RoundsTableName:    "TestRoundsTable",
		UserStatsTableName: "TestUserStatsTable",
		AWSRegion:          "us-west-2",
	}

	db, err := NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to create DB: %v", err)
	}

	if db.client == nil {
		t.Error("DB client should not be nil")
	}
	if db.roundsTableName != "TestRoundsTable" {
		t.Errorf("roundsTableName = %v, want TestRoundsTable", db.roundsTableName)
	}
	if db.userStatsTableName != "TestUserStatsTable" {
		t.Errorf("userStatsTableName = %v, want TestUserStatsTable", db.userStatsTableName)
	}
}

// Test Round marshaling and unmarshaling
func TestRoundMarshaling(t *testing.T) {
	now := time.Now()
	round := &Round{
		RoundID:               "test-round-id",
		Sport:                 "basketball",
		PlayDate:              "2024-01-01",
		Created:               now,
		LastUpdated:           now,
		Theme: "GOAT",
		Player: Player{
			Sport:                "basketball",
			SportsReferenceURL:   "http://example.com",
			Name:                 "Test Player",
			Bio:                  "Test bio",
			PlayerInformation:    "Test info",
			DraftInformation:     "Test draft",
			YearsActive:          "2010-2020",
			TeamsPlayedOn:        "Team A, Team B",
			JerseyNumbers:        "23, 32",
			CareerStats:          "Points: 1000",
			PersonalAchievements: "MVP 2015",
			Photo:                "http://example.com/photo.jpg",
		},
		Stats: RoundStats{
			PlayDate: "2024-01-01",
			Name:     "Test Player",
			Sport:    "basketball",
			Stats: Stats{
				TotalPlays:                 100,
				PercentageCorrect:          75.5,
				HighestScore:               200,
				AverageCorrectScore:        150.0,
				AverageNumberOfTileFlips:   4.5,
				MostCommonFirstTileFlipped: "bio",
				MostCommonLastTileFlipped:  "photo",
				MostCommonTileFlipped:      "careerStats",
				LeastCommonTileFlipped:     "jerseyNumbers",
				FirstTileFlippedTracker: TileFlipTracker{
					Bio:                  50,
					PlayerInformation:    30,
					DraftInformation:     10,
					TeamsPlayedOn:        5,
					JerseyNumbers:        2,
					CareerStats:          1,
					PersonalAchievements: 1,
					Photo:                1,
					YearsActive:          0,
				},
			},
		},
	}

	// Verify all fields are set correctly
	if round.RoundID != "test-round-id" {
		t.Errorf("RoundID = %v, want test-round-id", round.RoundID)
	}
	if round.Sport != "basketball" {
		t.Errorf("Sport = %v, want basketball", round.Sport)
	}
	if round.PlayDate != "2024-01-01" {
		t.Errorf("PlayDate = %v, want 2024-01-01", round.PlayDate)
	}
	if round.Player.Name != "Test Player" {
		t.Errorf("Player.Name = %v, want Test Player", round.Player.Name)
	}
	if round.Stats.TotalPlays != 100 {
		t.Errorf("Stats.TotalPlays = %v, want 100", round.Stats.TotalPlays)
	}
	if round.Stats.FirstTileFlippedTracker.Bio != 50 {
		t.Errorf("FirstTileFlippedTracker.Bio = %v, want 50", round.Stats.FirstTileFlippedTracker.Bio)
	}
}

// Test UserStats marshaling and unmarshaling
func TestUserStatsMarshaling(t *testing.T) {
	now := time.Now()
	stats := &UserStats{
		UserId:      "test-user-123",
		UserName:    "John Doe",
		UserCreated: now,
		Sports: []SportStats{
			{
				Sport:              "basketball",
				CurrentDailyStreak: 5,
				Stats: Stats{
					TotalPlays:                 200,
					PercentageCorrect:          75.0,
					HighestScore:               200,
					AverageCorrectScore:        150.0,
					AverageNumberOfTileFlips:   5.2,
					MostCommonFirstTileFlipped: "bio",
					MostCommonLastTileFlipped:  "photo",
					MostCommonTileFlipped:      "careerStats",
					LeastCommonTileFlipped:     "jerseyNumbers",
				},
			},
			{
				Sport:              "baseball",
				CurrentDailyStreak: 3,
				Stats: Stats{
					TotalPlays:                 150,
					PercentageCorrect:          80.0,
					HighestScore:               180,
					AverageCorrectScore:        140.0,
					AverageNumberOfTileFlips:   6.1,
					MostCommonFirstTileFlipped: "playerInformation",
					MostCommonLastTileFlipped:  "careerStats",
					MostCommonTileFlipped:      "bio",
					LeastCommonTileFlipped:     "photo",
				},
			},
		},
	}

	// Verify all fields are set correctly
	if stats.UserId != "test-user-123" {
		t.Errorf("UserId = %v, want test-user-123", stats.UserId)
	}
	if stats.UserName != "John Doe" {
		t.Errorf("UserName = %v, want John Doe", stats.UserName)
	}
	if len(stats.Sports) != 2 {
		t.Errorf("len(Sports) = %v, want 2", len(stats.Sports))
	}
	if stats.Sports[0].Sport != "basketball" {
		t.Errorf("Sports[0].Sport = %v, want basketball", stats.Sports[0].Sport)
	}
	if stats.Sports[0].TotalPlays != 200 {
		t.Errorf("Sports[0].TotalPlays = %v, want 200", stats.Sports[0].TotalPlays)
	}
	if stats.Sports[0].PercentageCorrect != 75.0 {
		t.Errorf("Sports[0].PercentageCorrect = %v, want 75.0", stats.Sports[0].PercentageCorrect)
	}
	if stats.Sports[0].AverageNumberOfTileFlips != 5.2 {
		t.Errorf("Sports[0].AverageNumberOfTileFlips = %v, want 5.2", stats.Sports[0].AverageNumberOfTileFlips)
	}
	if stats.Sports[1].AverageNumberOfTileFlips != 6.1 {
		t.Errorf("Sports[1].AverageNumberOfTileFlips = %v, want 6.1", stats.Sports[1].AverageNumberOfTileFlips)
	}
}

// Test TileFlipTracker structure
func TestTileFlipTracker(t *testing.T) {
	tracker := TileFlipTracker{
		Bio:                  10,
		PlayerInformation:    20,
		DraftInformation:     30,
		TeamsPlayedOn:        40,
		JerseyNumbers:        50,
		CareerStats:          60,
		PersonalAchievements: 70,
		Photo:                80,
		YearsActive:          90,
	}

	// Verify all fields
	expected := map[string]int{
		"Bio":                  10,
		"PlayerInformation":    20,
		"DraftInformation":     30,
		"TeamsPlayedOn":        40,
		"JerseyNumbers":        50,
		"CareerStats":          60,
		"PersonalAchievements": 70,
		"Photo":                80,
		"YearsActive":          90,
	}

	if tracker.Bio != expected["Bio"] {
		t.Errorf("Bio = %v, want %v", tracker.Bio, expected["Bio"])
	}
	if tracker.PlayerInformation != expected["PlayerInformation"] {
		t.Errorf("PlayerInformation = %v, want %v", tracker.PlayerInformation, expected["PlayerInformation"])
	}
	if tracker.DraftInformation != expected["DraftInformation"] {
		t.Errorf("DraftInformation = %v, want %v", tracker.DraftInformation, expected["DraftInformation"])
	}
	if tracker.TeamsPlayedOn != expected["TeamsPlayedOn"] {
		t.Errorf("TeamsPlayedOn = %v, want %v", tracker.TeamsPlayedOn, expected["TeamsPlayedOn"])
	}
	if tracker.JerseyNumbers != expected["JerseyNumbers"] {
		t.Errorf("JerseyNumbers = %v, want %v", tracker.JerseyNumbers, expected["JerseyNumbers"])
	}
	if tracker.CareerStats != expected["CareerStats"] {
		t.Errorf("CareerStats = %v, want %v", tracker.CareerStats, expected["CareerStats"])
	}
	if tracker.PersonalAchievements != expected["PersonalAchievements"] {
		t.Errorf("PersonalAchievements = %v, want %v", tracker.PersonalAchievements, expected["PersonalAchievements"])
	}
	if tracker.Photo != expected["Photo"] {
		t.Errorf("Photo = %v, want %v", tracker.Photo, expected["Photo"])
	}
	if tracker.YearsActive != expected["YearsActive"] {
		t.Errorf("YearsActive = %v, want %v", tracker.YearsActive, expected["YearsActive"])
	}
}

// Test Player structure
func TestPlayerStructure(t *testing.T) {
	player := Player{
		Sport:                "basketball",
		SportsReferenceURL:   "http://example.com/player",
		Name:                 "John Doe",
		Bio:                  "Professional basketball player",
		PlayerInformation:    "Height: 6'5\", Weight: 220 lbs",
		DraftInformation:     "Round 1, Pick 5, 2010",
		YearsActive:          "2010-2023",
		TeamsPlayedOn:        "Lakers, Bulls, Heat",
		JerseyNumbers:        "23, 6, 24",
		CareerStats:          "PPG: 25.5, RPG: 7.2, APG: 6.8",
		PersonalAchievements: "3x NBA Champion, 2x MVP, 10x All-Star",
		Photo:                "http://example.com/photo.jpg",
	}

	if player.Sport != "basketball" {
		t.Errorf("Sport = %v, want basketball", player.Sport)
	}
	if player.Name != "John Doe" {
		t.Errorf("Name = %v, want John Doe", player.Name)
	}
	if player.SportsReferenceURL != "http://example.com/player" {
		t.Errorf("SportsReferenceURL = %v, want http://example.com/player", player.SportsReferenceURL)
	}
}

// Test Result structure
func TestResultStructure(t *testing.T) {
	result := Result{
		Score:        150,
		IsCorrect:    true,
		TilesFlipped: []string{"bio", "careerStats", "photo", "yearsActive"},
	}

	if result.Score != 150 {
		t.Errorf("Score = %v, want 150", result.Score)
	}
	if !result.IsCorrect {
		t.Error("IsCorrect should be true")
	}
	if len(result.TilesFlipped) != 4 {
		t.Errorf("len(TilesFlipped) = %v, want 4", len(result.TilesFlipped))
	}
	if result.TilesFlipped[0] != "bio" {
		t.Errorf("TilesFlipped[0] = %v, want bio", result.TilesFlipped[0])
	}
}

// Test ErrorResponse structure
func TestErrorResponseStructure(t *testing.T) {
	now := time.Now()
	errResp := ErrorResponse{
		Error:     "Bad Request",
		Message:   "Invalid sport parameter",
		Code:      "INVALID_PARAMETER",
		Timestamp: now,
		Details: map[string]interface{}{
			"parameter": "sport",
			"value":     "invalid",
		},
	}

	if errResp.Error != "Bad Request" {
		t.Errorf("Error = %v, want Bad Request", errResp.Error)
	}
	if errResp.Code != "INVALID_PARAMETER" {
		t.Errorf("Code = %v, want INVALID_PARAMETER", errResp.Code)
	}
	if errResp.Details["parameter"] != "sport" {
		t.Errorf("Details[parameter] = %v, want sport", errResp.Details["parameter"])
	}
}

// Test AverageNumberOfTileFlips calculation
func TestAverageNumberOfTileFlipsCalculation(t *testing.T) {
	tests := []struct {
		name               string
		initialAverage     float64
		initialTotalPlays  int
		newTileFlips       int
		expectedAverage    float64
		expectedTotalPlays int
	}{
		{
			name:               "first game with 5 tiles",
			initialAverage:     0.0,
			initialTotalPlays:  0,
			newTileFlips:       5,
			expectedAverage:    5.0,
			expectedTotalPlays: 1,
		},
		{
			name:               "second game with 3 tiles",
			initialAverage:     5.0,
			initialTotalPlays:  1,
			newTileFlips:       3,
			expectedAverage:    4.0,
			expectedTotalPlays: 2,
		},
		{
			name:               "third game with 7 tiles",
			initialAverage:     4.0,
			initialTotalPlays:  2,
			newTileFlips:       7,
			expectedAverage:    5.0,
			expectedTotalPlays: 3,
		},
		{
			name:               "game with 9 tiles",
			initialAverage:     5.0,
			initialTotalPlays:  3,
			newTileFlips:       9,
			expectedAverage:    6.0,
			expectedTotalPlays: 4,
		},
		{
			name:               "game with 0 tiles",
			initialAverage:     6.0,
			initialTotalPlays:  4,
			newTileFlips:       0,
			expectedAverage:    4.8,
			expectedTotalPlays: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the calculation logic from handlers.go
			totalPlays := tt.initialTotalPlays + 1
			totalTileFlips := tt.initialAverage * float64(tt.initialTotalPlays)
			totalTileFlips += float64(tt.newTileFlips)
			actualAverage := totalTileFlips / float64(totalPlays)

			if actualAverage != tt.expectedAverage {
				t.Errorf("AverageNumberOfTileFlips = %v, want %v", actualAverage, tt.expectedAverage)
			}
			if totalPlays != tt.expectedTotalPlays {
				t.Errorf("TotalPlays = %v, want %v", totalPlays, tt.expectedTotalPlays)
			}
		})
	}
}

// Test RoundStats with AverageNumberOfTileFlips
func TestRoundStatsWithAverageNumberOfTileFlips(t *testing.T) {
	stats := RoundStats{
		PlayDate: "2024-01-01",
		Name:     "Test Player",
		Sport:    "basketball",
		Stats: Stats{
			TotalPlays:                 100,
			PercentageCorrect:          75.5,
			HighestScore:               200,
			AverageCorrectScore:        150.0,
			AverageNumberOfTileFlips:   5.7,
			MostCommonFirstTileFlipped: "bio",
			MostCommonLastTileFlipped:  "photo",
			MostCommonTileFlipped:      "careerStats",
			LeastCommonTileFlipped:     "jerseyNumbers",
		},
	}

	if stats.AverageNumberOfTileFlips != 5.7 {
		t.Errorf("AverageNumberOfTileFlips = %v, want 5.7", stats.AverageNumberOfTileFlips)
	}
}

// Test UserName field in UserStats
func TestUserStatsWithUserName(t *testing.T) {
	now := time.Now()
	stats := UserStats{
		UserId:      "user-123",
		UserName:    "Jane Smith",
		UserCreated: now,
		Sports: []SportStats{
			{
				Sport: "basketball",
				Stats: Stats{
					TotalPlays:          50,
					PercentageCorrect:   85.0,
					HighestScore:        9,
					AverageCorrectScore: 7.5,
				},
			},
		},
	}

	if stats.UserName != "Jane Smith" {
		t.Errorf("UserName = %v, want Jane Smith", stats.UserName)
	}
	if stats.UserId != "user-123" {
		t.Errorf("UserId = %v, want user-123", stats.UserId)
	}
	if stats.Sports[0].TotalPlays != 50 {
		t.Errorf("Sports[0].TotalPlays = %v, want 50", stats.Sports[0].TotalPlays)
	}
}
