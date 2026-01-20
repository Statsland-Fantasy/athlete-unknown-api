# How to Personally Verify AI Upscaler Works

Complete guide to verify the upscaler is working correctly.

---

## ✅ Method 1: Visual Comparison (Best Method)

### Step 1: Run the Test Script

```bash
./test-upscaler-simple.sh
```

### Step 2: You'll See Output Like This

```
✅ SUCCESS! Upscaling completed

Upscaled URL: https://replicate.delivery/pbxt/abc123xyz/out.png

You can download and compare:
Original:  https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg
Upscaled:  https://replicate.delivery/pbxt/abc123xyz/out.png
```

### Step 3: Open Both URLs in Your Browser

**Original Image:**
- Open: `https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg`
- Typically: 200x200 pixels
- Lower quality, pixelated

**Upscaled Image:**
- Open: `https://replicate.delivery/pbxt/[your-id]/out.png`
- Should be: 400x400 pixels (2x)
- Higher quality, sharper details

### Step 4: Download and Compare Side-by-Side

```bash
# Download original
curl -o original.jpg "https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg"

# Download upscaled (replace URL with your actual upscaled URL)
curl -o upscaled.png "https://replicate.delivery/pbxt/YOUR_URL_HERE/out.png"

# Open both
open original.jpg upscaled.png
```

**What to Look For:**
- ✅ Upscaled image is **larger** (2x resolution)
- ✅ Upscaled image is **sharper** (less blurry)
- ✅ Upscaled image has **better detail** (clearer facial features)
- ✅ Upscaled image is **less pixelated**

---

## ✅ Method 2: Check Image Properties

### Using Command Line:

```bash
# Check original image size
curl -s "https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg" -o test.jpg
file test.jpg

# Output: JPEG image data, ... 200 x 200

# Check upscaled image size
curl -s "YOUR_UPSCALED_URL" -o upscaled.png
file upscaled.png

# Output: PNG image data, ... 400 x 400
```

### Using ImageMagick (if installed):

```bash
# Check original dimensions
identify -format "%wx%h\n" original.jpg
# Output: 200x200

# Check upscaled dimensions
identify -format "%wx%h\n" upscaled.png
# Output: 400x400 (or higher)
```

---

## ✅ Method 3: Test Through Your API

This verifies the upscaler is integrated correctly with your scraping workflow.

### Step 1: Start Your API Server

```bash
# Make sure you have .env with AI_UPSCALER_ENABLED=true
go run .
```

**Expected logs:**
```
AI upscaler enabled (API: https://api.replicate.com/v1/predictions)
Server starting on port 8080
```

### Step 2: Create a Test Round with Scraping

```bash
curl -X POST "http://localhost:8080/v1/round" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: dev-admin-key-change-me" \
  -d '{
    "sport": "basketball",
    "sportsReferenceURL": "https://www.basketball-reference.com/players/j/jamesle01.html",
    "playDate": "2025-01-21"
  }'
```

### Step 3: Check the Response

Look at the `player.photo` field in the response:

**Before Upscaling (if disabled):**
```json
{
  "player": {
    "photo": "https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg"
  }
}
```

**After Upscaling (if working):**
```json
{
  "player": {
    "photo": "https://replicate.delivery/pbxt/abc123xyz/out.png"
  }
}
```

**Key Indicator:** The URL should start with `replicate.delivery` instead of `basketball-reference.com`

### Step 4: Open the Photo URL

Copy the `photo` URL from the response and open it in your browser:
```
https://replicate.delivery/pbxt/abc123xyz/out.png
```

You should see a high-quality, upscaled version of the athlete photo!

---

## ✅ Method 4: Check Server Logs

When the upscaler runs, you'll see detailed logs:

```bash
# Run your API server
go run .

# In another terminal, create a test round
curl -X POST http://localhost:8080/v1/round ...
```

**Look for these log messages:**

```
[Upscaler] Upscaling image: https://www.basketball-reference.com/req/202106291/images/headshots/jamesle01.jpg
[Upscaler] Successfully upscaled image: https://replicate.delivery/pbxt/abc123xyz/out.png
```

**If you see errors:**
```
[Upscaler] Failed to upscale image: ... Using original.
```

This means upscaling failed and it fell back to the original photo (which is good - graceful degradation!)

---

## ✅ Method 5: Check Replicate Dashboard

### View Usage in Real-Time:

1. Go to: https://replicate.com/dashboard
2. You'll see recent predictions
3. Each upscale will show:
   - ✅ Status: "succeeded"
   - ✅ Model: "nightmareai/real-esrgan"
   - ✅ Cost: ~$0.0023
   - ✅ Input/Output images

### Click on a Prediction:
- See the original input image
- See the upscaled output image
- Compare them side-by-side!

---

## 🎯 Quick Verification Checklist

Run through this after adding billing:

- [ ] Run `./test-upscaler-simple.sh`
- [ ] See "✅ SUCCESS! Upscaling completed" message
- [ ] Get an upscaled URL starting with `replicate.delivery`
- [ ] Open both original and upscaled URLs in browser
- [ ] Confirm upscaled image is larger and sharper
- [ ] Check Replicate dashboard shows the prediction
- [ ] See the cost charged (~$0.0023)

If all boxes are checked → **It's working!** ✅

---

## 🚨 Troubleshooting

### Problem: Still Getting "Insufficient Credit"

**Solution:**
1. Go to https://replicate.com/account/billing
2. Add payment method or buy credits
3. Wait 2-3 minutes
4. Try again

### Problem: Upscaled URL Returns Original

**Check:**
```bash
# Verify upscaler is enabled
cat .env | grep AI_UPSCALER_ENABLED
# Should show: AI_UPSCALER_ENABLED=true
```

### Problem: "API not configured" in Logs

**Check:**
```bash
# Verify API key is set
cat .env | grep AI_UPSCALER_API_KEY
# Should show your actual API key
```

### Problem: Timeout Errors

**Solution:** Increase timeout in `image_upscaler.go`:
```go
httpClient: &http.Client{
    Timeout: 120 * time.Second, // Increase to 2 minutes
}
```

---

## 📊 Expected Results Summary

| Test | Success Indicator |
|------|------------------|
| **Script Test** | URL starts with `replicate.delivery` |
| **Visual Check** | Image is 2x larger and sharper |
| **File Size** | Upscaled file is larger (more KB) |
| **Dimensions** | Upscaled is 400x400 vs 200x200 original |
| **API Response** | `photo` field has replicate URL |
| **Server Logs** | "Successfully upscaled image" message |
| **Dashboard** | Prediction shows "succeeded" |
| **Cost** | ~$0.0023 charged per image |

---

## 💡 Pro Tips

1. **Compare Multiple Athletes**: Test with different sports to ensure consistency
2. **Check Frontend**: Verify upscaled photos display correctly in your app
3. **Monitor Costs**: Keep an eye on your Replicate dashboard
4. **Set Spending Limits**: Protect yourself from unexpected charges
5. **Test Error Handling**: Disable the upscaler and verify original photos still work

---

## 🎉 Success Criteria

**You'll know it's working when:**
- ✅ Test script returns an upscaled URL
- ✅ Visual comparison shows clear quality improvement
- ✅ API integration works seamlessly
- ✅ No errors in server logs
- ✅ Replicate dashboard shows successful predictions

**If any of these fail** → Check the troubleshooting section or let me know!
