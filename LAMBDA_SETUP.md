# AWS Lambda Setup Summary

Your Gin HTTP server has been successfully converted to support AWS Lambda deployment!

## What Was Changed

### 1. **New Files Created**

- **[main_lambda.go](main_lambda.go)** - Lambda-specific entry point (build tag: `lambda`)
- **[router.go](router.go)** - Shared router setup used by both HTTP and Lambda modes
- **[cmd/lambda/README.md](cmd/lambda/README.md)** - Comprehensive Lambda deployment guide
- **[cmd/lambda/build.sh](cmd/lambda/build.sh)** - Build script for Linux/Mac
- **[cmd/lambda/build.bat](cmd/lambda/build.bat)** - Build script for Windows
- **[template.yaml](template.yaml)** - AWS SAM deployment template
- **[Makefile](Makefile)** - Build automation with multiple targets

### 2. **Modified Files**

- **[main.go](main.go)** - Updated with build tag `!lambda` to exclude from Lambda builds
- **[main_test.go](main_test.go)** - Updated handler function names to use exported versions
- **[.gitignore](.gitignore)** - Added build artifacts
- **[go.mod](go.mod)** - Added Lambda dependencies:
  - `github.com/aws/aws-lambda-go`
  - `github.com/awslabs/aws-lambda-go-api-proxy`

## Build Tags Strategy

The project now uses Go build tags to maintain two separate entry points:

| Mode | Build Tag | File | Usage |
|------|-----------|------|-------|
| HTTP Server | `!lambda` | main.go | Regular HTTP server on port 8080 |
| Lambda | `lambda` | main_lambda.go | AWS Lambda handler |
| Shared | none | router.go | Router setup used by both |

## Quick Start

### Build for Lambda

**Recommended - Using Makefile:**
```bash
make build-lambda
```

**Manual build from root directory:**
```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda -o cmd/lambda/bootstrap -ldflags="-s -w" .
```

**Using build scripts:**
```bash
# Windows
cmd\lambda\build.bat

# Linux/Mac
cmd/lambda/build.sh
```

This creates:
- `cmd/lambda/bootstrap` - The Lambda binary (Linux AMD64)
- `cmd/lambda/bootstrap.zip` - Ready-to-deploy package

**Note:** All builds must be run from the **root directory**. The `-tags lambda` flag ensures:
- ✅ Includes [main_lambda.go](main_lambda.go) (tagged with `lambda`)
- ❌ Excludes [main.go](main.go) (tagged with `!lambda`)
- ✅ Includes all other source files

### Deploy to AWS Lambda

**Option 1: Using AWS SAM (Recommended)**
```bash
make sam-deploy
```

**Option 2: Using AWS CLI**
```bash
AWS_LAMBDA_FUNCTION_NAME=my-function make deploy-lambda
```

**Option 3: Manual**
```bash
aws lambda update-function-code \
  --function-name athlete-unknown-api \
  --zip-file fileb://cmd/lambda/bootstrap.zip
```

### Test Locally

**Run as HTTP server:**
```bash
make run
# or
go run .
```

**Run as Lambda (requires AWS SAM CLI):**
```bash
make sam-local
```

## How It Works

### 1. **Gin Lambda Adapter**

The `awslabs/aws-lambda-go-api-proxy/gin` adapter translates AWS API Gateway events to Gin HTTP requests:

```go
ginLambda = ginadapter.New(router)
lambda.Start(Handler)
```

### 2. **Handler Function**

The Lambda handler receives API Gateway events and proxies them to Gin:

```go
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return ginLambda.ProxyWithContext(ctx, req)
}
```

### 3. **Shared Router**

The `SetupRouter()` function in [router.go](router.go) is used by both modes:
- Initializes DynamoDB
- Configures CORS
- Sets up all routes and middleware
- Returns configured Gin router

## Environment Variables for Lambda

Set these in your Lambda configuration:

