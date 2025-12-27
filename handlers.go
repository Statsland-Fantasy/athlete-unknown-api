package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

// Global DB instance
var db *DB

// handleGetRound handles GET /v1/round
func handleGetRound(c *gin.Context) {
	sport := c.Query("sport")
	if sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Sport parameter is required",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	// Validate sport
	if sport != "basketball" && sport != "baseball" && sport != "football" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Invalid sport parameter. Must be basketball, baseball, or football",
			"code":      "INVALID_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	playDate := c.Query("playDate")
	if playDate == "" {
		playDate = time.Now().Format("2006-01-02")
	}

	ctx := context.Background()
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to retrieve round: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Not Found",
			"message":   "No round found for the specified sport and playDate",
			"code":      "ROUND_NOT_FOUND",
			"timestamp": time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, round)
}

// handleCreateRound handles PUT /v1/round
func handleCreateRound(c *gin.Context) {
	var round Round
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Invalid request body: " + err.Error(),
			"code":      "INVALID_REQUEST_BODY",
			"timestamp": time.Now(),
		})
		return
	}

	// Validate required fields
	if round.Sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Missing required field: sport",
			"code":      "MISSING_REQUIRED_FIELD",
			"timestamp": time.Now(),
		})
		return
	}
	if round.PlayDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Missing required field: playDate",
			"code":      "MISSING_REQUIRED_FIELD",
			"timestamp": time.Now(),
		})
		return
	}
	if round.Player.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Missing required field: player.name",
			"code":      "MISSING_REQUIRED_FIELD",
			"timestamp": time.Now(),
		})
		return
	}

	// Set Created and LastUpdated timestamps if not provided
	now := time.Now()
	if round.Created.IsZero() {
		round.Created = now
	}
	if round.LastUpdated.IsZero() {
		round.LastUpdated = now
	}

	// Generate round ID
	roundID, err := GenerateRoundID(round.Sport, round.PlayDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Invalid playDate format: " + err.Error(),
			"code":      "INVALID_PLAY_DATE",
			"timestamp": time.Now(),
		})
		return
	}
	round.RoundID = roundID

	ctx := context.Background()
	err = db.CreateRound(ctx, &round)
	if err != nil {
		if err.Error() == "round already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error":     "Conflict",
				"message":   "Round already exists for sport '" + round.Sport + "' on playDate '" + round.PlayDate + "'",
				"code":      "ROUND_ALREADY_EXISTS",
				"timestamp": time.Now(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to create round: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, round)
}

// handleDeleteRound handles DELETE /v1/round
func handleDeleteRound(c *gin.Context) {
	sport := c.Query("sport")
	if sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Missing required parameter: sport",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	playDate := c.Query("playDate")
	if playDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Missing required parameter: playDate",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	ctx := context.Background()

	// Check if the round exists first
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to check round existence: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}
	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Not Found",
			"message":   "Round not found for sport '" + sport + "' on playDate '" + playDate + "'",
			"code":      "ROUND_NOT_FOUND",
			"timestamp": time.Now(),
		})
		return
	}

	err = db.DeleteRound(ctx, sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to delete round: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// handleGetUpcomingRounds handles GET /v1/upcoming-rounds
func handleGetUpcomingRounds(c *gin.Context) {
	sport := c.Query("sport")
	if sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Sport parameter is required",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	ctx := context.Background()
	upcomingRounds, err := db.GetRoundsBySport(ctx, sport, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to retrieve rounds: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	if len(upcomingRounds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Not Found",
			"message":   "No upcoming rounds found for sport '" + sport + "' in the specified date range",
			"code":      "NO_UPCOMING_ROUNDS",
			"timestamp": time.Now(),
		})
		return
	}

	// Sort by playDate
	sort.Slice(upcomingRounds, func(i, j int) bool {
		return upcomingRounds[i].PlayDate < upcomingRounds[j].PlayDate
	})

	c.JSON(http.StatusOK, upcomingRounds)
}

