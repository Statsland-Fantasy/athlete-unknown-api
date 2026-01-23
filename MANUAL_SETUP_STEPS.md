# Manual Setup Steps (Do It Yourself)

If you prefer to add your API key manually instead of sharing it:

## Step 1: Create/Edit .env file

```bash
# In the athlete-unknown-api directory
nano .env
```

## Step 2: Add These Lines

```env
# AI Image Upscaler Configuration
AI_UPSCALER_ENABLED=true
AI_UPSCALER_API_KEY=r8_YOUR_ACTUAL_TOKEN_FROM_REPLICATE_HERE
AI_UPSCALER_API_URL=https://api.replicate.com/v1/predictions
```

**Replace `r8_YOUR_ACTUAL_TOKEN_FROM_REPLICATE_HERE` with your real token!**

## Step 3: Save and Exit

- Press `Ctrl+X` to exit
- Press `Y` to confirm save
- Press `Enter` to confirm filename

## Step 4: Test It

```bash
go run test_upscaler.go
```

## Expected Output

```
=== AI Upscaler Test ===
Enabled: true
API URL: https://api.replicate.com/v1/predictions
API Key: r8_abcd...xyz

Testing upscaler with sample athlete photos...

Test 1/3
Original: https://www.basketball-reference.com/...
[Upscaler] Upscaling image: https://www.basketball-reference.com/...
[Upscaler] Successfully upscaled image: https://replicate.delivery/...
✅ Success! Upscaled: https://replicate.delivery/pbxt/...

Test 2/3
...
```

That's it! 🎉
