# Quick Testing Steps

Follow these steps to test the web-scraper-endpoint branch locally:

## Step 1: Setup DynamoDB Local

Run the setup script to start DynamoDB Local and create tables:

```bash
./setup-local-db.sh
```

This will:
- Start DynamoDB Local in a Docker container
- Create the required tables
- Display instructions for next steps

## Step 2: Start the API Server

Open a terminal and run:

```bash
export DYNAMODB_ENDPOINT=http://localhost:8000
export ROUNDS_TABLE_NAME=AthleteUnknownRoundsDev
export USER_STATS_TABLE_NAME=AthleteUnknownUserStatsDev
export AWS_REGION=us-west-2

go run .
```

You should see:
```
DynamoDB client initialized (Rounds Table: AthleteUnknownRoundsDev, User Stats Table: AthleteUnknownUserStatsDev, Region: us-west-2)
Server starting on port 8080
API endpoints available at /v1/*
```

**Keep this terminal open!**

## Step 3: Run the Tests

Open a **NEW terminal** and run:

```bash
cd /Users/nathan/Desktop/athlete-unknown-api
./test-scraper-local.sh
```

This will automatically test:
- ✅ Health check
- ✅ Basketball player scraping (LeBron James)
- ✅ Baseball player scraping (Derek Jeter)
- ✅ Football player scraping (Tom Brady)
- ✅ Error handling (missing parameters, duplicates)
- ✅ Name-based search

Watch the output to see if all tests pass!

## What You're Testing

The main feature is the **POST /v1/round** endpoint that:
1. Takes a sports-reference.com player path
2. Scrapes the player's page for data
3. Formats all the data properly
4. Creates a round in DynamoDB
5. Returns the complete round with player info

## Manual Test (Optional)

If you want to test manually instead:

```bash
# Test with a specific player
curl -X POST "http://localhost:8080/v1/round?sport=basketball&playDate=2025-12-15&sportsReferencePath=/players/j/jamesle01.html&theme=GOAT" | jq .

# Verify it was created
curl "http://localhost:8080/v1/round?sport=basketball&playDate=2025-12-15" | jq .
```

## Check the Data

View all rounds in DynamoDB:

```bash
aws dynamodb scan \
  --table-name AthleteUnknownRoundsDev \
  --endpoint-url http://localhost:8000 \
  --region us-west-2
```

## When You're Done

1. Stop the API server (Ctrl+C in the first terminal)
2. Stop DynamoDB Local:
   ```bash
   docker stop dynamodb-local
   ```

3. (Optional) Remove the container:
   ```bash
   docker rm dynamodb-local
   ```

## Troubleshooting

**"Cannot connect to Docker daemon"**
- Start Docker Desktop

**"Server not running" error**
- Make sure you completed Step 2 and the server is running
- Check that it's running on port 8080

**"Table not found"**
- Run `./setup-local-db.sh` again
- Make sure you set `DYNAMODB_ENDPOINT=http://localhost:8000`

**Empty player data**
- The scraper might have failed
- Check the server logs for errors
- Verify the sportsReferencePath is correct

## Expected Results

✅ All tests should pass
✅ Player data should be complete (name, bio, stats, achievements, etc.)
✅ Each round should have a `theme` field
✅ NO `previouslyPlayedDates` field (removed in this branch)
✅ Round IDs should be like `Basketball20251215`

Need more details? See [START-LOCAL-TESTING.md](START-LOCAL-TESTING.md)
