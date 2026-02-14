package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// --- Mock Service ---

type mockService struct {
	getRoundFn             func(ctx context.Context, sport, playDate string) (*Round, error)
	createRoundFn          func(ctx context.Context, round *Round) (*Round, error)
	deleteRoundFn          func(ctx context.Context, sport, playDate string) error
	getRoundsBySportFn     func(ctx context.Context, sport, startDate, endDate string) ([]*RoundSummary, error)
	submitResultsFn        func(ctx context.Context, params SubmitResultsParams) (*ResultResponse, error)
	getRoundStatsFn        func(ctx context.Context, sport, playDate string) (*RoundStats, error)
	getUserFn              func(ctx context.Context, userId string) (*User, error)
	migrateUserFn          func(ctx context.Context, userId, username string, user *User) (*User, error)
	scrapeAndCreateRoundFn func(ctx context.Context, params *scrapeParams) (*Round, error)
	updateUsernameFn       func(ctx context.Context, userId, username string) error
}

func (m *mockService) GetRound(ctx context.Context, sport, playDate string) (*Round, error) {
	if m.getRoundFn != nil {
		return m.getRoundFn(ctx, sport, playDate)
	}
	return nil, nil
}

func (m *mockService) CreateRound(ctx context.Context, round *Round) (*Round, error) {
	if m.createRoundFn != nil {
		return m.createRoundFn(ctx, round)
	}
	return round, nil
}

func (m *mockService) DeleteRound(ctx context.Context, sport, playDate string) error {
	if m.deleteRoundFn != nil {
		return m.deleteRoundFn(ctx, sport, playDate)
	}
	return nil
}

func (m *mockService) GetRoundsBySport(ctx context.Context, sport, startDate, endDate string) ([]*RoundSummary, error) {
	if m.getRoundsBySportFn != nil {
		return m.getRoundsBySportFn(ctx, sport, startDate, endDate)
	}
	return nil, nil
}

func (m *mockService) SubmitResults(ctx context.Context, params SubmitResultsParams) (*ResultResponse, error) {
	if m.submitResultsFn != nil {
		return m.submitResultsFn(ctx, params)
	}
	return &ResultResponse{Result: params.Result}, nil
}

func (m *mockService) GetRoundStats(ctx context.Context, sport, playDate string) (*RoundStats, error) {
	if m.getRoundStatsFn != nil {
		return m.getRoundStatsFn(ctx, sport, playDate)
	}
	return nil, nil
}

func (m *mockService) GetUser(ctx context.Context, userId string) (*User, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, userId)
	}
	return nil, nil
}

func (m *mockService) MigrateUser(ctx context.Context, userId, username string, user *User) (*User, error) {
	if m.migrateUserFn != nil {
		return m.migrateUserFn(ctx, userId, username, user)
	}
	return user, nil
}

func (m *mockService) ScrapeAndCreateRound(ctx context.Context, params *scrapeParams) (*Round, error) {
	if m.scrapeAndCreateRoundFn != nil {
		return m.scrapeAndCreateRoundFn(ctx, params)
	}
	return nil, nil
}

func (m *mockService) UpdateUsername(ctx context.Context, userId, username string) error {
	if m.updateUsernameFn != nil {
		return m.updateUsernameFn(ctx, userId, username)
	}
	return nil
}

// --- Test Helpers ---

// getTestServer creates a Server with a nil-function mockService for input validation tests
func getTestServer() *Server {
	return NewServer(&mockService{})
}

// assertStatus checks the response status code
func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Errorf("expected status %d, got %d. Body: %s", expected, w.Code, w.Body.String())
	}
}

// assertErrorCode checks the error code in the response body
func assertErrorCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode string) {
	t.Helper()
	var errResp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&errResp)
	if code, ok := errResp["code"].(string); !ok || code != expectedCode {
		t.Errorf("expected error code %s, got %v", expectedCode, errResp["code"])
	}
}

// --- GetRound Handler Tests ---

