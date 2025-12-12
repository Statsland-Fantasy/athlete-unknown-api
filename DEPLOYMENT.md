# AWS Lambda Deployment Guide

This guide walks through deploying the Athlete Unknown API as an AWS Lambda function in the Statsland DEV environment.

## Prerequisites

1. **AWS CLI** installed and configured with appropriate credentials
2. **Terraform** (>= 1.0) installed
3. **Go** (1.25+) installed
4. **DynamoDB tables** created in AWS:
   - `AthleteUnknownRoundsDev`
   - `AthleteUnknownUserStatsDev`

## Step 1: Build the Lambda Function

The API includes a build script that compiles the Go application for AWS Lambda:

```bash
./build-lambda.sh
```

This will create:
- `bootstrap` - The Lambda binary
- `lambda-deployment-package.zip` - The deployment package ready for Lambda

## Step 2: Configure Terraform Variables

1. Navigate to the terraform directory:
```bash
cd terraform
```

2. Create a `terraform.tfvars` file from the example:
```bash
cp terraform.tfvars.example terraform.tfvars
```

3. Edit `terraform.tfvars` and customize the values:
```hcl
aws_region            = "us-west-2"  # Your AWS region
environment           = "dev"
function_name         = "athlete-unknown-api-dev"
rounds_table_name     = "AthleteUnknownRoundsDev"
user_stats_table_name = "AthleteUnknownUserStatsDev"
log_retention_days    = 7
```

## Step 3: Deploy with Terraform

1. Initialize Terraform:
```bash
terraform init
```

2. Review the deployment plan:
```bash
terraform plan
```

3. Apply the configuration:
```bash
terraform apply
```

Review the changes and type `yes` to confirm.

## Step 4: Verify Deployment

After successful deployment, Terraform will output:

```
Outputs:

cloudwatch_log_group = "/aws/lambda/athlete-unknown-api-dev"
function_url = "https://xxxxxxxxxxxxx.lambda-url.us-west-2.on.aws/"
lambda_function_arn = "arn:aws:lambda:us-west-2:ACCOUNT_ID:function:athlete-unknown-api-dev"
lambda_function_name = "athlete-unknown-api-dev"
lambda_role_arn = "arn:aws:iam::ACCOUNT_ID:role/athlete-unknown-api-dev-role"
```

Save the `function_url` - this is your API endpoint!

## Step 5: Test the API

Test the health endpoint:
```bash
FUNCTION_URL="your-function-url-here"
curl "${FUNCTION_URL}/health"
```

Expected response:
```json
{"status":"healthy"}
```

Test the root endpoint:
```bash
curl "${FUNCTION_URL}/"
```

Test getting a round (requires data in DynamoDB):
```bash
curl "${FUNCTION_URL}/v1/round?sport=basketball&playDate=2025-11-15"
```

## Step 6: Update Your Website

Update your website configuration to point to the Lambda function URL:

```javascript
const API_BASE_URL = 'https://xxxxxxxxxxxxx.lambda-url.us-west-2.on.aws';
```

## What's Deployed

The Terraform configuration creates:

1. **Lambda Function** (`athlete-unknown-api-dev`)
   - Runtime: Custom runtime on Amazon Linux 2023
   - Memory: 512 MB
   - Timeout: 30 seconds
   - Handler: bootstrap

2. **IAM Role** with permissions:
   - CloudWatch Logs (write logs)
   - DynamoDB (read/write access to both tables)

3. **Function URL**
   - Public endpoint with CORS enabled
   - No authentication (configure if needed)

4. **CloudWatch Log Group**
   - Retention: 7 days (configurable)
   - Path: `/aws/lambda/athlete-unknown-api-dev`

## Environment Variables

The Lambda function uses these environment variables (configured via Terraform):

| Variable | Description | Example |
|----------|-------------|---------|
| `ROUNDS_TABLE_NAME` | DynamoDB Rounds table name | `AthleteUnknownRoundsDev` |
| `USER_STATS_TABLE_NAME` | DynamoDB User Stats table name | `AthleteUnknownUserStatsDev` |
| `AWS_REGION` | AWS region | `us-west-2` |

These can be changed via:
- **Terraform**: Edit `terraform.tfvars` and run `terraform apply`
- **AWS Console**: Lambda → Configuration → Environment variables
- **AWS CLI**:
  ```bash
  aws lambda update-function-configuration \
    --function-name athlete-unknown-api-dev \
    --environment "Variables={ROUNDS_TABLE_NAME=NewTableName,USER_STATS_TABLE_NAME=NewStatsTable,AWS_REGION=us-west-2}"
  ```

## Updating the Lambda Function

When you make code changes:

1. Rebuild the Lambda package:
```bash
./build-lambda.sh
```

2. Update via Terraform:
```bash
cd terraform
terraform apply
```

Or update directly via AWS CLI:
```bash
aws lambda update-function-code \
  --function-name athlete-unknown-api-dev \
  --zip-file fileb://lambda-deployment-package.zip
```

## Monitoring

View logs in CloudWatch:
```bash
aws logs tail /aws/lambda/athlete-unknown-api-dev --follow
```

Or via AWS Console:
- CloudWatch → Log groups → `/aws/lambda/athlete-unknown-api-dev`

## Troubleshooting

### Lambda returns 500 errors
Check CloudWatch logs:
```bash
aws logs tail /aws/lambda/athlete-unknown-api-dev --follow
```

### Cannot access DynamoDB tables
Verify:
1. Table names are correct in environment variables
2. IAM role has DynamoDB permissions
3. Tables exist in the same region

### CORS errors
The Function URL is configured with:
- Allow Origins: `*`
- Allow Methods: `GET, POST, DELETE, OPTIONS`
- Allow Headers: `Content-Type, Authorization, X-API-Key`

Update in `terraform/main.tf` if needed.

## Cost Optimization

Lambda pricing (us-west-2):
- **Requests**: $0.20 per 1M requests
- **Duration**: $0.0000166667 per GB-second

Estimated monthly cost for low-medium traffic:
- 100K requests/month: ~$0.20
- 512MB memory, 200ms avg duration: ~$0.17
- **Total**: ~$0.37/month (plus DynamoDB costs)

## Security Considerations

### Current Setup (DEV)
- Function URL has **NO authentication**
- CORS allows all origins (`*`)

### Recommended for Production
1. Enable authentication on Function URL
2. Implement API Gateway with:
   - API keys
   - Usage plans
   - Rate limiting
3. Restrict CORS origins
4. Enable AWS WAF
5. Use Secrets Manager for sensitive configuration

## Clean Up

To destroy all resources:
```bash
cd terraform
terraform destroy
```

Type `yes` to confirm deletion.

## Next Steps

1. ✅ Lambda deployed to AWS
2. ✅ IAM role configured with DynamoDB permissions
3. ✅ Environment variables set up
4. Test all API endpoints
5. Update website to use Lambda URL
6. Monitor CloudWatch logs
7. Set up CI/CD pipeline (optional)

## Support

- AWS Lambda Docs: https://docs.aws.amazon.com/lambda/
- Terraform AWS Provider: https://registry.terraform.io/providers/hashicorp/aws/latest/docs
- Go AWS Lambda: https://github.com/aws/aws-lambda-go
