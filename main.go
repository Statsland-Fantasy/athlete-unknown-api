package main

import (
	"athlete-unknown-api/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if file doesn't exist, for production environments)
	_ = godotenv.Load()

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

	// Set Gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})

	// API v1 routes
	v1 := router.Group("/v1")

	// Public endpoints (no auth. Available for guest users too)
	public := v1.Group("")
	{
		public.GET("/round", handleGetRound)
		public.GET("/stats/round", handleGetRoundStats)
	}

	// Public endpoints (with JWT auth for authenticated users)
	public.Use(middleware.JWTMiddleware())
	{
		public.POST("/results", middleware.RequirePermission("submit:athlete-unknown:results"), handleSubmitResults)
		public.GET("/stats/user", middleware.RequirePermission("read:athlete-unknown:user-stats"), handleGetUserStats)
		public.GET("/upcoming-rounds", middleware.RequirePermission("read:athlete-unknown:upcoming-rounds"), handleGetUpcomingRounds)
	}

	// Admin endpoints (API key auth)
	admin := v1.Group("")
	admin.Use(middleware.APIKeyMiddleware())
	{
		admin.PUT("/round", handleCreateRound)
		admin.POST("/round", handleScrapeAndCreateRound)
		admin.DELETE("/round", handleDeleteRound)
	}

	// Health check
	router.GET("/health", handleHealth)

	// Root endpoint
	router.GET("/", handleHome)

	log.Printf("Server starting on port %s", port)
	log.Printf("API endpoints available at /v1/*")
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

// handleHome handles the root endpoint
func handleHome(c *gin.Context) {
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
	c.JSON(200, response)
}

// handleHealth handles health check endpoint
func handleHealth(c *gin.Context) {
	c.JSON(200, map[string]string{"status": "healthy"})
}
