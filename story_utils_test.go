package main

import (
	"testing"
)

func TestCreateEmptyStoryMissions(t *testing.T) {
	today := "2026-01-15"
	missions := createEmptyStoryMissions(today)

	if len(missions) != 28 {
		t.Fatalf("expected 28 missions, got %d", len(missions))
	}

	// First two missions should have today as DateAchieved
	if missions[0].Criteria != CriteriaStartGame {
		t.Errorf("mission 0: expected criteria %q, got %q", CriteriaStartGame, missions[0].Criteria)
	}
	if missions[0].StoryId != "story_1" {
		t.Errorf("mission 0: expected storyId %q, got %q", "story_1", missions[0].StoryId)
	}
	if missions[0].DateAchieved != today {
		t.Errorf("mission 0: expected dateAchieved %q, got %q", today, missions[0].DateAchieved)
	}

	if missions[1].Criteria != CriteriaPlay1Day {
		t.Errorf("mission 1: expected criteria %q, got %q", CriteriaPlay1Day, missions[1].Criteria)
	}
	if missions[1].DateAchieved != today {
		t.Errorf("mission 1: expected dateAchieved %q, got %q", today, missions[1].DateAchieved)
	}

	// Remaining missions should have empty DateAchieved
	for i := 2; i < len(missions); i++ {
		if missions[i].DateAchieved != "" {
			t.Errorf("mission %d: expected empty dateAchieved, got %q", i, missions[i].DateAchieved)
		}
	}

	// Verify each mission has a unique storyId
	storyIds := make(map[string]bool)
	for i, m := range missions {
		if storyIds[m.StoryId] {
			t.Errorf("mission %d: duplicate storyId %q", i, m.StoryId)
		}
		storyIds[m.StoryId] = true
	}
}

func TestDaysPlayedStoryMissions(t *testing.T) {
	tests := []struct {
		name       string
		daysPlayed int
		expected   Criteria
	}{
		{"1 day", 1, CriteriaPlay1Day},
		{"2 days", 2, CriteriaPlay2Days},
		{"3 days", 3, CriteriaPlay3Days},
		{"4 days", 4, CriteriaPlay4Days},
		{"5 days", 5, CriteriaPlay5Days},
		{"6 days", 6, CriteriaPlay6Days},
		{"7 days", 7, CriteriaPlay7Days},
		{"8 days", 8, CriteriaPlay8Days},
		{"9 days", 9, CriteriaPlay9Days},
		{"10 days", 10, CriteriaPlay10Days},
		{"11 days", 11, CriteriaPlay11Days},
		{"12 days", 12, CriteriaPlay12Days},
		{"13 days", 13, CriteriaPlay13Days},
		{"14 days", 14, CriteriaPlay14Days},
		{"0 days returns none", 0, CriteriaNone},
		{"15 days returns none", 15, CriteriaNone},
		{"negative returns none", -1, CriteriaNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := daysPlayedStoryMissions(tt.daysPlayed)
			if result != tt.expected {
				t.Errorf("daysPlayedStoryMissions(%d) = %q, want %q", tt.daysPlayed, result, tt.expected)
			}
		})
	}
}

func TestTotalWinsStoryMissions(t *testing.T) {
	tests := []struct {
		name      string
		totalWins int
		expected  Criteria
	}{
		{"1 win", 1, CriteriaSolve1Case},
		{"10 wins", 10, CriteriaSolve10Cases},
		{"20 wins", 20, CriteriaSolve20Cases},
		{"30 wins", 30, CriteriaSolve30Cases},
		{"0 wins returns none", 0, CriteriaNone},
		{"5 wins returns none", 5, CriteriaNone},
		{"31 wins returns none", 31, CriteriaNone},
		{"negative returns none", -1, CriteriaNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := totalWinsStoryMissions(tt.totalWins)
			if result != tt.expected {
				t.Errorf("totalWinsStoryMissions(%d) = %q, want %q", tt.totalWins, result, tt.expected)
			}
		})
	}
}

