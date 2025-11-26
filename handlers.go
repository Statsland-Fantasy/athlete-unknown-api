package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

// Global DB instance
var db *DB

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

	ctx := context.Background()
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	if round == nil {
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

	ctx := context.Background()
	err := db.CreateRound(ctx, &round)
	if err != nil {
		if err.Error() == "round already exists" {
			errorResponseWithCode(w, "Conflict", "Round already exists for sport '"+round.Sport+"' on playDate '"+round.PlayDate+"'", "ROUND_ALREADY_EXISTS", http.StatusConflict)
			return
		}
		errorResponseWithCode(w, "Internal Server Error", "Failed to create round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

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

	ctx := context.Background()

	// Check if the round exists first
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to check round existence: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}
	if round == nil {
		errorResponseWithCode(w, "Not Found", "Round not found for sport '"+sport+"' on playDate '"+playDate+"'", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	err = db.DeleteRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to delete round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

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

	ctx := context.Background()
	upcomingRounds, err := db.GetRoundsBySport(ctx, sport, startDate, endDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve rounds: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

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

	ctx := context.Background()

	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}
	if round == nil {
		errorResponseWithCode(w, "Not Found", "Round not found for sport '"+sport+"' on date '"+playDate+"'", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	// Update round statistics
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

	// Save the updated round
	err = db.UpdateRound(ctx, round)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to update round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

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

	ctx := context.Background()
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	if round == nil {
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

	ctx := context.Background()
	stats, err := db.GetUserStats(ctx, userID)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve user stats: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	if stats == nil {
		errorResponseWithCode(w, "Not Found", "No statistics found for user '"+userID+"'", "USER_STATS_NOT_FOUND", http.StatusNotFound)
		return
	}

	jsonResponse(w, stats, http.StatusOK)
}
