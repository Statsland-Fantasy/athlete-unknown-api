# AWS Lambda Deployment Scripts

Quick reference for deploying the Athlete Unknown API to AWS Lambda.

## Quick Start

```bash
# 1. Initial setup (one-time)
./aws/setup-infrastructure.sh dev

# 2. Deploy your API
./aws/deploy.sh dev
```

## Available Scripts

### `setup-infrastructure.sh [environment]`
Creates AWS infrastructure (IAM roles, Lambda, API Gateway)

**Usage:**
```bash
./aws/setup-infrastructure.sh dev   # Development
./aws/setup-infrastructure.sh prod  # Production
```

**What it creates:**
- IAM role with DynamoDB read/write permissions
- Lambda function (Go ARM64)
- API Gateway HTTP API
- Lambda-API Gateway integration

**Time:** ~2 minutes

**Output:**
- API Endpoint URL
- Lambda function ARN
- Configuration saved to `.cache/{env}-infrastructure.json`

---

### `deploy.sh [environment]`
Builds and deploys Go API to Lambda

**Usage:**
```bash
./aws/deploy.sh dev   # Deploy to development
./aws/deploy.sh prod  # Deploy to production
```

**What it does:**
1. Compiles Go binary for Linux ARM64
2. Creates deployment package (zip)
3. Uploads to Lambda
4. Updates environment variables from `.env.{env}`
5. Runs automated smoke tests

**Time:** ~3-5 minutes

**Prerequisites:**
- Go 1.21+ installed
- Infrastructure set up (`setup-infrastructure.sh`)
- `.env.dev` or `.env.prod` configured

---

## Configuration

### Environment Files

Create `.env.dev` (gitignored):
```bash
# CORS - Add your frontend URL
ALLOWED_ORIGINS=https://your-cloudfront-url.cloudfront.net,http://localhost:3000

# Auth0
AUTH0_DOMAIN=dev-2l3xftm16ho266qq.us.auth0.com
AUTH0_AUDIENCE=https://api.statslandfantasy.com

# Admin API Key
ADMIN_API_KEY=your-secure-api-key
```

### AWS Configuration

Edit `aws/config.sh` to customize:
- AWS region
- Lambda memory/timeout
- Table names
- Function names

---

## Common Tasks

### Deploy after code changes
```bash
./aws/deploy.sh dev
```

### Check deployment status
```bash
# View infrastructure details
cat aws/.cache/dev-infrastructure.json

# View last deployment
cat aws/.cache/dev-last-deployment.json

# Test API
API_ENDPOINT=$(cat aws/.cache/dev-infrastructure.json | grep api_endpoint | cut -d'"' -f4)
curl $API_ENDPOINT/health
```

### Update environment variables

**Option 1:** Edit `.env.dev` and redeploy
```bash
vim .env.dev
./aws/deploy.sh dev
```

**Option 2:** Update directly in Lambda console
```
https://console.aws.amazon.com/lambda/home?region=us-west-2
→ Select function → Configuration → Environment variables
```

### View logs
```bash
# Tail logs in real-time
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev --follow

# View recent errors
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev \
  --since 1h \
  --filter-pattern "ERROR"
```

### Test locally before deploying

```bash
# Run local HTTP server
go run main.go

# Or use Docker to simulate Lambda
# (requires SAM CLI or Docker)
```

---

## Troubleshooting

### "Infrastructure not found" error
```bash
# Run setup first
./aws/setup-infrastructure.sh dev
```

### Build fails
```bash
# Ensure you're in project root
cd /path/to/athlete-unknown-api

# Clean and rebuild
rm -f bootstrap lambda-deployment.zip
./aws/deploy.sh dev
```

### API returns 500 error
```bash
# Check logs for errors
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev --since 5m

# Common issues:
# - DynamoDB table doesn't exist
# - Environment variables not set
# - IAM permissions missing
```

### CORS errors
```bash
# Update ALLOWED_ORIGINS in .env.dev
# Include your frontend CloudFront URL
./aws/deploy.sh dev
```

---

## Architecture

```
API Gateway (HTTPS)
    ↓
Lambda Function (Go + Gin)
    ↓
DynamoDB (AthleteUnknownRoundsDev, AthleteUnknownUserStatsDev)
```

**Key Features:**
- HTTP API Gateway (cheaper than REST API)
- Lambda ARM64/Graviton2 (20% cheaper)
- Automatic scaling (0 to 1000s of requests)
- Pay-per-use pricing (~$1-2/month for dev)

---

## Documentation

- **[AWS_LAMBDA_DEPLOYMENT.md](../AWS_LAMBDA_DEPLOYMENT.md)** - Full deployment guide
- **[.env.example](../.env.example)** - Environment variable template

---

## Cost

**Development:**
- ~$0.50-1/month (mostly DynamoDB)
- First 12 months likely FREE (AWS free tier)

**Production (10k DAU):**
- ~$12-18/month

See [AWS_LAMBDA_DEPLOYMENT.md - Cost Estimate](../AWS_LAMBDA_DEPLOYMENT.md#cost-estimate) for details.

---

## Next Steps

1. ✅ Deploy Lambda: `./aws/deploy.sh dev`
2. ⬜ Update frontend with API endpoint
3. ⬜ Test end-to-end flow
4. ⬜ Set up CloudWatch alarms
5. ⬜ Configure custom domain (optional)
