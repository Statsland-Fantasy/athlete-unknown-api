package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Domain errors that handlers map to HTTP status codes
var (
	ErrRoundNotFound       = errors.New("round not found")
	ErrRoundAlreadyExists  = errors.New("round already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyMigrated = errors.New("user already migrated")
	ErrInvalidPlayDate     = errors.New("invalid play date")
)

// SubmitResultsParams holds the parsed input for SubmitResults
type SubmitResultsParams struct {
	Sport    string
	PlayDate string
	Result   Result
	UserID   string // empty if unauthenticated
	Username string
	Timezone *time.Location
}

// Service defines the business logic interface that handlers depend on
type Service interface {
	GetRound(ctx context.Context, sport, playDate string) (*Round, error)
	CreateRound(ctx context.Context, round *Round) (*Round, error)
	DeleteRound(ctx context.Context, sport, playDate string) error
	GetRoundsBySport(ctx context.Context, sport, startDate, endDate string) ([]*RoundSummary, error)
	SubmitResults(ctx context.Context, params SubmitResultsParams) (*ResultResponse, error)
	GetRoundStats(ctx context.Context, sport, playDate string) (*RoundStats, error)
	GetUser(ctx context.Context, userId string) (*User, error)
	MigrateUser(ctx context.Context, userId, username string, user *User) (*User, error)
	ScrapeAndCreateRound(ctx context.Context, params *scrapeParams) (*Round, error)
	UpdateUsername(ctx context.Context, userId, username string) error
}

// GameService contains the business logic for the application
type GameService struct {
	db    Database
	auth0 Auth0Client
	now   func() time.Time
}

// NewGameService creates a new GameService with the given dependencies
func NewGameService(db Database, auth0 Auth0Client) *GameService {
	return &GameService{
		db:    db,
		auth0: auth0,
		now:   time.Now,
	}
}

// GetRound retrieves a round by sport and playDate, defaulting playDate to today if empty
func (gs *GameService) GetRound(ctx context.Context, sport, playDate string) (*Round, error) {
	if playDate == "" {
		playDate = gs.now().Format(DateFormatYYYYMMDD)
	}

	round, err := gs.db.GetRound(ctx, sport, playDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve round: %w", err)
	}

	if round == nil {
		return nil, ErrRoundNotFound
	}

	return round, nil
}

// CreateRound validates, sets timestamps, generates ID, and persists a new round
func (gs *GameService) CreateRound(ctx context.Context, round *Round) (*Round, error) {
	now := gs.now()
	if round.Created.IsZero() {
		round.Created = now
	}
	if round.LastUpdated.IsZero() {
		round.LastUpdated = now
	}

	roundID, err := GenerateRoundID(round.Sport, round.PlayDate)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPlayDate, err)
	}
	round.RoundID = roundID

	err = gs.db.CreateRound(ctx, round)
	if err != nil {
		if err.Error() == "round already exists" {
			return nil, ErrRoundAlreadyExists
		}
		return nil, fmt.Errorf("failed to create round: %w", err)
	}

	return round, nil
}

// DeleteRound checks existence and deletes a round
func (gs *GameService) DeleteRound(ctx context.Context, sport, playDate string) error {
	round, err := gs.db.GetRound(ctx, sport, playDate)
	if err != nil {
		return fmt.Errorf("failed to check round existence: %w", err)
	}
	if round == nil {
		return ErrRoundNotFound
	}

	err = gs.db.DeleteRound(ctx, sport, playDate)
	if err != nil {
		return fmt.Errorf("failed to delete round: %w", err)
	}

	return nil
}

// GetRoundsBySport retrieves round summaries for a sport within a date range
func (gs *GameService) GetRoundsBySport(ctx context.Context, sport, startDate, endDate string) ([]*RoundSummary, error) {
	rounds, err := gs.db.GetRoundsBySport(ctx, sport, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rounds: %w", err)
	}

	if len(rounds) == 0 {
		return nil, ErrRoundNotFound
	}

	return rounds, nil
}

