# Backend Lambda Deployment - Next Steps

All infrastructure code has been created! Here's what to do next.

## ✅ What's Been Completed

1. ✅ Lambda handler wrapper (`lambda_main.go`)
2. ✅ IAM role creation with DynamoDB permissions
3. ✅ API Gateway HTTP API setup
4. ✅ Automated deployment scripts
5. ✅ Environment variable management
6. ✅ Comprehensive documentation
7. ✅ All changes committed and pushed to `deployBE` branch

**Commit:** `84bafa6`
**Branch:** `deployBE`

## 📋 Deployment Steps

### Step 1: Ensure DynamoDB Tables Exist

**Required tables:**
- `AthleteUnknownRoundsDev` (for dev environment)
- `AthleteUnknownUserStatsDev` (for dev environment)

**Check if tables exist:**
```bash
aws dynamodb list-tables --region us-west-2 | grep AthleteUnknown
```

**If tables don't exist,** you'll need to create them first (check your DynamoDB setup docs).

### Step 2: Install Go Dependencies

```bash
cd /Users/nathan/Desktop/athlete-unknown-api
go mod download
go mod tidy
```

### Step 3: Configure Environment Variables

**Create `.env.dev` file** (not committed to git):
```bash
# CORS - Add your frontend CloudFront URL from frontend deployment
ALLOWED_ORIGINS=https://YOUR-CLOUDFRONT-URL.cloudfront.net,http://localhost:3000

# Auth0 Configuration
AUTH0_DOMAIN=dev-2l3xftm16ho266qq.us.auth0.com
AUTH0_AUDIENCE=https://api.statslandfantasy.com

# Admin API Key (change this!)
ADMIN_API_KEY=dev-secure-api-key-$(openssl rand -hex 16)
```

### Step 4: Set Up AWS Infrastructure

```bash
cd /Users/nathan/Desktop/athlete-unknown-api
./aws/setup-infrastructure.sh dev
```

⏱️ **Time:** ~2 minutes

**This creates:**
- IAM role: `statsland-athlete-unknown-api-dev-role`
- Lambda function: `statsland-athlete-unknown-api-dev`
- API Gateway: `statsland-athlete-unknown-api-dev`
- Lambda permissions and integrations

**Output example:**
```
API Endpoint: https://abc123xyz.execute-api.us-west-2.amazonaws.com/dev
```

**Save this API endpoint!** You'll need it for the frontend.

### Step 5: Deploy API Code

```bash
./aws/deploy.sh dev
```

⏱️ **Time:** ~3-5 minutes

**This does:**
- Compiles Go binary for Linux ARM64
- Creates deployment package
- Uploads to Lambda
- Updates environment variables
- Runs smoke tests

**Expected output:**
```
[SUCCESS] Build completed successfully
[SUCCESS] Lambda code uploaded successfully
[SUCCESS] All tests passed!
API Endpoint: https://abc123xyz.execute-api.us-west-2.amazonaws.com/dev
```

### Step 6: Test the API

```bash
# Save your API endpoint
API_ENDPOINT=$(cat aws/.cache/dev-infrastructure.json | grep api_endpoint | cut -d'"' -f4)

# Test health check
curl $API_ENDPOINT/health
# Should return: {"status":"healthy"}

# Test root endpoint
curl $API_ENDPOINT/
# Should return API info

# Test API endpoint (may return 404 if no data in DynamoDB)
curl "$API_ENDPOINT/v1/round?sport=basketball&playDate=2025-01-01"
```

### Step 7: Update Frontend Configuration

**In the statsland-website repo**, update `.env.dev`:

```bash
cd /Users/nathan/Desktop/statsland-website

# Edit .env.dev
# Change REACT_APP_API_BASE_URL to your Lambda API endpoint
REACT_APP_API_BASE_URL=https://abc123xyz.execute-api.us-west-2.amazonaws.com/dev
```

Then redeploy frontend:
```bash
./aws/deploy.sh dev
```

### Step 8: Update Frontend CORS in Backend

**Once you have your frontend CloudFront URL**, update backend `.env.dev`:

```bash
cd /Users/nathan/Desktop/athlete-unknown-api

# Edit .env.dev
ALLOWED_ORIGINS=https://YOUR-CLOUDFRONT-URL.cloudfront.net,http://localhost:3000

# Redeploy backend with updated CORS
./aws/deploy.sh dev
```

### Step 9: End-to-End Testing

1. **Open frontend** in browser (your CloudFront URL)
2. **Play a round** (test guest user flow)
3. **Check browser console** for API calls
4. **Verify no CORS errors**
5. **Check DynamoDB** tables for new data:
   ```bash
   # Check rounds table
   aws dynamodb scan --table-name AthleteUnknownRoundsDev --max-items 5

   # Check user stats table
   aws dynamodb scan --table-name AthleteUnknownUserStatsDev --max-items 5
   ```

## ⚠️ Important Notes

### Environment Variables

