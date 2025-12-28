package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// getTestServer creates a Server instance for testing
// Note: This uses nil for the DB since these tests only validate input parameters
func getTestServer() *Server {
	return &Server{db: nil}
}

// TestHandleGetRound tests the handleGetRound function's input validation
func TestHandleGetRound(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing sport parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
		{
			name:           "invalid sport parameter",
			queryParams:    "sport=soccer",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/round?"+tt.queryParams, nil)

			getTestServer().GetRound(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp map[string]interface{}
				json.NewDecoder(w.Body).Decode(&errResp)
				if code, ok := errResp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %v", tt.expectedCode, errResp["code"])
				}
			}
		})
	}
}

// TestHandleCreateRound tests the handleCreateRound function's input validation
func TestHandleCreateRound(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "invalid request body",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_REQUEST_BODY",
		},
		{
			name: "missing sport field",
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
			name: "missing playDate field",
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
			name: "missing player name field",
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

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPut, "/v1/round", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			getTestServer().CreateRound(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp map[string]interface{}
				json.NewDecoder(w.Body).Decode(&errResp)
				if code, ok := errResp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %v", tt.expectedCode, errResp["code"])
				}
			}
		})
	}
}

// TestHandleDeleteRound tests the handleDeleteRound function's input validation
func TestHandleDeleteRound(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing sport parameter",
			queryParams:    "playDate=2024-01-01",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
		{
			name:           "missing playDate parameter",
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodDelete, "/v1/round?"+tt.queryParams, nil)

			getTestServer().DeleteRound(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp map[string]interface{}
				json.NewDecoder(w.Body).Decode(&errResp)
				if code, ok := errResp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %v", tt.expectedCode, errResp["code"])
				}
			}
		})
	}
}

// TestHandleGetUpcomingRounds tests the handleGetUpcomingRounds function's input validation
func TestHandleGetUpcomingRounds(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing sport parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/upcoming-rounds?"+tt.queryParams, nil)

			getTestServer().GetUpcomingRounds(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp map[string]interface{}
				json.NewDecoder(w.Body).Decode(&errResp)
				if code, ok := errResp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %v", tt.expectedCode, errResp["code"])
				}
			}
		})
	}
}

// TestHandleSubmitResults tests the handleSubmitResults function's input validation
func TestHandleSubmitResults(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing query parameters",
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Result{Score: 100, IsCorrect: true}
			bodyBytes, _ := json.Marshal(result)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/results?"+tt.queryParams, bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			getTestServer().SubmitResults(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp map[string]interface{}
				json.NewDecoder(w.Body).Decode(&errResp)
				if code, ok := errResp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %v", tt.expectedCode, errResp["code"])
				}
			}
		})
	}
}

// TestHandleGetRoundStats tests the handleGetRoundStats function's input validation
func TestHandleGetRoundStats(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing query parameters",
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/stats/round?"+tt.queryParams, nil)

			getTestServer().GetRoundStats(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp map[string]interface{}
				json.NewDecoder(w.Body).Decode(&errResp)
				if code, ok := errResp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %v", tt.expectedCode, errResp["code"])
				}
			}
		})
	}
}

// TestHandleGetUserStats tests the handleGetUserStats function's input validation
func TestHandleGetUserStats(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing userId parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_REQUIRED_PARAMETER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/stats/user?"+tt.queryParams, nil)

			getTestServer().GetUserStats(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedCode != "" {
				var errResp map[string]interface{}
				json.NewDecoder(w.Body).Decode(&errResp)
				if code, ok := errResp["code"].(string); !ok || code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %v", tt.expectedCode, errResp["code"])
				}
			}
		})
	}
}
