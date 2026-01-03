//go:build lambda
// +build lambda

package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func main() {
	// Setup the Gin router
	router := SetupRouter()

	// Initialize the gin lambda adapter
	ginLambda = ginadapter.New(router)

	log.Println("Starting AWS Lambda handler")
	log.Printf("API endpoints available at /v1/*")

	// Start Lambda handler
	lambda.Start(Handler)
}

// Handler is the Lambda function handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