// handleSubmitResults handles POST /v1/results
func handleSubmitResults(c *gin.Context) {
	sport := c.Query("sport")
	playDate := c.Query("playDate")

	if sport == "" || playDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Sport and playDate parameters are required",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	var result Result
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Invalid request body: " + err.Error(),
			"code":      "INVALID_REQUEST_BODY",
			"timestamp": time.Now(),
		})
		return
	}

	// potential hack catcher. Score cannot be higher than 100
	if result.Score > 100 || result.Score < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Invalid request body: Score cannot be greater than 100 or less than 0",
			"code":      "INVALID_REQUEST_BODY",
			"timestamp": time.Now(),
		})
		return		
	}

	ctx := context.Background()

	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to retrieve round: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}
	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Not Found",
			"message":   "Round not found for sport '" + sport + "' on date '" + playDate + "'",
			"code":      "ROUND_NOT_FOUND",
			"timestamp": time.Now(),
		})
		return
	}

	// Update round statistics
	updateStatsWithResult(&round.Stats.Stats, &result)

	// Save the updated round
	err = db.UpdateRound(ctx, round)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to update round: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	// Get user_id from bearer token (set by JWT middleware)
	userIdToken, exists := c.Get("userId")
	if exists && userIdToken != "" {
		userId, ok := userIdToken.(string)
		if ok && userId != "" {
			// Fetch existing user stats or create new ones
			userStats, err := db.GetUserStats(ctx, userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "Internal Server Error",
					"message":   "Failed to retrieve user stats: " + err.Error(),
					"code":      "DATABASE_ERROR",
					"timestamp": time.Now(),
				})
				return
			}

			// If user stats don't exist, create new user stats
			if userStats == nil {
				userStats = &UserStats{
					UserId:  userId,
					Sports:  []UserSportStats{},
					CurrentDailyStreak: 1,
					LastDayPlayed: playDate,
					UserName: "", // TODO: update with user's username as fetched from Auth0
				}
			} else {
				// Update daily streak based on play date
				updateDailyStreak(userStats, playDate)
			}

			// Find or create specific sport stats
			var sportStats *UserSportStats
			for i := range userStats.Sports {
				if userStats.Sports[i].Sport == sport {
					sportStats = &userStats.Sports[i]
					break
				}
			}

			// If sport stats don't exist, create new entry
			if sportStats == nil {
				newSportStats := UserSportStats{
					Sport: sport,
				}
				userStats.Sports = append(userStats.Sports, newSportStats)
				sportStats = &userStats.Sports[len(userStats.Sports)-1]
			}			

			// Update sport-specific stats
			updateStatsWithResult(&sportStats.Stats, &result)

			// Create round history entry
			roundHistory := RoundHistory{
				PlayDate: playDate,
				Result:   result,
			}
			sportStats.History = append(sportStats.History, roundHistory)

			// Save or update user stats in DynamoDB
			if userStats.UserCreated.IsZero() {
				err = db.CreateUserStats(ctx, userStats)
			} else {
				err = db.UpdateUserStats(ctx, userStats)
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":     "Internal Server Error",
					"message":   "Failed to update user stats: " + err.Error(),
					"code":      "DATABASE_ERROR",
					"timestamp": time.Now(),
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, result)
}

// handleGetRoundStats handles GET /v1/stats/round
func handleGetRoundStats(c *gin.Context) {
	sport := c.Query("sport")
	playDate := c.Query("playDate")

	if sport == "" || playDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Sport and playDate parameters are required",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	ctx := context.Background()
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to retrieve round: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Not Found",
			"message":   "No statistics found for sport '" + sport + "' on date '" + playDate + "'",
			"code":      "STATS_NOT_FOUND",
			"timestamp": time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, round.Stats)
}

