#!/bin/bash

# AWS Lambda Deployment Script
# Builds the Go binary and deploys to Lambda

set -e  # Exit on any error

# Load configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/config.sh" "$@"

PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

log_info "Deploying to ${ENVIRONMENT} environment..."

# Function to check if infrastructure exists
check_infrastructure() {
    log_info "Checking infrastructure..."

    INFRA_FILE="${SCRIPT_DIR}/.cache/${ENVIRONMENT}-infrastructure.json"

    if [ ! -f "$INFRA_FILE" ]; then
        log_error "Infrastructure not found for ${ENVIRONMENT} environment"
        log_info "Please run: ./aws/setup-infrastructure.sh ${ENVIRONMENT}"
        exit 1
    fi

    # Load infrastructure info
    LAMBDA_FUNCTION_NAME=$(grep -o '"lambda_function_name": "[^"]*' "$INFRA_FILE" | grep -o '[^"]*$')

    if [ -z "$LAMBDA_FUNCTION_NAME" ]; then
        log_error "Lambda function name not found"
        exit 1
    fi

    log_success "Infrastructure found"
}

# Function to build Go binary for Lambda
build_lambda() {
    log_info "Building Go binary for Lambda (Linux ARM64)..."

    cd "$PROJECT_ROOT"

    # Clean previous build
    if [ -f "bootstrap" ]; then
        rm bootstrap
    fi
    if [ -f "lambda-deployment.zip" ]; then
        rm lambda-deployment.zip
    fi

    # Install dependencies
    log_info "Installing Go dependencies..."
    go mod download
    go mod tidy

    # Build for Lambda (Amazon Linux 2023 ARM64)
    log_info "Compiling Go binary..."
    GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap lambda_main.go

    if [ ! -f "bootstrap" ]; then
        log_error "Build failed - bootstrap binary not found"
        exit 1
    fi

    # Make it executable
    chmod +x bootstrap

    # Verify binary
    file bootstrap

    log_success "Build completed successfully"
}

# Function to package Lambda deployment
package_lambda() {
    log_info "Packaging Lambda deployment..."

    cd "$PROJECT_ROOT"

    # Create deployment package
    zip -q lambda-deployment.zip bootstrap

    # Verify zip file
    if [ ! -f "lambda-deployment.zip" ]; then
        log_error "Deployment package not found"
        exit 1
    fi

    ZIP_SIZE=$(du -h lambda-deployment.zip | cut -f1)
    log_success "Deployment package created: lambda-deployment.zip (${ZIP_SIZE})"
}

# Function to upload to Lambda
upload_to_lambda() {
    log_info "Uploading to Lambda function: ${LAMBDA_FUNCTION_NAME}..."

    cd "$PROJECT_ROOT"

    # Update Lambda function code
    aws lambda update-function-code \
        --function-name "${LAMBDA_FUNCTION_NAME}" \
        --zip-file fileb://lambda-deployment.zip \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" \
        --architectures arm64 \
        > /dev/null

    log_success "Lambda code uploaded successfully"

    # Wait for update to complete
    log_info "Waiting for Lambda update to complete..."
    aws lambda wait function-updated \
        --function-name "${LAMBDA_FUNCTION_NAME}" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}"

    log_success "Lambda function updated"
}

# Function to update Lambda environment variables
update_lambda_config() {
    log_info "Updating Lambda environment variables..."

    # Load infrastructure info
    INFRA_FILE="${SCRIPT_DIR}/.cache/${ENVIRONMENT}-infrastructure.json"
    ROUNDS_TABLE=$(grep -o '"dynamodb_rounds_table": "[^"]*' "$INFRA_FILE" | grep -o '[^"]*$')
    STATS_TABLE=$(grep -o '"dynamodb_user_stats_table": "[^"]*' "$INFRA_FILE" | grep -o '[^"]*$')

    # Prompt for Auth0 configuration (or load from .env if exists)
    if [ -f "${PROJECT_ROOT}/.env.${ENVIRONMENT}" ]; then
        log_info "Loading environment variables from .env.${ENVIRONMENT}"
        export $(cat "${PROJECT_ROOT}/.env.${ENVIRONMENT}" | grep -v '^#' | xargs)
    fi

    # Set default values if not provided
    AUTH0_DOMAIN=${AUTH0_DOMAIN:-"your-domain.auth0.com"}
    AUTH0_AUDIENCE=${AUTH0_AUDIENCE:-"your-api-audience"}
    ADMIN_API_KEY=${ADMIN_API_KEY:-"change-this-in-lambda-console"}
    ALLOWED_ORIGINS=${ALLOWED_ORIGINS:-"*"}

    # Update Lambda environment variables
    aws lambda update-function-configuration \
        --function-name "${LAMBDA_FUNCTION_NAME}" \
        --environment "Variables={
            GIN_MODE=release,
            AWS_REGION=${AWS_REGION},
            ROUNDS_TABLE_NAME=${ROUNDS_TABLE},
            USER_STATS_TABLE_NAME=${STATS_TABLE},
            AUTH0_DOMAIN=${AUTH0_DOMAIN},
            AUTH0_AUDIENCE=${AUTH0_AUDIENCE},
            ADMIN_API_KEY=${ADMIN_API_KEY},
            ALLOWED_ORIGINS=${ALLOWED_ORIGINS}
        }" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" \
        > /dev/null

    log_success "Lambda configuration updated"

    # Wait for configuration update to complete
    log_info "Waiting for configuration update to complete..."
    aws lambda wait function-updated \
        --function-name "${LAMBDA_FUNCTION_NAME}" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}"
}

