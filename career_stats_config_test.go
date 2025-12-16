package main

import (
	"testing"
)

func TestGetCareerStatsConfig(t *testing.T) {
	tests := []struct {
		name          string
		sport         string
		playerInfo    string
		expectedStats int
		checkStat     func(StatsConfig) bool
	}{
		{
			name:          "baseball pitcher",
			sport:         "baseball",
			playerInfo:    "Position: Pitcher",
			expectedStats: 5,
			checkStat: func(config StatsConfig) bool {
				// Check for pitcher-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "W" || stat.StatLabel == "SV" || stat.StatLabel == "K" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "baseball hitter",
			sport:         "baseball",
			playerInfo:    "Position: Outfielder",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for hitter-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "AVG" || stat.StatLabel == "HR" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "basketball player",
			sport:         "basketball",
			playerInfo:    "Position: Guard",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for basketball stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "PPG" || stat.StatLabel == "RPG" || stat.StatLabel == "APG" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football quarterback",
			sport:         "football",
			playerInfo:    "Position: QB",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for QB-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "YDS" || stat.StatLabel == "TD" || stat.StatLabel == "INT" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football running back",
			sport:         "football",
			playerInfo:    "Position: RB",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for RB-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "RUSH" || stat.StatLabel == "YDS" || stat.StatLabel == "TD" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football wide receiver",
			sport:         "football",
			playerInfo:    "Position: WR",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for WR-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "REC" || stat.StatLabel == "YDS" || stat.StatLabel == "TD" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football defensive back",
			sport:         "football",
			playerInfo:    "Position: CB",
			expectedStats: 3,
			checkStat: func(config StatsConfig) bool {
				// Check for DB-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "INT" || stat.StatLabel == "G" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football linebacker",
			sport:         "football",
			playerInfo:    "Position: LB",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for LB-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "SK" || stat.StatLabel == "SOLO" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football kicker",
			sport:         "football",
			playerInfo:    "Position: K",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for kicker-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "FGM" || stat.StatLabel == "FGA" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football punter",
			sport:         "football",
			playerInfo:    "Position: P",
			expectedStats: 4,
			checkStat: func(config StatsConfig) bool {
				// Check for punter-specific stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "PNT" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "football offensive lineman (default)",
			sport:         "football",
			playerInfo:    "Position: OT",
			expectedStats: 2,
			checkStat: func(config StatsConfig) bool {
				// Check for default stats
				for _, stat := range config.Stats {
					if stat.StatLabel == "G" || stat.StatLabel == "AV" {
						return true
					}
				}
				return false
			},
		},
		{
			name:          "unknown sport",
			sport:         "soccer",
			playerInfo:    "Forward",
			expectedStats: 0,
			checkStat: func(config StatsConfig) bool {
				return true
			},
		},
		{
			name:          "case insensitive baseball pitcher",
			sport:         "BASEBALL",
			playerInfo:    "POSITION: PITCHER",
			expectedStats: 5,
			checkStat: func(config StatsConfig) bool {
				// Should still detect pitcher
				for _, stat := range config.Stats {
					if stat.StatLabel == "W" || stat.StatLabel == "SV" {
						return true
					}
				}
				return false
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCareerStatsConfig(tt.sport, tt.playerInfo)
			if len(got.Stats) != tt.expectedStats {
				t.Errorf("GetCareerStatsConfig(%q, %q) returned %d stats, want %d", tt.sport, tt.playerInfo, len(got.Stats), tt.expectedStats)
			}
			if !tt.checkStat(got) {
				t.Errorf("GetCareerStatsConfig(%q, %q) missing expected stat", tt.sport, tt.playerInfo)
			}
		})
	}
}