**Required:**
- `AWS_REGION` - AWS region (e.g., us-west-2)
- `ROUNDS_TABLE_NAME` - DynamoDB table name for rounds
- `USER_STATS_TABLE_NAME` - DynamoDB table name for user stats

**Optional:**
- `GIN_MODE` - Set to "release" for production
- `CORS_ALLOWED_ORIGINS` - Comma-separated allowed origins
- `API_KEY` - Admin API key
- `AUTH0_DOMAIN` - Auth0 domain for JWT validation
- `AUTH0_AUDIENCE` - Auth0 audience for JWT validation

## API Gateway Setup

Your Lambda function needs an API Gateway trigger with proxy integration:

1. Create API Gateway (HTTP API recommended)
2. Create route: `ANY /{proxy+}`
3. Integration: Lambda Function (proxy integration enabled)
4. Deploy API

All requests will be routed through the Lambda handler to your Gin routes.

## Required IAM Permissions

Lambda execution role needs:

```json
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
        "dynamodb:Scan"
      ],
      "Resource": [
        "arn:aws:dynamodb:*:*:table/AthleteUnknownRounds*",
        "arn:aws:dynamodb:*:*:table/AthleteUnknownUserStats*"
      ]
    }
  ]
}
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build HTTP server binary |
| `make build-lambda` | Build Lambda deployment package |
| `make clean` | Remove all build artifacts |
| `make test` | Run all tests |
| `make run` | Run HTTP server locally |
| `make sam-local` | Run Lambda locally with SAM |
| `make sam-deploy` | Deploy using AWS SAM |
| `make deploy-lambda` | Deploy to existing Lambda function |
| `make help` | Show all available commands |

## Architecture Diagram

```
┌─────────────────────────────────────────────┐
│  API Gateway (HTTP API / REST API)         │
│  Route: ANY /{proxy+}                       │
└─────────────────┬───────────────────────────┘
                  │
                  │ APIGatewayProxyRequest
                  ▼
┌─────────────────────────────────────────────┐
│  AWS Lambda Function                        │
│  ┌────────────────────────────────────┐    │
│  │  Handler (main_lambda.go)          │    │
│  │    ↓                                │    │
│  │  Gin Lambda Adapter                │    │
│  │    ↓                                │    │
│  │  SetupRouter() (router.go)         │    │
│  │    ↓                                │    │
│  │  Gin Router + Middleware           │    │
│  │    ├─ CORS                          │    │
│  │    ├─ JWT Auth                      │    │
│  │    └─ API Key Auth                  │    │
│  │    ↓                                │    │
│  │  Route Handlers                     │    │
│  │    ├─ /v1/round                     │    │
│  │    ├─ /v1/stats/*                   │    │
│  │    └─ /health                       │    │
│  └────────────────────────────────────┘    │
└─────────────────┬───────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────┐
│  DynamoDB Tables                            │
│  ├─ AthleteUnknownRounds                    │
│  └─ AthleteUnknownUserStats                 │
└─────────────────────────────────────────────┘
```

## Testing

The conversion maintains all existing functionality:

- ✅ All routes work the same
- ✅ All middleware functions identically
- ✅ All tests pass
- ✅ CORS configured
- ✅ JWT authentication works
- ✅ API key authentication works
- ✅ DynamoDB integration unchanged

## Next Steps

1. **Set up Lambda function** in AWS Console or using SAM
2. **Configure environment variables** in Lambda
3. **Create API Gateway** and connect to Lambda
4. **Set up IAM permissions** for DynamoDB access
5. **Deploy your code** using `make build-lambda` and deploy commands
6. **Test your endpoints** through API Gateway URL

## Troubleshooting

See [cmd/lambda/README.md](cmd/lambda/README.md) for:
- Binary size optimization
- Cold start improvements
- CORS issues
- DynamoDB connection problems
- Build troubleshooting

## Resources

- [AWS Lambda Go Documentation](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)
- [Gin Lambda Adapter](https://github.com/awslabs/aws-lambda-go-api-proxy)
- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)