# Function to run smoke tests
run_smoke_tests() {
    local api_endpoint=$1

    log_info "Running smoke tests..."

    # Wait a few seconds for deployment to propagate
    sleep 3

    # Test 1: Health check
    log_info "Test 1: Health check..."
    HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${api_endpoint}/health")

    if [ "$HEALTH_STATUS" -eq 200 ]; then
        log_success "✓ Health check passed (HTTP ${HEALTH_STATUS})"
    else
        log_error "✗ Health check failed (HTTP ${HEALTH_STATUS})"
        return 1
    fi

    # Test 2: Root endpoint
    log_info "Test 2: Root endpoint..."
    ROOT_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${api_endpoint}/")

    if [ "$ROOT_STATUS" -eq 200 ]; then
        log_success "✓ Root endpoint accessible (HTTP ${ROOT_STATUS})"
    else
        log_error "✗ Root endpoint failed (HTTP ${ROOT_STATUS})"
        return 1
    fi

    # Test 3: Get round endpoint (may return 400 without params, but should not 500)
    log_info "Test 3: API v1 endpoint..."
    API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "${api_endpoint}/v1/round?sport=basketball&playDate=2025-01-01")

    if [ "$API_STATUS" -eq 200 ] || [ "$API_STATUS" -eq 400 ] || [ "$API_STATUS" -eq 404 ]; then
        log_success "✓ API endpoint responding (HTTP ${API_STATUS})"
    else
        log_error "✗ API endpoint error (HTTP ${API_STATUS})"
        return 1
    fi

    # Test 4: CORS headers
    log_info "Test 4: CORS headers..."
    CORS_HEADERS=$(curl -s -I -X OPTIONS "${api_endpoint}/v1/round" | grep -i "access-control-allow")

    if [ -n "$CORS_HEADERS" ]; then
        log_success "✓ CORS headers present"
    else
        log_warning "⚠ CORS headers not found (may need configuration)"
    fi

    log_success "Smoke tests completed!"
}

# Function to clean up build artifacts
cleanup() {
    log_info "Cleaning up build artifacts..."

    cd "$PROJECT_ROOT"

    if [ -f "bootstrap" ]; then
        rm bootstrap
    fi
    if [ -f "lambda-deployment.zip" ]; then
        rm lambda-deployment.zip
    fi

    log_success "Cleanup complete"
}

# Main execution
main() {
    log_info "Starting deployment process..."

    # Check infrastructure
    check_infrastructure

    # Build Lambda binary
    build_lambda

    # Package deployment
    package_lambda

    # Upload to Lambda
    upload_to_lambda

    # Update Lambda configuration
    update_lambda_config

    # Get deployment URL
    INFRA_FILE="${SCRIPT_DIR}/.cache/${ENVIRONMENT}-infrastructure.json"
    API_ENDPOINT=$(grep -o '"api_endpoint": "[^"]*' "$INFRA_FILE" | grep -o '[^"]*$')

    log_success "=== Deployment Complete ==="
    log_info "Environment: ${ENVIRONMENT}"
    log_info "Lambda Function: ${LAMBDA_FUNCTION_NAME}"
    log_info "API Endpoint: ${API_ENDPOINT}"
    echo ""

    # Run smoke tests
    log_info "Running smoke tests..."
    if run_smoke_tests "${API_ENDPOINT}"; then
        log_success "All tests passed!"
    else
        log_warning "Some tests failed. Please check manually."
    fi

    # Clean up
    cleanup

    echo ""
    log_info "Next steps:"
    log_info "1. Visit ${API_ENDPOINT}/health to verify deployment"
    log_info "2. Update frontend .env.dev with: REACT_APP_API_BASE_URL=${API_ENDPOINT}"
    log_info "3. Configure environment variables in Lambda console if needed:"
    log_info "   - AUTH0_DOMAIN, AUTH0_AUDIENCE"
    log_info "   - ADMIN_API_KEY"
    log_info "   - ALLOWED_ORIGINS"
    log_info ""
    log_info "Lambda Console: https://console.aws.amazon.com/lambda/home?region=${AWS_REGION}#/functions/${LAMBDA_FUNCTION_NAME}"

    # Save deployment info
    cat > "${SCRIPT_DIR}/.cache/${ENVIRONMENT}-last-deployment.json" <<EOF
{
    "environment": "${ENVIRONMENT}",
    "deployed_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
    "api_endpoint": "${API_ENDPOINT}",
    "lambda_function_name": "${LAMBDA_FUNCTION_NAME}"
}
EOF
}

# Run main function
main
