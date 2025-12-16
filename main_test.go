package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHome(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "GET request to root",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "GET request to non-existent path",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			checkResponse:  false,
		},
		{
			name:           "POST request to root",
			method:         http.MethodPost,
			path:           "/",
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			handleHome(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.checkResponse {
				var response map[string]interface{}
				if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
					t.Errorf("Failed to decode response: %v", err)
					return
				}

				if message, ok := response["message"].(string); !ok || message == "" {
					t.Errorf("Expected message in response")
				}

				if version, ok := response["version"].(string); !ok || version == "" {
					t.Errorf("Expected version in response")
				}
			}
		})
	}
}

func TestHandleHealth(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		{
			name:           "POST request not allowed",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
		{
			name:           "DELETE request not allowed",
			method:         http.MethodDelete,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   "METHOD_NOT_ALLOWED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			rec := httptest.NewRecorder()

			handleHealth(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedCode == "" {
				// Check for healthy status
				var response map[string]string
				if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
					t.Errorf("Failed to decode response: %v", err)
					return
				}
				if status, ok := response["status"]; !ok || status != "healthy" {
					t.Errorf("Expected status 'healthy', got %v", status)
				}
			} else {
				// Check for error code
				var errResp ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
					t.Errorf("Failed to decode error response: %v", err)
					return
				}
				if errResp.Code != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, errResp.Code)
				}
			}
		})
	}
}

func TestHandleRoundRouter(t *testing.T) {
	// Note: This test only validates routing logic, not the actual handler implementations
	tests := []struct {
		name           string
		method         string
		expectedStatus int // We expect bad request because we don't have DB initialized
	}{
		{
			name:           "GET method routes to handleGetRound",
			method:         http.MethodGet,
			expectedStatus: http.StatusBadRequest, // Missing required params
		},
		{
			name:           "PUT method routes to handleCreateRound",
			method:         http.MethodPut,
			expectedStatus: http.StatusBadRequest, // Invalid body
		},
		{
			name:           "POST method routes to handleScrapeAndCreateRound",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest, // Missing required params
		},
		{
			name:           "DELETE method routes to handleDeleteRound",
			method:         http.MethodDelete,
			expectedStatus: http.StatusBadRequest, // Missing required params
		},
		{
			name:           "PATCH method not allowed",
			method:         http.MethodPatch,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/v1/round", nil)
			rec := httptest.NewRecorder()

			handleRoundRouter(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestJsonResponse(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		status         int
		expectedStatus int
	}{
		{
			name:           "simple map response",
			data:           map[string]string{"key": "value"},
			status:         http.StatusOK,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "created status",
			data:           map[string]string{"message": "created"},
			status:         http.StatusCreated,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "struct response",
			data:           Player{Name: "Test Player", Sport: "basketball"},
			status:         http.StatusOK,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			jsonResponse(rec, tt.data, tt.status)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			contentType := rec.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Verify the response can be decoded as JSON
			var result interface{}
			if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
				t.Errorf("Failed to decode JSON response: %v", err)
			}
		})
	}
}

func TestErrorResponseWithCode(t *testing.T) {
	tests := []struct {
		name           string
		error          string
		message        string
		code           string
		status         int
		expectedStatus int
	}{
		{
			name:           "bad request error",
			error:          "Bad Request",
			message:        "Missing required parameter",
			code:           "MISSING_PARAMETER",
			status:         http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found error",
			error:          "Not Found",
			message:        "Resource not found",
			code:           "NOT_FOUND",
			status:         http.StatusNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "internal server error",
			error:          "Internal Server Error",
			message:        "Database connection failed",
			code:           "DATABASE_ERROR",
			status:         http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			errorResponseWithCode(rec, tt.error, tt.message, tt.code, tt.status)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			var errResp ErrorResponse
			if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
				t.Errorf("Failed to decode error response: %v", err)
				return
			}

			if errResp.Error != tt.error {
				t.Errorf("Expected error %q, got %q", tt.error, errResp.Error)
			}
			if errResp.Message != tt.message {
				t.Errorf("Expected message %q, got %q", tt.message, errResp.Message)
			}
			if errResp.Code != tt.code {
				t.Errorf("Expected code %q, got %q", tt.code, errResp.Code)
			}
			if errResp.Timestamp.IsZero() {
				t.Errorf("Expected non-zero timestamp")
			}
		})
	}
}

func TestCorsMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "OPTIONS request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			rec := httptest.NewRecorder()

			corsMiddleware(handler).ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Check CORS headers
			if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
				t.Errorf("Expected Access-Control-Allow-Origin *, got %s", origin)
			}
			if methods := rec.Header().Get("Access-Control-Allow-Methods"); methods == "" {
				t.Errorf("Expected Access-Control-Allow-Methods header to be set")
			}
			if headers := rec.Header().Get("Access-Control-Allow-Headers"); headers == "" {
				t.Errorf("Expected Access-Control-Allow-Headers header to be set")
			}
		})
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test?param=value", nil)
	rec := httptest.NewRecorder()

	loggingMiddleware(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// The middleware should pass through to the handler
	if rec.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got %s", rec.Body.String())
	}
}
