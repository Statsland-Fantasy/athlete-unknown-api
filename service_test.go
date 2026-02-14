package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

// --- Mock Database ---

type mockDB struct {
	getRoundFn         func(ctx context.Context, sport, playDate string) (*Round, error)
	createRoundFn      func(ctx context.Context, round *Round) error
	updateRoundFn      func(ctx context.Context, round *Round) error
	deleteRoundFn      func(ctx context.Context, sport, playDate string) error
	getRoundsBySportFn func(ctx context.Context, sport, startDate, endDate string) ([]*RoundSummary, error)
	getUserFn          func(ctx context.Context, userId string) (*User, error)
	createUserFn       func(ctx context.Context, user *User) error
	updateUserFn       func(ctx context.Context, user *User) error
}

func (m *mockDB) GetRound(ctx context.Context, sport, playDate string) (*Round, error) {
	if m.getRoundFn != nil {
		return m.getRoundFn(ctx, sport, playDate)
	}
	return nil, nil
}

func (m *mockDB) CreateRound(ctx context.Context, round *Round) error {
	if m.createRoundFn != nil {
		return m.createRoundFn(ctx, round)
	}
	return nil
}

func (m *mockDB) UpdateRound(ctx context.Context, round *Round) error {
	if m.updateRoundFn != nil {
		return m.updateRoundFn(ctx, round)
	}
	return nil
}

func (m *mockDB) DeleteRound(ctx context.Context, sport, playDate string) error {
	if m.deleteRoundFn != nil {
		return m.deleteRoundFn(ctx, sport, playDate)
	}
	return nil
}

func (m *mockDB) GetRoundsBySport(ctx context.Context, sport, startDate, endDate string) ([]*RoundSummary, error) {
	if m.getRoundsBySportFn != nil {
		return m.getRoundsBySportFn(ctx, sport, startDate, endDate)
	}
	return nil, nil
}

func (m *mockDB) GetUser(ctx context.Context, userId string) (*User, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, userId)
	}
	return nil, nil
}

func (m *mockDB) CreateUser(ctx context.Context, user *User) error {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, user)
	}
	return nil
}

func (m *mockDB) UpdateUser(ctx context.Context, user *User) error {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, user)
	}
	return nil
}

// --- Mock Auth0 Client ---

type mockAuth0 struct {
	getManagementTokenFn func() (string, error)
	updateUserMetadataFn func(userId, username, managementToken string) error
}

func (m *mockAuth0) GetManagementToken() (string, error) {
	if m.getManagementTokenFn != nil {
		return m.getManagementTokenFn()
	}
	return "mock-token", nil
}

func (m *mockAuth0) UpdateUserMetadata(userId, username, managementToken string) error {
	if m.updateUserMetadataFn != nil {
		return m.updateUserMetadataFn(userId, username, managementToken)
	}
	return nil
}

// --- Helper ---

func fixedTime() time.Time {
	return time.Date(2026, 2, 13, 12, 0, 0, 0, time.UTC)
}

func newTestService(db *mockDB, auth0 *mockAuth0) *GameService {
	svc := NewGameService(db, auth0)
	svc.now = fixedTime
	return svc
}

func testRound() *Round {
	return &Round{
		RoundID:  "basketball#5",
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Player:   Player{Name: "Test Player"},
		Stats: RoundStats{
			PlayDate: "2026-02-13",
			Name:     "Test Player",
			Sport:    "basketball",
		},
	}
}

func testUser() *User {
	return &User{
		UserId:             "user-123",
		UserName:           "testuser",
		UserCreated:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		CurrentDailyStreak: 3,
		LastDayPlayed:      "2026-02-12",
		TotalPlays:         10,
		TotalWins:          5,
		Sports:             []UserSportStats{},
		StoryMissions:      createEmptyStoryMissions("2026-01-01"),
	}
}

// --- GetRound Tests ---

func TestServiceGetRound_Success(t *testing.T) {
	expected := testRound()
	db := &mockDB{
		getRoundFn: func(_ context.Context, sport, playDate string) (*Round, error) {
			if sport == "basketball" && playDate == "2026-02-13" {
				return expected, nil
			}
			return nil, nil
		},
	}
	svc := newTestService(db, nil)

	round, err := svc.GetRound(context.Background(), "basketball", "2026-02-13")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if round.RoundID != expected.RoundID {
		t.Errorf("expected roundID %s, got %s", expected.RoundID, round.RoundID)
	}
}

