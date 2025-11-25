package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

// In-memory storage
var (
	rounds      = make(map[string]*Round)      // key: sport_playDate
	userStats   = make(map[string]*UserStats)  // key: userId
	roundsMutex sync.RWMutex
	statsMutex  sync.RWMutex
)

func init() {
	// Initialize with sample data
	sampleRound := &Round{
		RoundID:               "Basketball100",
		Sport:                 "basketball",
		PlayDate:              "2025-11-24",
		Created:               time.Now(),
		LastUpdated:           time.Now(),
		PreviouslyPlayedDates: []string{"2025-11-01", "2025-11-08"},
		Player: Player{
			Sport:                "basketball",
			SportsReferenceURL:   "https://www.basketball-reference.com/players/j/jamesle01.html",
			Name:                 "LeBron James",
			Bio:                  "DOB: December 30, 1984 in Akron, Ohio",
			PlayerInformation:    `6'9", 250 lbs, Forward, Shoots Right`,
			DraftInformation:     "Round 1 (1st overall) from St. Vincent-St. Mary High School",
			YearsActive:          "2003-Present",
			TeamsPlayedOn:        "CLE, MIA, LAL",
			JerseyNumbers:        "#23, #6",
			CareerStats:          "PPG: 27.2, RPG: 7.5, APG: 7.3, WS: 273.5",
			PersonalAchievements: "4x NBA Champion, 4x NBA MVP, 19x NBA All-Star, 2x Olympic Gold Medalist",
			Photo:                "https://cdn.triviagame.com/players/lebron-james.jpg",
		},
		Stats: RoundStats{
			PlayDate:                   "2025-11-24",
			Name:                       "LeBron James",
			Sport:                      "basketball",
			TotalPlays:                 1247,
			PercentageCorrect:          68.5,
			HighestScore:               9,
			AverageCorrectScore:        7.8,
			MostCommonFirstTileFlipped: "tile1",
			MostCommonLastTileFlipped:  "tile9",
			MostCommonTileFlipped:      "tile5",
			LeastCommonTileFlipped:     "tile3",
		},
	}
	rounds["basketball_2025-11-24"] = sampleRound
}