// handleGetUserStats handles GET /v1/stats/user
func handleGetUserStats(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		// if userId is not part in query param, extract from bearer token instead
		userIdToken, exists := c.Get("userId")
		if exists && userIdToken != "" {
			userIdStr, ok := userIdToken.(string)
			if ok && userIdStr != "" {
				userId = userIdStr
			}
		}

		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Bad Request",
				"message":   "userId parameter is required",
				"code":      "MISSING_REQUIRED_PARAMETER",
				"timestamp": time.Now(),
			})
			return
		}
	}

	ctx := context.Background()
	stats, err := db.GetUserStats(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to retrieve user stats: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Not Found",
			"message":   "No statistics found for user '" + userId + "'",
			"code":      "USER_STATS_NOT_FOUND",
			"timestamp": time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// handleScrapeAndCreateRound handles POST /v1/round - scrapes player data and creates a round
func handleScrapeAndCreateRound(c *gin.Context) {
	// Get required parameters
	sport := c.Query("sport")
	playDate := c.Query("playDate")

	// Validate required parameters
	if sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Sport parameter is required",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}
	if playDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "playDate parameter is required",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	// Validate sport
	if sport != "basketball" && sport != "baseball" && sport != "football" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Invalid sport parameter. Must be basketball, baseball, or football",
			"code":      "INVALID_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	// Get optional parameters
	name := c.Query("name")
	sportsReferenceURL := c.Query("sportsReferenceURL")
	theme := c.Query("theme")

	// Validate that at least one optional parameter is provided
	if name == "" && sportsReferenceURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Either 'name' or 'sportsReferenceURL' parameter must be provided",
			"code":      "MISSING_REQUIRED_PARAMETER",
			"timestamp": time.Now(),
		})
		return
	}

	// Get the hostname for the sport
	hostname := GetSportsReferenceHostname(sport)
	if hostname == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Unable to determine hostname for sport: " + sport,
			"code":      "CONFIGURATION_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	// If sportsReferenceURL is provided, go directly to the player page
	if sportsReferenceURL != "" {
		// Use the URL directly (it should be a full URL)
		playerURL := sportsReferenceURL
		fmt.Printf("Player page URL: %s\n", playerURL)

		// Scrape player page data
		player, err := scrapePlayerData(playerURL, hostname, sport)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal Server Error",
				"message":   "Failed to scrape player data: " + err.Error(),
				"code":      "SCRAPING_ERROR",
				"timestamp": time.Now(),
			})
			return
		}

		// Create round with scraped player data
		now := time.Now()
		roundID, err := GenerateRoundID(sport, playDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Bad Request",
				"message":   "Invalid playDate format: " + err.Error(),
				"code":      "INVALID_PLAY_DATE",
				"timestamp": time.Now(),
			})
			return
		}

		round := Round{
			RoundID:     roundID,
			Sport:       sport,
			PlayDate:    playDate,
			Player:      *player,
			Created:     now,
			LastUpdated: now,
			Theme:       theme,
			Stats: RoundStats{
				PlayDate: playDate,
				Name:     player.Name,
				Sport:    sport,
			},
		}

		// Store the round in DynamoDB
		err = db.CreateRound(c.Request.Context(), &round)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal Server Error",
				"message":   "Failed to create round: " + err.Error(),
				"code":      "DATABASE_ERROR",
				"timestamp": time.Now(),
			})
			return
		}

		c.JSON(http.StatusCreated, round)
		return
	}

	// Search for the player by name
	encodedName := url.QueryEscape(name)
	searchURL := fmt.Sprintf("https://www.%s/search/search.fcgi?search=%s", hostname, encodedName)

	// Initialize colly collector
	// Allow both www and non-www versions of the domain
	collector := colly.NewCollector(
		colly.AllowedDomains(hostname, "www."+hostname),
		colly.MaxDepth(1),
	)

	// Allow redirects between www and non-www versions
	collector.AllowURLRevisit = false

	// Variable to capture the final URL after redirects
	var finalURL string
	var scrapeError error
	var playerSearchItems []string

	// Set up error handling
	collector.OnError(func(r *colly.Response, err error) {
		scrapeError = err
		fmt.Printf("Scraping error: %v\n", err)
	})

	// Log request
	collector.OnRequest(func(r *colly.Request) {
		// fmt.Printf("Visiting: %s\n", r.URL.String())
	})

	// Capture the response and final URL
	collector.OnResponse(func(r *colly.Response) {
		finalURL = r.Request.URL.String()
		// fmt.Printf("Final URL after redirects: %s\n", finalURL)
	})

	// Extract player search results from #players div
	collector.OnHTML("div#players div.search-item", func(e *colly.HTMLElement) {
		// Get the player URL path from the search-item-url div text
		playerURLPath := strings.TrimSpace(e.ChildText("div.search-item-url"))
		if playerURLPath != "" {
			playerSearchItems = append(playerSearchItems, playerURLPath)
		}
	})

	// Visit the search URL
	err := collector.Visit(searchURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to initiate scraping: " + err.Error(),
			"code":      "SCRAPING_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	// Check if there was a scraping error
	if scrapeError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to scrape player data: " + scrapeError.Error(),
			"code":      "SCRAPING_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	// Check if the final URL contains "/players" (meaning it redirected to a specific player page)
	if !strings.Contains(finalURL, "/players") {
		// Check if there's exactly one player result in the search results
		if len(playerSearchItems) == 1 {
			// Use the single player result
			finalURL = fmt.Sprintf("https://www.%s%s", hostname, playerSearchItems[0])
			// fmt.Printf("Found single player result, using URL: %s\n", finalURL)
		} else if len(playerSearchItems) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Bad Request",
				"message":   "No players found with the name '" + name + "'. Please check the player name and sport again.",
				"code":      "NO_PLAYERS_FOUND",
				"timestamp": time.Now(),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Bad Request",
				"message":   "Multiple players found with the name '" + name + "'. Please provide the sportsReferenceURL parameter to specify the exact player.",
				"code":      "MULTIPLE_PLAYERS_FOUND",
				"timestamp": time.Now(),
			})
			return
		}
	}

	// Successfully redirected to a player page
	// fmt.Printf("Player page URL: %s\n", finalURL)

	// Scrape player page data
	player, err := scrapePlayerData(finalURL, hostname, sport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to scrape player data: " + err.Error(),
			"code":      "SCRAPING_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	// Create round with scraped player data
	now := time.Now()
	roundID, err := GenerateRoundID(sport, playDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Bad Request",
			"message":   "Invalid playDate format: " + err.Error(),
			"code":      "INVALID_PLAY_DATE",
			"timestamp": time.Now(),
		})
		return
	}

	round := Round{
		RoundID:     roundID,
		Sport:       sport,
		PlayDate:    playDate,
		Player:      *player,
		Created:     now,
		LastUpdated: now,
		Theme:       theme,
		Stats: RoundStats{
			PlayDate: playDate,
			Name:     player.Name,
			Sport:    sport,
		},
	}

	// Store the round in DynamoDB
	err = db.CreateRound(c.Request.Context(), &round)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Internal Server Error",
			"message":   "Failed to create round: " + err.Error(),
			"code":      "DATABASE_ERROR",
			"timestamp": time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, round)
}

