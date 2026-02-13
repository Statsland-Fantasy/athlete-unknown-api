package main

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// ValidateSportsReferenceURL validates that a URL is safe to scrape
// Returns an error if the URL is invalid or not whitelisted
func ValidateSportsReferenceURL(urlStr string) error {
	// Parse the URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Ensure URL has a scheme (http or https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https protocol, got: %s", parsedURL.Scheme)
	}

	// Ensure URL has a host
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a hostname")
	}

	// Extract hostname (without port)
	hostname := parsedURL.Hostname()

	// Check if hostname is in the whitelist
	isAllowed := false
	for _, allowedDomain := range AllowedScrapingDomains {
		if hostname == allowedDomain {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("URL hostname '%s' is not in the allowed whitelist. Allowed domains: %v", hostname, AllowedScrapingDomains)
	}

	// Prevent SSRF by ensuring the hostname doesn't resolve to a private IP
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve hostname: %w", err)
	}

	for _, ip := range ips {
		if isPrivateIP(ip) {
			return fmt.Errorf("URL resolves to a private IP address: %s", ip.String())
		}
	}

	return nil
}

// isPrivateIP checks if an IP address is private/internal
func isPrivateIP(ip net.IP) bool {
	// Check for loopback
	if ip.IsLoopback() {
		return true
	}

	// Check for private networks
	privateNetworks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16", // Link-local
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link-local
	}

	for _, cidr := range privateNetworks {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// getUserTimezone extracts and validates the user's timezone from the request header
// Returns the timezone location or UTC as fallback if header is missing/invalid
func getUserTimezone(c *gin.Context) *time.Location {
	// Get timezone from X-User-Timezone header
	tzHeader := c.GetHeader("X-User-Timezone")

	if tzHeader == "" {
		// No header provided, use UTC as fallback
		return time.UTC
	}

	// Validate and load the timezone
	loc, err := time.LoadLocation(tzHeader)
	if err != nil {
		// Invalid timezone string, use UTC as fallback
		// Could log this for monitoring: fmt.Printf("Invalid timezone header: %s\n", tzHeader)
		return time.UTC
	}

	return loc
}

// updateDailyStreak updates the user's daily streak based on real-life calendar days (engagement-based)
// This function tracks consecutive days the user plays ANY round, regardless of which round's playDate they choose
// For existing users, it:
// - Increments the streak if currentDate is exactly 1 day after the last real-life day they played
// - Resets the streak to 1 if more than 1 day has passed since the last real-life day they played
// - Keeps the streak unchanged if they already played on currentDate (prevents multiple increments on same day)
// Always updates LastDayPlayed to currentDate
// The currentDate parameter should be the real-life date (typically today's date from time.Now())
func updateDailyStreak(user *User, currentDate string) string {
	storyIdAchieved := ""
	if user == nil {
		return storyIdAchieved
	}

	// Check if we have a previous play date to compare against
	if user.LastDayPlayed != "" {
		// If they already played on this date, don't update the streak (prevents multiple increments)
		if user.LastDayPlayed == currentDate {
			return storyIdAchieved
		}

		lastPlayed, err := time.Parse(DateFormatYYYYMMDD, user.LastDayPlayed)
		currentParsed, err2 := time.Parse(DateFormatYYYYMMDD, currentDate)

		if err == nil && err2 == nil {
			daysDiff := int(currentParsed.Sub(lastPlayed).Hours() / 24)

			if daysDiff == 1 {
				// Consecutive day - increment streak
				user.CurrentDailyStreak++
				storyIdAchieved = currentDailyStreakStoryMissions(user.CurrentDailyStreak)
			} else if daysDiff > 1 {
				// Missed a day - reset streak to 1
				user.CurrentDailyStreak = 1
			}
			// If daysDiff == 0, same day - don't change streak (shouldn't happen due to check above)
		}
	}

	// Update last day played to the current date
	user.LastDayPlayed = currentDate
	return storyIdAchieved
}
