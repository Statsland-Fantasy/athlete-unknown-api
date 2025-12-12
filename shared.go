package main

import (
	"encoding/json"
	"net/http"
	"time"
)

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
