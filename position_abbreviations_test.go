package main

import (
	"testing"
)

func TestAbbreviatePositions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basketball - Single position",
			input:    "Position: Point Guard",
			expected: "Position: PG",
		},
		{
			name:     "Basketball - Multiple positions with 'and'",
			input:    "Position: Shooting Guard and Point Guard",
			expected: "Position: SG, PG", // " and " is replaced with ", "
		},
		{
			name:     "Basketball - Position with other attributes",
			input:    " ▪ Height: 6-3 ▪ Weight: 185lb Position: Small Forward Shoots: Right",
			expected: " ▪ Height: 6-3 ▪ Weight: 185lb Position: SF Shoots: Right", // "-" not replaced in abbreviatePositions
		},
		{
			name:     "Basketball - Hyphenated position (Guard-Forward)",
			input:    "Position: Guard-Forward",
			expected: "Position: G-F", // "-" not replaced in abbreviatePositions
		},
		{
			name:     "Basketball - Generic Guard position",
			input:    "Position: Guard",
			expected: "Position: G",
		},
		{
			name:     "No position mentioned",
			input:    " ▪ Height: 6-8 ▪ Weight: 220lb",
			expected: " ▪ Height: 6-8 ▪ Weight: 220lb", // "-" not replaced in abbreviatePositions
		},
		// Baseball position tests
		{
			name:     "Baseball - Designated Hitter",
			input:    "Position: Designated Hitter",
			expected: "Position: DH",
		},
		{
			name:     "Baseball - First Baseman",
			input:    "Position: First Baseman",
			expected: "Position: 1B",
		},
		{
			name:     "Baseball - Second Baseman",
			input:    "Position: Second Baseman",
			expected: "Position: 2B",
		},
		{
			name:     "Baseball - Third Baseman",
			input:    "Position: Third Baseman",
			expected: "Position: 3B",
		},
		{
			name:     "Baseball - Shortstop",
			input:    "Position: Shortstop",
			expected: "Position: SS",
		},
		{
			name:     "Baseball - Catcher",
			input:    "Position: Catcher",
			expected: "Position: C",
		},
		{
			name:     "Baseball - Pitcher",
			input:    "Position: Pitcher",
			expected: "Position: P",
		},
		{
			name:     "Baseball - Centerfielder",
			input:    "Position: Centerfielder",
			expected: "Position: CF",
		},
		{
			name:     "Baseball - Rightfielder",
			input:    "Position: Rightfielder",
			expected: "Position: RF",
		},
		{
			name:     "Baseball - Leftfielder",
			input:    "Position: Leftfielder",
			expected: "Position: LF",
		},
		{
			name:     "Baseball - Outfielder",
			input:    "Position: Outfielder",
			expected: "Position: OF",
		},
		{
			name:     "Baseball - Multiple positions with 'and'",
			input:    "Position: First Baseman and Outfielder",
			expected: "Position: 1B, OF", // " and " is replaced with ", "
		},
		{
			name:     "Baseball - Position with other attributes",
			input:    " ▪ Height: 6-2 ▪ Weight: 205lb Position: Pitcher Bats: Right Throws: Right",
			expected: " ▪ Height: 6-2 ▪ Weight: 205lb Position: P Bats: Right Throws: Right", // "-" not replaced in abbreviatePositions
		},
		// TODO: Add football position tests
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abbreviatePositions(tt.input)
			if result != tt.expected {
				t.Errorf("abbreviatePositions() = %q, want %q", result, tt.expected)
			}
		})
	}
}
