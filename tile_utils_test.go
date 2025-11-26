package main

import "testing"

func TestIncrementTileTracker(t *testing.T) {
	tests := []struct {
		name       string
		tracker    *TileFlipTracker
		tileName   string
		wantNil    bool
		checkField func(*TileFlipTracker) int
		expected   int
	}{
		{
			name:     "nil tracker",
			tracker:  nil,
			tileName: "bio",
			wantNil:  true,
		},
		{
			name:     "empty tile name",
			tracker:  &TileFlipTracker{},
			tileName: "",
			checkField: func(t *TileFlipTracker) int {
				return t.Bio
			},
			expected: 0,
		},
		{
			name:     "increment bio",
			tracker:  &TileFlipTracker{},
			tileName: "bio",
			checkField: func(t *TileFlipTracker) int {
				return t.Bio
			},
			expected: 1,
		},
		{
			name:     "increment playerInformation",
			tracker:  &TileFlipTracker{},
			tileName: "playerInformation",
			checkField: func(t *TileFlipTracker) int {
				return t.PlayerInformation
			},
			expected: 1,
		},
		{
			name:     "increment draftInformation",
			tracker:  &TileFlipTracker{},
			tileName: "draftInformation",
			checkField: func(t *TileFlipTracker) int {
				return t.DraftInformation
			},
			expected: 1,
		},
		{
			name:     "increment teamsPlayedOn",
			tracker:  &TileFlipTracker{},
			tileName: "teamsPlayedOn",
			checkField: func(t *TileFlipTracker) int {
				return t.TeamsPlayedOn
			},
			expected: 1,
		},
		{
			name:     "increment jerseyNumbers",
			tracker:  &TileFlipTracker{},
			tileName: "jerseyNumbers",
			checkField: func(t *TileFlipTracker) int {
				return t.JerseyNumbers
			},
			expected: 1,
		},
		{
			name:     "increment careerStats",
			tracker:  &TileFlipTracker{},
			tileName: "careerStats",
			checkField: func(t *TileFlipTracker) int {
				return t.CareerStats
			},
			expected: 1,
		},
		{
			name:     "increment personalAchievements",
			tracker:  &TileFlipTracker{},
			tileName: "personalAchievements",
			checkField: func(t *TileFlipTracker) int {
				return t.PersonalAchievements
			},
			expected: 1,
		},
		{
			name:     "increment photo",
			tracker:  &TileFlipTracker{},
			tileName: "photo",
			checkField: func(t *TileFlipTracker) int {
				return t.Photo
			},
			expected: 1,
		},
		{
			name:     "increment yearsActive",
			tracker:  &TileFlipTracker{},
			tileName: "yearsActive",
			checkField: func(t *TileFlipTracker) int {
				return t.YearsActive
			},
			expected: 1,
		},
		{
			name:     "invalid tile name",
			tracker:  &TileFlipTracker{},
			tileName: "invalidTile",
			checkField: func(t *TileFlipTracker) int {
				return t.Bio
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			incrementTileTracker(tt.tracker, tt.tileName)
			if tt.wantNil {
				return
			}
			if got := tt.checkField(tt.tracker); got != tt.expected {
				t.Errorf("incrementTileTracker() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFindMostCommonTile(t *testing.T) {
	tests := []struct {
		name    string
		tracker *TileFlipTracker
		want    string
	}{
		{
			name:    "nil tracker",
			tracker: nil,
			want:    "",
		},
		{
			name:    "empty tracker",
			tracker: &TileFlipTracker{},
			want:    "",
		},
		{
			name: "bio is most common",
			tracker: &TileFlipTracker{
				Bio:                  10,
				PlayerInformation:    5,
				DraftInformation:     3,
				TeamsPlayedOn:        2,
				JerseyNumbers:        1,
				CareerStats:          4,
				PersonalAchievements: 6,
				Photo:                7,
				YearsActive:          8,
			},
			want: "bio",
		},
		{
			name: "yearsActive is most common",
			tracker: &TileFlipTracker{
				Bio:                  1,
				PlayerInformation:    2,
				DraftInformation:     3,
				TeamsPlayedOn:        4,
				JerseyNumbers:        5,
				CareerStats:          6,
				PersonalAchievements: 7,
				Photo:                8,
				YearsActive:          15,
			},
			want: "yearsActive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findMostCommonTile(tt.tracker); got != tt.want {
				t.Errorf("findMostCommonTile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindLeastCommonTile(t *testing.T) {
	tests := []struct {
		name    string
		tracker *TileFlipTracker
		want    string
	}{
		{
			name:    "nil tracker",
			tracker: nil,
			want:    "",
		},
		{
			name:    "empty tracker",
			tracker: &TileFlipTracker{},
			want:    "",
		},
		{
			name: "jerseyNumbers is least common",
			tracker: &TileFlipTracker{
				Bio:                  10,
				PlayerInformation:    5,
				DraftInformation:     3,
				TeamsPlayedOn:        2,
				JerseyNumbers:        1,
				CareerStats:          4,
				PersonalAchievements: 6,
				Photo:                7,
				YearsActive:          8,
			},
			want: "jerseyNumbers",
		},
		{
			name: "bio is least common (ignoring zeros)",
			tracker: &TileFlipTracker{
				Bio:                  1,
				PlayerInformation:    2,
				DraftInformation:     3,
				TeamsPlayedOn:        4,
				JerseyNumbers:        5,
				CareerStats:          6,
				PersonalAchievements: 7,
				Photo:                8,
				YearsActive:          0,
			},
			want: "bio",
		},
		{
			name: "some fields are zero",
			tracker: &TileFlipTracker{
				Bio:                  0,
				PlayerInformation:    0,
				DraftInformation:     5,
				TeamsPlayedOn:        0,
				JerseyNumbers:        3,
				CareerStats:          0,
				PersonalAchievements: 0,
				Photo:                0,
				YearsActive:          0,
			},
			want: "jerseyNumbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findLeastCommonTile(tt.tracker); got != tt.want {
				t.Errorf("findLeastCommonTile() = %v, want %v", got, tt.want)
			}
		})
	}
}