func TestHandleGetRound(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing sport parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
		{
			name:           "invalid sport parameter",
			queryParams:    "sport=soccer",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorInvalidParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/round?"+tt.queryParams, nil)

			getTestServer().GetRound(c)

			assertStatus(t, w, tt.expectedStatus)
			assertErrorCode(t, w, tt.expectedCode)
		})
	}
}

func TestHandleGetRound_Success(t *testing.T) {
	svc := &mockService{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return &Round{RoundID: "basketball#5", Sport: "basketball", PlayDate: "2026-02-13"}, nil
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/round?sport=basketball&playDate=2026-02-13", nil)

	server.GetRound(c)

	assertStatus(t, w, http.StatusOK)

	var resp Round
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.RoundID != "basketball#5" {
		t.Errorf("expected roundID basketball#5, got %s", resp.RoundID)
	}
}

func TestHandleGetRound_NotFound(t *testing.T) {
	svc := &mockService{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, ErrRoundNotFound
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/round?sport=basketball&playDate=2026-02-13", nil)

	server.GetRound(c)

	assertStatus(t, w, http.StatusNotFound)
	assertErrorCode(t, w, ErrorRoundNotFound)
}

func TestHandleGetRound_DBError(t *testing.T) {
	svc := &mockService{
		getRoundFn: func(_ context.Context, _, _ string) (*Round, error) {
			return nil, errors.New("connection refused")
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/round?sport=basketball&playDate=2026-02-13", nil)

	server.GetRound(c)

	assertStatus(t, w, http.StatusInternalServerError)
	assertErrorCode(t, w, ErrorDatabaseError)
}

// --- CreateRound Handler Tests ---

func TestHandleCreateRound(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "invalid request body",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorInvalidRequestBody,
		},
		{
			name: "missing sport field",
			body: Round{
				PlayDate: "2024-01-01",
				Player:   Player{Name: "Test Player"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredField,
		},
		{
			name: "missing playDate field",
			body: Round{
				Sport:  "basketball",
				Player: Player{Name: "Test Player"},
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredField,
		},
		{
			name: "missing player name field",
			body: Round{
				Sport:    "basketball",
				PlayDate: "2024-01-01",
				Player:   Player{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if str, ok := tt.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPut, "/v1/round", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			getTestServer().CreateRound(c)

			assertStatus(t, w, tt.expectedStatus)
			assertErrorCode(t, w, tt.expectedCode)
		})
	}
}

func TestHandleCreateRound_Success(t *testing.T) {
	svc := &mockService{
		createRoundFn: func(_ context.Context, round *Round) (*Round, error) {
			round.RoundID = "basketball#5"
			return round, nil
		},
	}
	server := NewServer(svc)

	round := Round{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Player:   Player{Name: "Test Player"},
	}
	bodyBytes, _ := json.Marshal(round)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/round", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	server.CreateRound(c)

	assertStatus(t, w, http.StatusCreated)

	var resp Round
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.RoundID == "" {
		t.Error("expected roundID to be generated")
	}
}

func TestHandleCreateRound_AlreadyExists(t *testing.T) {
	svc := &mockService{
		createRoundFn: func(_ context.Context, _ *Round) (*Round, error) {
			return nil, ErrRoundAlreadyExists
		},
	}
	server := NewServer(svc)

	round := Round{
		Sport:    "basketball",
		PlayDate: "2026-02-13",
		Player:   Player{Name: "Test Player"},
	}
	bodyBytes, _ := json.Marshal(round)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/round", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	server.CreateRound(c)

	assertStatus(t, w, http.StatusConflict)
	assertErrorCode(t, w, ErrorRoundAlreadyExists)
}

// --- DeleteRound Handler Tests ---

func TestHandleDeleteRound(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing sport parameter",
			queryParams:    "playDate=2024-01-01",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
		{
			name:           "missing playDate parameter",
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodDelete, "/v1/round?"+tt.queryParams, nil)

			getTestServer().DeleteRound(c)

			assertStatus(t, w, tt.expectedStatus)
			assertErrorCode(t, w, tt.expectedCode)
		})
	}
}

func TestHandleDeleteRound_Success(t *testing.T) {
	svc := &mockService{
		deleteRoundFn: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/v1/round?sport=basketball&playDate=2026-02-13", nil)

	server.DeleteRound(c)

	assertStatus(t, w, http.StatusNoContent)
}

func TestHandleDeleteRound_NotFound(t *testing.T) {
	svc := &mockService{
		deleteRoundFn: func(_ context.Context, _, _ string) error {
			return ErrRoundNotFound
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/v1/round?sport=basketball&playDate=2026-02-13", nil)

	server.DeleteRound(c)

	assertStatus(t, w, http.StatusNotFound)
	assertErrorCode(t, w, ErrorRoundNotFound)
}

// --- GetRounds Handler Tests ---

func TestHandleGetRounds_MissingSport(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/rounds", nil)

	getTestServer().GetRounds(c)

	assertStatus(t, w, http.StatusBadRequest)
	assertErrorCode(t, w, ErrorMissingRequiredParameter)
}

func TestHandleGetRounds_Success(t *testing.T) {
	svc := &mockService{
		getRoundsBySportFn: func(_ context.Context, _, _, _ string) ([]*RoundSummary, error) {
			return []*RoundSummary{
				{RoundID: "basketball#5", Sport: "basketball", PlayDate: "2026-02-13"},
			}, nil
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/rounds?sport=basketball", nil)

	server.GetRounds(c)

	assertStatus(t, w, http.StatusOK)
}

func TestHandleGetRounds_NoRoundsFound(t *testing.T) {
	svc := &mockService{
		getRoundsBySportFn: func(_ context.Context, _, _, _ string) ([]*RoundSummary, error) {
			return nil, ErrRoundNotFound
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/rounds?sport=basketball", nil)

	server.GetRounds(c)

	assertStatus(t, w, http.StatusNotFound)
	assertErrorCode(t, w, ErrorNoUpcomingRounds)
}

// --- GetUpcomingRounds Handler Tests ---

func TestHandleGetUpcomingRounds(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing sport parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/upcoming-rounds?"+tt.queryParams, nil)

			getTestServer().GetUpcomingRounds(c)

			assertStatus(t, w, tt.expectedStatus)
			assertErrorCode(t, w, tt.expectedCode)
		})
	}
}

// --- SubmitResults Handler Tests ---

func TestHandleSubmitResults(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		body           interface{}
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing query parameters",
			queryParams:    "sport=basketball",
			body:           Result{Score: 100, IsCorrect: true},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
		{
			name:           "missing both parameters",
			queryParams:    "",
			body:           Result{Score: 100, IsCorrect: true},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/results?"+tt.queryParams, bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			getTestServer().SubmitResults(c)

			assertStatus(t, w, tt.expectedStatus)
			assertErrorCode(t, w, tt.expectedCode)
		})
	}
}

func TestHandleSubmitResults_InvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/results?sport=basketball&playDate=2026-02-13", bytes.NewReader([]byte("bad json")))
	c.Request.Header.Set("Content-Type", "application/json")

	getTestServer().SubmitResults(c)

	assertStatus(t, w, http.StatusBadRequest)
	assertErrorCode(t, w, ErrorInvalidRequestBody)
}

func TestHandleSubmitResults_ScoreTooHigh(t *testing.T) {
	result := Result{Score: 101, IsCorrect: true}
	bodyBytes, _ := json.Marshal(result)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/results?sport=basketball&playDate=2026-02-13", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	getTestServer().SubmitResults(c)

	assertStatus(t, w, http.StatusBadRequest)
	assertErrorCode(t, w, ErrorInvalidRequestBody)
}

func TestHandleSubmitResults_ScoreNegative(t *testing.T) {
	result := Result{Score: -1, IsCorrect: false}
	bodyBytes, _ := json.Marshal(result)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/results?sport=basketball&playDate=2026-02-13", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	getTestServer().SubmitResults(c)

	assertStatus(t, w, http.StatusBadRequest)
	assertErrorCode(t, w, ErrorInvalidRequestBody)
}

func TestHandleSubmitResults_Success(t *testing.T) {
	svc := &mockService{
		submitResultsFn: func(_ context.Context, params SubmitResultsParams) (*ResultResponse, error) {
			return &ResultResponse{Result: params.Result}, nil
		},
	}
	server := NewServer(svc)

	result := Result{Score: 80, IsCorrect: true, PlayerName: "Test Player"}
	bodyBytes, _ := json.Marshal(result)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/results?sport=basketball&playDate=2026-02-13", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	server.SubmitResults(c)

	assertStatus(t, w, http.StatusOK)

	var resp ResultResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Score != 80 {
		t.Errorf("expected score 80, got %d", resp.Score)
	}
}

func TestHandleSubmitResults_RoundNotFound(t *testing.T) {
	svc := &mockService{
		submitResultsFn: func(_ context.Context, _ SubmitResultsParams) (*ResultResponse, error) {
			return nil, ErrRoundNotFound
		},
	}
	server := NewServer(svc)

	result := Result{Score: 80, IsCorrect: true}
	bodyBytes, _ := json.Marshal(result)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/results?sport=basketball&playDate=2026-02-13", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	server.SubmitResults(c)

	assertStatus(t, w, http.StatusNotFound)
	assertErrorCode(t, w, ErrorRoundNotFound)
}

// --- GetRoundStats Handler Tests ---

func TestHandleGetRoundStats(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing query parameters",
			queryParams:    "sport=basketball",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
		{
			name:           "missing both parameters",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/stats/round?"+tt.queryParams, nil)

			getTestServer().GetRoundStats(c)

			assertStatus(t, w, tt.expectedStatus)
			assertErrorCode(t, w, tt.expectedCode)
		})
	}
}

func TestHandleGetRoundStats_Success(t *testing.T) {
	svc := &mockService{
		getRoundStatsFn: func(_ context.Context, _, _ string) (*RoundStats, error) {
			return &RoundStats{Sport: "basketball", PlayDate: "2026-02-13", Name: "Test Player"}, nil
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/stats/round?sport=basketball&playDate=2026-02-13", nil)

	server.GetRoundStats(c)

	assertStatus(t, w, http.StatusOK)

	var resp RoundStats
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Sport != "basketball" {
		t.Errorf("expected sport basketball, got %s", resp.Sport)
	}
}

func TestHandleGetRoundStats_NotFound(t *testing.T) {
	svc := &mockService{
		getRoundStatsFn: func(_ context.Context, _, _ string) (*RoundStats, error) {
			return nil, ErrRoundNotFound
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/stats/round?sport=basketball&playDate=2026-02-13", nil)

	server.GetRoundStats(c)

	assertStatus(t, w, http.StatusNotFound)
	assertErrorCode(t, w, ErrorStatsNotFound)
}

// --- GetUser Handler Tests ---

func TestHandleGetUser(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "missing userId parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrorMissingRequiredParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/v1/user?"+tt.queryParams, nil)

			getTestServer().GetUser(c)

			assertStatus(t, w, tt.expectedStatus)
			assertErrorCode(t, w, tt.expectedCode)
		})
	}
}

func TestHandleGetUser_SuccessWithQueryParam(t *testing.T) {
	svc := &mockService{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return &User{UserId: "user-123", UserName: "testuser"}, nil
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/user?userId=user-123", nil)

	server.GetUser(c)

	assertStatus(t, w, http.StatusOK)

	var resp User
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.UserId != "user-123" {
		t.Errorf("expected userId user-123, got %s", resp.UserId)
	}
}

func TestHandleGetUser_SuccessWithJWTToken(t *testing.T) {
	svc := &mockService{
		getUserFn: func(_ context.Context, userId string) (*User, error) {
			return &User{UserId: userId, UserName: "testuser"}, nil
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/user", nil)
	c.Set(ConstantUserId, "jwt-user-id")

	server.GetUser(c)

	assertStatus(t, w, http.StatusOK)
}

func TestHandleGetUser_NotFound(t *testing.T) {
	svc := &mockService{
		getUserFn: func(_ context.Context, _ string) (*User, error) {
			return nil, ErrUserNotFound
		},
	}
	server := NewServer(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/user?userId=nonexistent", nil)

	server.GetUser(c)

	assertStatus(t, w, http.StatusNotFound)
	assertErrorCode(t, w, ErrorUserNotFound)
}

// --- MigrateUserStats Handler Tests ---

func TestHandleMigrateUserStats_MissingToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/stats/user/migrate", bytes.NewReader([]byte("{}")))
	c.Request.Header.Set("Content-Type", "application/json")

	getTestServer().MigrateUserStats(c)

	assertStatus(t, w, http.StatusUnauthorized)
}

func TestHandleMigrateUserStats_Success(t *testing.T) {
	svc := &mockService{
		migrateUserFn: func(_ context.Context, userId, username string, user *User) (*User, error) {
			user.UserId = userId
			user.UserName = username
			return user, nil
		},
	}
	server := NewServer(svc)

	user := User{Sports: []UserSportStats{{Sport: "basketball"}}}
	bodyBytes, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/stats/user/migrate", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(ConstantUserId, "user-123")
	c.Set(ConstantUsername, "testuser")

	server.MigrateUserStats(c)

	assertStatus(t, w, http.StatusCreated)

	var resp User
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.UserId != "user-123" {
		t.Errorf("expected userId user-123, got %s", resp.UserId)
	}
}

func TestHandleMigrateUserStats_AlreadyMigrated(t *testing.T) {
	svc := &mockService{
		migrateUserFn: func(_ context.Context, _, _ string, _ *User) (*User, error) {
			return nil, ErrUserAlreadyMigrated
		},
	}
	server := NewServer(svc)

	bodyBytes, _ := json.Marshal(User{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/stats/user/migrate", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(ConstantUserId, "user-123")
	c.Set(ConstantUsername, "testuser")

	server.MigrateUserStats(c)

	assertStatus(t, w, http.StatusConflict)
	assertErrorCode(t, w, ErrorUserAlreadyMigrated)
}

func TestHandleMigrateUserStats_InvalidBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/stats/user/migrate", bytes.NewReader([]byte("bad json")))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(ConstantUserId, "user-123")
	c.Set(ConstantUsername, "testuser")

	getTestServer().MigrateUserStats(c)

	assertStatus(t, w, http.StatusBadRequest)
	assertErrorCode(t, w, ErrorInvalidRequestBody)
}

// --- UpdateUsername Handler Tests ---

func TestHandleUpdateUsername_MissingToken(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/user/username", bytes.NewReader([]byte(`{"username":"test"}`)))
	c.Request.Header.Set("Content-Type", "application/json")

	getTestServer().UpdateUsername(c)

	assertStatus(t, w, http.StatusUnauthorized)
}

func TestHandleUpdateUsername_MissingUsername(t *testing.T) {
	bodyBytes, _ := json.Marshal(map[string]string{"username": ""})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/user/username", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(ConstantUserId, "user-123")

	getTestServer().UpdateUsername(c)

	assertStatus(t, w, http.StatusBadRequest)
	assertErrorCode(t, w, ErrorMissingRequiredField)
}

func TestHandleUpdateUsername_Success(t *testing.T) {
	svc := &mockService{
		updateUsernameFn: func(_ context.Context, _, _ string) error {
			return nil
		},
	}
	server := NewServer(svc)

	bodyBytes, _ := json.Marshal(map[string]string{"username": "newname"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/user/username", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(ConstantUserId, "user-123")

	server.UpdateUsername(c)

	assertStatus(t, w, http.StatusOK)

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "newname" {
		t.Errorf("expected username newname, got %v", resp["username"])
	}
}

func TestHandleUpdateUsername_Auth0Error(t *testing.T) {
	svc := &mockService{
		updateUsernameFn: func(_ context.Context, _, _ string) error {
			return errors.New("auth0 unavailable")
		},
	}
	server := NewServer(svc)

	bodyBytes, _ := json.Marshal(map[string]string{"username": "newname"})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/v1/user/username", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(ConstantUserId, "user-123")

	server.UpdateUsername(c)

	assertStatus(t, w, http.StatusInternalServerError)
	assertErrorCode(t, w, ErrorConfigurationError)
}