**Currently set via `.env.dev`:**
- `ALLOWED_ORIGINS` - CORS origins (must include frontend URL)
- `AUTH0_DOMAIN` - Auth0 tenant
- `AUTH0_AUDIENCE` - Auth0 API identifier
- `ADMIN_API_KEY` - API key for admin endpoints

**Auto-configured by deployment:**
- `ROUNDS_TABLE_NAME` - DynamoDB rounds table
- `USER_STATS_TABLE_NAME` - DynamoDB user stats table
- `AWS_REGION` - AWS region
- `GIN_MODE` - Gin framework mode (release)

### Updating Environment Variables

**Method 1:** Edit `.env.dev` and redeploy
```bash
./aws/deploy.sh dev
```

**Method 2:** Update in Lambda console
1. Go to: https://console.aws.amazon.com/lambda/
2. Select `statsland-athlete-unknown-api-dev`
3. Configuration → Environment variables → Edit

### Monitoring

**View Lambda logs:**
```bash
# Real-time logs
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev --follow

# Recent errors
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev \
  --since 1h \
  --filter-pattern "ERROR"
```

**Lambda Console:**
https://console.aws.amazon.com/lambda/home?region=us-west-2#/functions/statsland-athlete-unknown-api-dev

**API Gateway Console:**
https://console.aws.amazon.com/apigateway/main/apis?region=us-west-2

## 🚨 Troubleshooting

### "DynamoDB table not found"

**Cause:** Tables don't exist or wrong table names

**Solution:**
```bash
# List tables
aws dynamodb list-tables --region us-west-2

# If missing, create them (check DynamoDB setup docs)
```

### "Access Denied" on DynamoDB

**Cause:** IAM role lacks permissions

**Solution:**
```bash
# Re-run infrastructure setup
./aws/setup-infrastructure.sh dev
```

### CORS errors in browser

**Cause:** Frontend origin not in ALLOWED_ORIGINS

**Solution:**
```bash
# Update .env.dev with frontend CloudFront URL
ALLOWED_ORIGINS=https://YOUR-CF-URL.cloudfront.net,http://localhost:3000

# Redeploy
./aws/deploy.sh dev
```

### Lambda returns 500 error

**Cause:** Runtime error or missing dependencies

**Solution:**
```bash
# Check logs
aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev --since 10m

# Common issues:
# - DynamoDB table doesn't exist
# - Environment variables not set
# - Auth0 configuration invalid
```

### Build fails

**Cause:** Go dependencies not installed or wrong directory

**Solution:**
```bash
cd /Users/nathan/Desktop/athlete-unknown-api
go mod download
go mod tidy
./aws/deploy.sh dev
```

## 📊 Expected Costs

**Development (light usage):**
- Lambda: $0 (free tier)
- API Gateway: $0.01 (free tier covers most)
- DynamoDB: ~$0.50
- **Total: ~$0.50-1/month**

**Free tier coverage (first 12 months):**
- Lambda: 1M requests/month
- API Gateway: 1M requests/month
- DynamoDB: 25 GB storage, 25 WCU, 25 RCU

## 📚 Documentation

- **[AWS_LAMBDA_DEPLOYMENT.md](AWS_LAMBDA_DEPLOYMENT.md)** - Full deployment guide
- **[aws/README.md](aws/README.md)** - Quick script reference

## ✅ Deployment Checklist

- [ ] DynamoDB tables exist (AthleteUnknownRoundsDev, AthleteUnknownUserStatsDev)
- [ ] Go dependencies installed (`go mod download`)
- [ ] `.env.dev` created with Auth0 and CORS config
- [ ] Infrastructure set up (`./aws/setup-infrastructure.sh dev`)
- [ ] API code deployed (`./aws/deploy.sh dev`)
- [ ] API endpoint tested (`curl $API_ENDPOINT/health`)
- [ ] Frontend updated with API endpoint
- [ ] Backend CORS updated with frontend CloudFront URL
- [ ] End-to-end flow tested (browser → frontend → API → DynamoDB)
- [ ] No errors in browser console
- [ ] Data appears in DynamoDB tables

## 🎯 Success Criteria

When everything is working:
1. ✅ API health check returns 200
2. ✅ Frontend can fetch rounds from API
3. ✅ Guest users can submit results
4. ✅ No CORS errors in browser console
5. ✅ Data persists in DynamoDB
6. ✅ CloudWatch logs show successful requests
7. ✅ Authenticated users can fetch their stats (if Auth0 configured)

## 🔄 Regular Deployment (After Initial Setup)

For future code updates:

```bash
cd /Users/nathan/Desktop/athlete-unknown-api
./aws/deploy.sh dev
```

That's it! The deployment process is fully automated.

---

**Need help?** Check:
- [AWS_LAMBDA_DEPLOYMENT.md - Troubleshooting](AWS_LAMBDA_DEPLOYMENT.md#troubleshooting)
- Lambda logs: `aws logs tail /aws/lambda/statsland-athlete-unknown-api-dev --follow`
- AWS Lambda Console: https://console.aws.amazon.com/lambda/