// scrapePlayerData scrapes player information from a sports reference page
func scrapePlayerData(playerURL, hostname, sport string) (*Player, error) {
	player := &Player{
		Sport:              sport,
		SportsReferenceURL: playerURL,
	}

	// Initialize colly collector
	c := colly.NewCollector(
		colly.AllowedDomains(hostname, "www."+hostname),
		colly.MaxDepth(1),
	)

	// Allow redirects between www and non-www versions
	c.AllowURLRevisit = false

	var scrapeError error

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		scrapeError = err
		fmt.Printf("Scraping error: %v\n", err)
	})

	// Extract player name from the page title or h1
	c.OnHTML("h1[itemprop='name'], h1 span", func(e *colly.HTMLElement) {
		if player.Name == "" {
			player.Name = strings.TrimSpace(e.Text)
		}
	})

	// Extract bio information (birth date and location)
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		// Look for birth information
		if strings.Contains(text, "Born:") || strings.Contains(text, "born") {
			// Remove newlines and extra spaces
			text = strings.ReplaceAll(text, "\n", " ")
			text = strings.Join(strings.Fields(text), " ")
			// Remove country code (last 3 characters: space + 2-char code)
			if sport != "football" && len(text) > 3 {
				text = strings.TrimSpace(text[:len(text)-3])
			}
			player.Bio = text
		}
	})

	// Extract player information (height, weight, position, handedness)
	// Collect all physical attributes into a slice first
	var physicalAttrs []string
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		// Look for physical attributes
		if strings.Contains(text, "Position:") || strings.Contains(text, "Positions:") || strings.Contains(text, "Bats:") ||
		   strings.Contains(text, "Throws:") || strings.Contains(text, "Shoots:") {
			// Remove newlines and extra spaces
			text = strings.ReplaceAll(text, "\n", " ")
			text = strings.ReplaceAll(text, "-", ", ") // football uses - instead of ,
			if sport == "football" {
				text = strings.ReplaceAll(text, "Throws:", " ▪ Throws:") // football has Throws in same line as position. Need extra separator
			}
			text = strings.Join(strings.Fields(text), " ")
			physicalAttrs = append(physicalAttrs, text)
		}
	})

	// Extract height and weight from span elements
	c.OnHTML("div#meta span[itemprop='height'], div#meta span[itemprop='weight']", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if text != "" {
			// Remove metric measurements in parentheses (e.g., "(193cm, 88kg)")
			// Use regex to remove anything in parentheses
			re := regexp.MustCompile(`\s*\([^)]*\)`)
			text = re.ReplaceAllString(text, "")
			text = strings.TrimSpace(text)

			// Get the label (Height or Weight) with bullet separator
			if e.Attr("itemprop") == "height" {
				physicalAttrs = append(physicalAttrs, " ▪ Height: "+text)
			} else if e.Attr("itemprop") == "weight" {
				physicalAttrs = append(physicalAttrs, " ▪ Weight: "+text)
			}
		}
	})

	// Alternative: Look for height/weight in paragraphs
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if (strings.Contains(text, "cm") && strings.Contains(text, "kg")) ||
		   (strings.Contains(text, "lb") && (strings.Contains(text, "-") || strings.Contains(text, "'"))) {
			// This likely contains height and weight
			text = strings.ReplaceAll(text, "\n", " ")
			text = strings.Join(strings.Fields(text), " ")

			// Remove metric measurements in parentheses
			re := regexp.MustCompile(`\s*\([^)]*\)`)
			text = re.ReplaceAllString(text, "")
			text = strings.TrimSpace(text)

			physicalAttrs = append(physicalAttrs, text)
		}
	})

	// Extract college information for football/basketball
	// Prioritize college over high school
	var college string
	var highSchool string
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if strings.Contains(text, "College:") || strings.Contains(text, "Colleges:") {
			// Extract college name after the colon
			parts := strings.SplitN(text, ":", 2)
			if len(parts) == 2 {
				college = strings.TrimSpace(parts[1])
				// Remove "(College Stats)" suffix if present
				college = strings.TrimSpace(strings.ReplaceAll(college, "(College Stats)", ""))

				// If multiple colleges separated by comma, choose the latter (most recent)
				if strings.Contains(college, ",") {
					colleges := strings.Split(college, ",")
					college = strings.TrimSpace(colleges[len(colleges)-1])
				}
			}
		} else if strings.Contains(text, "High School:") {
			// Extract high school name after the colon (as fallback)
			parts := strings.SplitN(text, ":", 2)
			if len(parts) == 2 {
				highSchool = strings.TrimSpace(parts[1])
			}
		}
	})

	// Extract draft information
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if strings.Contains(text, "Draft:") || strings.Contains(text, "Drafted") {
			// Remove extra spaces - normalize to single space between words
			text = strings.ReplaceAll(text, "\n", " ")
			text = strings.ReplaceAll(text, "\t", " ")
			text = strings.Join(strings.Fields(text), " ")

			// Use college if available, otherwise fall back to high school
			school := college
			if school == "" {
				school = highSchool
			}
			school = strings.ReplaceAll(school, "\n", "")
			school = strings.ReplaceAll(school, "\t", "")

			player.DraftInformation = formatDraftInformation(text, sport, school)
		}

		// Set default for draft information if not found
		if player.DraftInformation == "" {
			player.DraftInformation = "Undrafted"
		}
	})

	// Extract years active from the first table
	var years []string
	var injuredYears []string
	var firstTableProcessed bool
	c.OnHTML("table", func(e *colly.HTMLElement) {
		// Skip tables with id="last5" (last 5 games tables)
		if e.Attr("id") == "last5" {
			return
		}

		// Only process the very first table on the page
		if firstTableProcessed {
			return
		}
		firstTableProcessed = true

		// Extract all th elements from this table
		e.ForEach("th[data-stat='year_id']", func(_ int, el *colly.HTMLElement) {
			year := strings.TrimSpace(el.Text)
			if isValidYear(year) && !contains(years, year) {
				years = append(years, year)
			}
		})
	})	

	// Extract team information - just get ALL td with team data
	var teams []string
	c.OnHTML("td[data-stat='team_name_abbr']", func(e *colly.HTMLElement) {
		team := strings.TrimSpace(e.Text)

		// Check if player was injured (didn't play)
		if strings.Contains(team, "Did not play") || strings.Contains(team, "Injured") {
			// Find the corresponding year from the parent row (tr)
			yearElement := e.DOM.Parent().Find("th[data-stat='year_ID'], th[data-stat='year_id']")
			if yearElement.Length() > 0 {
				injuredYear := strings.TrimSpace(yearElement.Text())
				if isValidYear(injuredYear) && !contains(injuredYears, injuredYear) {
					injuredYears = append(injuredYears, injuredYear)
				}
			}
			return
		}

		// Skip entries with "TM" (total/multi-team rows) and empty entries
		if team != "" && team != "TM" && !strings.Contains(team, "TM") && !contains(teams, team) {
			teams = append(teams, team)
		}
	})

	// Capture the stats_pullout element for later processing
	var statsPulloutElement *colly.HTMLElement
	c.OnHTML(".stats_pullout", func(e *colly.HTMLElement) {
		statsPulloutElement = e
	})

	// Extract jersey numbers using uni_holder class
	var jerseyNumbersMap = make(map[string]bool)
	c.OnHTML(".uni_holder", func(e *colly.HTMLElement) {
		// Remove all newlines, tabs, and extra spaces
		text := e.Text
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ReplaceAll(text, "\t", "")

		// Deduplicate fields before joining
		fields := strings.Fields(text)
		uniqueFields := make(map[string]bool)
		var dedupedFields []string
		for _, field := range fields {
			if !uniqueFields[field] {
				uniqueFields[field] = true
				dedupedFields = append(dedupedFields, field)
			}
		}
		// Only join with comma if there are multiple fields
		if len(dedupedFields) > 1 {
			text = strings.Join(dedupedFields, ", ")
		} else if len(dedupedFields) == 1 {
			text = dedupedFields[0]
		}

		// Remove footnote markers like +1, +2, etc.
		// These appear as superscripts in the HTML
		for i := 0; i <= 9; i++ {
			text = strings.ReplaceAll(text, fmt.Sprintf(", +%d", i), "") // remove accompanying ", " as well
		}
		text = strings.TrimSpace(text)

		if text != "" {
			jerseyNumbersMap[text] = true
		}
	})

	// After scraping, join teams, years, jersey numbers, and physical attributes
	c.OnScraped(func(r *colly.Response) {
		if len(teams) > 0 {
			player.TeamsPlayedOn = strings.Join(teams, ", ")
		}
		if len(years) > 0 {
			// Filter out injured years from active years
			var activeYears []string
			for _, year := range years {
				if !contains(injuredYears, year) {
					activeYears = append(activeYears, year)
				}
			}
			if len(activeYears) > 0 {
				player.YearsActive = formatYearsAsRanges(activeYears, sport)
			}
		}

		// Convert map to slice for jersey numbers (automatically deduplicates)
		var jerseyNumbers []string
		for num := range jerseyNumbersMap {
			jerseyNumbers = append(jerseyNumbers, num)
		}
		if len(jerseyNumbers) > 0 {
			player.JerseyNumbers = strings.Join(jerseyNumbers, ", ")
		}

		if len(physicalAttrs) > 0 {
			player.PlayerInformation = strings.Join(physicalAttrs, " ▪ ")
			// Abbreviate positions in the player information
			player.PlayerInformation = abbreviatePositions(player.PlayerInformation)
		}

		// Extract career stats using the config based on sport and position
		if statsPulloutElement != nil {
			careerStatsConfig := GetCareerStatsConfig(sport, player.PlayerInformation)

			// baseball pitcher-only
			winsOrSavesValue := 0
			var winsOrSavesLabel string

			// Extract each stat from the configuration
			var careerStats []string
			for _, statConfig := range careerStatsConfig.Stats {
				statValue := strings.TrimSpace(statsPulloutElement.DOM.Find(statConfig.HTMLPath).Text())
				if statValue != "" {
					// logic to correctly display the higher of Wins or Saves
					if statConfig.StatLabel == "W" || statConfig.StatLabel == "SV" {
						intStatValue, err := strconv.Atoi(statValue)
						if err != nil {
							// If conversion fails, skip this stat and continue
							fmt.Printf("Warning: Failed to convert %s value '%s' to integer: %v\n", statConfig.StatLabel, statValue, err)
							continue
						}
						if intStatValue > winsOrSavesValue {
							winsOrSavesValue = intStatValue
							winsOrSavesLabel = statConfig.StatLabel
						}
					} else {
						careerStats = append(careerStats, fmt.Sprintf("%s %s", statValue, statConfig.StatLabel))
					}
				}
			}

			// baseball pitcher-only
			if winsOrSavesLabel != "" {
				careerStats = append([]string{fmt.Sprintf("%d %s", winsOrSavesValue, winsOrSavesLabel)}, careerStats...)
			}

			// Join all stats into the CareerStats field
			if len(careerStats) > 0 {
				player.CareerStats = strings.Join(careerStats, ", ")
			}
		}


	})

	// Extract career stats based on sport
	if sport == "baseball" {
		// For baseball, we need to determine if hitter or pitcher
		c.OnHTML("div#info strong", func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if strings.Contains(text, "Career:") {
				player.CareerStats = text
			}
		})

		// Alternative: extract from career stats summary
		c.OnHTML("div#meta div.p1", func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if player.CareerStats == "" {
				player.CareerStats = text
			}
		})
	} else if sport == "basketball" {
		// For basketball, look for career stats
		c.OnHTML("div#info div", func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if strings.Contains(text, "Career:") || strings.Contains(text, "Points") {
				player.CareerStats = text
			}
		})
	} else if sport == "football" {
		// For football, look for career stats
		c.OnHTML("div#meta div", func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if strings.Contains(text, "Career:") {
				player.CareerStats = text
			}
		})
	}

	// Extract personal achievements (awards, honors, championships)
	var achievements []string
	c.OnHTML("ul#bling li", func(e *colly.HTMLElement) {
		achievement := strings.TrimSpace(e.Text)
		if achievement != "" {
			achievements = append(achievements, achievement)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		if len(achievements) > 0 {
			maxLength := 100
			processedAchievements := ProcessAchievements(sport, achievements, maxLength) 
			player.PersonalAchievements = processedAchievements
		} else {
			player.PersonalAchievements = "N/A"
		}
	})

	// Extract photo URL - try multiple selectors
	c.OnHTML("div#meta img", func(e *colly.HTMLElement) {
		if player.Photo == "" {
			src := e.Attr("src")
			// Only capture actual image URLs, not placeholder icons
			if src != "" && (strings.Contains(src, ".jpg") || strings.Contains(src, ".png") || strings.Contains(src, ".jpeg")) {
				player.Photo = src
			}
		}
	})

	// Alternative photo selector
	c.OnHTML("img[itemProp='image']", func(e *colly.HTMLElement) {
		if player.Photo == "" {
			player.Photo = e.Attr("src")
		}
	})

	// Another alternative - media-item class
	c.OnHTML("img.media-item", func(e *colly.HTMLElement) {
		if player.Photo == "" {
			player.Photo = e.Attr("src")
		}
	})

	// Visit the player page
	err := c.Visit(playerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to visit player page: %w", err)
	}

	if scrapeError != nil {
		return nil, fmt.Errorf("scraping error: %w", scrapeError)
	}

	return player, nil
}
