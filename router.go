package main

import (
	"athlete-unknown-api/middleware"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// SetupRouter initializes and configures the Gin router with all routes and middleware
func SetupRouter() *gin.Engine {
	// Load .env file (ignore error if file doesn't exist, for production environments)
	_ = godotenv.Load()

	// Load configuration
	cfg := LoadConfig()

	// Initialize DynamoDB client
	var err error
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}
	log.Printf("DynamoDB client initialized (Rounds Table: %s, User Stats Table: %s, Region: %s)",
		cfg.RoundsTableName, cfg.UserStatsTableName, cfg.AWSRegion)

	// Create server with database dependency injection
	server := NewServer(db)

	// Set Gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware with environment-based configuration
	allowedOrigins := GetAllowedCORSOrigins()
	log.Printf("CORS allowed origins: %v", allowedOrigins)

	corsConfig := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// API v1 routes
	v1 := router.Group("/v1")

	// Public endpoints (no auth. Available for guest users too)
	public := v1.Group("")
	public.Use(middleware.OptionalJWTMiddleware())
	{
		public.GET("/round", server.GetRound)
		public.GET("/stats/round", server.GetRoundStats)
		public.POST("/results", server.SubmitResults)
	}

	// Public endpoints (with required JWT auth for authenticated users)
	publicAuth := v1.Group("")
	publicAuth.Use(middleware.JWTMiddleware())
	{
		publicAuth.GET("/stats/user", middleware.RequirePermission("read:athlete-unknown:user-stats"), server.GetUserStats)
		publicAuth.GET("/upcoming-rounds", middleware.RequirePermission("read:athlete-unknown:upcoming-rounds"), server.GetUpcomingRounds)
	}

	// Admin endpoints (API key auth)
	admin := v1.Group("")
	admin.Use(middleware.APIKeyMiddleware())
	{
		admin.PUT("/round", server.CreateRound)
		admin.POST("/round", server.ScrapeAndCreateRound)
		admin.DELETE("/round", server.DeleteRound)
	}

	// Health check
	router.GET("/health", HandleHealth)

	// Root endpoint
	router.GET("/", HandleHome)

	return router
}

// HandleHome handles the root endpoint
func HandleHome(c *gin.Context) {
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

// HandleHealth handles health check endpoint
func HandleHealth(c *gin.Context) {
	c.JSON(200, map[string]string{"status": "healthy"})
}
