package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server holds dependencies for HTTP handlers
type Server struct {
	service Service
}

// NewServer creates a new Server with the given service
func NewServer(service Service) *Server {
	return &Server{service: service}
}

// writeError writes a standardized error response
func writeError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		JSONFieldError:     http.StatusText(status),
		JSONFieldMessage:   message,
		JSONFieldCode:      code,
		JSONFieldTimestamp: time.Now(),
	})
}

// writeServiceError maps domain errors to HTTP responses
func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrRoundNotFound):
		writeError(c, http.StatusNotFound, ErrorRoundNotFound, err.Error())
	case errors.Is(err, ErrRoundAlreadyExists):
		writeError(c, http.StatusConflict, ErrorRoundAlreadyExists, err.Error())
	case errors.Is(err, ErrUserNotFound):
		writeError(c, http.StatusNotFound, ErrorUserNotFound, err.Error())
	case errors.Is(err, ErrUserAlreadyMigrated):
		writeError(c, http.StatusConflict, ErrorUserAlreadyMigrated, err.Error())
	case errors.Is(err, ErrInvalidPlayDate):
		writeError(c, http.StatusBadRequest, ErrorInvalidPlayDate, err.Error())
	default:
		// Check if it's a scrapeError
		var se *scrapeError
		if errors.As(err, &se) {
			c.JSON(se.StatusCode, gin.H{
				JSONFieldError:     getStatusText(se.StatusCode),
				JSONFieldMessage:   se.Message,
				JSONFieldCode:      se.ErrorCode,
				JSONFieldTimestamp: time.Now(),
			})
			return
		}
		writeError(c, http.StatusInternalServerError, ErrorDatabaseError, err.Error())
	}
}

// GetRound handles GET /v1/round
func (s *Server) GetRound(c *gin.Context) {
	sport := c.Query(QueryParamSport)
	if sport == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredParameter, "Sport parameter is required")
		return
	}

	if !IsValidSport(sport) {
		writeError(c, http.StatusBadRequest, ErrorInvalidParameter, "Invalid sport parameter. Must be basketball, baseball, or football")
		return
	}

	round, err := s.service.GetRound(c.Request.Context(), sport, c.Query(QueryParamPlayDate))
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, round)
}

// CreateRound handles PUT /v1/round
func (s *Server) CreateRound(c *gin.Context) {
	var round Round
	if err := c.ShouldBindJSON(&round); err != nil {
		writeError(c, http.StatusBadRequest, ErrorInvalidRequestBody, "Invalid request body: "+err.Error())
		return
	}

	if round.Sport == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredField, "Missing required field: sport")
		return
	}
	if round.PlayDate == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredField, "Missing required field: playDate")
		return
	}
	if round.Player.Name == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredField, "Missing required field: player.name")
		return
	}

	created, err := s.service.CreateRound(c.Request.Context(), &round)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

// DeleteRound handles DELETE /v1/round
func (s *Server) DeleteRound(c *gin.Context) {
	sport := c.Query(QueryParamSport)
	if sport == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredParameter, "Missing required parameter: sport")
		return
	}

	playDate := c.Query(QueryParamPlayDate)
	if playDate == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredParameter, "Missing required parameter: playDate")
		return
	}

	err := s.service.DeleteRound(c.Request.Context(), sport, playDate)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}

// dateRangeProvider is a function that computes the date range based on query parameters
type dateRangeProvider func(startDateQuery, endDateQuery string) (startDate, endDate string)

// getRoundsWithDateProvider handles common logic for retrieving rounds with custom date range logic
func (s *Server) getRoundsWithDateProvider(c *gin.Context, dateProvider dateRangeProvider) {
	sport := c.Query(QueryParamSport)
	if sport == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredParameter, "Sport parameter is required")
		return
	}

	startDate, endDate := dateProvider(c.Query(QueryParamStartDate), c.Query(QueryParamEndDate))
	rounds, err := s.service.GetRoundsBySport(c.Request.Context(), sport, startDate, endDate)
	if err != nil {
		if errors.Is(err, ErrRoundNotFound) {
			writeError(c, http.StatusNotFound, ErrorNoUpcomingRounds, "No rounds found for sport '"+sport+"' in the specified date range")
			return
		}
		writeServiceError(c, err)
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
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredParameter, "Sport and playDate parameters are required")
		return
	}

	var result Result
	if err := c.ShouldBindJSON(&result); err != nil {
		writeError(c, http.StatusBadRequest, ErrorInvalidRequestBody, "Invalid request body: "+err.Error())
		return
	}

	if result.Score > 100 || result.Score < 0 {
		writeError(c, http.StatusBadRequest, ErrorInvalidRequestBody, "Invalid request body: Score cannot be greater than 100 or less than 0")
		return
	}

	// Extract user info from JWT context
	params := SubmitResultsParams{
		Sport:    sport,
		PlayDate: playDate,
		Result:   result,
		Timezone: getUserTimezone(c),
	}

	userIdToken, exists := c.Get(ConstantUserId)
	if exists && userIdToken != "" {
		userId, ok := userIdToken.(string)
		if ok && userId != "" {
			params.UserID = userId
			usernameToken, exists := c.Get(ConstantUsername)
			if exists && usernameToken != "" {
				params.Username = usernameToken.(string)
			}
		}
	}

	resultResponse, err := s.service.SubmitResults(c.Request.Context(), params)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, resultResponse)
}

