package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server holds dependencies for HTTP handlers
type Server struct {
	db *DB
}

// NewServer creates a new Server with the given database
func NewServer(db *DB) *Server {
	return &Server{db: db}
}

// GetRound handles GET /v1/round
func (s *Server) GetRound(c *gin.Context) {
	sport := c.Query(QueryParamSport)
	if sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Sport parameter is required",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Validate sport
	if !IsValidSport(sport) {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid sport parameter. Must be basketball, baseball, or football",
			JSONFieldCode:      ErrorInvalidParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	playDate := c.Query(QueryParamPlayDate)
	if playDate == "" {
		playDate = time.Now().Format(DateFormatYYYYMMDD)
	}

	round, err := s.db.GetRound(c.Request.Context(), sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to retrieve round: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "No round found for the specified sport and playDate",
			JSONFieldCode:      ErrorRoundNotFound,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, round)
}

// CreateRound handles PUT /v1/round
func (s *Server) CreateRound(c *gin.Context) {
	var round Round
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid request body: " + err.Error(),
			JSONFieldCode:      ErrorInvalidRequestBody,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Validate required fields
	if round.Sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Missing required field: sport",
			JSONFieldCode:      ErrorMissingRequiredField,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}
	if round.PlayDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Missing required field: playDate",
			JSONFieldCode:      ErrorMissingRequiredField,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}
	if round.Player.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Missing required field: player.name",
			JSONFieldCode:      ErrorMissingRequiredField,
			JSONFieldTimestamp: time.Now(),
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
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid playDate format: " + err.Error(),
			JSONFieldCode:      ErrorInvalidPlayDate,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}
	round.RoundID = roundID

	err = s.db.CreateRound(c.Request.Context(), &round)
	if err != nil {
		if err.Error() == "round already exists" {
			c.JSON(http.StatusConflict, gin.H{
				JSONFieldError:     StatusConflict,
				JSONFieldMessage:   "Round already exists for sport '" + round.Sport + "' on playDate '" + round.PlayDate + "'",
				JSONFieldCode:      ErrorRoundAlreadyExists,
				JSONFieldTimestamp: time.Now(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to create round: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, round)
}

// DeleteRound handles DELETE /v1/round
func (s *Server) DeleteRound(c *gin.Context) {
	sport := c.Query(QueryParamSport)
	if sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Missing required parameter: sport",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	playDate := c.Query(QueryParamPlayDate)
	if playDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Missing required parameter: playDate",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Check if the round exists first
	round, err := s.db.GetRound(c.Request.Context(), sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to check round existence: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}
	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "Round not found for sport '" + sport + "' on playDate '" + playDate + "'",
			JSONFieldCode:      ErrorRoundNotFound,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	err = s.db.DeleteRound(c.Request.Context(), sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to delete round: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// dateRangeProvider is a function that computes the date range based on query parameters
type dateRangeProvider func(startDateQuery, endDateQuery string) (startDate, endDate string)

// getRoundsWithDateProvider handles common logic for retrieving rounds with custom date range logic
func (s *Server) getRoundsWithDateProvider(c *gin.Context, dateProvider dateRangeProvider) {
	sport := c.Query(QueryParamSport)
	if sport == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Sport parameter is required",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	startDate, endDate := dateProvider(c.Query(QueryParamStartDate), c.Query(QueryParamEndDate))
	rounds, err := s.db.GetRoundsBySport(c.Request.Context(), sport, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to retrieve rounds: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	if len(rounds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "No rounds found for sport '" + sport + "' in the specified date range",
			JSONFieldCode:      ErrorNoUpcomingRounds,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, rounds)
}

// GetRounds handles GET /v1/rounds
func (s *Server) GetRounds(c *gin.Context) {
	s.getRoundsWithDateProvider(c, func(startDateQuery, endDateQuery string) (string, string) {
		startDate := startDateQuery
		if startDate == "" {
			startDate = FIRST_ROUND_DATE_STRING
		}

		endDate := endDateQuery
		if endDate == "" {
			endDate = time.Now().Format(DateFormatYYYYMMDD)
		}

		return startDate, endDate
	})
}

// GetUpcomingRounds handles GET /v1/upcoming-rounds
func (s *Server) GetUpcomingRounds(c *gin.Context) {
	s.getRoundsWithDateProvider(c, func(startDateQuery, endDateQuery string) (string, string) {
		startDate := startDateQuery
		if startDate == "" {
			startDate = FIRST_ROUND_DATE_STRING
		}

		endDate := endDateQuery
		if endDate == "" {
			endDateTime := FIRST_ROUND_DATE
			endDateTime2 := time.Now()

			// set endDateTime as the later date for maximum range
			if endDateTime.Before(endDateTime2) {
				endDateTime = endDateTime2
			}
			endDate = endDateTime.AddDate(0, 0, 30).Format(DateFormatYYYYMMDD)
		}

		return startDate, endDate
	})
}

// SubmitResults handles POST /v1/results
func (s *Server) SubmitResults(c *gin.Context) {
	sport := c.Query(QueryParamSport)
	playDate := c.Query(QueryParamPlayDate)

	if sport == "" || playDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Sport and playDate parameters are required",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	var result Result
	if err := c.ShouldBindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid request body: " + err.Error(),
			JSONFieldCode:      ErrorInvalidRequestBody,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// potential hack catcher. Score cannot be higher than 100 or lower than 0
	if result.Score > 100 || result.Score < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid request body: Score cannot be greater than 100 or less than 0",
			JSONFieldCode:      ErrorInvalidRequestBody,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	round, err := s.db.GetRound(c.Request.Context(), sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to retrieve round: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}
	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "Round not found for sport '" + sport + "' on date '" + playDate + "'",
			JSONFieldCode:      ErrorRoundNotFound,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Update round statistics
	updateStatsWithResult(&round.Stats.Stats, &result)

	// Save the updated round
	err = s.db.UpdateRound(c.Request.Context(), round)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to update round: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Build response with optional storyId
	resultResponse := ResultResponse{Result: result}

	// Get user_id from bearer token (set by JWT middleware)
	userIdToken, exists := c.Get(ConstantUserId)
	if exists && userIdToken != "" {
		userId, ok := userIdToken.(string)
		if ok && userId != "" {
			var overwriteUsername string
			usernameToken, exists := c.Get(ConstantUsername)
			if exists && usernameToken != "" {
				overwriteUsername = usernameToken.(string)
			}

			// Fetch existing user or create new ones
			user, err := s.db.GetUser(c.Request.Context(), userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					JSONFieldError:     StatusInternalServerError,
					JSONFieldMessage:   "Failed to retrieve user: " + err.Error(),
					JSONFieldCode:      ErrorDatabaseError,
					JSONFieldTimestamp: time.Now(),
				})
				return
			}

			// initialize empty storyId, populate if necessary
			currentDailyStreakStoryId := ""
			totalWinsStoryId := ""

			// Get user's timezone from header (or UTC fallback)
			userLoc := getUserTimezone(c)
			// Calculate today's date in the user's local timezone
			today := time.Now().In(userLoc).Format(DateFormatYYYYMMDD)

			// If user stats don't exist, create new user stats
			if user == nil {
				totalWins := 0
				if result.Score > 0 {
					totalWins = 1
				}
				user = &User{
					UserId:             userId,
					UserName:           overwriteUsername,
					UserCreated:        time.Now(),
					CurrentDailyStreak: 1,
					TotalPlays:         1,
					TotalWins:          totalWins,
					LastDayPlayed:      today, // Track real-life date in user's timezone, not round playDate
					Sports:             []UserSportStats{},
					StoryMissions:      createEmptyStoryMissions(today),
				}
			} else {
				// Update daily streak based on real-life date in user's timezone (engagement-based tracking)
				currentDailyStreakStoryId = updateDailyStreak(user, today)
			}

			// update username
			if overwriteUsername != "" {
				user.UserName = overwriteUsername
			}

			// Find or create specific sport stats
			var sportStats *UserSportStats
			for i := range user.Sports {
				if user.Sports[i].Sport == sport {
					sportStats = &user.Sports[i]
					break
				}
			}

			// If sport stats don't exist, create new entry
			if sportStats == nil {
				newSportStats := UserSportStats{
					Sport: sport,
				}
				user.Sports = append(user.Sports, newSportStats)
				sportStats = &user.Sports[len(user.Sports)-1]
			}

			// Update root level stats
			user.TotalPlays++
			if result.Score > 0 {
				user.TotalWins++
				totalWinsStoryId = totalWinsStoryMissions(user.TotalWins)
			}

			// Update sport-specific stats
			updateStatsWithResult(&sportStats.Stats, &result)

			// Check if history entry for this playDate already exists
			historyExists := false
			for i := range sportStats.History {
				if sportStats.History[i].PlayDate == playDate {
					// Update existing history entry instead of creating duplicate
					sportStats.History[i].Result = result
					historyExists = true
					break
				}
			}

			// Only append if this playDate doesn't already exist in history
			if !historyExists {
				roundHistory := RoundHistory{
					PlayDate: playDate,
					Result:   result,
				}
				sportStats.History = append(sportStats.History, roundHistory)
			}

			// Update storyMissions with any new completed ones
			updateStoryMissions(&user.StoryMissions, &currentDailyStreakStoryId, &totalWinsStoryId, today, result.PlayerName)

			// Set storyId on response if a total wins story mission was achieved
			resultResponse.StoryId = totalWinsStoryId

			// Save or update user stats in DynamoDB
			if user.UserCreated.IsZero() {
				err = s.db.CreateUser(c.Request.Context(), user)
			} else {
				err = s.db.UpdateUser(c.Request.Context(), user)
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					JSONFieldError:     StatusInternalServerError,
					JSONFieldMessage:   "Failed to update user stats: " + err.Error(),
					JSONFieldCode:      ErrorDatabaseError,
					JSONFieldTimestamp: time.Now(),
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, resultResponse)
}

// GetRoundStats handles GET /v1/stats/round
func (s *Server) GetRoundStats(c *gin.Context) {
	sport := c.Query(QueryParamSport)
	playDate := c.Query(QueryParamPlayDate)

	if sport == "" || playDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Sport and playDate parameters are required",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	round, err := s.db.GetRound(c.Request.Context(), sport, playDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to retrieve round: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	if round == nil {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "No statistics found for sport '" + sport + "' on date '" + playDate + "'",
			JSONFieldCode:      ErrorStatsNotFound,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, round.Stats)
}

// GetUser handles GET /v1/user
func (s *Server) GetUser(c *gin.Context) {
	userId := c.Query(QueryParamUserId)
	if userId == "" {
		// if userId is not part in query param, extract from bearer token instead
		userIdToken, exists := c.Get(ConstantUserId)
		if exists && userIdToken != "" {
			userIdStr, ok := userIdToken.(string)
			if ok && userIdStr != "" {
				userId = userIdStr
			}
		}

		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				JSONFieldError:     StatusBadRequest,
				JSONFieldMessage:   "userId parameter is required",
				JSONFieldCode:      ErrorMissingRequiredParameter,
				JSONFieldTimestamp: time.Now(),
			})
			return
		}
	}

	user, err := s.db.GetUser(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to retrieve user: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "User not found for '" + userId + "'",
			JSONFieldCode:      ErrorUserNotFound,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// TODO: check most advanced currentDailyStreak story and return that

	c.JSON(http.StatusOK, user)
}

// MigrateUserStats handles POST /v1/stats/user/migrate - migrates user stats from local storage to backend
func (s *Server) MigrateUserStats(c *gin.Context) {
	// Get userId from JWT token (set by JWT middleware)
	userIdToken, exists := c.Get(ConstantUserId)
	if !exists || userIdToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			JSONFieldError:     "Unauthorized",
			JSONFieldMessage:   "User ID not found in token",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	userId, ok := userIdToken.(string)
	if !ok || userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			JSONFieldError:     "Unauthorized",
			JSONFieldMessage:   "Invalid user ID in token",
			JSONFieldCode:      ErrorInvalidParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	usernameToken, exists := c.Get(ConstantUsername)
	if !exists {
		usernameToken = ""
	}

	username, ok := usernameToken.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			JSONFieldError:     "Unauthorized",
			JSONFieldMessage:   "Invalid username in token",
			JSONFieldCode:      ErrorInvalidParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Parse User from request body
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid request body: " + err.Error(),
			JSONFieldCode:      ErrorInvalidRequestBody,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Check if user already exist in the database
	existingUser, err := s.db.GetUser(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to check existing user stats: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// If user stats already exist, return 409 Conflict
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			JSONFieldError:     StatusConflict,
			JSONFieldMessage:   "User stats already exist. Migration not allowed for user '" + userId + "'",
			JSONFieldCode:      ErrorUserAlreadyMigrated,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Override userId and username from payload with userId from JWT token for security
	user.UserId = userId
	user.UserName = username

	// Save user stats to DynamoDB
	err = s.db.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to migrate user stats: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ScrapeAndCreateRound handles POST /v1/round - scrapes player data and creates a round
func (s *Server) ScrapeAndCreateRound(c *gin.Context) {
	// 1. Parse and validate input
	params, err := parseAndValidateScrapeParams(c)
	if err != nil {
		respondWithScrapeError(c, err)
		return
	}

	// 2. Resolve player URL (search or direct)
	playerURL, err := resolvePlayerURL(params)
	if err != nil {
		respondWithScrapeError(c, err)
		return
	}

	// 3. Scrape player data
	player, scrapeErr := scrapePlayerData(playerURL, params.Hostname, params.Sport)
	if scrapeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to scrape player data: " + scrapeErr.Error(),
			JSONFieldCode:      ErrorScrapingError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// 4. Build and save round
	round, err := s.createRoundFromPlayer(c.Request.Context(), player, params)
	if err != nil {
		respondWithScrapeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, round)
}

// UpdateUsername handles PUT /v1/user/username - updates user's username in Auth0
func (s *Server) UpdateUsername(c *gin.Context) {
	// Get userId from JWT token (set by JWT middleware)
	userIdToken, exists := c.Get(ConstantUserId)
	if !exists || userIdToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			JSONFieldError:     "Unauthorized",
			JSONFieldMessage:   "User ID not found in token",
			JSONFieldCode:      ErrorMissingRequiredParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	userId, ok := userIdToken.(string)
	if !ok || userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			JSONFieldError:     "Unauthorized",
			JSONFieldMessage:   "Invalid user ID in token",
			JSONFieldCode:      ErrorInvalidParameter,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Parse request body
	var requestBody struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid request body: " + err.Error(),
			JSONFieldCode:      ErrorInvalidRequestBody,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Validate username field
	if requestBody.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Missing required field: username",
			JSONFieldCode:      ErrorMissingRequiredField,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Get Auth0 Management API access token
	managementToken, err := getAuth0ManagementToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to obtain Auth0 Management API token: " + err.Error(),
			JSONFieldCode:      ErrorConfigurationError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Update user metadata in Auth0
	err = updateAuth0UserMetadata(userId, requestBody.Username, managementToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to update username in Auth0: " + err.Error(),
			JSONFieldCode:      ErrorConfigurationError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Username updated successfully",
		"userId":   userId,
		"username": requestBody.Username,
	})
}
