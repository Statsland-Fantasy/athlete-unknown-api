package main

import (
	"net/http"
)

// handleRoundRouter routes /v1/round based on HTTP method
func handleRoundRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetRound(w, r)
	case http.MethodPost:
		handleCreateRound(w, r)
	case http.MethodDelete:
		handleDeleteRound(w, r)
	default:
		errorResponseWithCode(w, "Method Not Allowed", "Method "+r.Method+" is not allowed for this endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
	}
}

// handleHome handles the root endpoint
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "" {
		http.NotFound(w, r)
		return
	}

	response := map[string]interface{}{
		"message": "Welcome to the Athlete Unknown Trivia Game API",
		"version": "1.0.0",
		"endpoints": []string{
			"GET /health",
			"GET /v1/round?sport={sport}&playDate={date}",
			"POST /v1/round",
			"DELETE /v1/round?sport={sport}&playDate={date}",
			"GET /v1/upcoming-rounds?sport={sport}&startDate={date}&endDate={date}",
			"POST /v1/results?sport={sport}&playDate={date}",
			"GET /v1/stats/round?sport={sport}&playDate={date}",
			"GET /v1/stats/user?userId={userId}",
		},
	}
	jsonResponse(w, response, http.StatusOK)
}

// handleHealth handles health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	jsonResponse(w, map[string]string{"status": "healthy"}, http.StatusOK)
}
