package main

// createInitialStoryMissions returns the default set of story missions for a new user
func createEmptyStoryMissions(today string) []StoryMission {
	return []StoryMission{
		{
			Criteria:     CriteriaStartGame,
			StoryId:      "story_1",
			Title:        "A Smashing Debut",
			DateAchieved: today,
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay1Day,
			StoryId:      "story_2",
			Title:        "A Smashing Debut",
			DateAchieved: today,
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay2Days,
			StoryId:      "story_3",
			Title:        "A Smashing Debut",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay3Days,
			StoryId:      "story_4",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay4Days,
			StoryId:      "story_5",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay5Days,
			StoryId:      "story_6",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay6Days,
			StoryId:      "story_7",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay7Days,
			StoryId:      "story_8",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay8Days,
			StoryId:      "story_9",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay9Days,
			StoryId:      "story_10",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay10Days,
			StoryId:      "story_11",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay11Days,
			StoryId:      "story_12",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay12Days,
			StoryId:      "story_13",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay13Days,
			StoryId:      "story_14",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay14Days,
			StoryId:      "story_15",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaSolve1Case,
			StoryId:      "story_16",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaSolve10Cases,
			StoryId:      "story_17",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaSolve20Cases,
			StoryId:      "story_18",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaSolve30Cases,
			StoryId:      "story_19",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay3ConsecutiveDays,
			StoryId:      "story_20",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay5ConsecutiveDays,
			StoryId:      "story_21",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay7ConsecutiveDays,
			StoryId:      "story_22",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaPlay10ConsecutiveDays,
			StoryId:      "story_23",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaScore100,
			StoryId:      "story_24",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaScore95,
			StoryId:      "story_25",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaLose,
			StoryId:      "story_26",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaNone,
			StoryId:      "story_27",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
		{
			Criteria:     CriteriaNone,
			StoryId:      "story_28",
			Title:        "",
			DateAchieved: "",
			PlayerName:   "",
		},
	}
}

func daysPlayedStoryMissions(daysPlayed int) Criteria {
	switch daysPlayed {
	case 1:
		return CriteriaPlay1Day
	case 2:
		return CriteriaPlay2Days
	case 3:
		return CriteriaPlay3Days
	case 4:
		return CriteriaPlay4Days
	case 5:
		return CriteriaPlay5Days
	case 6:
		return CriteriaPlay6Days
	case 7:
		return CriteriaPlay7Days
	case 8:
		return CriteriaPlay8Days
	case 9:
		return CriteriaPlay9Days
	case 10:
		return CriteriaPlay10Days
	case 11:
		return CriteriaPlay11Days
	case 12:
		return CriteriaPlay12Days
	case 13:
		return CriteriaPlay13Days
	case 14:
		return CriteriaPlay14Days
	default:
		return CriteriaNone
	}
}

func totalWinsStoryMissions(totalWins int) Criteria {
	switch totalWins {
	case 1:
		return CriteriaSolve1Case
	case 10:
		return CriteriaSolve10Cases
	case 20:
		return CriteriaSolve20Cases
	case 30:
		return CriteriaSolve30Cases
	default:
		return CriteriaNone
	}
}

func currentDailyStreakStoryMissions(currentStreak int) Criteria {
	switch currentStreak {
	case 3:
		return CriteriaPlay3ConsecutiveDays
	case 5:
		return CriteriaPlay5ConsecutiveDays
	case 7:
		return CriteriaPlay7ConsecutiveDays
	case 10:
		return CriteriaPlay10ConsecutiveDays
	default:
		return CriteriaNone
	}
}

func scoreStoryMissions(score int) Criteria {
	switch {
	case score == 100:
		return CriteriaScore100
	case score >= 95:
		return CriteriaScore95
	case score == 0:
		return CriteriaLose
	default:
		return CriteriaNone
	}
}

func calculateAchievedStoryMissions(user *User, score int) []Criteria {
	if user == nil {
		return nil
	}

	var earnedStoryMissionCriteria []Criteria

	// Score-earned story missions
	scoreStoryMissionsCriteria := scoreStoryMissions(score)
	if scoreStoryMissionsCriteria != CriteriaNone {
		earnedStoryMissionCriteria = append(earnedStoryMissionCriteria, scoreStoryMissionsCriteria)
	}

	// total days played-earned story missions
	daysPlayedStoryMissionCriteria := daysPlayedStoryMissions(user.TotalDaysPlayed)
	if daysPlayedStoryMissionCriteria != CriteriaNone {
		earnedStoryMissionCriteria = append(earnedStoryMissionCriteria, daysPlayedStoryMissionCriteria)
	}

	// current daily streak-earned story missions
	currentDailyStreakStoryMissionsCriteria := currentDailyStreakStoryMissions(user.CurrentDailyStreak)
	if currentDailyStreakStoryMissionsCriteria != CriteriaNone {
		earnedStoryMissionCriteria = append(earnedStoryMissionCriteria, currentDailyStreakStoryMissionsCriteria)
	}

	// total wins-earned story missions
	totalWinsStoryMissionsCriteria := totalWinsStoryMissions(user.TotalWins)
	if totalWinsStoryMissionsCriteria != CriteriaNone {
		earnedStoryMissionCriteria = append(earnedStoryMissionCriteria, totalWinsStoryMissionsCriteria)
	}

	return earnedStoryMissionCriteria
}

func updateStoryMissions(user *User, todayDate string, playerName string, earnedStoryMissionsCriteria []Criteria) []Criteria {
	var filteredCritiera []Criteria

	for _, criteria := range earnedStoryMissionsCriteria {
		for i, mission := range user.StoryMissions {
			if mission.Criteria == criteria && mission.DateAchieved == "" {
				user.StoryMissions[i].DateAchieved = todayDate
				user.StoryMissions[i].PlayerName = playerName

				if !isDaysPlayedCriteria(criteria) {
					filteredCritiera = append(filteredCritiera, criteria)
				}
				break
			}
		}
	}

	return filteredCritiera
}

func isDaysPlayedCriteria(c Criteria) bool {
	switch c {
	case CriteriaPlay1Day, CriteriaPlay2Days, CriteriaPlay3Days, CriteriaPlay4Days,
		CriteriaPlay5Days, CriteriaPlay6Days, CriteriaPlay7Days, CriteriaPlay8Days,
		CriteriaPlay9Days, CriteriaPlay10Days, CriteriaPlay11Days, CriteriaPlay12Days,
		CriteriaPlay13Days, CriteriaPlay14Days:
		return true
	default:
		return false
	}
}