func TestCurrentDailyStreakStoryMissions(t *testing.T) {
	tests := []struct {
		name          string
		currentStreak int
		expected      Criteria
	}{
		{"3 day streak", 3, CriteriaPlay3ConsecutiveDays},
		{"5 day streak", 5, CriteriaPlay5ConsecutiveDays},
		{"7 day streak", 7, CriteriaPlay7ConsecutiveDays},
		{"10 day streak", 10, CriteriaPlay10ConsecutiveDays},
		{"0 streak returns none", 0, CriteriaNone},
		{"1 streak returns none", 1, CriteriaNone},
		{"4 streak returns none", 4, CriteriaNone},
		{"11 streak returns none", 11, CriteriaNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := currentDailyStreakStoryMissions(tt.currentStreak)
			if result != tt.expected {
				t.Errorf("currentDailyStreakStoryMissions(%d) = %q, want %q", tt.currentStreak, result, tt.expected)
			}
		})
	}
}

func TestScoreStoryMissions(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected Criteria
	}{
		{"perfect score", 100, CriteriaScore100},
		{"score 95", 95, CriteriaScore95},
		{"score 96", 96, CriteriaScore95},
		{"score 99", 99, CriteriaScore95},
		{"score 0 is a loss", 0, CriteriaLose},
		{"score 50 returns none", 50, CriteriaNone},
		{"score 94 returns none", 94, CriteriaNone},
		{"score 1 returns none", 1, CriteriaNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scoreStoryMissions(tt.score)
			if result != tt.expected {
				t.Errorf("scoreStoryMissions(%d) = %q, want %q", tt.score, result, tt.expected)
			}
		})
	}
}