// handleGetRound handles GET /v1/round
func handleGetRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	if sport == "" {
		errorResponseWithCode(w, "Bad Request", "Sport parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	// Validate sport
	if sport != "basketball" && sport != "baseball" && sport != "football" {
		errorResponseWithCode(w, "Bad Request", "Invalid sport parameter. Must be basketball, baseball, or football", "INVALID_PARAMETER", http.StatusBadRequest)
		return
	}

	playDate := r.URL.Query().Get("playDate")
	if playDate == "" {
		playDate = time.Now().Format("2006-01-02")
	}
	fmt.Println("playdateeee", playDate)

	key := sport + "_" + playDate

	roundsMutex.RLock()
	round, exists := rounds[key]
	roundsMutex.RUnlock()

	if !exists {
		errorResponseWithCode(w, "Not Found", "No round found for the specified sport and playDate", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	jsonResponse(w, round, http.StatusOK)
}

// handleCreateRound handles POST /v1/round
func handleCreateRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponseWithCode(w, "Method Not Allowed", "Only POST method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	var round Round
	if err := json.NewDecoder(r.Body).Decode(&round); err != nil {
		errorResponseWithCode(w, "Bad Request", "Invalid request body: "+err.Error(), "INVALID_REQUEST_BODY", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if round.Sport == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required field: sport", "MISSING_REQUIRED_FIELD", http.StatusBadRequest)
		return
	}
	if round.PlayDate == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required field: playDate", "MISSING_REQUIRED_FIELD", http.StatusBadRequest)
		return
	}
	if round.Player.Name == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required field: player.name", "MISSING_REQUIRED_FIELD", http.StatusBadRequest)
		return
	}

	key := round.Sport + "_" + round.PlayDate

	roundsMutex.Lock()
	if _, exists := rounds[key]; exists {
		roundsMutex.Unlock()
		errorResponseWithCode(w, "Conflict", "Round already exists for sport '"+round.Sport+"' on playDate '"+round.PlayDate+"'", "ROUND_ALREADY_EXISTS", http.StatusConflict)
		return
	}

	// Set timestamps
	now := time.Now()
	round.Created = now
	round.LastUpdated = now

	// Initialize stats if not provided
	if round.Stats.PlayDate == "" {
		round.Stats.PlayDate = round.PlayDate
		round.Stats.Name = round.Player.Name
		round.Stats.Sport = round.Sport
	}

	rounds[key] = &round
	roundsMutex.Unlock()

	jsonResponse(w, round, http.StatusCreated)
}

// handleDeleteRound handles DELETE /v1/round
func handleDeleteRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorResponseWithCode(w, "Method Not Allowed", "Only DELETE method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	if sport == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required parameter: sport", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	playDate := r.URL.Query().Get("playDate")
	if playDate == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required parameter: playDate", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	key := sport + "_" + playDate

	roundsMutex.Lock()
	if _, exists := rounds[key]; !exists {
		roundsMutex.Unlock()
		errorResponseWithCode(w, "Not Found", "Round not found for sport '"+sport+"' on playDate '"+playDate+"'", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	delete(rounds, key)
	roundsMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// handleGetUpcomingRounds handles GET /v1/upcoming-rounds
func handleGetUpcomingRounds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	if sport == "" {
		errorResponseWithCode(w, "Bad Request", "Sport parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	roundsMutex.RLock()
	var upcomingRounds []*Round
	for _, round := range rounds {
		if round.Sport != sport {
			continue
		}

		// Filter by date range if provided
		if startDate != "" && round.PlayDate < startDate {
			continue
		}
		if endDate != "" && round.PlayDate > endDate {
			continue
		}

		upcomingRounds = append(upcomingRounds, round)
	}
	roundsMutex.RUnlock()

	if len(upcomingRounds) == 0 {
		errorResponseWithCode(w, "Not Found", "No upcoming rounds found for sport '"+sport+"' in the specified date range", "NO_UPCOMING_ROUNDS", http.StatusNotFound)
		return
	}

	// Sort by playDate
	sort.Slice(upcomingRounds, func(i, j int) bool {
		return upcomingRounds[i].PlayDate < upcomingRounds[j].PlayDate
	})

	jsonResponse(w, upcomingRounds, http.StatusOK)
}

// handleSubmitResults handles POST /v1/results
func handleSubmitResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponseWithCode(w, "Method Not Allowed", "Only POST method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	playDate := r.URL.Query().Get("playDate")

	if sport == "" || playDate == "" {
		errorResponseWithCode(w, "Bad Request", "Sport and playDate parameters are required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	var result Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		errorResponseWithCode(w, "Bad Request", "Invalid request body: "+err.Error(), "INVALID_REQUEST_BODY", http.StatusBadRequest)
		return
	}

	key := sport + "_" + playDate

	roundsMutex.RLock()
	round, exists := rounds[key]
	roundsMutex.RUnlock()

	if !exists {
		errorResponseWithCode(w, "Not Found", "Round not found for sport '"+sport+"' on date '"+playDate+"'", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	// Update round statistics
	roundsMutex.Lock()
	round.Stats.TotalPlays++
	if result.IsCorrect {
		correctCount := int(round.Stats.PercentageCorrect * float64(round.Stats.TotalPlays-1) / 100)
		correctCount++
		round.Stats.PercentageCorrect = float64(correctCount) * 100 / float64(round.Stats.TotalPlays)

		// Update average correct score
		totalCorrectScore := round.Stats.AverageCorrectScore * float64(correctCount-1)
		totalCorrectScore += float64(result.Score)
		round.Stats.AverageCorrectScore = totalCorrectScore / float64(correctCount)
	}

	if result.Score > round.Stats.HighestScore {
		round.Stats.HighestScore = result.Score
	}

	round.LastUpdated = time.Now()
	roundsMutex.Unlock()

	jsonResponse(w, result, http.StatusOK)
}

// handleGetRoundStats handles GET /v1/stats/round
func handleGetRoundStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	playDate := r.URL.Query().Get("playDate")

	if sport == "" || playDate == "" {
		errorResponseWithCode(w, "Bad Request", "Sport and playDate parameters are required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	key := sport + "_" + playDate

	roundsMutex.RLock()
	round, exists := rounds[key]
	roundsMutex.RUnlock()

	if !exists {
		errorResponseWithCode(w, "Not Found", "No statistics found for sport '"+sport+"' on date '"+playDate+"'", "STATS_NOT_FOUND", http.StatusNotFound)
		return
	}

	jsonResponse(w, round.Stats, http.StatusOK)
}

// handleGetUserStats handles GET /v1/stats/user
func handleGetUserStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		errorResponseWithCode(w, "Bad Request", "userId parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	statsMutex.RLock()
	stats, exists := userStats[userID]
	statsMutex.RUnlock()

	if !exists {
		errorResponseWithCode(w, "Not Found", "No statistics found for user '"+userID+"'", "USER_STATS_NOT_FOUND", http.StatusNotFound)
		return
	}

	jsonResponse(w, stats, http.StatusOK)
}
