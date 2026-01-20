# AI Upscaler Testing Guide

Quick guide to test the Replicate AI image upscaler integration.

---

## Prerequisites

1. **Replicate Account**: Sign up at https://replicate.com
2. **API Token**: Get from https://replicate.com/account/api-tokens
3. **Go installed**: Version 1.25+

---

## Setup (5 minutes)

### Step 1: Add API Key to .env

Create or edit `.env` in the project root:

```bash
# Copy from example
cp .env.example .env

# Or edit directly
nano .env
```

Add these lines to `.env`:

```env
# AI Image Upscaler Configuration
AI_UPSCALER_ENABLED=true
AI_UPSCALER_API_KEY=r8_YOUR_ACTUAL_TOKEN_HERE
AI_UPSCALER_API_URL=https://api.replicate.com/v1/predictions
```

**Replace `r8_YOUR_ACTUAL_TOKEN_HERE` with your actual Replicate API token!**

---

## Testing Methods

### Method 1: Quick Test Script (Recommended)

Run the standalone test script:

```bash
go run test_upscaler.go
```

**What it does:**
- Tests upscaler with 3 real athlete photos
- Shows original and upscaled URLs
- Validates configuration
- Takes ~60 seconds total

**Expected output:**
```
=== AI Upscaler Test ===
Enabled: true
API URL: https://api.replicate.com/v1/predictions
API Key: r8_abcd...xyz

Testing upscaler with sample athlete photos...

Test 1/3
Original: https://www.basketball-reference.com/...
✅ Success! Upscaled: https://replicate.delivery/pbxt/...

Test 2/3
Original: https://www.pro-football-reference.com/...
✅ Success! Upscaled: https://replicate.delivery/pbxt/...

Test 3/3
Original: https://www.baseball-reference.com/...
✅ Success! Upscaled: https://replicate.delivery/pbxt/...

=== Test Complete ===
```

---

### Method 2: Test via API Endpoint

**Step 1:** Start the API server:

```bash
# Make sure local DynamoDB is running (if testing locally)
# Otherwise it will try to connect to AWS DynamoDB

go run .
```

**Step 2:** Create a test round with scraping:

```bash
curl -X POST "http://localhost:8080/v1/round" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "sport": "basketball",
    "sportsReferenceURL": "https://www.basketball-reference.com/players/j/jamesle01.html",
    "playDate": "2025-01-20"
  }'
```

**Step 3:** Check the response - the `photo` field should have an upscaled URL:

```json
{
  "player": {
    "photo": "https://replicate.delivery/pbxt/..."
  }
}
```

The URL starting with `replicate.delivery` confirms upscaling worked!

---

### Method 3: Unit Test the Upscaler Directly

Create a simple Go test:

```bash
# Create test file
cat > test_manual.go << 'EOF'
package main

import "fmt"

func main() {
    upscaler := NewImageUpscaler(
        "YOUR_API_KEY",
        "https://api.replicate.com/v1/predictions",
        true,
    )

    testURL := "https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg"
    result := upscaler.UpscaleImage(testURL)

    fmt.Printf("Original: %s\n", testURL)
    fmt.Printf("Upscaled: %s\n", result)
}
EOF

# Run it
go run test_manual.go
```

---

## Troubleshooting

### Error: "API not configured"

**Cause:** Environment variables not loaded

**Fix:**
```bash
# Verify .env exists and has correct values
cat .env | grep AI_UPSCALER
```

### Error: "401 Unauthorized"

**Cause:** Invalid or expired API token

**Fix:**
1. Go to https://replicate.com/account/api-tokens
2. Generate a new token
3. Update `.env` with new token

### Error: "prediction status: failed"

**Cause:** Invalid image URL or unsupported format

**Fix:**
- Verify the image URL is publicly accessible
- Try a different image
- Check Replicate dashboard for error details

### Error: "request failed: timeout"

**Cause:** Upscaling takes longer than 60 seconds

**Fix:** Increase timeout in `image_upscaler.go`:
```go
httpClient: &http.Client{
    Timeout: 120 * time.Second, // Increase to 2 minutes
},
```

### Images Not Actually Upscaled

**Symptoms:** Returns original URL without error

**Check:**
1. Is `AI_UPSCALER_ENABLED=true` in `.env`?
2. Is the API key valid?
3. Check logs for `[Upscaler]` messages

---

## Cost Estimates

**Replicate Real-ESRGAN Pricing:**
- ~$0.0023 per image
- ~$2.30 per 1,000 images
- ~$23 per 10,000 images

**Budget planning:**
- 1 player photo per day = ~$0.84/year
- 3 player photos per day (multi-sport) = ~$2.52/year
- 10 player photos per day = ~$8.40/year

Very affordable! 🎉

---

## Verify Upscaling Quality

Compare the original and upscaled images:

**Original (typically 200x200px):**
```
https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg
```

**Upscaled (400x400px or better):**
```
https://replicate.delivery/pbxt/[generated-id]/out.png
```

**Visual check:**
- Download both images
- Open side-by-side
- Upscaled should have:
  - Higher resolution (2x)
  - Sharper details
  - Less pixelation
  - Better clarity

---

## Next Steps

Once testing is successful:

1. ✅ **Enable in production**: Update production `.env` with real API key
2. ✅ **Monitor costs**: Check Replicate dashboard regularly
3. ✅ **Test with real scraping**: Create a few test rounds
4. ✅ **Verify in frontend**: Check that upscaled photos display correctly
5. ✅ **Set up alerts**: Monitor API usage and costs

---

## Quick Commands Reference

```bash
# Test upscaler standalone
go run test_upscaler.go

# Start API server with upscaler enabled
go run .

# Check environment variables
cat .env | grep AI_UPSCALER

# Test scraping with upscaling
curl -X POST http://localhost:8080/v1/round \
  -H "X-API-Key: dev-admin-key-change-me" \
  -H "Content-Type: application/json" \
  -d '{"sport":"basketball","sportsReferenceURL":"https://www.basketball-reference.com/players/j/jamesle01.html","playDate":"2025-01-20"}'
```

---

## Support

If you encounter issues:
1. Check Replicate status: https://replicate.com/status
2. View logs in terminal (look for `[Upscaler]` messages)
3. Check Replicate dashboard: https://replicate.com/dashboard
4. Review API docs: https://replicate.com/docs/reference/http

Happy upscaling! 🚀
