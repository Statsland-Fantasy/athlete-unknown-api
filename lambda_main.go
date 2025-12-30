package main

import (
	"athlete-unknown-api/middleware"
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

// init is called once during Lambda cold start
func init() {
	log.Println("Lambda cold start - initializing...")

	// Load configuration
	cfg := LoadConfig()

	// Initialize DynamoDB client
	db, err := NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}
	log.Printf("DynamoDB client initialized (Rounds Table: %s, User Stats Table: %s, Region: %s)",
		cfg.RoundsTableName, cfg.UserStatsTableName, cfg.AWSRegion)

	// Create server with database dependency injection
	server := NewServer(db)

	// Set Gin mode to release for Lambda
	gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	router := gin.New()

	// Use recovery middleware
	router.Use(gin.Recovery())

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
	router.GET("/health", handleHealth)

	// Root endpoint
	router.GET("/", handleHome)

	// Create Lambda adapter for Gin
	ginLambda = ginadapter.New(router)

	log.Println("Lambda initialization complete")
}

// Handler is the Lambda function handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log request for debugging (remove in production or use structured logging)
	log.Printf("Lambda invoked: %s %s", req.HTTPMethod, req.Path)

	// Proxy the request to Gin
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Check if running in Lambda environment
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Running in Lambda
		lambda.Start(Handler)
	} else {
		// Running locally - use standard HTTP server
		log.Println("Not running in Lambda environment, starting HTTP server...")
		log.Println("To run as Lambda locally, use AWS SAM or Lambda emulator")

		// You could start the regular HTTP server here for local testing
		// Or just exit with a message
		log.Fatal("Please use 'go run main.go' for local HTTP server or deploy to Lambda")
	}
}
