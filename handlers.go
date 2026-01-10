package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// Server holds dependencies for HTTP handlers
type Server struct {
	db       *DB
	upscaler *ImageUpscaler
}

// NewServer creates a new Server with the given database and upscaler
func NewServer(db *DB, upscaler *ImageUpscaler) *Server {
	return &Server{
		db:       db,
		upscaler: upscaler,
	}
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

// GetUpcomingRounds handles GET /v1/upcoming-rounds
func (s *Server) GetUpcomingRounds(c *gin.Context) {
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

	startDate := c.Query(QueryParamStartDate)
	endDate := c.Query(QueryParamEndDate)

	upcomingRounds, err := s.db.GetRoundsBySport(c.Request.Context(), sport, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to retrieve rounds: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	if len(upcomingRounds) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "No upcoming rounds found for sport '" + sport + "' in the specified date range",
			JSONFieldCode:      ErrorNoUpcomingRounds,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Sort by playDate
	sort.Slice(upcomingRounds, func(i, j int) bool {
		return upcomingRounds[i].PlayDate < upcomingRounds[j].PlayDate
	})

	c.JSON(http.StatusOK, upcomingRounds)
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

	// potential hack catcher. Score cannot be higher than 100
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

	// Get user_id from bearer token (set by JWT middleware)
	userIdToken, exists := c.Get(ConstantUserId)
	if exists && userIdToken != "" {
		userId, ok := userIdToken.(string)
		if ok && userId != "" {
			// Fetch existing user stats or create new ones
			userStats, err := s.db.GetUserStats(c.Request.Context(), userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					JSONFieldError:     StatusInternalServerError,
					JSONFieldMessage:   "Failed to retrieve user stats: " + err.Error(),
					JSONFieldCode:      ErrorDatabaseError,
					JSONFieldTimestamp: time.Now(),
				})
				return
			}

			// If user stats don't exist, create new user stats
			if userStats == nil {
				userStats = &UserStats{
					UserId:             userId,
					Sports:             []UserSportStats{},
					CurrentDailyStreak: 1,
					LastDayPlayed:      playDate,
					UserName:           "", // TODO: update with user's username as fetched from Auth0
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
				err = s.db.CreateUserStats(c.Request.Context(), userStats)
			} else {
				err = s.db.UpdateUserStats(c.Request.Context(), userStats)
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

	c.JSON(http.StatusOK, result)
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

// GetUserStats handles GET /v1/stats/user
func (s *Server) GetUserStats(c *gin.Context) {
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

	stats, err := s.db.GetUserStats(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to retrieve user stats: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{
			JSONFieldError:     StatusNotFound,
			JSONFieldMessage:   "No statistics found for user '" + userId + "'",
			JSONFieldCode:      ErrorUserStatsNotFound,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
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

	// Parse UserStats from request body
	var userStats UserStats
	if err := c.ShouldBindJSON(&userStats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			JSONFieldError:     StatusBadRequest,
			JSONFieldMessage:   "Invalid request body: " + err.Error(),
			JSONFieldCode:      ErrorInvalidRequestBody,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Check if user stats already exist in the database
	existingStats, err := s.db.GetUserStats(c.Request.Context(), userId)
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
	if existingStats != nil {
		c.JSON(http.StatusConflict, gin.H{
			JSONFieldError:     StatusConflict,
			JSONFieldMessage:   "User stats already exist. Migration not allowed for user '" + userId + "'",
			JSONFieldCode:      ErrorUserAlreadyMigrated,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	// Override userId from payload with userId from JWT token for security
	userStats.UserId = userId

	// Save user stats to DynamoDB
	err = s.db.CreateUserStats(c.Request.Context(), &userStats)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			JSONFieldError:     StatusInternalServerError,
			JSONFieldMessage:   "Failed to migrate user stats: " + err.Error(),
			JSONFieldCode:      ErrorDatabaseError,
			JSONFieldTimestamp: time.Now(),
		})
		return
	}

	c.JSON(http.StatusCreated, userStats)
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

	// 3.5. Upscale player photo
	if player.Photo != "" {
		player.Photo = s.upscaler.UpscaleImage(player.Photo)
	}

	// 4. Build and save round
	round, err := s.createRoundFromPlayer(c.Request.Context(), player, params)
	if err != nil {
		respondWithScrapeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, round)
}
