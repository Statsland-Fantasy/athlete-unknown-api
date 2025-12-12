# Quick Start: Deploy to AWS Lambda

This guide will get your API deployed to AWS Lambda in minutes.

## Prerequisites Checklist

- [ ] AWS CLI installed and configured (`aws configure`)
- [ ] Terraform installed (`terraform --version`)
- [ ] Go 1.25+ installed (`go version`)
- [ ] DynamoDB tables exist in AWS:
  - [ ] `AthleteUnknownRoundsDev`
  - [ ] `AthleteUnknownUserStatsDev`

## Deploy in 3 Steps

### Step 1: Configure Variables

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your settings:
```hcl
aws_region            = "us-west-2"  # Your AWS region
environment           = "dev"
function_name         = "athlete-unknown-api-dev"
rounds_table_name     = "AthleteUnknownRoundsDev"      # Your table name
user_stats_table_name = "AthleteUnknownUserStatsDev"   # Your table name
log_retention_days    = 7
```

### Step 2: Deploy

```bash
cd ..  # Back to project root
./deploy.sh
```

The script will:
1. Build the Lambda function
2. Create the deployment package
3. Deploy everything to AWS using Terraform

### Step 3: Test

After deployment, Terraform will output your function URL:

```
function_url = "https://xxxxxxxxxxxxx.lambda-url.us-west-2.on.aws/"
```

Test it:
```bash
# Save your function URL
FUNCTION_URL="paste-your-url-here"

# Test health endpoint
curl "${FUNCTION_URL}/health"

# Expected: {"status":"healthy"}

# Test root endpoint
curl "${FUNCTION_URL}/"

# Test a round (if you have data)
curl "${FUNCTION_URL}/v1/round?sport=basketball&playDate=2025-11-15"
```

## Update Your Website

Update your frontend to use the Lambda URL:

```javascript
// Before (local)
const API_BASE_URL = 'http://localhost:8080';

// After (Lambda)
const API_BASE_URL = 'https://xxxxxxxxxxxxx.lambda-url.us-west-2.on.aws';
```

## Environment Variables

The Lambda function automatically gets these environment variables from Terraform:

- `ROUNDS_TABLE_NAME` → Your DynamoDB Rounds table
- `USER_STATS_TABLE_NAME` → Your DynamoDB User Stats table
- `AWS_REGION` → Your AWS region

To change them later:

**Option 1: Via Terraform** (Recommended)
```bash
cd terraform
# Edit terraform.tfvars
terraform apply
```

**Option 2: Via AWS Console**
1. Go to Lambda → Functions → `athlete-unknown-api-dev`
2. Configuration → Environment variables
3. Edit and save

**Option 3: Via AWS CLI**
```bash
aws lambda update-function-configuration \
  --function-name athlete-unknown-api-dev \
  --environment "Variables={ROUNDS_TABLE_NAME=NewTable,USER_STATS_TABLE_NAME=NewStatsTable,AWS_REGION=us-west-2}"
```

## What Gets Created

- ✅ Lambda Function (`athlete-unknown-api-dev`)
- ✅ IAM Role with DynamoDB read/write permissions
- ✅ Function URL (public HTTPS endpoint)
- ✅ CloudWatch Log Group (7-day retention)

## Common Issues

### "AccessDeniedException" from DynamoDB
**Cause:** IAM role doesn't have DynamoDB permissions
**Fix:** Check that table names in `terraform.tfvars` match your actual tables

### "Table not found"
**Cause:** DynamoDB tables don't exist or wrong region
**Fix:** Create tables or update `AWS_REGION` in `terraform.tfvars`

### Build fails
**Cause:** Go not installed or wrong version
**Fix:** Install Go 1.25+: https://go.dev/dl/

## Monitoring

View live logs:
```bash
aws logs tail /aws/lambda/athlete-unknown-api-dev --follow
```

Or in AWS Console:
CloudWatch → Log groups → `/aws/lambda/athlete-unknown-api-dev`

## Making Updates

When you change code:

```bash
./deploy.sh
```

That's it! The script rebuilds and redeploys automatically.

## Cost

Estimated monthly cost (low-medium traffic):
- 100K requests: ~$0.20
- Compute time: ~$0.17
- **Total: ~$0.37/month** (plus DynamoDB)

## Need Help?

See detailed documentation:
- [DEPLOYMENT.md](DEPLOYMENT.md) - Complete deployment guide
- [README.md](README.md) - API documentation
- AWS Support: https://console.aws.amazon.com/support/

## Clean Up

To remove everything:
```bash
cd terraform
terraform destroy
```

Type `yes` to confirm.
