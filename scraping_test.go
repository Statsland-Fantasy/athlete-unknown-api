package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// TestScrapeError tests the scrapeError type
func TestScrapeError(t *testing.T) {
	tests := []struct {
		name     string
		err      *scrapeError
		expected string
	}{
		{
			name: "error with wrapped error",
			err: &scrapeError{
				StatusCode: 500,
				Message:    "Failed to scrape",
				ErrorCode:  ErrorScrapingError,
				Err:        fmt.Errorf("connection timeout"),
			},
			expected: "Failed to scrape: connection timeout",
		},
		{
			name: "error without wrapped error",
			err: &scrapeError{
				StatusCode: 400,
				Message:    "Invalid parameter",
				ErrorCode:  ErrorInvalidParameter,
				Err:        nil,
			},
			expected: "Invalid parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("scrapeError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestParseAndValidateScrapeParams tests parameter extraction and validation
func TestParseAndValidateScrapeParams(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
		shouldSucceed  bool
	}{
		{
			name:           "valid params with name",
			queryParams:    "sport=basketball&playDate=2024-01-15&name=LeBron+James",
			expectedStatus: 0,
			expectedCode:   "",
			shouldSucceed:  true,
		},
		{
			name:           "valid params with sportsReferenceURL",
			queryParams:    "sport=baseball&playDate=2024-06-01&sportsReferenceURL=https://www.baseball-reference.com/players/t/troutmi01.shtml",
			expectedStatus: 0,
			expectedCode:   "",
			shouldSucceed:  true,
		},
		{
			name:           "missing sport parameter",
			queryParams:    "playDate=2024-01-15&name=Test+Player",
			expectedStatus: 400,
			expectedCode:   ErrorMissingRequiredParameter,
			shouldSucceed:  false,
		},
		{
			name:           "missing playDate parameter",
			queryParams:    "sport=basketball&name=Test+Player",
			expectedStatus: 400,
			expectedCode:   ErrorMissingRequiredParameter,
			shouldSucceed:  false,
		},
		{
			name:           "invalid sport parameter",
			queryParams:    "sport=soccer&playDate=2024-01-15&name=Test+Player",
			expectedStatus: 400,
			expectedCode:   ErrorInvalidParameter,
			shouldSucceed:  false,
		},
		{
			name:           "missing both name and sportsReferenceURL",
			queryParams:    "sport=basketball&playDate=2024-01-15",
			expectedStatus: 400,
			expectedCode:   ErrorMissingRequiredParameter,
			shouldSucceed:  false,
		},
		{
			name:           "valid params with theme",
			queryParams:    "sport=football&playDate=2024-09-01&name=Patrick+Mahomes&theme=dark",
			expectedStatus: 0,
			expectedCode:   "",
			shouldSucceed:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/scrape?"+tt.queryParams, nil)

			params, err := parseAndValidateScrapeParams(c)

			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if params == nil {
					t.Errorf("Expected params to be non-nil")
				}
			} else {
				if err == nil {
					t.Errorf("Expected error with code %s, got nil", tt.expectedCode)
				} else {
					if err.StatusCode != tt.expectedStatus {
						t.Errorf("Expected status %d, got %d", tt.expectedStatus, err.StatusCode)
					}
					if err.ErrorCode != tt.expectedCode {
						t.Errorf("Expected error code %s, got %s", tt.expectedCode, err.ErrorCode)
					}
				}
			}
		})
	}
}

// TestResolvePlayerURL tests URL resolution logic
func TestResolvePlayerURL(t *testing.T) {
	tests := []struct {
		name          string
		params        *scrapeParams
		shouldSucceed bool
		expectedCode  string
	}{
		{
			name: "valid direct URL",
			params: &scrapeParams{
				Sport:              SportBaseball,
				PlayDate:           "2024-06-01",
				Name:               "",
				SportsReferenceURL: "https://www.baseball-reference.com/players/t/troutmi01.shtml",
				Hostname:           "baseball-reference.com",
			},
			shouldSucceed: true,
			expectedCode:  "",
		},
		{
			name: "invalid direct URL - wrong domain",
			params: &scrapeParams{
				Sport:              SportBaseball,
				PlayDate:           "2024-06-01",
				Name:               "",
				SportsReferenceURL: "https://malicious-site.com/players/test",
				Hostname:           "baseball-reference.com",
			},
			shouldSucceed: false,
			expectedCode:  ErrorInvalidURL,
		},
		{
			name: "invalid direct URL - wrong scheme",
			params: &scrapeParams{
				Sport:              SportBaseball,
				PlayDate:           "2024-06-01",
				Name:               "",
				SportsReferenceURL: "ftp://baseball-reference.com/players/test",
				Hostname:           "baseball-reference.com",
			},
			shouldSucceed: false,
			expectedCode:  ErrorInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := resolvePlayerURL(tt.params)

			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if url == "" {
					t.Errorf("Expected non-empty URL")
				}
			} else {
				if err == nil {
					t.Errorf("Expected error with code %s, got nil", tt.expectedCode)
				} else if err.ErrorCode != tt.expectedCode {
					t.Errorf("Expected error code %s, got %s", tt.expectedCode, err.ErrorCode)
				}
			}
		})
	}
}

// TestGetStatusText tests HTTP status code to text conversion
func TestGetStatusText(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   string
	}{
		{
			name:       "bad request",
			statusCode: 400,
			expected:   StatusBadRequest,
		},
		{
			name:       "not found",
			statusCode: 404,
			expected:   StatusNotFound,
		},
		{
			name:       "conflict",
			statusCode: 409,
			expected:   StatusConflict,
		},
		{
			name:       "internal server error",
			statusCode: 500,
			expected:   StatusInternalServerError,
		},
		{
			name:       "unknown status code defaults to internal server error",
			statusCode: 418,
			expected:   StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStatusText(tt.statusCode)
			if got != tt.expected {
				t.Errorf("getStatusText(%d) = %v, want %v", tt.statusCode, got, tt.expected)
			}
		})
	}
}
