package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Load configuration
	cfg := LoadConfig()

	// Initialize DynamoDB client
	var err error
	db, err = NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}
	log.Printf("DynamoDB client initialized (Rounds Table: %s, User Stats Table: %s, Region: %s)",
		cfg.RoundsTableName, cfg.UserStatsTableName, cfg.AWSRegion)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// API v1 routes
	mux.HandleFunc("/v1/round", handleRoundRouter)
	mux.HandleFunc("/v1/upcoming-rounds", handleGetUpcomingRounds)
	mux.HandleFunc("/v1/results", handleSubmitResults)
	mux.HandleFunc("/v1/stats/round", handleGetRoundStats)
	mux.HandleFunc("/v1/stats/user", handleGetUserStats)

	// Health check
	mux.HandleFunc("/health", handleHealth)

	// Root endpoint
	mux.HandleFunc("/", handleHome)

	// Add CORS and logging middleware
	handler := corsMiddleware(loggingMiddleware(mux))

	log.Printf("Server starting on port %s", port)
	log.Printf("API endpoints available at /v1/*")
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

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
	if r.URL.Path != "/" {
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

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.URL.RawQuery)
		next.ServeHTTP(w, r)
	})
}

// Middleware for CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Response helper
func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Error response helper with code
func errorResponseWithCode(w http.ResponseWriter, error string, message string, code string, status int) {
	errorResp := ErrorResponse{
		Error:     error,
		Message:   message,
		Code:      code,
		Timestamp: time.Now(),
	}
	jsonResponse(w, errorResp, status)
}
