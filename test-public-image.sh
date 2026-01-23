#!/bin/bash

# Test with a publicly accessible image

echo "=== Testing AI Upscaler with Public Image ==="
echo ""

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

API_KEY="${AI_UPSCALER_API_KEY}"

# Use a public test image that allows external access
TEST_IMAGE="https://replicate.delivery/pbxt/GtQawIBiT0IfImNsvTKbDxikv3MxKlYDmCUEUFT0r3RlPakS/ComfyUI_00929_.png"

echo "Testing with public image:"
echo "URL: $TEST_IMAGE"
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
    echo "✅ SUCCESS! Your Replicate API is working!"
    echo ""
    echo "Upscaled URL: $OUTPUT"
    echo ""
    echo "The issue with Basketball Reference photos is hotlink protection."
    echo "We need to download images first, then upscale them."
elif [ "$ERROR" != "null" ]; then
    echo "❌ ERROR: $ERROR"
else
    echo "⚠️  Status: $STATUS"
fi
