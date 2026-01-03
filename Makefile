.PHONY: build build-lambda clean test run run-lambda deploy-lambda sam-local sam-deploy

# Build the regular HTTP server
build:
	@echo "Building HTTP server..."
	go build -o bin/server .
	@echo "Build complete: bin/server"

# Build the Lambda function
build-lambda:
	@echo "Building Lambda function for Linux AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda -o lambda/bootstrap -ldflags="-s -w" .
	@echo "Creating deployment package..."
	cd lambda && zip -9 bootstrap.zip bootstrap
	@echo "Lambda deployment package created: lambda/bootstrap.zip"
	@ls -lh lambda/bootstrap lambda/bootstrap.zip

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/server
	rm -f lambda/bootstrap
	rm -f lambda/bootstrap.zip
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run the server locally (HTTP mode)
run:
	@echo "Starting server locally..."
	go run .

# Run the Lambda function locally using AWS SAM CLI
sam-local: build-lambda
	@echo "Starting Lambda function locally with SAM..."
	AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_REGION=us-west-2 sam local start-api

# Deploy using AWS SAM
sam-deploy: build-lambda
	@echo "Deploying with AWS SAM..."
	sam deploy --guided

# Deploy Lambda function using AWS CLI (requires AWS_LAMBDA_FUNCTION_NAME env var)
deploy-lambda: build-lambda
	@if [ -z "$$AWS_LAMBDA_FUNCTION_NAME" ]; then \
		echo "Error: AWS_LAMBDA_FUNCTION_NAME environment variable not set"; \
		echo "Usage: AWS_LAMBDA_FUNCTION_NAME=my-function make deploy-lambda"; \
		exit 1; \
	fi
	@echo "Deploying to Lambda function: $$AWS_LAMBDA_FUNCTION_NAME"
	aws lambda update-function-code \
		--function-name $$AWS_LAMBDA_FUNCTION_NAME \
		--zip-file fileb://lambda/bootstrap.zip
	@echo "Deployment complete"

# Local DynamoDB helpers
dynamodb-start: ## Start local DynamoDB (Docker)
	@echo "Starting local DynamoDB on :8000..."
	docker run -d -p 8000:8000 --name dynamodb-local amazon/dynamodb-local

dynamodb-stop: ## Stop local DynamoDB
	@echo "Stopping local DynamoDB..."
	docker stop dynamodb-local
	docker rm dynamodb-local

# Create local DynamoDB tables
create-local-tables: ## Create DynamoDB tables locally
	@echo "Creating local DynamoDB tables..."
	AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_REGION=us-west-2 \
	aws dynamodb create-table \
		--table-name AthleteUnknownRoundsDev \
		--attribute-definitions \
			AttributeName=sport,AttributeType=S \
			AttributeName=playDate,AttributeType=S \
		--key-schema \
			AttributeName=sport,KeyType=HASH \
			AttributeName=playDate,KeyType=RANGE \
		--billing-mode PAY_PER_REQUEST \
		--endpoint-url http://localhost:8000

	AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_REGION=us-west-2 \
	aws dynamodb create-table \
		--table-name AthleteUnknownUserStatsDev \
		--attribute-definitions \
			AttributeName=userId,AttributeType=S \
		--key-schema \
			AttributeName=userId,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST \
		--endpoint-url http://localhost:8000

# Help command
help:
	@echo "Available targets:"
	@echo "  build               - Build the HTTP server"
	@echo "  build-lambda        - Build the Lambda deployment package"
	@echo "  clean               - Remove build artifacts"
	@echo "  test                - Run tests"
	@echo "  run                 - Run server locally (HTTP mode)"
	@echo "  sam-local           - Run Lambda locally using AWS SAM CLI"
	@echo "  sam-deploy          - Deploy using AWS SAM (guided)"
	@echo "  deploy-lambda       - Deploy to existing Lambda (requires AWS_LAMBDA_FUNCTION_NAME)"
	@echo "  dynamodb-start      - Start local instance of DynamoDB on port 8000"
	@echo "  dynamodb-stop       - Stop local instance of DynamoDB"
	@echo "  create-local-tables - Create Rounds and UserStats local DynamoDB tables"
	@echo "  help                - Show this help message"
