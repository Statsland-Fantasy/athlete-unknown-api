#!/bin/bash

# AWS Infrastructure Setup Script for Lambda API
# Creates IAM roles, Lambda function, and API Gateway

set -e  # Exit on any error

# Load configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/config.sh" "$@"

log_info "Setting up AWS infrastructure for ${ENVIRONMENT} environment..."
log_info "Region: ${AWS_REGION}"
log_info "Lambda Function: ${LAMBDA_FUNCTION_NAME}"

# Function to create IAM role for Lambda
create_lambda_role() {
    log_info "Creating IAM role for Lambda..."

    # Check if role already exists
    if aws iam get-role --role-name "${LAMBDA_ROLE_NAME}" --profile "${AWS_PROFILE}" 2>/dev/null; then
        log_warning "IAM role ${LAMBDA_ROLE_NAME} already exists"
        ROLE_ARN=$(aws iam get-role --role-name "${LAMBDA_ROLE_NAME}" --profile "${AWS_PROFILE}" --query "Role.Arn" --output text)
        echo "${ROLE_ARN}"
        return 0
    fi

    # Create trust policy for Lambda
    cat > /tmp/lambda-trust-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF

    # Create IAM role
    aws iam create-role \
        --role-name "${LAMBDA_ROLE_NAME}" \
        --assume-role-policy-document file:///tmp/lambda-trust-policy.json \
        --description "IAM role for ${LAMBDA_FUNCTION_NAME}" \
        --tags Key=Project,Value="${TAG_PROJECT}" Key=Environment,Value="${TAG_ENVIRONMENT}" \
        --profile "${AWS_PROFILE}"

    rm /tmp/lambda-trust-policy.json

    log_success "IAM role created: ${LAMBDA_ROLE_NAME}"

    # Attach AWS managed policy for basic Lambda execution
    aws iam attach-role-policy \
        --role-name "${LAMBDA_ROLE_NAME}" \
        --policy-arn "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole" \
        --profile "${AWS_PROFILE}"

    log_success "Attached AWSLambdaBasicExecutionRole policy"

    # Create and attach custom DynamoDB policy
    create_dynamodb_policy

    # Get role ARN
    ROLE_ARN=$(aws iam get-role --role-name "${LAMBDA_ROLE_NAME}" --profile "${AWS_PROFILE}" --query "Role.Arn" --output text)

    log_success "Lambda role ARN: ${ROLE_ARN}"
    echo "${ROLE_ARN}"
}

# Function to create DynamoDB access policy
create_dynamodb_policy() {
    log_info "Creating DynamoDB access policy..."

    POLICY_NAME="${LAMBDA_ROLE_NAME}-dynamodb-policy"

    # Check if policy already exists
    ACCOUNT_ID=$(aws sts get-caller-identity --profile "${AWS_PROFILE}" --query "Account" --output text)
    POLICY_ARN="arn:aws:iam::${ACCOUNT_ID}:policy/${POLICY_NAME}"

    if aws iam get-policy --policy-arn "${POLICY_ARN}" --profile "${AWS_PROFILE}" 2>/dev/null; then
        log_warning "DynamoDB policy already exists"
    else
        # Create DynamoDB access policy
        cat > /tmp/dynamodb-policy.json <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem",
        "dynamodb:PutItem",
        "dynamodb:UpdateItem",
        "dynamodb:DeleteItem",
        "dynamodb:Query",
        "dynamodb:Scan",
        "dynamodb:BatchGetItem",
        "dynamodb:BatchWriteItem"
      ],
      "Resource": [
        "arn:aws:dynamodb:${AWS_REGION}:${ACCOUNT_ID}:table/${DYNAMODB_ROUNDS_TABLE}",
        "arn:aws:dynamodb:${AWS_REGION}:${ACCOUNT_ID}:table/${DYNAMODB_USER_STATS_TABLE}",
        "arn:aws:dynamodb:${AWS_REGION}:${ACCOUNT_ID}:table/${DYNAMODB_ROUNDS_TABLE}/index/*",
        "arn:aws:dynamodb:${AWS_REGION}:${ACCOUNT_ID}:table/${DYNAMODB_USER_STATS_TABLE}/index/*"
      ]
    }
  ]
}
EOF

        aws iam create-policy \
            --policy-name "${POLICY_NAME}" \
            --policy-document file:///tmp/dynamodb-policy.json \
            --description "DynamoDB access policy for ${LAMBDA_FUNCTION_NAME}" \
            --profile "${AWS_PROFILE}"

        rm /tmp/dynamodb-policy.json
        log_success "DynamoDB policy created: ${POLICY_NAME}"
    fi

    # Attach policy to role
    aws iam attach-role-policy \
        --role-name "${LAMBDA_ROLE_NAME}" \
        --policy-arn "${POLICY_ARN}" \
        --profile "${AWS_PROFILE}"

    log_success "Attached DynamoDB policy to Lambda role"
}

