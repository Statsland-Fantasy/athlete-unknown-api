package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandleGetRound tests the handleGetRound function's input validation
func TestHandleGetRound(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "missing sport parameter",
			method:         http.MethodGet,
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
		{
			name:           "invalid sport parameter",
			method:         http.MethodGet,
			queryParams:    "sport=soccer",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/round?"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			handleGetRound(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}

// TestHandleCreateRound tests the handleCreateRound function's input validation
func TestHandleCreateRound(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			body:           map[string]string{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "invalid request body",
			method:         http.MethodPost,
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_REQUEST_BODY",
		},
		{
			name:   "missing sport field",
			method: http.MethodPost,
			body: Round{
				PlayDate: "2024-01-01",
				Player: Player{
					Name: "Test Player",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_FIELD",
		},
		{
			name:   "missing playDate field",
			method: http.MethodPost,
			body: Round{
				Sport: "basketball",
				Player: Player{
					Name: "Test Player",
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_FIELD",
		},
		{
			name:   "missing player name field",
			method: http.MethodPost,
			body: Round{
				Sport:    "basketball",
				PlayDate: "2024-01-01",
				Player:   Player{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_FIELD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/v1/round", bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handleCreateRound(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}

// TestHandleDeleteRound tests the handleDeleteRound function's input validation
func TestHandleDeleteRound(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			queryParams:    "sport=basketball&playDate=2024-01-01",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "missing sport parameter",
			method:         http.MethodDelete,
			queryParams:    "playDate=2024-01-01",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
		{
			name:           "missing playDate parameter",
			method:         http.MethodDelete,
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/round?"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			handleDeleteRound(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}

// TestHandleGetUpcomingRounds tests the handleGetUpcomingRounds function's input validation
func TestHandleGetUpcomingRounds(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "missing sport parameter",
			method:         http.MethodGet,
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/upcoming-rounds?"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			handleGetUpcomingRounds(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}

// TestHandleSubmitResults tests the handleSubmitResults function's input validation
func TestHandleSubmitResults(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			queryParams:    "sport=basketball&playDate=2024-01-01",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "missing query parameters",
			method:         http.MethodPost,
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Result{Score: 100, IsCorrect: true}
			bodyBytes, _ := json.Marshal(result)

			req := httptest.NewRequest(tt.method, "/v1/results?"+tt.queryParams, bytes.NewReader(bodyBytes))
			rec := httptest.NewRecorder()

			handleSubmitResults(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}

// TestHandleGetRoundStats tests the handleGetRoundStats function's input validation
func TestHandleGetRoundStats(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			queryParams:    "sport=basketball&playDate=2024-01-01",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "missing query parameters",
			method:         http.MethodGet,
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/stats/round?"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			handleGetRoundStats(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}

// TestHandleGetUserStats tests the handleGetUserStats function's input validation
func TestHandleGetUserStats(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "method not allowed",
			method:         http.MethodPost,
			queryParams:    "userId=user123",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "missing userId parameter",
			method:         http.MethodGet,
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/stats/user?"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			handleGetUserStats(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode != "" {
				var errResp ErrorResponse
				json.NewDecoder(rec.Body).Decode(&errResp)
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}