func TestServiceGetRound_DefaultsPlayDateToToday(t *testing.T) {
	var capturedPlayDate string
	db := &mockDB{
		getRoundFn: func(_ context.Context, sport, playDate string) (*Round, error) {
			capturedPlayDate = playDate
			return testRound(), nil
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetRound(context.Background(), "basketball", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := fixedTime().Format(DateFormatYYYYMMDD)
	if capturedPlayDate != expected {
		t.Errorf("expected playDate %s, got %s", expected, capturedPlayDate)
	}
}

func TestServiceGetRound_NotFound(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, nil
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetRound(context.Background(), "basketball", "2026-02-13")
	if !errors.Is(err, ErrRoundNotFound) {
		t.Errorf("expected ErrRoundNotFound, got %v", err)
	}
}

func TestServiceGetRound_DBError(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, errors.New("connection refused")
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetRound(context.Background(), "basketball", "2026-02-13")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, ErrRoundNotFound) {
		t.Error("should not be ErrRoundNotFound for a DB error")
	}
}

// --- CreateRound Tests ---

func TestServiceCreateRound_Success(t *testing.T) {
	db := &mockDB{}
	svc := newTestService(db, nil)

	round := &Round{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Player:   Player{Name: "Test Player"},
	}

	created, err := svc.CreateRound(context.Background(), round)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.RoundID == "" {
		t.Error("expected roundID to be generated")
	}
	if created.Created.IsZero() {
		t.Error("expected Created timestamp to be set")
	}
}

func TestServiceCreateRound_AlreadyExists(t *testing.T) {
	db := &mockDB{
		createRoundFn: func(_ context.Context, _ *Round) error {
			return errors.New("round already exists")
		},
	}
	svc := newTestService(db, nil)

	round := &Round{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Player:   Player{Name: "Test Player"},
	}

	_, err := svc.CreateRound(context.Background(), round)
	if !errors.Is(err, ErrRoundAlreadyExists) {
		t.Errorf("expected ErrRoundAlreadyExists, got %v", err)
	}
}

func TestServiceCreateRound_InvalidPlayDate(t *testing.T) {
	db := &mockDB{}
	svc := newTestService(db, nil)

	round := &Round{
		Sport:    "basketball",
		PlayDate: "not-a-date",
		Player:   Player{Name: "Test Player"},
	}

	_, err := svc.CreateRound(context.Background(), round)
	if !errors.Is(err, ErrInvalidPlayDate) {
		t.Errorf("expected ErrInvalidPlayDate, got %v", err)
	}
}

// --- DeleteRound Tests ---

func TestServiceDeleteRound_Success(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return testRound(), nil
		},
	}
	svc := newTestService(db, nil)

	err := svc.DeleteRound(context.Background(), "basketball", "2026-02-13")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServiceDeleteRound_NotFound(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, nil
		},
	}
	svc := newTestService(db, nil)

	err := svc.DeleteRound(context.Background(), "basketball", "2026-02-13")
	if !errors.Is(err, ErrRoundNotFound) {
		t.Errorf("expected ErrRoundNotFound, got %v", err)
	}
}

func TestServiceDeleteRound_DBErrorOnGet(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, errors.New("timeout")
		},
	}
	svc := newTestService(db, nil)

	err := svc.DeleteRound(context.Background(), "basketball", "2026-02-13")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestServiceDeleteRound_DBErrorOnDelete(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return testRound(), nil
		},
		deleteRoundFn: func(_ context.Context, _, _ string) error {
			return errors.New("delete failed")
		},
	}
	svc := newTestService(db, nil)

	err := svc.DeleteRound(context.Background(), "basketball", "2026-02-13")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- GetRoundsBySport Tests ---

func TestServiceGetRoundsBySport_Success(t *testing.T) {
	expected := []*RoundSummary{
		{RoundID: "basketball#5", Sport: "basketball", PlayDate: "2026-02-13"},
	}
	db := &mockDB{
		getRoundsBySportFn: func(_ context.Context, _, _, _ string) ([]*RoundSummary, error) {
			return expected, nil
		},
	}
	svc := newTestService(db, nil)

	rounds, err := svc.GetRoundsBySport(context.Background(), "basketball", "2026-02-01", "2026-02-28")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rounds) != 1 {
		t.Errorf("expected 1 round, got %d", len(rounds))
	}
}

