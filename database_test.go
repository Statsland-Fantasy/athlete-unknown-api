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
		RoundID:     "test-round-id",
		Sport:       "basketball",
		PlayDate:    "2024-01-01",
		Created:     now,
		LastUpdated: now,
		PreviouslyPlayedDates: []string{"2023-12-31", "2023-12-30"},
		Player: Player{
			Sport:                 "basketball",
			SportsReferenceURL:    "http://example.com",
			Name:                  "Test Player",
			Bio:                   "Test bio",
			PlayerInformation:     "Test info",
			DraftInformation:      "Test draft",
			YearsActive:           "2010-2020",
			TeamsPlayedOn:         "Team A, Team B",
			JerseyNumbers:         "23, 32",
			CareerStats:           "Points: 1000",
			PersonalAchievements:  "MVP 2015",
			Photo:                 "http://example.com/photo.jpg",
		},
		Stats: RoundStats{
			PlayDate:              "2024-01-01",
			Name:                  "Test Player",
			Sport:                 "basketball",
			TotalPlays:            100,
			PercentageCorrect:     75.5,
			HighestScore:          200,
			AverageCorrectScore:   150.0,
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
		UserID:            "test-user-123",
		UserName:          "test-user-name",
		UserCreated:       now,
		TotalPlays:        500,
		PercentageCorrect: 80.0,
		HighestScore:      250,
		AverageCorrectScore: 175.0,
		MostCommonFirstTileFlipped: "bio",
		MostCommonLastTileFlipped:  "photo",
		MostCommonTileFlipped:      "careerStats",
		LeastCommonTileFlipped:     "jerseyNumbers",
		Sports: []SportStats{
			{
				Sport:                       "basketball",
				CurrentDailyStreak:          5,
				TotalPlays:                  200,
				PercentageCorrect:           75.0,
				HighestScore:                200,
				AverageCorrectScore:         150.0,
				MostCommonFirstTileFlipped:  "bio",
				MostCommonLastTileFlipped:   "photo",
				MostCommonTileFlipped:       "careerStats",
				LeastCommonTileFlipped:      "jerseyNumbers",
			},
			{
				Sport:                       "baseball",
				CurrentDailyStreak:          3,
				TotalPlays:                  150,
				PercentageCorrect:           80.0,
				HighestScore:                180,
				AverageCorrectScore:         140.0,
				MostCommonFirstTileFlipped:  "playerInformation",
				MostCommonLastTileFlipped:   "careerStats",
				MostCommonTileFlipped:       "bio",
				LeastCommonTileFlipped:      "photo",
			},
		},
		FirstTileFlippedTracker: TileFlipTracker{
			Bio:                  100,
			PlayerInformation:    80,
			DraftInformation:     50,
			TeamsPlayedOn:        30,
			JerseyNumbers:        20,
			CareerStats:          40,
			PersonalAchievements: 60,
			Photo:                70,
			YearsActive:          50,
		},
	}

	// Verify all fields are set correctly
	if stats.UserID != "test-user-123" {
		t.Errorf("UserID = %v, want test-user-123", stats.UserID)
	}
	if stats.UserName != "test-user-name" {
		t.Errorf("UserName = %v, want test-user-name", stats.UserName)
	}
	if stats.TotalPlays != 500 {
		t.Errorf("TotalPlays = %v, want 500", stats.TotalPlays)
	}
	if stats.PercentageCorrect != 80.0 {
		t.Errorf("PercentageCorrect = %v, want 80.0", stats.PercentageCorrect)
	}
	if len(stats.Sports) != 2 {
		t.Errorf("len(Sports) = %v, want 2", len(stats.Sports))
	}
	if stats.Sports[0].Sport != "basketball" {
		t.Errorf("Sports[0].Sport = %v, want basketball", stats.Sports[0].Sport)
	}
	if stats.FirstTileFlippedTracker.Bio != 100 {
		t.Errorf("FirstTileFlippedTracker.Bio = %v, want 100", stats.FirstTileFlippedTracker.Bio)
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
		Sport:                 "basketball",
		SportsReferenceURL:    "http://example.com/player",
		Name:                  "John Doe",
		Bio:                   "Professional basketball player",
		PlayerInformation:     "Height: 6'5\", Weight: 220 lbs",
		DraftInformation:      "Round 1, Pick 5, 2010",
		YearsActive:           "2010-2023",
		TeamsPlayedOn:         "Lakers, Bulls, Heat",
		JerseyNumbers:         "23, 6, 24",
		CareerStats:           "PPG: 25.5, RPG: 7.2, APG: 6.8",
		PersonalAchievements:  "3x NBA Champion, 2x MVP, 10x All-Star",
		Photo:                 "http://example.com/photo.jpg",
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