// GetRoundStats handles GET /v1/stats/round
func (s *Server) GetRoundStats(c *gin.Context) {
	sport := c.Query(QueryParamSport)
	playDate := c.Query(QueryParamPlayDate)

	if sport == "" || playDate == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredParameter, "Sport and playDate parameters are required")
		return
	}

	stats, err := s.service.GetRoundStats(c.Request.Context(), sport, playDate)
	if err != nil {
		if errors.Is(err, ErrRoundNotFound) {
			writeError(c, http.StatusNotFound, ErrorStatsNotFound, "No statistics found for sport '"+sport+"' on date '"+playDate+"'")
			return
		}
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetUser handles GET /v1/user
func (s *Server) GetUser(c *gin.Context) {
	userId := c.Query(QueryParamUserId)
	if userId == "" {
		userIdToken, exists := c.Get(ConstantUserId)
		if exists && userIdToken != "" {
			userIdStr, ok := userIdToken.(string)
			if ok && userIdStr != "" {
				userId = userIdStr
			}
		}

		if userId == "" {
			writeError(c, http.StatusBadRequest, ErrorMissingRequiredParameter, "userId parameter is required")
			return
		}
	}

	user, err := s.service.GetUser(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			writeError(c, http.StatusNotFound, ErrorUserNotFound, "User not found for '"+userId+"'")
			return
		}
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// MigrateUserStats handles POST /v1/stats/user/migrate
func (s *Server) MigrateUserStats(c *gin.Context) {
	userIdToken, exists := c.Get(ConstantUserId)
	if !exists || userIdToken == "" {
		writeError(c, http.StatusUnauthorized, ErrorMissingRequiredParameter, "User ID not found in token")
		return
	}

	userId, ok := userIdToken.(string)
	if !ok || userId == "" {
		writeError(c, http.StatusUnauthorized, ErrorInvalidParameter, "Invalid user ID in token")
		return
	}

	usernameToken, exists := c.Get(ConstantUsername)
	if !exists {
		usernameToken = ""
	}

	username, ok := usernameToken.(string)
	if !ok {
		writeError(c, http.StatusUnauthorized, ErrorInvalidParameter, "Invalid username in token")
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		writeError(c, http.StatusBadRequest, ErrorInvalidRequestBody, "Invalid request body: "+err.Error())
		return
	}

	migratedUser, err := s.service.MigrateUser(c.Request.Context(), userId, username, &user)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyMigrated) {
			writeError(c, http.StatusConflict, ErrorUserAlreadyMigrated, "User stats already exist. Migration not allowed for user '"+userId+"'")
			return
		}
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, migratedUser)
}

// ScrapeAndCreateRound handles POST /v1/round - scrapes player data and creates a round
func (s *Server) ScrapeAndCreateRound(c *gin.Context) {
	params, err := parseAndValidateScrapeParams(c)
	if err != nil {
		respondWithScrapeError(c, err)
		return
	}

	round, svcErr := s.service.ScrapeAndCreateRound(c.Request.Context(), params)
	if svcErr != nil {
		writeServiceError(c, svcErr)
		return
	}

	c.JSON(http.StatusCreated, round)
}

// UpdateUsername handles PUT /v1/user/username
func (s *Server) UpdateUsername(c *gin.Context) {
	userIdToken, exists := c.Get(ConstantUserId)
	if !exists || userIdToken == "" {
		writeError(c, http.StatusUnauthorized, ErrorMissingRequiredParameter, "User ID not found in token")
		return
	}

	userId, ok := userIdToken.(string)
	if !ok || userId == "" {
		writeError(c, http.StatusUnauthorized, ErrorInvalidParameter, "Invalid user ID in token")
		return
	}

	var requestBody struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		writeError(c, http.StatusBadRequest, ErrorInvalidRequestBody, "Invalid request body: "+err.Error())
		return
	}

	if requestBody.Username == "" {
		writeError(c, http.StatusBadRequest, ErrorMissingRequiredField, "Missing required field: username")
		return
	}

	err := s.service.UpdateUsername(c.Request.Context(), userId, requestBody.Username)
	if err != nil {
		writeError(c, http.StatusInternalServerError, ErrorConfigurationError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Username updated successfully",
		"userId":   userId,
		"username": requestBody.Username,
	})
}