func TestServiceGetRoundsBySport_Empty(t *testing.T) {
	db := &mockDB{
		getRoundsBySportFn: func(_ context.Context, _, _, _ string) ([]*RoundSummary, error) {
			return []*RoundSummary{}, nil
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetRoundsBySport(context.Background(), "basketball", "2026-02-01", "2026-02-28")
	if !errors.Is(err, ErrRoundNotFound) {
		t.Errorf("expected ErrRoundNotFound for empty results, got %v", err)
	}
}

func TestServiceGetRoundsBySport_DBError(t *testing.T) {
	db := &mockDB{
		getRoundsBySportFn: func(_ context.Context, _, _, _ string) ([]*RoundSummary, error) {
			return nil, errors.New("query failed")
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetRoundsBySport(context.Background(), "basketball", "2026-02-01", "2026-02-28")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- GetRoundStats Tests ---

func TestServiceGetRoundStats_Success(t *testing.T) {
	round := testRound()
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return round, nil
		},
	}
	svc := newTestService(db, nil)

	stats, err := svc.GetRoundStats(context.Background(), "basketball", "2026-02-13")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Sport != "basketball" {
		t.Errorf("expected sport basketball, got %s", stats.Sport)
	}
}

func TestServiceGetRoundStats_NotFound(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, nil
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetRoundStats(context.Background(), "basketball", "2026-02-13")
	if !errors.Is(err, ErrRoundNotFound) {
		t.Errorf("expected ErrRoundNotFound, got %v", err)
	}
}

// --- GetUser Tests ---

func TestServiceGetUser_Success(t *testing.T) {
	expected := testUser()
	db := &mockDB{
		getUserFn: func(_ context.Context, userId string) (*User, error) {
			if userId == "user-123" {
				return expected, nil
			}
			return nil, nil
		},
	}
	svc := newTestService(db, nil)

	user, err := svc.GetUser(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.UserId != "user-123" {
		t.Errorf("expected userId user-123, got %s", user.UserId)
	}
}

func TestServiceGetUser_NotFound(t *testing.T) {
	db := &mockDB{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return nil, nil
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetUser(context.Background(), "nonexistent")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestServiceGetUser_DBError(t *testing.T) {
	db := &mockDB{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return nil, errors.New("connection error")
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.GetUser(context.Background(), "user-123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, ErrUserNotFound) {
		t.Error("should not be ErrUserNotFound for a DB error")
	}
}

// --- MigrateUser Tests ---

func TestServiceMigrateUser_Success(t *testing.T) {
	var createdUser *User
	db := &mockDB{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return nil, nil // no existing user
		},
		createUserFn: func(_ context.Context, user *User) error {
			createdUser = user
			return nil
		},
	}
	svc := newTestService(db, nil)

	user := &User{
		UserId:   "old-id",
		UserName: "old-name",
		Sports:   []UserSportStats{{Sport: "basketball"}},
	}

	result, err := svc.MigrateUser(context.Background(), "jwt-user-id", "jwt-username", user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify JWT values override payload values
	if result.UserId != "jwt-user-id" {
		t.Errorf("expected userId jwt-user-id, got %s", result.UserId)
	}
	if result.UserName != "jwt-username" {
		t.Errorf("expected userName jwt-username, got %s", result.UserName)
	}
	if createdUser == nil {
		t.Fatal("expected CreateUser to be called")
	}
}

func TestServiceMigrateUser_AlreadyMigrated(t *testing.T) {
	db := &mockDB{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return testUser(), nil // user already exists
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.MigrateUser(context.Background(), "user-123", "testuser", &User{})
	if !errors.Is(err, ErrUserAlreadyMigrated) {
		t.Errorf("expected ErrUserAlreadyMigrated, got %v", err)
	}
}

func TestServiceMigrateUser_DBErrorOnGet(t *testing.T) {
	db := &mockDB{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.MigrateUser(context.Background(), "user-123", "testuser", &User{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestServiceMigrateUser_DBErrorOnCreate(t *testing.T) {
	db := &mockDB{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return nil, nil
		},
		createUserFn: func(_ context.Context, _ *User) error {
			return errors.New("create failed")
		},
	}
	svc := newTestService(db, nil)

	_, err := svc.MigrateUser(context.Background(), "user-123", "testuser", &User{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- SubmitResults Tests ---

func TestServiceSubmitResults_GuestUser(t *testing.T) {
	round := testRound()
	var updatedRound *Round
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return round, nil
		},
		updateRoundFn: func(_ context.Context, r *Round) error {
			updatedRound = r
			return nil
		},
	}
	svc := newTestService(db, nil)

	params := SubmitResultsParams{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Result:   Result{Score: 80, IsCorrect: true, PlayerName: "Test Player"},
		Timezone: time.UTC,
	}

	resp, err := svc.SubmitResults(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Score != 80 {
		t.Errorf("expected score 80, got %d", resp.Score)
	}
	if updatedRound == nil {
		t.Fatal("expected round to be updated")
	}
}

func TestServiceSubmitResults_AuthenticatedNewUser(t *testing.T) {
	round := testRound()
	var savedUser *User
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return round, nil
		},
		updateRoundFn: func(_ context.Context, _ *Round) error {
			return nil
		},
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return nil, nil // new user
		},
		// New user gets UserCreated set to non-zero, so UpdateUser is called
		updateUserFn: func(_ context.Context, user *User) error {
			savedUser = user
			return nil
		},
	}
	svc := newTestService(db, nil)

	params := SubmitResultsParams{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Result:   Result{Score: 80, IsCorrect: true, PlayerName: "Test Player"},
		UserID:   "user-123",
		Username: "testuser",
		Timezone: time.UTC,
	}

	resp, err := svc.SubmitResults(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedUser == nil {
		t.Fatal("expected user to be saved")
	}
	if savedUser.UserId != "user-123" {
		t.Errorf("expected userId user-123, got %s", savedUser.UserId)
	}
	if savedUser.CurrentDailyStreak != 1 {
		t.Errorf("expected daily streak 1, got %d", savedUser.CurrentDailyStreak)
	}
	if savedUser.TotalPlays != 2 { // 1 from init + 1 from increment
		t.Errorf("expected totalPlays 2, got %d", savedUser.TotalPlays)
	}
	if savedUser.TotalWins != 2 { // 1 from init + 1 from increment (score > 0)
		t.Errorf("expected totalWins 2, got %d", savedUser.TotalWins)
	}
	if len(savedUser.Sports) != 1 || savedUser.Sports[0].Sport != "basketball" {
		t.Error("expected basketball sport stats to be created")
	}
	if resp.Score != 80 {
		t.Errorf("expected score 80, got %d", resp.Score)
	}
}

func TestServiceSubmitResults_AuthenticatedExistingUser(t *testing.T) {
	round := testRound()
	existingUser := testUser()
	var updatedUser *User
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return round, nil
		},
		updateRoundFn: func(_ context.Context, _ *Round) error {
			return nil
		},
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return existingUser, nil
		},
		updateUserFn: func(_ context.Context, user *User) error {
			updatedUser = user
			return nil
		},
	}
	svc := newTestService(db, nil)

	params := SubmitResultsParams{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Result:   Result{Score: 90, IsCorrect: true, PlayerName: "Test Player"},
		UserID:   "user-123",
		Username: "testuser",
		Timezone: time.UTC,
	}

	_, err := svc.SubmitResults(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updatedUser == nil {
		t.Fatal("expected user to be updated")
	}
	// Streak should increment (lastDayPlayed was 2026-02-12, now is 2026-02-13)
	if updatedUser.CurrentDailyStreak != 4 {
		t.Errorf("expected daily streak 4, got %d", updatedUser.CurrentDailyStreak)
	}
	if updatedUser.TotalPlays != 11 {
		t.Errorf("expected totalPlays 11, got %d", updatedUser.TotalPlays)
	}
	if updatedUser.TotalWins != 6 {
		t.Errorf("expected totalWins 6, got %d", updatedUser.TotalWins)
	}
	if updatedUser.LastDayPlayed != "2026-02-13" {
		t.Errorf("expected LastDayPlayed 2026-02-13, got %s", updatedUser.LastDayPlayed)
	}
	if updatedUser.Sports[0].Sport != "basketball" {
		t.Errorf("expected sport stats for sport basketball, got %s", updatedUser.Sports[0].Sport)
	}
	if updatedUser.Sports[0].Stats.TotalPlays != 1 {
		t.Errorf("expected for sport basketball totalPlays 1, got %d", updatedUser.Sports[0].Stats.TotalPlays)

	}
	// TODO add more unit tests here for validating new story notification logic and story missions
}

func TestServiceSubmitResults_RoundNotFound(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, nil
		},
	}
	svc := newTestService(db, nil)

	params := SubmitResultsParams{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Result:   Result{Score: 80},
		Timezone: time.UTC,
	}

	_, err := svc.SubmitResults(context.Background(), params)
	if !errors.Is(err, ErrRoundNotFound) {
		t.Errorf("expected ErrRoundNotFound, got %v", err)
	}
}

func TestServiceSubmitResults_ZeroScoreNoWin(t *testing.T) {
	round := testRound()
	var savedUser *User
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return round, nil
		},
		updateRoundFn: func(_ context.Context, _ *Round) error { return nil },
		getUserFn:     func(_ context.Context, _ string) (*User, error) { return nil, nil },
		updateUserFn: func(_ context.Context, user *User) error {
			savedUser = user
			return nil
		},
	}
	svc := newTestService(db, nil)

	params := SubmitResultsParams{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Result:   Result{Score: 0, IsCorrect: false, PlayerName: "Test Player"},
		UserID:   "user-123",
		Timezone: time.UTC,
	}

	_, err := svc.SubmitResults(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if savedUser == nil {
		t.Fatal("expected user to be saved")
	}
	// TotalWins from init is 0 (score == 0), plus increment only happens if score > 0
	if savedUser.TotalWins != 0 {
		t.Errorf("expected totalWins 0 for zero score, got %d", savedUser.TotalWins)
	}
}

func TestServiceSubmitResults_ExistingHistoryUpdated(t *testing.T) {
	round := testRound()
	existingUser := testUser()
	existingUser.Sports = []UserSportStats{
		{
			Sport: "basketball",
			History: []RoundHistory{
				{PlayDate: "2026-02-13", Result: Result{Score: 50}},
			},
		},
	}
	var updatedUser *User
	db := &mockDB{
		getRoundFn:    func(_ context.Context, _, _ string) (*Round, error) { return round, nil },
		updateRoundFn: func(_ context.Context, _ *Round) error { return nil },
		getUserFn:     func(_ context.Context, _ string) (*User, error) { return existingUser, nil },
		updateUserFn: func(_ context.Context, user *User) error {
			updatedUser = user
			return nil
		},
	}
	svc := newTestService(db, nil)

	params := SubmitResultsParams{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Result:   Result{Score: 90, IsCorrect: true, PlayerName: "Test Player"},
		UserID:   "user-123",
		Timezone: time.UTC,
	}

	_, err := svc.SubmitResults(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should update existing history, not add a duplicate
	if len(updatedUser.Sports[0].History) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(updatedUser.Sports[0].History))
	}
	if updatedUser.Sports[0].History[0].Result.Score != 90 {
		t.Errorf("expected updated score 90, got %d", updatedUser.Sports[0].History[0].Result.Score)
	}
}

func TestServiceSubmitResults_UpdateRoundError(t *testing.T) {
	db := &mockDB{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return testRound(), nil
		},
		updateRoundFn: func(_ context.Context, _ *Round) error {
			return errors.New("update failed")
		},
	}
	svc := newTestService(db, nil)

	params := SubmitResultsParams{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Result:   Result{Score: 80},
		Timezone: time.UTC,
	}

	_, err := svc.SubmitResults(context.Background(), params)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- UpdateUsername Tests ---

func TestServiceUpdateUsername_Success(t *testing.T) {
	var capturedUserId, capturedUsername, capturedToken string
	auth0 := &mockAuth0{
		getManagementTokenFn: func() (string, error) {
			return "mgmt-token-123", nil
		},
		updateUserMetadataFn: func(userId, username, managementToken string) error {
			capturedUserId = userId
			capturedUsername = username
			capturedToken = managementToken
			return nil
		},
	}
	svc := newTestService(&mockDB{}, auth0)

	err := svc.UpdateUsername(context.Background(), "user-123", "newname")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedUserId != "user-123" {
		t.Errorf("expected userId user-123, got %s", capturedUserId)
	}
	if capturedUsername != "newname" {
		t.Errorf("expected username newname, got %s", capturedUsername)
	}
	if capturedToken != "mgmt-token-123" {
		t.Errorf("expected token mgmt-token-123, got %s", capturedToken)
	}
}

func TestServiceUpdateUsername_TokenError(t *testing.T) {
	auth0 := &mockAuth0{
		getManagementTokenFn: func() (string, error) {
			return "", errors.New("auth0 unavailable")
		},
	}
	svc := newTestService(&mockDB{}, auth0)

	err := svc.UpdateUsername(context.Background(), "user-123", "newname")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestServiceUpdateUsername_UpdateMetadataError(t *testing.T) {
	auth0 := &mockAuth0{
		getManagementTokenFn: func() (string, error) {
			return "token", nil
		},
		updateUserMetadataFn: func(_, _, _ string) error {
			return errors.New("update failed")
		},
	}
	svc := newTestService(&mockDB{}, auth0)

	err := svc.UpdateUsername(context.Background(), "user-123", "newname")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
