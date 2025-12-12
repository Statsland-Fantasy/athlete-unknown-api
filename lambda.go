// +build lambda

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// lambdaHandler wraps the HTTP handlers for Lambda execution
func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Create HTTP request from API Gateway event
	httpReq, err := convertAPIGatewayRequestToHTTPRequest(request)
	if err != nil {
		log.Printf("Error converting request: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"Internal Server Error","message":"Failed to process request"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Create a custom response writer that captures the response
	w := &lambdaResponseWriter{
		headers:    make(http.Header),
		statusCode: http.StatusOK,
	}

	// Route the request to the appropriate handler
	routeRequest(w, httpReq)

	// Convert response to API Gateway format
	return events.APIGatewayProxyResponse{
		StatusCode: w.statusCode,
		Body:       w.body.String(),
		Headers:    flattenHeaders(w.headers),
	}, nil
}

// convertAPIGatewayRequestToHTTPRequest converts API Gateway request to http.Request
func convertAPIGatewayRequestToHTTPRequest(request events.APIGatewayProxyRequest) (*http.Request, error) {
	// Build the full path with query string
	path := request.Path
	if len(request.QueryStringParameters) > 0 {
		queryParts := []string{}
		for key, value := range request.QueryStringParameters {
			queryParts = append(queryParts, key+"="+value)
		}
		path += "?" + strings.Join(queryParts, "&")
	}

	// Create the HTTP request
	httpReq, err := http.NewRequest(
		request.HTTPMethod,
		path,
		strings.NewReader(request.Body),
	)
	if err != nil {
		return nil, err
	}

	// Copy headers
	for key, value := range request.Headers {
		httpReq.Header.Set(key, value)
	}

	return httpReq, nil
}

// lambdaResponseWriter implements http.ResponseWriter for Lambda
type lambdaResponseWriter struct {
	headers    http.Header
	body       strings.Builder
	statusCode int
}

func (w *lambdaResponseWriter) Header() http.Header {
	return w.headers
}

func (w *lambdaResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *lambdaResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// flattenHeaders converts http.Header to map[string]string
func flattenHeaders(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}

// routeRequest routes the HTTP request to the appropriate handler
func routeRequest(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route based on path
	path := r.URL.Path

	// Remove leading/trailing slashes and normalize
	path = strings.Trim(path, "/")

	log.Printf("Lambda routing: %s %s", r.Method, path)

	switch {
	case path == "health":
		handleHealth(w, r)
	case path == "v1/round":
		handleRoundRouter(w, r)
	case path == "v1/upcoming-rounds":
		handleGetUpcomingRounds(w, r)
	case path == "v1/results":
		handleSubmitResults(w, r)
	case path == "v1/stats/round":
		handleGetRoundStats(w, r)
	case path == "v1/stats/user":
		handleGetUserStats(w, r)
	case path == "" || path == "/":
		handleHome(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
		})
	}
}

func main() {
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

	// Start Lambda handler
	lambda.Start(lambdaHandler)
}