# Function to create Lambda function (placeholder - will be updated on first deploy)
create_lambda_function() {
    local role_arn=$1

    log_info "Creating Lambda function: ${LAMBDA_FUNCTION_NAME}"

    # Check if Lambda function already exists
    if aws lambda get-function --function-name "${LAMBDA_FUNCTION_NAME}" --region "${AWS_REGION}" --profile "${AWS_PROFILE}" 2>/dev/null; then
        log_warning "Lambda function ${LAMBDA_FUNCTION_NAME} already exists"
        FUNCTION_ARN=$(aws lambda get-function --function-name "${LAMBDA_FUNCTION_NAME}" --region "${AWS_REGION}" --profile "${AWS_PROFILE}" --query "Configuration.FunctionArn" --output text)
        echo "${FUNCTION_ARN}"
        return 0
    fi

    # Create a minimal bootstrap file for initial creation
    mkdir -p /tmp/lambda-bootstrap
    cat > /tmp/lambda-bootstrap/bootstrap <<'BOOTSTRAPEOF'
#!/bin/sh
echo "Lambda function not yet deployed. Please run deploy.sh"
BOOTSTRAPEOF
    chmod +x /tmp/lambda-bootstrap/bootstrap

    # Create zip file
    cd /tmp/lambda-bootstrap
    zip -q bootstrap.zip bootstrap
    cd -

    # Wait for IAM role to propagate
    log_info "Waiting for IAM role to propagate (10 seconds)..."
    sleep 10

    # Create Lambda function with placeholder code
    FUNCTION_ARN=$(aws lambda create-function \
        --function-name "${LAMBDA_FUNCTION_NAME}" \
        --runtime "${LAMBDA_RUNTIME}" \
        --role "${role_arn}" \
        --handler "${LAMBDA_HANDLER}" \
        --zip-file fileb:///tmp/lambda-bootstrap/bootstrap.zip \
        --timeout "${LAMBDA_TIMEOUT}" \
        --memory-size "${LAMBDA_MEMORY}" \
        --environment "Variables={GIN_MODE=release,AWS_REGION=${AWS_REGION},ROUNDS_TABLE_NAME=${DYNAMODB_ROUNDS_TABLE},USER_STATS_TABLE_NAME=${DYNAMODB_USER_STATS_TABLE}}" \
        --description "${PROJECT_NAME} Athlete Unknown API - ${ENVIRONMENT}" \
        --tags Project="${TAG_PROJECT}",Environment="${TAG_ENVIRONMENT}",ManagedBy="${TAG_MANAGED_BY}" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" \
        --query "FunctionArn" \
        --output text)

    # Clean up
    rm -rf /tmp/lambda-bootstrap

    log_success "Lambda function created: ${FUNCTION_ARN}"
    log_warning "Lambda function contains placeholder code. Run deploy.sh to deploy actual code."

    echo "${FUNCTION_ARN}"
}

# Function to create API Gateway
create_api_gateway() {
    local lambda_arn=$1

    log_info "Creating API Gateway..."

    # Check if API already exists
    EXISTING_API=$(aws apigatewayv2 get-apis \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" \
        --query "Items[?Name=='${API_GATEWAY_NAME}'].ApiId" \
        --output text)

    if [ -n "$EXISTING_API" ]; then
        log_warning "API Gateway ${API_GATEWAY_NAME} already exists: ${EXISTING_API}"
        API_ID="${EXISTING_API}"
    else
        # Create HTTP API (cheaper and simpler than REST API)
        API_ID=$(aws apigatewayv2 create-api \
            --name "${API_GATEWAY_NAME}" \
            --protocol-type HTTP \
            --description "${PROJECT_NAME} Athlete Unknown API - ${ENVIRONMENT}" \
            --cors-configuration "AllowOrigins=*,AllowMethods=GET\\,POST\\,PUT\\,DELETE\\,OPTIONS,AllowHeaders=Content-Type\\,Authorization\\,X-API-Key" \
            --region "${AWS_REGION}" \
            --profile "${AWS_PROFILE}" \
            --query "ApiId" \
            --output text)

        log_success "API Gateway created: ${API_ID}"
    fi

    # Create Lambda integration
    log_info "Creating Lambda integration..."

    INTEGRATION_ID=$(aws apigatewayv2 create-integration \
        --api-id "${API_ID}" \
        --integration-type AWS_PROXY \
        --integration-uri "${lambda_arn}" \
        --payload-format-version 2.0 \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" \
        --query "IntegrationId" \
        --output text)

    log_success "Integration created: ${INTEGRATION_ID}"

    # Create default route (catch-all)
    aws apigatewayv2 create-route \
        --api-id "${API_ID}" \
        --route-key '$default' \
        --target "integrations/${INTEGRATION_ID}" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" >/dev/null || log_warning "Default route may already exist"

    # Create specific routes for better routing
    for method in GET POST PUT DELETE; do
        for path in "/{proxy+}" "/v1/{proxy+}" "/health" "/"; do
            aws apigatewayv2 create-route \
                --api-id "${API_ID}" \
                --route-key "${method} ${path}" \
                --target "integrations/${INTEGRATION_ID}" \
                --region "${AWS_REGION}" \
                --profile "${AWS_PROFILE}" >/dev/null 2>&1 || true
        done
    done

    log_success "Routes created"

    # Create deployment/stage
    log_info "Creating API stage..."

    STAGE_NAME="${API_GATEWAY_STAGE}"

    aws apigatewayv2 create-stage \
        --api-id "${API_ID}" \
        --stage-name "${STAGE_NAME}" \
        --auto-deploy \
        --description "${ENVIRONMENT} stage" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" >/dev/null 2>&1 || log_warning "Stage may already exist"

    # Grant API Gateway permission to invoke Lambda
    log_info "Granting API Gateway permission to invoke Lambda..."

    ACCOUNT_ID=$(aws sts get-caller-identity --profile "${AWS_PROFILE}" --query "Account" --output text)
    SOURCE_ARN="arn:aws:execute-api:${AWS_REGION}:${ACCOUNT_ID}:${API_ID}/*/*"

    aws lambda add-permission \
        --function-name "${LAMBDA_FUNCTION_NAME}" \
        --statement-id "apigateway-invoke-${API_ID}" \
        --action lambda:InvokeFunction \
        --principal apigateway.amazonaws.com \
        --source-arn "${SOURCE_ARN}" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" >/dev/null 2>&1 || log_warning "Permission may already exist"

    log_success "API Gateway configured"

    # Get API endpoint
    API_ENDPOINT=$(aws apigatewayv2 get-api \
        --api-id "${API_ID}" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" \
        --query "ApiEndpoint" \
        --output text)

    log_success "API Endpoint: ${API_ENDPOINT}/${STAGE_NAME}"

    echo "${API_ID}"
}

# Main execution
main() {
    log_info "Starting infrastructure setup..."

    # Create IAM role
    ROLE_ARN=$(create_lambda_role)

    # Create Lambda function
    LAMBDA_ARN=$(create_lambda_function "${ROLE_ARN}")

    # Create API Gateway
    API_ID=$(create_api_gateway "${LAMBDA_ARN}")

    # Get API endpoint
    API_ENDPOINT=$(aws apigatewayv2 get-api \
        --api-id "${API_ID}" \
        --region "${AWS_REGION}" \
        --profile "${AWS_PROFILE}" \
        --query "ApiEndpoint" \
        --output text)

    # Save infrastructure info
    mkdir -p "${SCRIPT_DIR}/.cache"
    cat > "${SCRIPT_DIR}/.cache/${ENVIRONMENT}-infrastructure.json" <<EOF
{
    "environment": "${ENVIRONMENT}",
    "region": "${AWS_REGION}",
    "lambda_function_name": "${LAMBDA_FUNCTION_NAME}",
    "lambda_arn": "${LAMBDA_ARN}",
    "lambda_role_arn": "${ROLE_ARN}",
    "api_gateway_id": "${API_ID}",
    "api_endpoint": "${API_ENDPOINT}/${API_GATEWAY_STAGE}",
    "dynamodb_rounds_table": "${DYNAMODB_ROUNDS_TABLE}",
    "dynamodb_user_stats_table": "${DYNAMODB_USER_STATS_TABLE}",
    "created_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
}
EOF

    log_success "Infrastructure setup complete!"
    log_info "Configuration saved to: ${SCRIPT_DIR}/.cache/${ENVIRONMENT}-infrastructure.json"

    echo ""
    log_success "=== Deployment Information ==="
    log_info "Environment: ${ENVIRONMENT}"
    log_info "Lambda Function: ${LAMBDA_FUNCTION_NAME}"
    log_info "Lambda ARN: ${LAMBDA_ARN}"
    log_info "API Gateway ID: ${API_ID}"
    log_info "API Endpoint: ${API_ENDPOINT}/${API_GATEWAY_STAGE}"
    log_info "DynamoDB Tables: ${DYNAMODB_ROUNDS_TABLE}, ${DYNAMODB_USER_STATS_TABLE}"
    echo ""
    log_info "Next steps:"
    log_info "1. Run './aws/deploy.sh ${ENVIRONMENT}' to deploy your API code"
    log_info "2. Update your frontend .env.dev with: REACT_APP_API_BASE_URL=${API_ENDPOINT}/${API_GATEWAY_STAGE}"
}

# Run main function
main
