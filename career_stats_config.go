package main

import (
	"strings"
)

// CareerStatsConfig defines the HTML path and stat name for a career stat
type CareerStatsConfig struct {
	HTMLPath  string // CSS selector or HTML path to the stat element
	StatLabel string // The name/label of the stat (e.g., "Points", "Games", "WAR")
}

// StatsConfig defines which stats to extract for a sport/position combination
type StatsConfig struct {
	Stats []CareerStatsConfig
}

// GetCareerStatsConfig returns the stats configuration for a given sport and position
func GetCareerStatsConfig(sport, playerInfo string) StatsConfig {
	// Normalize inputs to lowercase for case-insensitive matching
	sport = strings.ToLower(sport)
	playerInfo = strings.ToLower(playerInfo)

	switch sport {
	case SportBaseball:
		// Check if pitcher or hitter
		if strings.Contains(playerInfo, "position: p") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1 > :nth-of-type(2) > p:last-of-type", StatLabel: "W"},
					{HTMLPath: "div.p2 > :nth-of-type(3) > p:last-of-type", StatLabel: "SV"},
					{HTMLPath: "div.p3 > :nth-of-type(2) > p:last-of-type", StatLabel: "K"},
					{HTMLPath: "div.p1 > :nth-of-type(4) > p:last-of-type", StatLabel: "ERA"},
					{HTMLPath: "div.p1 > :nth-of-type(1) > p:last-of-type", StatLabel: "WAR"},
				},
			}
		}
		// Hitter stats
		return StatsConfig{
			Stats: []CareerStatsConfig{
				{HTMLPath: "div.p1 > :nth-of-type(5) > p:last-of-type", StatLabel: "AVG"},
				{HTMLPath: "div.p1 > :nth-of-type(4) > p:last-of-type", StatLabel: "HR"},
				{HTMLPath: "div.p2 > :nth-of-type(3) > p:last-of-type", StatLabel: "SB"},
				{HTMLPath: "div.p1 > :nth-of-type(1) > p:last-of-type", StatLabel: "WAR"},
			},
		}

	case SportBasketball:
		return StatsConfig{
			Stats: []CareerStatsConfig{
				{HTMLPath: "div.p1 > :nth-of-type(2) > p:last-of-type", StatLabel: "PPG"},
				{HTMLPath: "div.p1 > :nth-of-type(3) > p:last-of-type", StatLabel: "RPG"},
				{HTMLPath: "div.p1 > :nth-of-type(4) > p:last-of-type", StatLabel: "APG"},
				{HTMLPath: "div.p3 > :nth-of-type(2) > p:last-of-type", StatLabel: "WS"},
			},
		}

	case SportFootball:
		if strings.Contains(playerInfo, "qb") || strings.Contains(playerInfo, "quarterback") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(3) > p:last-of-type", StatLabel: "YDS"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(5) > p:last-of-type", StatLabel: "TD"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p2:nth-of-type(6) > p:last-of-type", StatLabel: "INT"},
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
				},
			}
		} else if strings.Contains(playerInfo, "rb") || strings.Contains(playerInfo, "running") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "RUSH"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "YDS"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(4) > p:last-of-type", StatLabel: "TD"},
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
				},
			}
		} else if strings.Contains(playerInfo, "wr") || strings.Contains(playerInfo, "te") || strings.Contains(playerInfo, "receiver") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "REC"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "YDS"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(4) > p:last-of-type", StatLabel: "TD"},
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
				},
			}
		} else if strings.Contains(playerInfo, "db") || strings.Contains(playerInfo, "cb") || strings.Contains(playerInfo, "fs") || strings.Contains(playerInfo, "ss") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "G"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "INT"},
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
				},
			}
		} else if strings.Contains(playerInfo, "nt") || strings.Contains(playerInfo, "dt") || strings.Contains(playerInfo, "de") || strings.Contains(playerInfo, ": lb") || strings.Contains(playerInfo, ": olb") || strings.Contains(playerInfo, ": ilb") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "G"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p2:nth-of-type(1) > p:last-of-type", StatLabel: "SK"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p2:nth-of-type(2) > p:last-of-type", StatLabel: "SOLO"},
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
				},
			}
		} else if strings.Contains(playerInfo, ": k") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "G"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p2:nth-of-type(1) > p:last-of-type", StatLabel: "FGM"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p2:nth-of-type(2) > p:last-of-type", StatLabel: "FGA"},
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
				},
			}
		} else if strings.Contains(playerInfo, ": p") {
			return StatsConfig{
				Stats: []CareerStatsConfig{
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "G"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p2:nth-of-type(1) > p:last-of-type", StatLabel: "PNT"},
					{HTMLPath: "div.p1:nth-of-type(3) > div.p2:nth-of-type(2) > p:last-of-type", StatLabel: "YDS"},
					{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
				},
			}
		}
		// Default for other positions (ie offensive linemen)
		return StatsConfig{
			Stats: []CareerStatsConfig{
				{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(1) > p:last-of-type", StatLabel: "G"},
				{HTMLPath: "div.p1:nth-of-type(2) > div.p1:nth-of-type(2) > p:last-of-type", StatLabel: "AV"},
			},
		}

	default:
		return StatsConfig{}
	}
}
