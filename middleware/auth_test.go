package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// TestJWTMiddleware tests the JWT authentication middleware
func TestJWTMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name:           "invalid authorization format - no Bearer prefix",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name:           "invalid JWT token",
			authHeader:     "Bearer invalid.jwt.token",
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			// Set authorization header if provided
			if tt.authHeader != "" {
				c.Request.Header.Set("Authorization", tt.authHeader)
			}

			// Create middleware and execute
			middleware := JWTMiddleware()
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

// TestRequirePermission tests the permission checking middleware
func TestRequirePermission(t *testing.T) {
	tests := []struct {
		name           string
		permissions    []string
		requiredPerm   string
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name:           "has required permission",
			permissions:    []string{"read:data", "write:data"},
			requiredPerm:   "read:data",
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name:           "missing required permission",
			permissions:    []string{"read:data"},
			requiredPerm:   "write:data",
			expectedStatus: http.StatusForbidden,
			shouldAbort:    true,
		},
		{
			name:           "no permissions in context",
			permissions:    nil,
			requiredPerm:   "read:data",
			expectedStatus: http.StatusForbidden,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			// Set permissions in context if provided
			if tt.permissions != nil {
				c.Set("permissions", tt.permissions)
			}

			// Create middleware and execute
			middleware := RequirePermission(tt.requiredPerm)
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

// TestRequireRole tests the role checking middleware
func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		roles          []string
		requiredRole   string
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name:           "has required role",
			roles:          []string{"admin", "user"},
			requiredRole:   "admin",
			expectedStatus: http.StatusOK,
			shouldAbort:    false,
		},
		{
			name:           "missing required role",
			roles:          []string{"user"},
			requiredRole:   "admin",
			expectedStatus: http.StatusForbidden,
			shouldAbort:    true,
		},
		{
			name:           "no roles in context",
			roles:          nil,
			requiredRole:   "admin",
			expectedStatus: http.StatusForbidden,
			shouldAbort:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			// Set roles in context if provided
			if tt.roles != nil {
				c.Set("roles", tt.roles)
			}

			// Create middleware and execute
			middleware := RequireRole(tt.requiredRole)
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

// TestOptionalJWTMiddleware tests the optional JWT authentication middleware
func TestOptionalJWTMiddleware(t *testing.T) {
	tests := []struct {
		name                  string
		authHeader            string
		expectedStatus        int
		shouldAbort           bool
		shouldHaveUserContext bool
	}{
		{
			name:                  "missing authorization header - allows through",
			authHeader:            "",
			expectedStatus:        http.StatusOK,
			shouldAbort:           false,
			shouldHaveUserContext: false,
		},
		{
			name:                  "invalid authorization format - allows through",
			authHeader:            "InvalidToken",
			expectedStatus:        http.StatusOK,
			shouldAbort:           false,
			shouldHaveUserContext: false,
		},
		{
			name:                  "invalid JWT token - allows through",
			authHeader:            "Bearer invalid.jwt.token",
			expectedStatus:        http.StatusOK,
			shouldAbort:           false,
			shouldHaveUserContext: false,
		},
		{
			name:                  "malformed JWT - allows through",
			authHeader:            "Bearer notajwt",
			expectedStatus:        http.StatusOK,
			shouldAbort:           false,
			shouldHaveUserContext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test context with a handler that returns 200 OK
			w := httptest.NewRecorder()
			c, router := gin.CreateTestContext(w)

			// Add a test endpoint that the middleware will protect
			router.Use(OptionalJWTMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			// Create request
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			// Set authorization header if provided
			if tt.authHeader != "" {
				c.Request.Header.Set("Authorization", tt.authHeader)
			}

			// Execute the request through the router
			router.ServeHTTP(w, c.Request)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check if request was aborted (should never abort for optional middleware)
			if tt.shouldAbort && !c.IsAborted() {
				t.Errorf("Expected request to be aborted, but it wasn't")
			}

			// Check if user context was set
			_, hasUserId := c.Get("userId")
			if tt.shouldHaveUserContext && !hasUserId {
				t.Errorf("Expected userId to be set in context, but it wasn't")
			}
			if !tt.shouldHaveUserContext && hasUserId {
				t.Errorf("Expected userId NOT to be set in context, but it was")
			}
		})
	}
}
