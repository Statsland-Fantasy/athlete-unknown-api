package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// ImageUpscaler handles AI-powered image upscaling using Replicate's Real-ESRGAN model
// Pricing: ~$0.0023 per image
// Model: nightmareai/real-esrgan (2x upscaling)
// API Docs: https://replicate.com/docs/reference/http
type ImageUpscaler struct {
	apiKey     string
	apiURL     string
	enabled    bool
	httpClient *http.Client
}

// NewImageUpscaler creates a new image upscaler service
func NewImageUpscaler(apiKey, apiURL string, enabled bool) *ImageUpscaler {
	return &ImageUpscaler{
		apiKey:  apiKey,
		apiURL:  apiURL,
		enabled: enabled,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // AI upscaling can take time
		},
	}
}

// UpscaleImage takes an original image URL and returns an upscaled version
// If upscaling fails or is disabled, returns the original URL
func (u *ImageUpscaler) UpscaleImage(originalURL string) string {
	// Return original if upscaling is disabled
	if !u.enabled {
		return originalURL
	}

	// Return original if no URL provided
	if originalURL == "" {
		return originalURL
	}

	// Return original if API not configured
	if u.apiURL == "" || u.apiKey == "" {
		log.Printf("[Upscaler] API not configured, using original photo")
		return originalURL
	}

	log.Printf("[Upscaler] Upscaling image: %s", originalURL)

	// Attempt upscaling with retry logic
	upscaledURL, err := u.upscaleWithRetry(originalURL, 3)
	if err != nil {
		log.Printf("[Upscaler] Failed to upscale image: %v. Using original.", err)
		return originalURL
	}

	log.Printf("[Upscaler] Successfully upscaled image: %s", upscaledURL)
	return upscaledURL
}

// upscaleWithRetry attempts to upscale with exponential backoff retry
func (u *ImageUpscaler) upscaleWithRetry(imageURL string, maxRetries int) (string, error) {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		upscaledURL, err := u.callUpscaleAPI(imageURL)
		if err == nil {
			return upscaledURL, nil
		}

		lastErr = err
		if attempt < maxRetries {
			waitTime := time.Duration(attempt) * time.Second
			log.Printf("[Upscaler] Retry %d/%d after %v", attempt, maxRetries, waitTime)
			time.Sleep(waitTime)
		}
	}

	return "", fmt.Errorf("all %d attempts failed: %w", maxRetries, lastErr)
}

// callUpscaleAPI makes the actual API call to Replicate's Real-ESRGAN upscaling service
func (u *ImageUpscaler) callUpscaleAPI(imageURL string) (string, error) {
	// Download the image first to bypass hotlink protection
	log.Printf("[Upscaler] Downloading image from: %s", imageURL)
	imageData, contentType, err := u.downloadImage(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}

	// Convert to base64 data URL
	base64Data := base64.StdEncoding.EncodeToString(imageData)
	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
	log.Printf("[Upscaler] Converted image to base64 data URL (%d bytes)", len(imageData))

	// Replicate API payload for Real-ESRGAN model
	// Model: nightmareai/real-esrgan (popular 2x upscaling model)
	payload := map[string]interface{}{
		"version": "42fed1c4974146d4d2414e2be2c5277c7fcf05fcc3a73abf41610695738c1d7b",
		"input": map[string]interface{}{
			"image":        dataURL, // Use base64 data URL instead of direct URL
			"scale":        2,       // 2x upscaling
			"face_enhance": false,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create prediction
	req, err := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "wait") // Wait for result instead of polling

	// Make request
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Replicate response
	var result struct {
		Status string      `json:"status"`
		Output interface{} `json:"output"`
		Error  string      `json:"error,omitempty"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for errors
	if result.Error != "" {
		return "", fmt.Errorf("Replicate error: %s", result.Error)
	}

	// Check status
	if result.Status != "succeeded" {
		return "", fmt.Errorf("prediction status: %s", result.Status)
	}

	// Extract output URL
	// Replicate returns output as a string URL or array of URLs
	if outputStr, ok := result.Output.(string); ok {
		return outputStr, nil
	}

	if outputArray, ok := result.Output.([]interface{}); ok && len(outputArray) > 0 {
		if url, ok := outputArray[0].(string); ok {
			return url, nil
		}
	}

	return "", fmt.Errorf("invalid output format in Replicate response")
}

// downloadImage downloads an image from a URL and returns the image data and content type
func (u *ImageUpscaler) downloadImage(imageURL string) ([]byte, string, error) {
	// Create request with user agent to avoid bot detection
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
	req.Header.Set("Referer", "https://www.basketball-reference.com/")

	// Download the image
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Get content type (default to image/jpeg if not specified)
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	return imageData, contentType, nil
}
