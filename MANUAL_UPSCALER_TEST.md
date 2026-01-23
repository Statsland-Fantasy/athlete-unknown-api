# Manual AI Upscaler Test - Step by Step

Follow these steps to verify your AI upscaler is working.

---

## ✅ What We've Done So Far

1. ✅ Added your Replicate API key to `.env`
2. ✅ Set `AI_UPSCALER_ENABLED=true`
3. ✅ Fixed the code to download images first (bypasses hotlink protection)
4. ✅ Verified you have credits in Replicate

---

## 🧪 Test Method: Watch It Work Live

### Step 1: Start the Server in One Terminal

```bash
cd /Users/nathan/Desktop/athlete-unknown-api
go run .
```

**Expected output:**
```
DynamoDB client initialized...
AI upscaler enabled (API: https://api.replicate.com/v1/predictions)
Server starting on port 8080
```

**✅ Confirm you see "AI upscaler enabled"**

---

### Step 2: Open a Second Terminal

Keep the first terminal open (server running), and open a new terminal window.

---

### Step 3: Test by Getting an Existing Round

In the second terminal:

```bash
curl "http://localhost:8080/v1/round?sport=basketball&playDate=2026-01-12"
```

This will try to get today's round. If it doesn't exist, that's fine - we'll create one next.

---

### Step 4: Create a Simple Test Round (No Scraping)

Create a file called `test-round-simple.json`:

```json
{
  "roundId": "Basketball#999",
  "sport": "basketball",
  "playDate": "2026-01-20",
  "theme": "Test",
  "player": {
    "sport": "basketball",
    "sportsReferenceURL": "https://www.basketball-reference.com/players/j/jamesle01.html",
    "name": "LeBron James",
    "bio": "Test",
    "playerInformation": "Test",
    "draftInformation": "Test",
    "yearsActive": "2003-Present",
    "teamsPlayedOn": "CLE, MIA, LAL",
    "jerseyNumbers": "#23, #6",
    "careerStats": "Test stats",
    "personalAchievements": "4x NBA Champion",
    "photo": "https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg"
  },
  "stats": {
    "playDate": "2026-01-20",
    "name": "LeBron James",
    "sport": "basketball",
    "totalPlays": 0,
    "percentageCorrect": 0,
    "highestScore": 0,
    "averageCorrectScore": 0,
    "averageNumberOfTileFlips": 0,
    "mostCommonFirstTileFlipped": "",
    "mostCommonLastTileFlipped": "",
    "mostCommonTileFlipped": "",
    "leastCommonTileFlipped": "",
    "mostTileFlippedTracker": {"bio":0,"playerInformation":0,"draftInformation":0,"yearsActive":0,"teamsPlayedOn":0,"jerseyNumbers":0,"careerStats":0,"personalAchievements":0,"photo":0},
    "firstTileFlippedTracker": {"bio":0,"playerInformation":0,"draftInformation":0,"yearsActive":0,"teamsPlayedOn":0,"jerseyNumbers":0,"careerStats":0,"personalAchievements":0,"photo":0},
    "lastTileFlippedTracker": {"bio":0,"playerInformation":0,"draftInformation":0,"yearsActive":0,"teamsPlayedOn":0,"jerseyNumbers":0,"careerStats":0,"personalAchievements":0,"photo":0}
  }
}
```

Wait - this won't trigger the upscaler because it's using the PUT endpoint which stores the data as-is.

---

## 🎯 BETTER TEST: Use curl + watch logs

Actually, the easiest way is to watch what happens when the upscaler is called directly in the code. Let me give you a simpler approach:

### Direct Test Script

Create `direct-test-upscale.sh`:

```bash
#!/bin/bash

echo "Testing AI Upscaler..."
echo ""

# Test photo URL
PHOTO_URL="https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg"

echo "Photo to upscale: $PHOTO_URL"
echo ""
echo "Starting test... (this will take 20-30 seconds)"
echo ""

# Load API key from .env
export $(cat .env | grep AI_UPSCALER_API_KEY | xargs)
export $(cat .env | grep AI_UPSCALER_ENABLED | xargs)

if [ "$AI_UPSCALER_ENABLED" != "true" ]; then
    echo "❌ Error: AI_UPSCALER_ENABLED is not true"
    exit 1
fi

# Download image first
echo "1. Downloading image..."
curl -s "$PHOTO_URL" -o /tmp/test-photo.jpg \
  -H "User-Agent: Mozilla/5.0" \
  -H "Referer: https://www.basketball-reference.com/"

if [ ! -f /tmp/test-photo.jpg ]; then
    echo "❌ Failed to download image"
    exit 1
fi

echo "✅ Downloaded $(wc -c < /tmp/test-photo.jpg) bytes"
echo ""

# Convert to base64
echo "2. Converting to base64..."
BASE64_DATA=$(base64 -i /tmp/test-photo.jpg)
DATA_URL="data:image/jpeg;base64,$BASE64_DATA"
echo "✅ Converted to base64 ($(echo -n "$BASE64_DATA" | wc -c) chars)"
echo ""

# Call Replicate API
echo "3. Sending to Replicate for upscaling..."
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

echo ""
echo "=== RESULT ==="
echo "$RESPONSE" | jq '.'

STATUS=$(echo "$RESPONSE" | jq -r '.status')
OUTPUT=$(echo "$RESPONSE" | jq -r '.output')
ERROR=$(echo "$RESPONSE" | jq -r '.error')

echo ""
if [ "$STATUS" = "succeeded" ]; then
    echo "✅ SUCCESS! Upscaling worked!"
    echo ""
    echo "Original:  $PHOTO_URL"
    echo "Upscaled:  $OUTPUT"
    echo ""
    echo "Open both URLs in your browser to compare!"
elif [ "$ERROR" != "null" ]; then
    echo "❌ ERROR: $ERROR"
else
    echo "⚠️  Status: $STATUS"
fi
```

Make it executable and run it:

```bash
chmod +x direct-test-upscale.sh
./direct-test-upscale.sh
```

---

## 📊 What Success Looks Like

You should see:

```
✅ SUCCESS! Upscaling worked!

Original:  https://www.basketball-reference.com/req/.../jamesle01.jpg
Upscaled:  https://replicate.delivery/pbxt/abc123xyz/out.png

Open both URLs in your browser to compare!
```

---

## 🔍 Visual Verification

1. Copy the **Original** URL
2. Copy the **Upscaled** URL
3. Open both in browser tabs
4. **Compare them:**
   - Upscaled should be **2x larger** (400x400 vs 200x200)
   - Upscaled should be **sharper** and **clearer**
   - Upscaled should have **better detail** in facial features

---

## ✅ Integration Test (When Ready for Production)

Once the direct test works, test it integrated with your API:

1. Make sure you have a scraping endpoint that calls the upscaler
2. Watch the server logs for `[Upscaler]` messages
3. Check that the photo URL in the database starts with `replicate.delivery`

---

## 💰 Cost Check

After testing:
1. Go to https://replicate.com/dashboard
2. Click on your test prediction
3. Verify cost was ~$0.0023
4. See the before/after images side-by-side

---

## Need Help?

If you see errors, check:
- ✅ API key is correct in `.env`
- ✅ AI_UPSCALER_ENABLED=true
- ✅ You have credits in Replicate
- ✅ Internet connection is working
- ✅ Image URL is accessible

Run `./direct-test-upscale.sh` and share the output if you need help!
