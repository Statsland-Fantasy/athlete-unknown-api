#!/bin/bash

echo "=== AI Upscaler Direct Test ==="
echo ""

# Test photo URL
PHOTO_URL="https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg"

echo "Photo to upscale: $PHOTO_URL"
echo ""
echo "Starting test... (this will take 20-30 seconds)"
echo ""

# Load API key from .env
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

if [ "$AI_UPSCALER_ENABLED" != "true" ]; then
    echo "❌ Error: AI_UPSCALER_ENABLED is not true"
    exit 1
fi

if [ -z "$AI_UPSCALER_API_KEY" ]; then
    echo "❌ Error: AI_UPSCALER_API_KEY is not set"
    exit 1
fi

# Download image first
echo "1. Downloading image..."
curl -s "$PHOTO_URL" -o /tmp/test-photo.jpg \
  -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36" \
  -H "Referer: https://www.basketball-reference.com/"

if [ ! -f /tmp/test-photo.jpg ]; then
    echo "❌ Failed to download image"
    exit 1
fi

SIZE=$(wc -c < /tmp/test-photo.jpg | tr -d ' ')
echo "✅ Downloaded $SIZE bytes"
echo ""

# Convert to base64
echo "2. Converting to base64..."
BASE64_DATA=$(base64 -i /tmp/test-photo.jpg)
# Remove any newlines from base64 output
BASE64_DATA=$(echo "$BASE64_DATA" | tr -d '\n')
DATA_URL="data:image/jpeg;base64,$BASE64_DATA"
echo "✅ Converted to base64 ($(echo -n "$BASE64_DATA" | wc -c | tr -d ' ') characters)"
echo ""

# Call Replicate API
echo "3. Sending to Replicate for upscaling..."
echo "   (This takes 15-30 seconds, please wait...)"
echo ""

RESPONSE=$(curl -s -X POST \
  https://api.replicate.com/v1/predictions \
  -H "Authorization: Bearer $AI_UPSCALER_API_KEY" \
  -H "Content-Type: application/json" \
  -H "Prefer: wait" \
  -d "{
    \"version\": \"42fed1c4974146d4d2414e2be2c5277c7fcf05fcc3a73abf41610695738c1d7b\",
    \"input\": {
      \"image\": \"$DATA_URL\",
      \"scale\": 2,
      \"face_enhance\": false
    }
  }")

echo "=== RESPONSE ==="
echo "$RESPONSE" | jq '.' || echo "$RESPONSE"

STATUS=$(echo "$RESPONSE" | jq -r '.status' 2>/dev/null)
OUTPUT=$(echo "$RESPONSE" | jq -r '.output' 2>/dev/null)
ERROR=$(echo "$RESPONSE" | jq -r '.error' 2>/dev/null)

echo ""
echo "=== RESULT ==="
if [ "$STATUS" = "succeeded" ]; then
    echo "✅ SUCCESS! Upscaling worked!"
    echo ""
    echo "Original:  $PHOTO_URL"
    echo "Upscaled:  $OUTPUT"
    echo ""
    echo "🎉 TO VERIFY:"
    echo "1. Open both URLs in your browser"
    echo "2. Compare - upscaled should be larger and sharper"
    echo ""
    echo "Original:"
    echo "$PHOTO_URL"
    echo ""
    echo "Upscaled:"
    echo "$OUTPUT"
elif [ "$ERROR" != "null" ] && [ -n "$ERROR" ]; then
    echo "❌ ERROR: $ERROR"
else
    echo "⚠️  Status: $STATUS"
    echo "Check the response above for details"
fi

echo ""
echo "Cleaning up..."
rm -f /tmp/test-photo.jpg
echo "Done!"
