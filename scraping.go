package main

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

// scrapeParams holds validated parameters for scraping operations
type scrapeParams struct {
	Sport              string
	PlayDate           string
	Name               string
	SportsReferenceURL string
	Theme              string
	Hostname           string
}

// scrapeError represents a scraping error with HTTP status and error details
type scrapeError struct {
	StatusCode int
	Message    string
	ErrorCode  string
	Err        error
}

const (
	ClueMaxLength = 55
)

func (e *scrapeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// parseAndValidateScrapeParams extracts and validates scraping parameters from the request
func parseAndValidateScrapeParams(c *gin.Context) (*scrapeParams, *scrapeError) {
	sport := c.Query(QueryParamSport)
	playDate := c.Query(QueryParamPlayDate)
	name := c.Query(QueryParamName)
	sportsReferenceURL := c.Query(QueryParamSportsReferenceURL)
	theme := c.Query(QueryParamTheme)

	// Validate required parameters
	if sport == "" {
		return nil, &scrapeError{
			StatusCode: 400,
			Message:    "Sport parameter is required",
			ErrorCode:  ErrorMissingRequiredParameter,
		}
	}

	if playDate == "" {
		return nil, &scrapeError{
			StatusCode: 400,
			Message:    "playDate parameter is required",
			ErrorCode:  ErrorMissingRequiredParameter,
		}
	}

	// Validate sport
	if !IsValidSport(sport) {
		return nil, &scrapeError{
			StatusCode: 400,
			Message:    "Invalid sport parameter. Must be basketball, baseball, or football",
			ErrorCode:  ErrorInvalidParameter,
		}
	}

	// Validate that at least one optional parameter is provided
	if name == "" && sportsReferenceURL == "" {
		return nil, &scrapeError{
			StatusCode: 400,
			Message:    "Either 'name' or 'sportsReferenceURL' parameter must be provided",
			ErrorCode:  ErrorMissingRequiredParameter,
		}
	}

	// Get the hostname for the sport
	hostname := GetSportsReferenceHostname(sport)
	if hostname == "" {
		return nil, &scrapeError{
			StatusCode: 500,
			Message:    "Unable to determine hostname for sport: " + sport,
			ErrorCode:  ErrorConfigurationError,
		}
	}

	return &scrapeParams{
		Sport:              sport,
		PlayDate:           playDate,
		Name:               name,
		SportsReferenceURL: sportsReferenceURL,
		Theme:              theme,
		Hostname:           hostname,
	}, nil
}

// resolvePlayerURL determines the final player URL either via direct URL or search
func resolvePlayerURL(params *scrapeParams) (string, *scrapeError) {
	// If direct URL provided, validate and return
	if params.SportsReferenceURL != "" {
		if err := ValidateSportsReferenceURL(params.SportsReferenceURL); err != nil {
			return "", &scrapeError{
				StatusCode: 400,
				Message:    "Invalid sportsReferenceURL: " + err.Error(),
				ErrorCode:  ErrorInvalidURL,
				Err:        err,
			}
		}
		fmt.Printf("Player page URL: %s\n", params.SportsReferenceURL)
		return params.SportsReferenceURL, nil
	}

	// Otherwise search by name
	return searchPlayerByName(params.Name, params.Hostname)
}

// searchPlayerByName performs player search and returns the player's URL
func searchPlayerByName(name, hostname string) (string, *scrapeError) {
	encodedName := url.QueryEscape(name)
	searchURL := fmt.Sprintf("https://www.%s/search/search.fcgi?search=%s", hostname, encodedName)

	// Initialize colly collector
	collector := colly.NewCollector(
		colly.AllowedDomains(hostname, "www."+hostname),
		colly.MaxDepth(1),
	)

	// Allow redirects between www and non-www versions
	collector.AllowURLRevisit = false

	// Variable to capture the final URL after redirects
	var finalURL string
	var collectorError error
	var playerSearchItems []string

	// Set up error handling
	collector.OnError(func(r *colly.Response, err error) {
		collectorError = err
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
		return "", &scrapeError{
			StatusCode: 500,
			Message:    "Failed to initiate scraping: " + err.Error(),
			ErrorCode:  ErrorScrapingError,
			Err:        err,
		}
	}

	// Check if there was a scraping error
	if collectorError != nil {
		return "", &scrapeError{
			StatusCode: 500,
			Message:    "Failed to scrape player data: " + collectorError.Error(),
			ErrorCode:  ErrorScrapingError,
			Err:        collectorError,
		}
	}

	// Check if the final URL contains "/players" (meaning it redirected to a specific player page)
	if !strings.Contains(finalURL, "/players") {
		// Check if there's exactly one player result in the search results
		if len(playerSearchItems) == 1 {
			// Use the single player result
			finalURL = fmt.Sprintf("https://www.%s%s", hostname, playerSearchItems[0])
			// fmt.Printf("Found single player result, using URL: %s\n", finalURL)
		} else if len(playerSearchItems) == 0 {
			return "", &scrapeError{
				StatusCode: 400,
				Message:    "No players found with the name '" + name + "'. Please check the player name and sport again.",
				ErrorCode:  ErrorNoPlayersFound,
			}
		} else {
			return "", &scrapeError{
				StatusCode: 400,
				Message:    "Multiple players found with the name '" + name + "'. Please provide the sportsReferenceURL parameter to specify the exact player.",
				ErrorCode:  ErrorMultiplePlayersFound,
			}
		}
	}

	// Successfully redirected to a player page
	// fmt.Printf("Player page URL: %s\n", finalURL)

	// Validate the URL before returning (finalURL comes from redirect/search results)
	if err := ValidateSportsReferenceURL(finalURL); err != nil {
		return "", &scrapeError{
			StatusCode: 500,
			Message:    "Invalid player URL from search results: " + err.Error(),
			ErrorCode:  ErrorInvalidSearchResultURL,
			Err:        err,
		}
	}

	return finalURL, nil
}

// createRoundFromPlayer builds a Round struct from Player data and params
func (s *Server) createRoundFromPlayer(ctx context.Context, player *Player, params *scrapeParams) (*Round, *scrapeError) {
	roundID, err := GenerateRoundID(params.Sport, params.PlayDate)
	if err != nil {
		return nil, &scrapeError{
			StatusCode: 400,
			Message:    "Invalid playDate format: " + err.Error(),
			ErrorCode:  ErrorInvalidPlayDate,
			Err:        err,
		}
	}

	now := time.Now()
	round := &Round{
		RoundID:     roundID,
		Sport:       params.Sport,
		PlayDate:    params.PlayDate,
		Player:      *player,
		Created:     now,
		LastUpdated: now,
		Theme:       params.Theme,
		Stats: RoundStats{
			PlayDate: params.PlayDate,
			Name:     player.Name,
			Sport:    params.Sport,
		},
	}

	// Store the round in DynamoDB
	if err := s.db.CreateRound(ctx, round); err != nil {
		return nil, &scrapeError{
			StatusCode: 500,
			Message:    "Failed to create round: " + err.Error(),
			ErrorCode:  ErrorDatabaseError,
			Err:        err,
		}
	}

	return round, nil
}

// respondWithScrapeError sends an error response based on scrapeError
func respondWithScrapeError(c *gin.Context, err *scrapeError) {
	c.JSON(err.StatusCode, gin.H{
		JSONFieldError:     getStatusText(err.StatusCode),
		JSONFieldMessage:   err.Message,
		JSONFieldCode:      err.ErrorCode,
		JSONFieldTimestamp: time.Now(),
	})
}

// getStatusText returns the HTTP status text for a status code
func getStatusText(statusCode int) string {
	switch statusCode {
	case 400:
		return StatusBadRequest
	case 404:
		return StatusNotFound
	case 409:
		return StatusConflict
	case 500:
		return StatusInternalServerError
	default:
		return StatusInternalServerError
	}
}

// scrapePlayerData orchestrates the scraping of all player information
func scrapePlayerData(playerURL, hostname, sport string) (*Player, error) {
	// Validate URL as an additional safety layer
	if err := ValidateSportsReferenceURL(playerURL); err != nil {
		return nil, fmt.Errorf("invalid player URL: %w", err)
	}

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

	var statsPulloutElement *colly.HTMLElement
	var rawAchievements []string

	// Register all scrapers
	scrapePlayerName(c, player) // includes initials
	scrapeBio(c, player, sport)
	scrapePlayerInformation(c, player, sport)
	scrapeDraftInformation(c, player, sport)
	scrapeYearsActiveAndTeamsPlayedOn(c, player, sport)
	scrapeJerseyNumbers(c, player)
	scrapeCareerStats(c, &statsPulloutElement)
	scrapePersonalAchievements(c, &rawAchievements)
	scrapePhoto(c, player)
	scrapeNicknames(c, player, sport)

	// Post-processing after scraping completes
	c.OnScraped(func(r *colly.Response) {
		// Set draft information default if not found
		if player.DraftInformation == "" {
			player.DraftInformation = "Undrafted"
		}

		// Set career stats. Need playerInformation to determine position for stats to get
		if statsPulloutElement != nil {
			careerStatsConfig := GetCareerStatsConfig(sport, player.PlayerInformation)

			winsOrSavesValue := 0
			var winsOrSavesLabel string

			var careerStats []string
			for _, statConfig := range careerStatsConfig.Stats {
				statValue := strings.TrimSpace(statsPulloutElement.DOM.Find(statConfig.HTMLPath).Text())
				if statValue != "" {
					// Logic to correctly display the higher of Wins or Saves (baseball pitcher-only)
					if statConfig.StatLabel == "W" || statConfig.StatLabel == "SV" {
						intStatValue, err := strconv.Atoi(statValue)
						if err != nil {
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

			if winsOrSavesLabel != "" {
				careerStats = append([]string{fmt.Sprintf("%d %s", winsOrSavesValue, winsOrSavesLabel)}, careerStats...)
			}

			if len(careerStats) > 0 {
				player.CareerStats = strings.Join(careerStats, ", ")
			}
		}

		// Set personal achievements
		if len(rawAchievements) > 0 {
			player.PersonalAchievements = ProcessAchievements(sport, rawAchievements, ClueMaxLength)
		} else {
			player.PersonalAchievements = "N/A"
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

// scrapePlayerName extracts the player's name
func scrapePlayerName(c *colly.Collector, player *Player) {
	c.OnHTML("h1[itemprop='name'], h1 span", func(e *colly.HTMLElement) {
		if player.Name == "" {
			playerName := strings.TrimSpace(e.Text)
			player.Name = playerName
			player.Initials = getPlayerInitials(playerName)
		}
	})
}

// scrapeBio extracts the player's biographical information (birth date and location)
func scrapeBio(c *colly.Collector, player *Player, sport string) {
	var dobText string
	var physicalAttributesText string
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		// Look for birth information
		if strings.Contains(text, "Born:") || strings.Contains(text, "born") {
			// Remove newlines and extra spaces
			dobText = text
			dobText = strings.ReplaceAll(dobText, "\n", " ")
			dobText = strings.Join(strings.Fields(dobText), " ")
			// Abbreviate month names to first 3 characters
			dobText = regexp.MustCompile(`\b(January|February|March|April|May|June|July|August|September|October|November|December)\b`).ReplaceAllStringFunc(dobText, func(month string) string {
				return month[:3]
			})
			// Remove country code (last 3 characters: space + 2-char code)
			if sport != SportFootball && len(dobText) > 3 {
				dobText = strings.TrimSpace(dobText[:len(dobText)-3])
			}
			// Abbreviate US state names to 2-letter codes
			dobText = abbreviateUSState(dobText)
		}

		if (strings.Contains(text, "cm") && strings.Contains(text, "kg")) ||
			(strings.Contains(text, "lb") && (strings.Contains(text, "-") || strings.Contains(text, "'"))) {
			// This likely contains height and weight
			physicalAttributesText = text
			physicalAttributesText = strings.ReplaceAll(physicalAttributesText, "\n", " ")
			physicalAttributesText = strings.Join(strings.Fields(physicalAttributesText), " ")

			// Remove metric measurements in parentheses
			re := regexp.MustCompile(`\s*\([^)]*\)`)
			physicalAttributesText = re.ReplaceAllString(physicalAttributesText, "")
			physicalAttributesText = strings.TrimSpace(physicalAttributesText)
		}

		player.Bio = dobText + " ▪ " + physicalAttributesText
	})
}

// scrapePlayerInformation extracts physical attributes (height, weight, position, handedness)
func scrapePlayerInformation(c *colly.Collector, player *Player, sport string) {
	playerInformation := &[]string{}

	// Extract from paragraph elements
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		// Look for physical attributes
		if strings.Contains(text, "Position:") || strings.Contains(text, "Positions:") || strings.Contains(text, "Bats:") ||
			strings.Contains(text, "Throws:") || strings.Contains(text, "Shoots:") {
			// Remove newlines and extra spaces
			text = strings.ReplaceAll(text, "\n", " ")
			text = strings.ReplaceAll(text, "-", ", ") // football uses - instead of ,
			if sport == SportFootball {
				text = strings.ReplaceAll(text, "Throws:", " ▪ Throws:") // football has Throws in same line as position
			}
			text = strings.Join(strings.Fields(text), " ")
			*playerInformation = append(*playerInformation, text)
		}

		playerInformationString := strings.Join(*playerInformation, " ▪ ")
		player.PlayerInformation = abbreviatePositions(playerInformationString)
	})
}

// scrapeDraftInformation extracts draft information and returns with college/high school names
func scrapeDraftInformation(c *colly.Collector, player *Player, sport string) {
	var college string
	var highSchool string

	// Extract college information for football/basketball
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		schoolText := strings.TrimSpace(e.Text)
		if strings.Contains(schoolText, "College:") || strings.Contains(schoolText, "Colleges:") {
			// Extract college name after the colon
			parts := strings.SplitN(schoolText, ":", 2)
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
		} else if strings.Contains(schoolText, "High School:") {
			// Extract high school name after the colon (as fallback)
			parts := strings.SplitN(schoolText, ":", 2)
			if len(parts) == 2 {
				highSchool = strings.TrimSpace(parts[1])
			}
		}

		// extract draft information
		draftText := strings.TrimSpace(e.Text)
		if strings.Contains(draftText, "Draft:") || strings.Contains(draftText, "Drafted") {
			// Remove extra spaces - normalize to single space between words
			draftText = strings.ReplaceAll(draftText, "\n", " ")
			draftText = strings.ReplaceAll(draftText, "\t", " ")
			draftText = strings.Join(strings.Fields(draftText), " ")

			// Use college if available, otherwise fall back to high school
			school := college
			if school == "" {
				school = highSchool
			}
			school = strings.ReplaceAll(school, "\n", "")
			school = strings.ReplaceAll(school, "\t", "")
			player.DraftInformation = formatDraftInformation(draftText, sport, school)
		}

		if len(player.DraftInformation) > ClueMaxLength {
			// should only occur for baseball
			// Remove parentheses that don't contain "OVR" (e.g., city/state locations)
			re := regexp.MustCompile(`\s*\([^)]*\)`)
			player.DraftInformation = re.ReplaceAllStringFunc(player.DraftInformation, func(match string) string {
				if strings.Contains(strings.ToUpper(match), "OVR") {
					return match // Keep parentheses containing "OVR"
				}
				return "" // Remove other parentheses (locations, etc.)
			})
			// Clean up any extra spaces left behind
			player.DraftInformation = strings.Join(strings.Fields(player.DraftInformation), " ")
		}
	})
}

// scrapeYearsActive extracts years active and teams played on accounting for injury/unplayed years
func scrapeYearsActiveAndTeamsPlayedOn(c *colly.Collector, player *Player, sport string) {
	var firstTableProcessed bool

	c.OnHTML("table", func(e *colly.HTMLElement) {
		// Skip tables with id="last5" (last 5 games tables at the top of the page)
		if e.Attr("id") == "last5" {
			return
		}

		// Only process the very first table on the page
		if firstTableProcessed {
			return
		}
		firstTableProcessed = true

		var years []string
		var teams []string
		// Extract all tr elements from tbody
		e.ForEach("tbody tr", func(_ int, row *colly.HTMLElement) {
			year := strings.TrimSpace(row.ChildText("th[data-stat='year_id']"))
			teamNameAbbr := strings.TrimSpace(row.ChildText("td[data-stat='team_name_abbr']"))
			teamNameAbbrLower := strings.ToLower(teamNameAbbr)

			isActiveYear := !(strings.Contains(teamNameAbbrLower, "did not play"))
			if isValidYear(year) && !contains(years, year) && isActiveYear {
				years = append(years, year)
				// only process non-duplicate teams and rows without "TM" (Total/multi-team)
				if !contains(teams, teamNameAbbr) && !strings.Contains(teamNameAbbrLower, "tm") {
					teams = append(teams, teamNameAbbr)
				}
			}
		})

		player.YearsActive = formatYearsAsRanges(years, sport)
		player.TeamsPlayedOn = strings.Join(teams, ", ")
	})
}

// scrapeJerseyNumbers extracts jersey numbers using uni_holder class
func scrapeJerseyNumbers(c *colly.Collector, player *Player) {

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
		for i := 20; i > 0; i-- {
			text = strings.ReplaceAll(text, fmt.Sprintf(", +%d", i), "")
		}
		player.JerseyNumbers = strings.TrimSpace(text)
	})
}

// scrapeCareerStats captures the stats pullout element for processing
func scrapeCareerStats(c *colly.Collector, statsPulloutElement **colly.HTMLElement) {
	c.OnHTML(".stats_pullout", func(e *colly.HTMLElement) {
		*statsPulloutElement = e
	})
}

// scrapePersonalAchievements extracts awards, honors, and championships
func scrapePersonalAchievements(c *colly.Collector, rawAchievements *[]string) {
	c.OnHTML("ul#bling li", func(e *colly.HTMLElement) {
		achievement := strings.TrimSpace(e.Text)
		if achievement != "" {
			*rawAchievements = append(*rawAchievements, achievement)
		}
	})
}

// scrapePhoto extracts the player photo URL from multiple possible selectors
func scrapePhoto(c *colly.Collector, player *Player) {
	c.OnHTML("div#meta", func(e *colly.HTMLElement) {
		mediaItemFind := e.DOM.Find("div.media-item")
		if mediaItemFind.Length() > 0 {
			mediaItemImages := mediaItemFind.Children()
			if mediaItemImages.Length() > 0 {
				mediaItemImage := mediaItemImages.First()
				src, exists := mediaItemImage.Attr("src")
				if exists && player.Photo == "" {
					player.Photo = src
				}
			}
		} else {
			// Fallback to any img in div#meta
			imgFind := e.DOM.Find("img")
			if imgFind.Length() > 0 {
				src, exists := imgFind.Attr("src")
				if exists && player.Photo == "" {
					player.Photo = src
				}
			}
		}
	})
}

// scrapeNicknames extracts the nicknames from various places given the sport page
func scrapeNicknames(c *colly.Collector, player *Player, sport string) {
	isFirstP := true // for football
	var nicknamesText string
	c.OnHTML("div#meta p", func(e *colly.HTMLElement) {
		divMetaText := strings.TrimSpace(e.Text)
		switch sport {
		case "baseball":
			if strings.Contains(divMetaText, "Nicknames:") {
				divMetaText = strings.ReplaceAll(divMetaText, "\n", " ")
				divMetaText = strings.ReplaceAll(divMetaText, "\t", " ")
				parts := strings.SplitN(divMetaText, ":", 2)
				if len(parts) == 2 {
					nicknamesText = strings.TrimSpace(parts[1])
				} else {
					nicknamesText = strings.TrimSpace(parts[0])
				}
			}
		case "basketball":
			// Extract text between parentheses if string starts with one
			parenRegex := regexp.MustCompile(`^\(([^)]+)\)`)
			if matches := parenRegex.FindStringSubmatch(divMetaText); len(matches) > 1 {
				nicknamesText = matches[1]
			}
		case "football":
			if isFirstP {
				// Extract text between parentheses from first p element
				parenRegex := regexp.MustCompile(`\(([^)]+)\)`)
				if matches := parenRegex.FindStringSubmatch(divMetaText); len(matches) > 1 {
					nicknamesText = matches[1]
					nicknamesText = strings.ReplaceAll(nicknamesText, " or", ",")
				}
				isFirstP = false
			}
		}

		player.Nicknames = nicknamesText
	})
}
