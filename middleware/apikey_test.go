package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// TestAPIKeyMiddleware tests the API key authentication middleware
func TestAPIKeyMiddleware(t *testing.T) {
	// Save original env var and restore after test
	originalKey := os.Getenv("ADMIN_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("ADMIN_API_KEY", originalKey)
		} else {
			os.Unsetenv("ADMIN_API_KEY")
		}
	}()

	tests := []struct {
		name           string
		envKeyValue    string
		headerKey      string
		expectedStatus int
		expectedError  string
		shouldAbort    bool
	}{
		{
			name:           "valid API key",
			envKeyValue:    "test-valid-key",
			headerKey:      "test-valid-key",
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name:           "missing API key header",
			envKeyValue:    "test-valid-key",
			headerKey:      "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Missing API key",
			shouldAbort:    true,
		},
		{
			name:           "invalid API key",
			envKeyValue:    "test-valid-key",
			headerKey:      "wrong-key",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid API key",
			shouldAbort:    true,
		},
		{
			name:           "missing env variable",
			envKeyValue:    "",
			headerKey:      "some-key",
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Server configuration error",
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variable
			if tt.envKeyValue != "" {
				os.Setenv("ADMIN_API_KEY", tt.envKeyValue)
			} else {
				os.Unsetenv("ADMIN_API_KEY")
			}

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			// Set header if provided
			if tt.headerKey != "" {
				c.Request.Header.Set("X-API-Key", tt.headerKey)
			}

			// Create middleware and execute
			middleware := APIKeyMiddleware()
			middleware(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check if request was aborted
			if tt.shouldAbort && !c.IsAborted() {
				t.Errorf("Expected request to be aborted, but it wasn't")
			}
		})
	}
}
