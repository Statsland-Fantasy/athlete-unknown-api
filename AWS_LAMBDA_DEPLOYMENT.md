######  AWS Lambda Deployment Guide

Comprehensive guide for deploying the Athlete Unknown API to AWS Lambda with API Gateway.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Detailed Setup](#detailed-setup)
- [Environment Variables](#environment-variables)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Cost Estimate](#cost-estimate)

## Architecture Overview

```
Users/Frontend
  ↓
API Gateway (HTTP API)
  ↓
Lambda Function (Go)
  ↓
DynamoDB Tables (AthleteUnknownRounds*, AthleteUnknownUserStats*)
```

**Components:**
- **API Gateway HTTP API**: RESTful API endpoint with CORS, routing
- **Lambda Function**: Serverless Go application (Gin framework)
- **IAM Role**: Lambda execution role with DynamoDB permissions
- **DynamoDB**: NoSQL database (must be created separately)

## Prerequisites

1. **AWS Account** with appropriate permissions
2. **AWS CLI** installed and configured
   ```bash
   aws configure
   # Enter Access Key ID, Secret Access Key, Region (us-west-2)
   ```

3. **Go 1.21+** installed
   ```bash
   go version
   ```

4. **DynamoDB Tables** already created:
   - `AthleteUnknownRoundsDev` (for dev)
   - `AthleteUnknownUserStatsDev` (for dev)
   - See backend DynamoDB setup documentation

5. **Auth0 Account** with API configured (optional, for JWT auth)

## Quick Start

### 1. Configure Environment

```bash
# Copy environment file
cp .env.dev.example .env.dev

# Edit with your settings
# - Update ALLOWED_ORIGINS with your frontend CloudFront URL
# - Update AUTH0_DOMAIN and AUTH0_AUDIENCE
# - Set a secure ADMIN_API_KEY
```

### 2. Set Up AWS Infrastructure

```bash
chmod +x aws/setup-infrastructure.sh aws/deploy.sh
./aws/setup-infrastructure.sh dev
```

⏱️ **Time**: ~2 minutes

**This creates:**
- IAM role with DynamoDB permissions
- Lambda function (with placeholder code)
- API Gateway HTTP API
- Lambda-API Gateway integration

**Output:**
```
API Endpoint: https://abc123.execute-api.us-west-2.amazonaws.com/dev
```

### 3. Deploy API Code

```bash
./aws/deploy.sh dev
```

⏱️ **Time**: ~3-5 minutes

**This does:**
- Downloads Go dependencies
- Compiles Go binary for Linux ARM64
- Creates deployment package (zip)
- Uploads to Lambda
- Updates environment variables
- Runs smoke tests

### 4. Test Deployment

```bash
# Get your API endpoint from output or:
API_ENDPOINT=$(cat aws/.cache/dev-infrastructure.json | grep api_endpoint | cut -d'"' -f4)

# Test health check
curl $API_ENDPOINT/health

# Test root endpoint
curl $API_ENDPOINT/

# Test API endpoint
curl "$API_ENDPOINT/v1/round?sport=basketball&playDate=2025-01-01"
```

### 5. Update Frontend

Update your frontend `.env.dev`:
```bash
REACT_APP_API_BASE_URL=https://abc123.execute-api.us-west-2.amazonaws.com/dev
```

## Detailed Setup

### Environment Configuration

Create `.env.dev` file (gitignored):

```bash
# CORS - Add your frontend CloudFront URL
ALLOWED_ORIGINS=https://d111111abcdef8.cloudfront.net,http://localhost:3000

# Auth0
AUTH0_DOMAIN=dev-2l3xftm16ho266qq.us.auth0.com
AUTH0_AUDIENCE=https://api.statslandfantasy.com

# Admin API Key (change to secure random value)
ADMIN_API_KEY=your-secure-api-key-here
```

### Infrastructure Setup Details

The `setup-infrastructure.sh` script creates:

1. **IAM Role** (`statsland-athlete-unknown-api-dev-role`)
   - Assume role policy for Lambda
   - AWSLambdaBasicExecutionRole (for CloudWatch logs)
   - Custom DynamoDB policy (read/write/update/delete)

2. **Lambda Function** (`statsland-athlete-unknown-api-dev`)
   - Runtime: `provided.al2023` (custom runtime for Go)
   - Architecture: ARM64 (Graviton2, cheaper than x86)
   - Memory: 512 MB
   - Timeout: 30 seconds
   - Environment variables: Table names, region, Auth0 config

3. **API Gateway HTTP API** (`statsland-athlete-unknown-api-dev`)
   - Protocol: HTTP (cheaper than REST API)
   - CORS: Configured for all origins by default
   - Routes: Catch-all routes for all HTTP methods
   - Stage: `dev` or `prod`

4. **Lambda Permission**
   - Allows API Gateway to invoke Lambda function

### Deployment Process

The `deploy.sh` script:

1. **Validates** infrastructure exists
2. **Installs** Go dependencies (`go mod download`)
3. **Compiles** Go binary:
   ```bash
   GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap lambda_main.go
   ```
4. **Packages** binary into `lambda-deployment.zip`
5. **Uploads** to Lambda
6. **Updates** Lambda environment variables from `.env.dev`
7. **Runs** smoke tests
8. **Cleans up** build artifacts

### Lambda Code Structure

**Key files:**
- `lambda_main.go` - Lambda entry point with Gin adapter
- `main.go` - Original HTTP server (for local development)
- `handlers.go` - API route handlers
- `database.go` - DynamoDB client
- `config.go` - Configuration loading

**Lambda handler flow:**
1. Cold start: `init()` function runs once
   - Initializes DynamoDB client
   - Sets up Gin router and routes
   - Creates Lambda adapter
2. Warm invocations: `Handler()` function
   - Proxies API Gateway events to Gin
   - Returns API Gateway response

## Environment Variables

**Managed by deployment scripts** (auto-configured):
| Variable | Description | Example |
|----------|-------------|---------|
| `ROUNDS_TABLE_NAME` | DynamoDB rounds table | `AthleteUnknownRoundsDev` |
| `USER_STATS_TABLE_NAME` | DynamoDB user stats table | `AthleteUnknownUserStatsDev` |
| `AWS_REGION` | AWS region | `us-west-2` |
| `GIN_MODE` | Gin framework mode | `release` |

**Configured via .env files:**
| Variable | Description | Required |
|----------|-------------|----------|
| `ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | Yes |
| `AUTH0_DOMAIN` | Auth0 tenant domain | Yes (for JWT auth) |
| `AUTH0_AUDIENCE` | Auth0 API audience | Yes (for JWT auth) |
| `ADMIN_API_KEY` | API key for admin endpoints | Yes |

**Updating environment variables:**

**Option 1:** Edit `.env.dev` and redeploy:
```bash
./aws/deploy.sh dev
```

**Option 2:** Update directly in Lambda console:
1. Go to: https://console.aws.amazon.com/lambda/
2. Select your function
3. Configuration → Environment variables
4. Edit and Save

**Option 3:** Use AWS CLI:
```bash
aws lambda update-function-configuration \
  --function-name statsland-athlete-unknown-api-dev \
  --environment "Variables={AUTH0_DOMAIN=new-domain.auth0.com,...}"
```

## Testing

### Automated Smoke Tests

Included in `deploy.sh`:
- ✅ Health check endpoint
- ✅ Root endpoint
- ✅ API v1 endpoint
- ✅ CORS headers

### Manual Testing

```bash
API_ENDPOINT="https://your-api-id.execute-api.us-west-2.amazonaws.com/dev"

# Health check
curl $API_ENDPOINT/health

# Get round (may return 404 if no data)
curl "$API_ENDPOINT/v1/round?sport=basketball&playDate=2025-01-01"

# Submit results (guest user - no auth required)
curl -X POST "$API_ENDPOINT/v1/results?sport=basketball&playDate=2025-01-01" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "guest-123",
    "guesses": [...],
    "solved": false,
    "revealedAt": 6,
    "duration": 120
  }'

# Get user stats (requires JWT token)
curl "$API_ENDPOINT/v1/stats/user?userId=user-123" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Create round (requires admin API key)
curl -X PUT "$API_ENDPOINT/v1/round" \
  -H "X-API-Key: your-admin-api-key" \
  -H "Content-Type: application/json" \
  -d '{...round data...}'
```

### End-to-End Testing

1. **Deploy backend** to Lambda
2. **Deploy frontend** to CloudFront
3. **Update frontend** `.env.dev` with Lambda API endpoint
4. **Test from browser**:
   - Load frontend
   - Play a round (tests guest user flow)
   - Check browser console for API calls
   - Verify DynamoDB tables are updated

### Viewing Logs

**CloudWatch Logs:**
```bash
# View recent logs
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev --follow

# View logs in console
# https://console.aws.amazon.com/cloudwatch/home?region=us-west-2#logsV2:log-groups/log-group/$252Faws$252Flambda$252Fstatsland-athlete-unknown-api-dev
```

## Troubleshooting

### Build fails with "go: cannot find main module"

**Cause:** Running from wrong directory

**Solution:**
```bash
cd /path/to/athlete-unknown-api
./aws/deploy.sh dev
```

### "Infrastructure not found" error

**Cause:** Haven't run setup-infrastructure.sh yet

**Solution:**
```bash
./aws/setup-infrastructure.sh dev
```

### Lambda returns 500 Internal Server Error

**Cause:** Runtime error, missing dependencies, or configuration issue

**Solution:**
```bash
# Check CloudWatch logs
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev --since 5m

# Common issues:
# - DynamoDB table doesn't exist
# - IAM role lacks DynamoDB permissions
# - Environment variables not set correctly
```

### CORS errors in browser

**Cause:** Frontend origin not in ALLOWED_ORIGINS

**Solution:**
```bash
# Update .env.dev
ALLOWED_ORIGINS=https://d111111abcdef8.cloudfront.net,http://localhost:3000

# Redeploy
./aws/deploy.sh dev
```

### Lambda cold start is slow

**Cause:** Initial invocation initializes DynamoDB client

**Solution:**
- This is normal (1-2 seconds on cold start)
- Subsequent requests are fast (<100ms)
- For production, consider:
  - Provisioned concurrency (costs extra)
  - Keep-alive via CloudWatch Events
  - Accept cold starts (most cost-effective)

### DynamoDB "Table not found" error

**Cause:** DynamoDB tables not created or wrong table name

**Solution:**
```bash
# Verify tables exist
aws dynamodb list-tables --region us-west-2

# Check environment variables in Lambda
aws lambda get-function-configuration \
  --function-name statsland-athlete-unknown-api-dev \
  --query "Environment.Variables"
```

### "Access Denied" on DynamoDB operations

**Cause:** IAM role lacks permissions

**Solution:**
```bash
# Check IAM role policies
ROLE_NAME="statsland-athlete-unknown-api-dev-role"
aws iam list-attached-role-policies --role-name $ROLE_NAME

# Verify DynamoDB policy is attached
# If not, re-run setup-infrastructure.sh
```

### API Gateway returns 403 Forbidden

**Cause:** Lambda permission not configured

**Solution:**
```bash
# Check Lambda permissions
aws lambda get-policy --function-name statsland-athlete-unknown-api-dev

# Re-run infrastructure setup to fix
./aws/setup-infrastructure.sh dev
```

## Cost Estimate

### Development Environment

**Lambda:**
- Free tier: 1M requests/month, 400,000 GB-seconds
- After free tier: $0.20 per 1M requests
- Compute: $0.0000133334 per GB-second (ARM64)
- For 10k requests/month @ 512MB, ~1s duration:
  - Requests: $0.00 (under free tier)
  - Compute: $0.00 (under free tier)

**API Gateway HTTP API:**
- Free tier: 1M requests/month (first 12 months)
- After free tier: $1.00 per million requests
- For 10k requests/month: $0.01

**DynamoDB:**
- On-demand pricing: $1.25 per million write requests, $0.25 per million read requests
- Storage: $0.25 per GB/month
- For light usage: ~$0.50/month

**Total (dev environment):**
- First 12 months: ~$0.50/month (mostly DynamoDB)
- After 12 months: ~$1.50-2/month

### Production Environment (10k daily active users)

- Lambda requests: ~300k/month → $0.06
- Lambda compute: ~$2-3
- API Gateway: ~300k requests → $0.30
- DynamoDB: ~$10-15 (depends on data volume)

**Total (production):** ~$12-18/month

### Cost Optimization Tips

1. **Use ARM64 architecture** (Graviton2) - 20% cheaper than x86
   - Already configured in deployment scripts
2. **Use HTTP API** instead of REST API - 70% cheaper
   - Already configured
3. **Optimize Lambda memory** - Balance speed vs cost
   - Current: 512MB (good balance)
   - Monitor and adjust based on metrics
4. **Enable DynamoDB on-demand** - Pay per request, no idle costs
   - Better for variable traffic
5. **Set appropriate Lambda timeout** - Don't pay for hanging requests
   - Current: 30s (reasonable for API)

## Additional Resources

- [AWS Lambda Documentation](https://docs.aws.amazon.com/lambda/)
- [API Gateway HTTP API](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api.html)
- [Go on Lambda](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)
- [Lambda Pricing](https://aws.amazon.com/lambda/pricing/)
- [API Gateway Pricing](https://aws.amazon.com/api-gateway/pricing/)

## Next Steps

1. ✅ Deploy backend to Lambda
2. ✅ Update frontend with API endpoint
3. ⬜ Set up custom domain for API (optional)
4. ⬜ Configure CloudWatch alarms for errors
5. ⬜ Set up CI/CD pipeline (GitHub Actions)
6. ⬜ Enable API Gateway access logs
7. ⬜ Add WAF for production (optional security)

---

**Support:**
- Backend GitHub: https://github.com/Statsland-Fantasy/athlete-unknown-api
- AWS Lambda Console: https://console.aws.amazon.com/lambda/
- API Gateway Console: https://console.aws.amazon.com/apigateway/
