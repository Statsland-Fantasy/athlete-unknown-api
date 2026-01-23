# Local Testing Guide - Web Scraper Endpoint

## Quick Start (3 Steps)

### Option A: Use AWS DynamoDB (If Tables Exist)

If you already have DynamoDB tables in AWS:

```bash
# Step 1: Set environment variables
export ROUNDS_TABLE_NAME=AthleteUnknownRoundsDev
export USER_STATS_TABLE_NAME=AthleteUnknownUserStatsDev
export AWS_REGION=us-west-2

# Step 2: Start the server
go run .

# Step 3: In a NEW terminal, run the test script
./test-scraper-local.sh
```

### Option B: Use DynamoDB Local

If you want to test locally without AWS:

**Step 1: Start DynamoDB Local**

If you have Docker:
```bash
docker run -d -p 8000:8000 amazon/dynamodb-local
```

OR if you have DynamoDB Local JAR:
```bash
java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb -port 8000
```

**Step 2: Create Tables**
```bash
# Create Rounds table
aws dynamodb create-table \
    --table-name AthleteUnknownRoundsDev \
    --attribute-definitions \
        AttributeName=playDate,AttributeType=S \
        AttributeName=sport,AttributeType=S \
    --key-schema \
        AttributeName=playDate,KeyType=HASH \
        AttributeName=sport,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000

# Create User Stats table
aws dynamodb create-table \
    --table-name AthleteUnknownUserStatsDev \
    --attribute-definitions \
        AttributeName=userId,AttributeType=S \
    --key-schema \
        AttributeName=userId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://localhost:8000
```

**Step 3: Set Environment Variables & Start Server**
```bash
export DYNAMODB_ENDPOINT=http://localhost:8000
export ROUNDS_TABLE_NAME=AthleteUnknownRoundsDev
export USER_STATS_TABLE_NAME=AthleteUnknownUserStatsDev
export AWS_REGION=us-west-2

go run .
```

**Step 4: Run Tests (in a NEW terminal)**
```bash
./test-scraper-local.sh
```

## Manual Testing

If you prefer to test manually:

### 1. Test Health Endpoint
```bash
curl http://localhost:8080/health
```

Expected: `{"status":"healthy"}`

### 2. Test Basketball Scraper (LeBron James)
```bash
curl -X POST "http://localhost:8080/v1/round?sport=basketball&playDate=2025-12-15&sportsReferencePath=/players/j/jamesle01.html&theme=GOAT" | jq .
```

Expected: 201 Created with full player data

### 3. Verify Round Was Created
```bash
curl "http://localhost:8080/v1/round?sport=basketball&playDate=2025-12-15" | jq .
```

Expected: 200 OK with the round data

### 4. Test Baseball Scraper (Derek Jeter)
```bash
curl -X POST "http://localhost:8080/v1/round?sport=baseball&playDate=2025-12-16&sportsReferencePath=/players/j/jeterde01.shtml&theme=Captain" | jq .
```

### 5. Test Football Scraper (Tom Brady)
```bash
curl -X POST "http://localhost:8080/v1/round?sport=football&playDate=2025-12-17&sportsReferencePath=/players/B/BradTo00.htm&theme=GOAT" | jq .
```

## What to Look For

When testing the scraper, verify:

✅ **Status Code**: Should be `201 Created`
✅ **Player Name**: Should be correctly scraped
✅ **All Fields Populated**: bio, playerInformation, draftInformation, etc.
✅ **Theme Field**: Should match what you provided
✅ **NO previouslyPlayedDates**: This field was removed
✅ **Round ID Format**: Should be like `Basketball20251215`
✅ **Stats Initialized**: totalPlays: 0, percentageCorrect: 0.0, etc.

## Troubleshooting

### "Connection refused" when starting server
- **Cause**: DynamoDB is not running
- **Fix**: Start DynamoDB Local or check AWS credentials

### "Table not found"
- **Cause**: DynamoDB tables don't exist
- **Fix**: Create tables using commands in Step 2 above

### Empty player fields in response
- **Cause**: Scraper couldn't find data on the page
- **Fix**: Verify the sportsReferencePath is correct

### "Round already exists"
- **Cause**: You already created a round for that sport+date
- **Fix**: Use a different date or delete the existing round:
  ```bash
  curl -X DELETE "http://localhost:8080/v1/round?sport=basketball&playDate=2025-12-15"
  ```

## View Data in DynamoDB

**DynamoDB Local:**
```bash
aws dynamodb scan \
  --table-name AthleteUnknownRoundsDev \
  --endpoint-url http://localhost:8000
```

**AWS DynamoDB:**
```bash
aws dynamodb scan --table-name AthleteUnknownRoundsDev
```

## Stop Testing

When done:
1. Stop the Go server: `Ctrl+C`
2. Stop DynamoDB Local: `docker stop <container_id>` or `Ctrl+C`

## Next Steps

After local testing, you can:
1. Merge the PR if everything works
2. Deploy to AWS Lambda using the `backendLambda` branch
3. Test the scraper in production
