#!/bin/bash

# Test script for web-scraper-endpoint branch
# This script tests the new POST /v1/round scraper endpoint

BASE_URL="http://localhost:8080"
DATE=$(date +%Y-%m-%d)

echo "=========================================="
echo "Web Scraper Endpoint Test Suite"
echo "=========================================="
echo ""

# Function to print test header
print_test() {
    echo ""
    echo "=========================================="
    echo "Test $1: $2"
    echo "=========================================="
}

# Function to check if server is running
check_server() {
    if ! curl -s "${BASE_URL}/health" > /dev/null 2>&1; then
        echo "❌ ERROR: Server not running at ${BASE_URL}"
        echo ""
        echo "Please start the server first:"
        echo "  export DYNAMODB_ENDPOINT=http://localhost:8000"
        echo "  export ROUNDS_TABLE_NAME=AthleteUnknownRoundsDev"
        echo "  export USER_STATS_TABLE_NAME=AthleteUnknownUserStatsDev"
        echo "  go run ."
        exit 1
    fi
    echo "✅ Server is running"
}

# Check server
check_server

# Test 1: Health check
print_test "1" "Health Check"
curl -s "${BASE_URL}/health" | jq .

# Test 2: Basketball scraper with direct path
print_test "2" "Basketball - LeBron James (Direct Path)"
curl -s -X POST "${BASE_URL}/v1/round?sport=basketball&playDate=${DATE}&sportsReferencePath=/players/j/jamesle01.html&theme=GOAT" | jq .

# Give server a moment
sleep 2

# Test 3: Verify the round was created
print_test "3" "Verify Basketball Round Created"
curl -s "${BASE_URL}/v1/round?sport=basketball&playDate=${DATE}" | jq .

# Test 4: Baseball scraper
print_test "4" "Baseball - Derek Jeter"
BASEBALL_DATE=$(date -v+1d +%Y-%m-%d 2>/dev/null || date -d "+1 day" +%Y-%m-%d)
curl -s -X POST "${BASE_URL}/v1/round?sport=baseball&playDate=${BASEBALL_DATE}&sportsReferencePath=/players/j/jeterde01.shtml&theme=Captain" | jq .

sleep 2

# Test 5: Football scraper
print_test "5" "Football - Tom Brady"
FOOTBALL_DATE=$(date -v+2d +%Y-%m-%d 2>/dev/null || date -d "+2 days" +%Y-%m-%d)
curl -s -X POST "${BASE_URL}/v1/round?sport=football&playDate=${FOOTBALL_DATE}&sportsReferencePath=/players/B/BradTo00.htm&theme=GOAT" | jq .

sleep 2

# Test 6: Test error - missing parameters
print_test "6" "Error Test - Missing Parameters"
curl -s -X POST "${BASE_URL}/v1/round?sport=basketball" | jq .

# Test 7: Test error - duplicate round
print_test "7" "Error Test - Duplicate Round"
curl -s -X POST "${BASE_URL}/v1/round?sport=basketball&playDate=${DATE}&sportsReferencePath=/players/j/jamesle01.html" | jq .

# Test 8: Test name search (may be less reliable)
print_test "8" "Name Search - Stephen Curry"
CURRY_DATE=$(date -v+3d +%Y-%m-%d 2>/dev/null || date -d "+3 days" +%Y-%m-%d)
curl -s -X POST "${BASE_URL}/v1/round?sport=basketball&playDate=${CURRY_DATE}&name=Stephen+Curry&theme=Splash+Brother" | jq .

echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo "All tests completed!"
echo ""
echo "To view created rounds in DynamoDB:"
echo "  aws dynamodb scan --table-name AthleteUnknownRoundsDev --endpoint-url http://localhost:8000"
echo ""
