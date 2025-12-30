#!/bin/bash

# AWS Deployment Configuration for Lambda API
# This file contains environment-specific settings for deploying to AWS

# Environment (dev or prod)
ENVIRONMENT=${1:-dev}

# AWS Configuration
AWS_REGION="us-west-2"
AWS_PROFILE="default"  # Change if using named AWS profiles

# Project Configuration
PROJECT_NAME="statsland"
APP_NAME="athlete-unknown-api"

# Lambda Configuration
if [ "$ENVIRONMENT" == "prod" ]; then
    LAMBDA_FUNCTION_NAME="${PROJECT_NAME}-${APP_NAME}-prod"
    API_GATEWAY_NAME="${PROJECT_NAME}-${APP_NAME}-prod"
    DYNAMODB_ROUNDS_TABLE="AthleteUnknownRounds"
    DYNAMODB_USER_STATS_TABLE="AthleteUnknownUserStats"
else
    LAMBDA_FUNCTION_NAME="${PROJECT_NAME}-${APP_NAME}-dev"
    API_GATEWAY_NAME="${PROJECT_NAME}-${APP_NAME}-dev"
    DYNAMODB_ROUNDS_TABLE="AthleteUnknownRoundsDev"
    DYNAMODB_USER_STATS_TABLE="AthleteUnknownUserStatsDev"
fi

# Lambda Runtime Configuration
LAMBDA_RUNTIME="provided.al2023"  # Custom runtime for Go
LAMBDA_HANDLER="bootstrap"  # Standard name for Go Lambda
LAMBDA_MEMORY="512"  # MB
LAMBDA_TIMEOUT="30"  # seconds

# IAM Role Names
LAMBDA_ROLE_NAME="${LAMBDA_FUNCTION_NAME}-role"

# API Gateway Configuration
API_GATEWAY_STAGE="${ENVIRONMENT}"

# Tags for all resources
TAG_PROJECT="${PROJECT_NAME}"
TAG_ENVIRONMENT="${ENVIRONMENT}"
TAG_MANAGED_BY="script"

# Output colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}
