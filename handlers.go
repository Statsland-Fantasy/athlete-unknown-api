package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Global DB instance
var db *DB

// handleGetRound handles GET /v1/round
func handleGetRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	if sport == "" {
		errorResponseWithCode(w, "Bad Request", "Sport parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	// Validate sport
	if sport != "basketball" && sport != "baseball" && sport != "football" {
		errorResponseWithCode(w, "Bad Request", "Invalid sport parameter. Must be basketball, baseball, or football", "INVALID_PARAMETER", http.StatusBadRequest)
		return
	}

	playDate := r.URL.Query().Get("playDate")
	if playDate == "" {
		playDate = time.Now().Format("2006-01-02")
	}

	ctx := context.Background()
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	if round == nil {
		errorResponseWithCode(w, "Not Found", "No round found for the specified sport and playDate", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	jsonResponse(w, round, http.StatusOK)
}

// handleCreateRound handles PUT /v1/round
func handleCreateRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errorResponseWithCode(w, "Method Not Allowed", "Only PUT method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	var round Round
	if err := json.NewDecoder(r.Body).Decode(&round); err != nil {
		errorResponseWithCode(w, "Bad Request", "Invalid request body: "+err.Error(), "INVALID_REQUEST_BODY", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if round.Sport == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required field: sport", "MISSING_REQUIRED_FIELD", http.StatusBadRequest)
		return
	}
	if round.PlayDate == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required field: playDate", "MISSING_REQUIRED_FIELD", http.StatusBadRequest)
		return
	}
	if round.Player.Name == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required field: player.name", "MISSING_REQUIRED_FIELD", http.StatusBadRequest)
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
		errorResponseWithCode(w, "Bad Request", "Invalid playDate format: "+err.Error(), "INVALID_PLAY_DATE", http.StatusBadRequest)
		return
	}
	round.RoundID = roundID

	ctx := context.Background()
	err = db.CreateRound(ctx, &round)
	if err != nil {
		if err.Error() == "round already exists" {
			errorResponseWithCode(w, "Conflict", "Round already exists for sport '"+round.Sport+"' on playDate '"+round.PlayDate+"'", "ROUND_ALREADY_EXISTS", http.StatusConflict)
			return
		}
		errorResponseWithCode(w, "Internal Server Error", "Failed to create round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, round, http.StatusCreated)
}

// handleDeleteRound handles DELETE /v1/round
func handleDeleteRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errorResponseWithCode(w, "Method Not Allowed", "Only DELETE method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	if sport == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required parameter: sport", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	playDate := r.URL.Query().Get("playDate")
	if playDate == "" {
		errorResponseWithCode(w, "Bad Request", "Missing required parameter: playDate", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Check if the round exists first
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to check round existence: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}
	if round == nil {
		errorResponseWithCode(w, "Not Found", "Round not found for sport '"+sport+"' on playDate '"+playDate+"'", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	err = db.DeleteRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to delete round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleGetUpcomingRounds handles GET /v1/upcoming-rounds
func handleGetUpcomingRounds(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	if sport == "" {
		errorResponseWithCode(w, "Bad Request", "Sport parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	ctx := context.Background()
	upcomingRounds, err := db.GetRoundsBySport(ctx, sport, startDate, endDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve rounds: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	if len(upcomingRounds) == 0 {
		errorResponseWithCode(w, "Not Found", "No upcoming rounds found for sport '"+sport+"' in the specified date range", "NO_UPCOMING_ROUNDS", http.StatusNotFound)
		return
	}

	// Sort by playDate
	sort.Slice(upcomingRounds, func(i, j int) bool {
		return upcomingRounds[i].PlayDate < upcomingRounds[j].PlayDate
	})

	jsonResponse(w, upcomingRounds, http.StatusOK)
}

// handleSubmitResults handles POST /v1/results
func handleSubmitResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponseWithCode(w, "Method Not Allowed", "Only POST method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	playDate := r.URL.Query().Get("playDate")

	if sport == "" || playDate == "" {
		errorResponseWithCode(w, "Bad Request", "Sport and playDate parameters are required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	var result Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		errorResponseWithCode(w, "Bad Request", "Invalid request body: "+err.Error(), "INVALID_REQUEST_BODY", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}
	if round == nil {
		errorResponseWithCode(w, "Not Found", "Round not found for sport '"+sport+"' on date '"+playDate+"'", "ROUND_NOT_FOUND", http.StatusNotFound)
		return
	}

	// Update round statistics
	round.Stats.TotalPlays++
	if result.IsCorrect {
		correctCount := int(round.Stats.PercentageCorrect * float64(round.Stats.TotalPlays-1) / 100)
		correctCount++
		round.Stats.PercentageCorrect = float64(correctCount) * 100 / float64(round.Stats.TotalPlays)

		// Update average correct score
		totalCorrectScore := round.Stats.AverageCorrectScore * float64(correctCount-1)
		totalCorrectScore += float64(result.Score)
		round.Stats.AverageCorrectScore = totalCorrectScore / float64(correctCount)
	}

	if result.Score > round.Stats.HighestScore {
		round.Stats.HighestScore = result.Score
	}

	// Update average number of tile flips
	totalTileFlips := round.Stats.AverageNumberOfTileFlips * float64(round.Stats.TotalPlays-1)
	totalTileFlips += float64(len(result.TilesFlipped))
	round.Stats.AverageNumberOfTileFlips = totalTileFlips / float64(round.Stats.TotalPlays)

	// Track tile flips
	if len(result.TilesFlipped) > 0 {
		// Track first tile flipped
		incrementTileTracker(&round.Stats.FirstTileFlippedTracker, result.TilesFlipped[0])

		// Track last tile flipped
		incrementTileTracker(&round.Stats.LastTileFlippedTracker, result.TilesFlipped[len(result.TilesFlipped)-1])

		// Track all tiles flipped
		for _, tile := range result.TilesFlipped {
			incrementTileTracker(&round.Stats.MostTileFlippedTracker, tile)
		}

		// Recalculate most/least common tiles
		round.Stats.MostCommonFirstTileFlipped = findMostCommonTile(&round.Stats.FirstTileFlippedTracker)
		round.Stats.MostCommonLastTileFlipped = findMostCommonTile(&round.Stats.LastTileFlippedTracker)
		round.Stats.MostCommonTileFlipped = findMostCommonTile(&round.Stats.MostTileFlippedTracker)
		round.Stats.LeastCommonTileFlipped = findLeastCommonTile(&round.Stats.MostTileFlippedTracker)
	}

	// Save the updated round
	err = db.UpdateRound(ctx, round)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to update round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, result, http.StatusOK)
}

// handleGetRoundStats handles GET /v1/stats/round
func handleGetRoundStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	sport := r.URL.Query().Get("sport")
	playDate := r.URL.Query().Get("playDate")

	if sport == "" || playDate == "" {
		errorResponseWithCode(w, "Bad Request", "Sport and playDate parameters are required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	round, err := db.GetRound(ctx, sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	if round == nil {
		errorResponseWithCode(w, "Not Found", "No statistics found for sport '"+sport+"' on date '"+playDate+"'", "STATS_NOT_FOUND", http.StatusNotFound)
		return
	}

	jsonResponse(w, round.Stats, http.StatusOK)
}

// handleGetUserStats handles GET /v1/stats/user
func handleGetUserStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponseWithCode(w, "Method Not Allowed", "Only GET method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userId")
	if userID == "" {
		errorResponseWithCode(w, "Bad Request", "userId parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	stats, err := db.GetUserStats(ctx, userID)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to retrieve user stats: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	if stats == nil {
		errorResponseWithCode(w, "Not Found", "No statistics found for user '"+userID+"'", "USER_STATS_NOT_FOUND", http.StatusNotFound)
		return
	}

	jsonResponse(w, stats, http.StatusOK)
}

// handleScrapeAndCreateRound handles POST /v1/round - scrapes player data and creates a round
func handleScrapeAndCreateRound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponseWithCode(w, "Method Not Allowed", "Only POST method is allowed", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}

	// Get required parameters
	sport := r.URL.Query().Get("sport")
	playDate := r.URL.Query().Get("playDate")

	// Validate required parameters
	if sport == "" {
		errorResponseWithCode(w, "Bad Request", "Sport parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}
	if playDate == "" {
		errorResponseWithCode(w, "Bad Request", "playDate parameter is required", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	// Validate sport
	if sport != "basketball" && sport != "baseball" && sport != "football" {
		errorResponseWithCode(w, "Bad Request", "Invalid sport parameter. Must be basketball, baseball, or football", "INVALID_PARAMETER", http.StatusBadRequest)
		return
	}

	// Get optional parameters
	name := r.URL.Query().Get("name")
	sportsReferencePath := r.URL.Query().Get("sportsReferencePath")

	// Validate that at least one optional parameter is provided
	if name == "" && sportsReferencePath == "" {
		errorResponseWithCode(w, "Bad Request", "Either 'name' or 'sportsReferencePath' parameter must be provided", "MISSING_REQUIRED_PARAMETER", http.StatusBadRequest)
		return
	}

	// Get the hostname for the sport
	hostname := GetSportsReferenceHostname(sport)
	if hostname == "" {
		errorResponseWithCode(w, "Internal Server Error", "Unable to determine hostname for sport: "+sport, "CONFIGURATION_ERROR", http.StatusInternalServerError)
		return
	}

	// If sportsReferencePath is provided, go directly to the player page
	if sportsReferencePath != "" {
		// Use the direct path to the player page (always use www subdomain)
		playerURL := fmt.Sprintf("https://www.%s%s", hostname, sportsReferencePath)
		fmt.Printf("Player page URL: %s\n", playerURL)

		// Scrape player page data
		player, err := scrapePlayerData(playerURL, hostname, sport)
		if err != nil {
			errorResponseWithCode(w, "Internal Server Error", "Failed to scrape player data: "+err.Error(), "SCRAPING_ERROR", http.StatusInternalServerError)
			return
		}

		// Create round with scraped player data
		now := time.Now()
		roundID, err := GenerateRoundID(sport, playDate)
		if err != nil {
			errorResponseWithCode(w, "Bad Request", "Invalid playDate format: "+err.Error(), "INVALID_PLAY_DATE", http.StatusBadRequest)
			return
		}

		round := Round{
			RoundID:     roundID,
			Sport:       sport,
			PlayDate:    playDate,
			Player:      *player,
			Created:     now,
			LastUpdated: now,
			Stats: RoundStats{
				PlayDate: playDate,
				Name:     player.Name,
				Sport:    sport,
			},
		}

		// Store the round in DynamoDB
		err = db.CreateRound(r.Context(), &round)
		if err != nil {
			errorResponseWithCode(w, "Internal Server Error", "Failed to create round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
			return
		}

		jsonResponse(w, round, http.StatusCreated)
		return
	}

	// Search for the player by name
	encodedName := url.QueryEscape(name)
	searchURL := fmt.Sprintf("https://www.%s/search/search.fcgi?search=%s", hostname, encodedName)

	// Initialize colly collector
	// Allow both www and non-www versions of the domain
	c := colly.NewCollector(
		colly.AllowedDomains(hostname, "www."+hostname),
		colly.MaxDepth(1),
	)

	// Allow redirects between www and non-www versions
	c.AllowURLRevisit = false

	// Variable to capture the final URL after redirects
	var finalURL string
	var scrapeError error
	var playerSearchItems []string

	// Set up error handling
	c.OnError(func(r *colly.Response, err error) {
		scrapeError = err
		fmt.Printf("Scraping error: %v\n", err)
	})

	// Log request
	c.OnRequest(func(r *colly.Request) {
		// fmt.Printf("Visiting: %s\n", r.URL.String())
	})

	// Capture the response and final URL
	c.OnResponse(func(r *colly.Response) {
		finalURL = r.Request.URL.String()
		// fmt.Printf("Final URL after redirects: %s\n", finalURL)
	})

	// Extract player search results from #players div
	c.OnHTML("div#players div.search-item", func(e *colly.HTMLElement) {
		// Get the player URL path from the search-item-url div text
		playerURLPath := strings.TrimSpace(e.ChildText("div.search-item-url"))
		if playerURLPath != "" {
			playerSearchItems = append(playerSearchItems, playerURLPath)
		}
	})

	// Visit the search URL
	err := c.Visit(searchURL)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to initiate scraping: "+err.Error(), "SCRAPING_ERROR", http.StatusInternalServerError)
		return
	}

	// Check if there was a scraping error
	if scrapeError != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to scrape player data: "+scrapeError.Error(), "SCRAPING_ERROR", http.StatusInternalServerError)
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
			errorResponseWithCode(w, "Bad Request", "No players found with the name '"+name+"'. Please check the player name and sport again.", "NO_PLAYERS_FOUND", http.StatusBadRequest)
			return			
		} else {
			errorResponseWithCode(w, "Bad Request", "Multiple players found with the name '"+name+"'. Please provide the sportsReferencePath parameter to specify the exact player.", "MULTIPLE_PLAYERS_FOUND", http.StatusBadRequest)
			return
		}
	}

	// Successfully redirected to a player page
	// fmt.Printf("Player page URL: %s\n", finalURL)

	// Scrape player page data
	player, err := scrapePlayerData(finalURL, hostname, sport)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to scrape player data: "+err.Error(), "SCRAPING_ERROR", http.StatusInternalServerError)
		return
	}

	// Create round with scraped player data
	now := time.Now()
	roundID, err := GenerateRoundID(sport, playDate)
	if err != nil {
		errorResponseWithCode(w, "Bad Request", "Invalid playDate format: "+err.Error(), "INVALID_PLAY_DATE", http.StatusBadRequest)
		return
	}

	round := Round{
		RoundID:     roundID,
		Sport:       sport,
		PlayDate:    playDate,
		Player:      *player,
		Created:     now,
		LastUpdated: now,
		Stats: RoundStats{
			PlayDate: playDate,
			Name:     player.Name,
			Sport:    sport,
		},
	}

	// Store the round in DynamoDB
	err = db.CreateRound(r.Context(), &round)
	if err != nil {
		errorResponseWithCode(w, "Internal Server Error", "Failed to create round: "+err.Error(), "DATABASE_ERROR", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, round, http.StatusCreated)
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

			physicalAttrs = append(physicalAttrs, "▪ "+text)
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
			player.PlayerInformation = strings.Join(physicalAttrs, " ")
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
			ttt := ProcessAchievements(sport, achievements, maxLength)
			player.PersonalAchievements = ttt
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
