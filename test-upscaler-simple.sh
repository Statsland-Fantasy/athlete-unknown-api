#!/bin/bash

# Simple shell script to test AI Upscaler via curl

echo "=== AI Upscaler Test via Replicate API ==="
echo ""

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

API_KEY="${AI_UPSCALER_API_KEY}"

if [ -z "$API_KEY" ] || [ "$API_KEY" = "r8_your_replicate_api_token_here" ]; then
    echo "❌ Error: AI_UPSCALER_API_KEY not set in .env"
    exit 1
fi

echo "✅ API Key configured: ${API_KEY:0:8}...${API_KEY: -4}"
echo ""

# Test image URL
TEST_IMAGE="https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg"

echo "Testing with LeBron James photo:"
echo "Original: $TEST_IMAGE"
echo ""
echo "Sending request to Replicate..."

# Make the API request
RESPONSE=$(curl -s -X POST \
  https://api.replicate.com/v1/predictions \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -H "Prefer: wait" \
  -d "{
    \"version\": \"42fed1c4974146d4d2414e2be2c5277c7fcf05fcc3a73abf41610695738c1d7b\",
    \"input\": {
      \"image\": \"$TEST_IMAGE\",
      \"scale\": 2,
      \"face_enhance\": false
    }
  }")

echo ""
echo "Response:"
echo "$RESPONSE" | jq '.'

# Extract status and output
STATUS=$(echo "$RESPONSE" | jq -r '.status')
OUTPUT=$(echo "$RESPONSE" | jq -r '.output')
ERROR=$(echo "$RESPONSE" | jq -r '.error')

echo ""
if [ "$STATUS" = "succeeded" ]; then
    echo "✅ SUCCESS! Upscaling completed"
    echo ""
    echo "Upscaled URL: $OUTPUT"
    echo ""
    echo "You can download and compare:"
    echo "Original:  $TEST_IMAGE"
    echo "Upscaled:  $OUTPUT"
elif [ "$ERROR" != "null" ]; then
    echo "❌ ERROR: $ERROR"
else
    echo "⚠️  Status: $STATUS"
fi