// SubmitResults processes a game result: updates round stats, user stats, streaks, and story missions
func (gs *GameService) SubmitResults(ctx context.Context, params SubmitResultsParams) (*ResultResponse, error) {
	round, err := gs.db.GetRound(ctx, params.Sport, params.PlayDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve round: %w", err)
	}
	if round == nil {
		return nil, ErrRoundNotFound
	}

	// Update round statistics
	updateStatsWithResult(&round.Stats.Stats, &params.Result)

	err = gs.db.UpdateRound(ctx, round)
	if err != nil {
		return nil, fmt.Errorf("failed to update round: %w", err)
	}

	resultResponse := ResultResponse{Result: params.Result}

	// Process user stats if authenticated
	if params.UserID != "" {
		user, err := gs.db.GetUser(ctx, params.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user: %w", err)
		}

		today := gs.now().In(params.Timezone).Format(DateFormatYYYYMMDD)

		if user == nil {
			user = &User{
				UserId:             params.UserID,
				UserName:           params.Username,
				UserCreated:        gs.now(),
				CurrentDailyStreak: 1,
				TotalPlays:         1,
				TotalWins:          0,
				TotalDaysPlayed:    1,
				LastDayPlayed:      today,
				Sports:             []UserSportStats{},
				StoryMissions:      createEmptyStoryMissions(today),
			}
		} else {
			updateDailyStreak(user, today)
		}

		// Update root level stats
		user.TotalPlays++
		if params.Result.Score > 0 {
			user.TotalWins++
		}

		if params.Username != "" {
			user.UserName = params.Username
		}

		// Find or create specific sport stats
		var sportStats *UserSportStats
		for i := range user.Sports {
			if user.Sports[i].Sport == params.Sport {
				sportStats = &user.Sports[i]
				break
			}
		}

		if sportStats == nil {
			newSportStats := UserSportStats{
				Sport: params.Sport,
			}
			user.Sports = append(user.Sports, newSportStats)
			sportStats = &user.Sports[len(user.Sports)-1]
		}

		// Update sport-specific stats
		updateStatsWithResult(&sportStats.Stats, &params.Result)

		// Check if history entry for this playDate already exists
		historyExists := false
		for i := range sportStats.History {
			if sportStats.History[i].PlayDate == params.PlayDate {
				sportStats.History[i].Result = params.Result
				historyExists = true
				break
			}
		}

		if !historyExists {
			roundHistory := RoundHistory{
				PlayDate: params.PlayDate,
				Result:   params.Result,
			}
			sportStats.History = append(sportStats.History, roundHistory)
		}

		earnedStoryMissionsCriteria := calculateAchievedStoryMissions(user, params.Result.Score)

		// Update storyMissions
		filteredEarnedStoryMissionsCriteria := updateStoryMissions(user, today, params.Result.PlayerName, earnedStoryMissionsCriteria)
		resultResponse.EarnedStoryMissionsCriteria = filteredEarnedStoryMissionsCriteria

		// Save or update user
		err = gs.db.UpdateUser(ctx, user)

		if err != nil {
			return nil, fmt.Errorf("failed to update user stats: %w", err)
		}
	}

	return &resultResponse, nil
}

// GetRoundStats retrieves the stats for a specific round
func (gs *GameService) GetRoundStats(ctx context.Context, sport, playDate string) (*RoundStats, error) {
	round, err := gs.db.GetRound(ctx, sport, playDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve round: %w", err)
	}

	if round == nil {
		return nil, ErrRoundNotFound
	}

	return &round.Stats, nil
}

// GetUser retrieves a user by ID
func (gs *GameService) GetUser(ctx context.Context, userId string) (*User, error) {
	user, err := gs.db.GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// MigrateUser checks for conflicts, overrides IDs from JWT, and creates the user
func (gs *GameService) MigrateUser(ctx context.Context, userId, username string, user *User) (*User, error) {
	existingUser, err := gs.db.GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user stats: %w", err)
	}

	if existingUser != nil {
		return nil, ErrUserAlreadyMigrated
	}

	// Override userId and username from payload with values from JWT for security
	user.UserId = userId
	user.UserName = username

	err = gs.db.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate user stats: %w", err)
	}

	return user, nil
}

// ScrapeAndCreateRound resolves the player URL, scrapes data, and persists the round
func (gs *GameService) ScrapeAndCreateRound(ctx context.Context, params *scrapeParams) (*Round, error) {
	// Resolve player URL (search or direct)
	playerURL, scrapeErr := resolvePlayerURL(params)
	if scrapeErr != nil {
		return nil, scrapeErr
	}

	// Scrape player data
	player, err := scrapePlayerData(playerURL, params.Hostname, params.Sport)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape player data: %w", err)
	}

	// Build and save round
	round, scrapeErr := gs.createRoundFromPlayer(ctx, player, params)
	if scrapeErr != nil {
		return nil, scrapeErr
	}

	return round, nil
}

// createRoundFromPlayer builds a Round struct from Player data and params
func (gs *GameService) createRoundFromPlayer(ctx context.Context, player *Player, params *scrapeParams) (*Round, *scrapeError) {
	roundID, err := GenerateRoundID(params.Sport, params.PlayDate)
	if err != nil {
		return nil, &scrapeError{
			StatusCode: 400,
			Message:    "Invalid playDate format: " + err.Error(),
			ErrorCode:  ErrorInvalidPlayDate,
			Err:        err,
		}
	}

	now := gs.now()
	round := &Round{
		RoundID:     roundID,
		Sport:       params.Sport,
		PlayDate:    params.PlayDate,
		Player:      *player,
		Created:     now,
		LastUpdated: now,
		Title:       params.Title,
		Stats: RoundStats{
			PlayDate: params.PlayDate,
			Name:     player.Name,
			Sport:    params.Sport,
		},
	}

	if err := gs.db.CreateRound(ctx, round); err != nil {
		return nil, &scrapeError{
			StatusCode: 500,
			Message:    "Failed to create round: " + err.Error(),
			ErrorCode:  ErrorDatabaseError,
			Err:        err,
		}
	}

	return round, nil
}

// UpdateUsername updates the user's username in Auth0 & in User DB record
func (gs *GameService) UpdateUsername(ctx context.Context, userId, username string) error {
	managementToken, err := gs.auth0.GetManagementToken()
	if err != nil {
		return fmt.Errorf("failed to obtain Auth0 Management API token: %w", err)
	}

	err = gs.auth0.UpdateUserMetadata(userId, username, managementToken)
	if err != nil {
		return fmt.Errorf("failed to update username in Auth0: %w", err)
	}

	user, err := gs.db.GetUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}

	user.UserName = username
	if err := gs.db.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update username in DB: %w", err)
	}

	return nil
}
