// middleware/apikey.go
package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        if apiKey == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
            c.Abort()
            return
        }

        // In production, check against database or secure store
        validKey := os.Getenv("ADMIN_API_KEY")
        if apiKey != validKey {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
            c.Abort()
            return
        }

        // Mark as admin access
        c.Set("isAdmin", true)
        c.Next()
    }
}