func TestCalculateAchievedStoryMissions(t *testing.T) {
	t.Run("nil user returns nil", func(t *testing.T) {
		result := calculateAchievedStoryMissions(nil, 100)
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})

	t.Run("no milestones hit returns empty", func(t *testing.T) {
		user := &User{
			TotalDaysPlayed:    0,
			TotalWins:          0,
			CurrentDailyStreak: 0,
		}
		result := calculateAchievedStoryMissions(user, 50)
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("perfect score only", func(t *testing.T) {
		user := &User{
			TotalDaysPlayed:    0,
			TotalWins:          0,
			CurrentDailyStreak: 0,
		}
		result := calculateAchievedStoryMissions(user, 100)
		if len(result) != 1 || result[0] != CriteriaScore100 {
			t.Errorf("expected [%q], got %v", CriteriaScore100, result)
		}
	})

	t.Run("multiple milestones at once", func(t *testing.T) {
		user := &User{
			TotalDaysPlayed:    3,
			TotalWins:          1,
			CurrentDailyStreak: 3,
		}
		result := calculateAchievedStoryMissions(user, 100)

		expected := map[Criteria]bool{
			CriteriaScore100:             true,
			CriteriaPlay3Days:            true,
			CriteriaPlay3ConsecutiveDays: true,
			CriteriaSolve1Case:           true,
		}

		if len(result) != len(expected) {
			t.Fatalf("expected %d criteria, got %d: %v", len(expected), len(result), result)
		}
		for _, c := range result {
			if !expected[c] {
				t.Errorf("unexpected criteria %q in result", c)
			}
		}
	})

	t.Run("days played milestone only", func(t *testing.T) {
		user := &User{
			TotalDaysPlayed:    7,
			TotalWins:          5,
			CurrentDailyStreak: 2,
		}
		result := calculateAchievedStoryMissions(user, 50)
		if len(result) != 1 || result[0] != CriteriaPlay7Days {
			t.Errorf("expected [%q], got %v", CriteriaPlay7Days, result)
		}
	})
}

func TestUpdateStoryMissions(t *testing.T) {

	t.Run("updates matching mission and returns filtered criteria", func(t *testing.T) {
		user := &User{
			StoryMissions: createEmptyStoryMissions("2026-01-01"),
		}
		earned := []Criteria{CriteriaSolve1Case}
		result := updateStoryMissions(user, "2026-01-15", "Mike Trout", earned)

		if len(result) != 1 || result[0] != CriteriaSolve1Case {
			t.Errorf("expected [%q], got %v", CriteriaSolve1Case, result)
		}

		// Verify the mission was updated on the user
		for _, m := range user.StoryMissions {
			if m.Criteria == CriteriaSolve1Case {
				if m.DateAchieved != "2026-01-15" {
					t.Errorf("expected dateAchieved %q, got %q", "2026-01-15", m.DateAchieved)
				}
				if m.PlayerName != "Mike Trout" {
					t.Errorf("expected playerName %q, got %q", "Mike Trout", m.PlayerName)
				}
			}
		}
	})

	t.Run("days played criteria updates mission but is excluded from returned slice", func(t *testing.T) {
		user := &User{
			StoryMissions: createEmptyStoryMissions("2026-01-01"),
		}
		earned := []Criteria{CriteriaPlay3Days}
		result := updateStoryMissions(user, "2026-01-15", "Babe Ruth", earned)

		if len(result) != 0 {
			t.Errorf("expected empty filtered result for days played criteria, got %v", result)
		}

		// But the mission itself should still be updated
		for _, m := range user.StoryMissions {
			if m.Criteria == CriteriaPlay3Days {
				if m.DateAchieved != "2026-01-15" {
					t.Errorf("expected dateAchieved %q, got %q", "2026-01-15", m.DateAchieved)
				}
				if m.PlayerName != "Babe Ruth" {
					t.Errorf("expected playerName %q, got %q", "Babe Ruth", m.PlayerName)
				}
			}
		}
	})

	t.Run("skips already achieved missions", func(t *testing.T) {
		user := &User{
			StoryMissions: createEmptyStoryMissions("2026-01-01"),
		}
		// First update
		updateStoryMissions(user, "2026-01-10", "Player 1", []Criteria{CriteriaSolve1Case})

		// Second update with same criteria should not overwrite
		result := updateStoryMissions(user, "2026-01-20", "Player 2", []Criteria{CriteriaSolve1Case})

		if len(result) != 0 {
			t.Errorf("expected empty result for already achieved mission, got %v", result)
		}

		for _, m := range user.StoryMissions {
			if m.Criteria == CriteriaSolve1Case {
				if m.DateAchieved != "2026-01-10" {
					t.Errorf("original dateAchieved should be preserved, got %q", m.DateAchieved)
				}
				if m.PlayerName != "Player 1" {
					t.Errorf("original playerName should be preserved, got %q", m.PlayerName)
				}
			}
		}
	})

	t.Run("mixed criteria returns only non-days-played", func(t *testing.T) {
		user := &User{
			StoryMissions: createEmptyStoryMissions("2026-01-01"),
		}
		earned := []Criteria{CriteriaPlay5Days, CriteriaSolve1Case, CriteriaScore100}
		result := updateStoryMissions(user, "2026-01-15", "Test", earned)

		if len(result) != 2 {
			t.Fatalf("expected 2 filtered criteria, got %d: %v", len(result), result)
		}

		expected := map[Criteria]bool{CriteriaSolve1Case: true, CriteriaScore100: true}
		for _, c := range result {
			if !expected[c] {
				t.Errorf("unexpected criteria %q in filtered result", c)
			}
		}
	})

	t.Run("empty earned criteria returns nil", func(t *testing.T) {
		user := &User{
			StoryMissions: createEmptyStoryMissions("2026-01-01"),
		}
		result := updateStoryMissions(user, "2026-01-15", "Test", []Criteria{})

		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

func TestIsDaysPlayedCriteria(t *testing.T) {
	daysPlayedCriteria := []Criteria{
		CriteriaPlay1Day, CriteriaPlay2Days, CriteriaPlay3Days, CriteriaPlay4Days,
		CriteriaPlay5Days, CriteriaPlay6Days, CriteriaPlay7Days, CriteriaPlay8Days,
		CriteriaPlay9Days, CriteriaPlay10Days, CriteriaPlay11Days, CriteriaPlay12Days,
		CriteriaPlay13Days, CriteriaPlay14Days,
	}

	for _, c := range daysPlayedCriteria {
		t.Run(string(c)+" is days played", func(t *testing.T) {
			if !isDaysPlayedCriteria(c) {
				t.Errorf("expected %q to be days played criteria", c)
			}
		})
	}

	nonDaysPlayedCriteria := []Criteria{
		CriteriaStartGame, CriteriaSolve1Case, CriteriaSolve10Cases,
		CriteriaPlay3ConsecutiveDays, CriteriaScore100, CriteriaLose, CriteriaNone,
	}

	for _, c := range nonDaysPlayedCriteria {
		t.Run(string(c)+" is not days played", func(t *testing.T) {
			if isDaysPlayedCriteria(c) {
				t.Errorf("expected %q to NOT be days played criteria", c)
			}
		})
	}
}
