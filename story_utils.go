package main

import "fmt"

// createInitialStoryMissions returns the default set of story missions for a new user
func createEmptyStoryMissions(today string) []StoryMission {
	return []StoryMission{
		{
			StoryId:      "story_1",
			Criteria:     "Start the Game",
			Title:        "A Smashing Debut",
			DateAchieved: today,
			PlayerName:   "",
		},
		{
			StoryId:      "story_2",
			Criteria:     "Start a Case",
			Title:        "A Smashing Debut",
			DateAchieved: today,
			PlayerName:   "",
		},
		{
			StoryId:      "story_3",
			Criteria:     "Solve 1 Case",
			Title:        "A Smashing Debut",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_4",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_5",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_6",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_7",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_8",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_9",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_10",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_11",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_12",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_13",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_14",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_15",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_16",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_17",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_18",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_19",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_20",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_21",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_22",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_23",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_24",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_25",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_26",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_27",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			StoryId:      "story_28",
			Criteria:     "",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
	}
}

func currentDailyStreakStoryMissions(dailyStreak int) string {
	switch dailyStreak {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	case 10:
		return "10"
	case 11:
		return "11"
	case 12:
		return "12"
	case 13:
		return "13"
	case 14:
		return "14"
	case 21:
		return "21"
	default:
		return ""
	}
}

func totalWinsStoryMissions(totalWins int) string {
	switch totalWins {
	case 1:
		return "1"
	case 5:
		return "5"
	case 10:
		return "10"
	case 15:
		return "15"
	case 20:
		return "20"
	case 25:
		return "25"
	case 30:
		return "30"
	case 35:
		return "35"
	case 40:
		return "40"
	case 45:
		return "45"
	case 50:
		return "50"
	case 75:
		return "75"
	default:
		return ""
	}
}

func updateStoryMissions(storyMissions *[]StoryMission, currentDailyStreakStoryId *string, totalWinsStoryId *string, todayDate string, playerName string) {
	if storyMissions == nil || currentDailyStreakStoryId == nil || totalWinsStoryId == nil {
		return
	}
	storyIds := []string{*currentDailyStreakStoryId, *totalWinsStoryId}

	for _, storyId := range storyIds {
		if storyId == "" {
			continue
		}

		for i, mission := range *storyMissions {
			if mission.StoryId == storyId {
				fmt.Printf("Matched story mission: index=%d, storyId=%s, title=%s", i, mission.StoryId, mission.Title)
				if mission.DateAchieved != "" {
					// If already achieved Clear the original pointer value
					if storyId == *currentDailyStreakStoryId {
						*currentDailyStreakStoryId = ""
					}
					if storyId == *totalWinsStoryId {
						*totalWinsStoryId = ""
					}
				}
				(*storyMissions)[i].DateAchieved = todayDate
				(*storyMissions)[i].PlayerName = playerName
				break
			}
		}
	}
}
