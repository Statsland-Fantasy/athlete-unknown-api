package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// ImageUpscaler handles AI-powered image upscaling
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

// callUpscaleAPI makes the actual API call to the upscaling service
func (u *ImageUpscaler) callUpscaleAPI(imageURL string) (string, error) {
	// Prepare request payload
	// NOTE: Adjust this structure based on your specific AI upscaling service
	// Examples: Replicate, DeepAI, Let's Enhance, etc.
	payload := map[string]interface{}{
		"image": imageURL,
		"scale": 2, // 2x upscaling
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", u.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	// NOTE: Adjust this based on your specific AI service's response format
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract upscaled URL
	// NOTE: This field name may vary by service - adjust accordingly
	upscaledURL, ok := result["output"].(string)
	if !ok {
		// Try alternative field names
		if url, exists := result["url"].(string); exists {
			return url, nil
		}
		if url, exists := result["result"].(string); exists {
			return url, nil
		}
		return "", fmt.Errorf("invalid response format: no URL found")
	}

	return upscaledURL, nil
}
