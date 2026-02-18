package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Auth0Client defines the interface for Auth0 operations
type Auth0Client interface {
	GetManagementToken() (string, error)
	UpdateUserMetadata(userId, username, managementToken string) error
}

// auth0Client is the concrete implementation of Auth0Client
type auth0Client struct{}

// NewAuth0Client creates a new Auth0Client
func NewAuth0Client() Auth0Client {
	return &auth0Client{}
}

// Auth0ManagementTokenResponse represents the response from Auth0 token endpoint
type Auth0ManagementTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetManagementToken obtains an access token for the Auth0 Management API
func (a *auth0Client) GetManagementToken() (string, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_MANAGEMENT_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_MANAGEMENT_CLIENT_SECRET")

	// The audience for Management API must be https://{domain}/api/v2/
	// This is different from AUTH0_AUDIENCE which is for user JWT validation
	audience := fmt.Sprintf("https://%s/api/v2/", domain)

	if domain == "" || clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("missing required Auth0 environment variables")
	}

	tokenURL := fmt.Sprintf("https://%s/oauth/token", domain)

	requestBody := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to request management token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Auth0 token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse Auth0ManagementTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResponse.AccessToken, nil
}

// UpdateUserMetadata updates the user_metadata for a user in Auth0
func (a *auth0Client) UpdateUserMetadata(userId, username, managementToken string) error {
	domain := os.Getenv("AUTH0_DOMAIN")
	if domain == "" {
		return fmt.Errorf("AUTH0_DOMAIN environment variable not set")
	}

	// Construct the Auth0 Management API URL
	updateURL := fmt.Sprintf("https://%s/api/v2/users/%s", domain, userId)

	// Prepare the update payload
	updatePayload := map[string]interface{}{
		"user_metadata": map[string]string{
			"athlete_unknown_username": username,
		},
	}

	jsonBody, err := json.Marshal(updatePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal update payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("PATCH", updateURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+managementToken)
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute update request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Auth0 user update failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
