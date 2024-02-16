package handler

import (
	"context"
	"encoding/json"
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

type MockGithubService struct {
	mock.Mock
}

func (m *MockGithubService) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)          //nolint:typecheck
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

func (m *MockGithubService) GetUserInfo(ctx context.Context, token string) (*github.UserInfo, error) {
	args := m.Called(ctx, token)                         //nolint:typecheck
	return args.Get(0).(*github.UserInfo), args.Error(1) //nolint:wrapcheck
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) CreateTokenString(userID string, expiresAt time.Time) (string, error) {
	args := m.Called(userID, expiresAt)  //nolint:typecheck
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

func (m *MockJWTService) ParseTokenString(token string) (string, error) {
	args := m.Called(token) //nolint:typecheck
	return args.String(0), args.Error(1)
}

func Test_GetHealth(t *testing.T) {
	e, _, _ := registerHandlers()

	res := testutil.NewRequest().Get("/health").GoWithHTTPHandler(t, e)

	var body client.HealthCheckResponse
	err := res.UnmarshalBodyToObject(&body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.Code())
	assert.Equal(t, "OK", body.Status)
}

func Test_PostLoginGithubAuthorize(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		e, mockGithubService, mockJwtService := registerHandlers()

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(&github.UserInfo{
				ID:    123,
				Login: "testUser",
			}, nil)

		mockJwtService.
			On("CreateTokenString", "123", mock.Anything).
			Return("someToken", nil)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.
			NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		var body client.JWTToken
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.NotEmpty(t, body.Token)
		assert.Equal(t, "someToken", body.Token)
	})

	t.Run("ValidationError", func(t *testing.T) {
		e, _, _ := registerHandlers()

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "", // empty code
		})

		res := testutil.
			NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errBodyValidation, body.Code)
		assert.Equal(t, "Code field is required", body.Message)
	})

	t.Run("GitHub ExchangeCodeForToken error", func(t *testing.T) {
		e, mockGithubService, _ := registerHandlers()

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("", assert.AnError)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.
			NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusInternalServerError, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errExchangeCode, body.Code)
		assert.Equal(t, "Error while exchanging GitHub code for token", body.Message)
	})

	t.Run("GitHub GetUserInfo error", func(t *testing.T) {
		e, mockGithubService, _ := registerHandlers()

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(&github.UserInfo{}, assert.AnError)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.
			NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusInternalServerError, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errGetUserInfo, body.Code)
		assert.Equal(t, "Error while getting user info from GitHub", body.Message)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		e, _, _ := registerHandlers()

		res := testutil.
			NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody([]byte("invalid json")).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errRequestBodyBinding, body.Code)
		assert.Equal(t, "Error binding request body", body.Message)
	})

	t.Run("JWTService CreateTokenString error", func(t *testing.T) {
		e, mockGithubService, mockJwtService := registerHandlers()

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(&github.UserInfo{
				ID:    123,
				Login: "testUser",
			}, nil)

		mockJwtService.
			On("CreateTokenString", "123", mock.Anything).
			Return("", assert.AnError)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.
			NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusInternalServerError, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errCreateToken, body.Code)
		assert.Equal(t, "Error while creating JWT token", body.Message)
	})
}

func Test_PostLoginRefresh(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		e, _, _ := registerHandlers()

		res := testutil.
			NewRequest().
			Post("/login/refresh").
			WithJWSAuth("token").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		var body client.JWTToken
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, "", body.Token) // TODO: check token
	})

	t.Run("Forbidden", func(t *testing.T) {
		e, _, _ := registerHandlers()

		res := testutil.NewRequest().Post("/login/refresh").GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusUnauthorized, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errForbidden, body.Code)
		assert.Equal(t, "Authorization header is required", body.Message)
	})
}

func registerHandlers() (server *echo.Echo, githubService *MockGithubService, jwtService *MockJWTService) {
	// Create mocks
	g := new(MockGithubService)
	j := new(MockJWTService)

	// Create echo instance
	e := echo.New()
	h := NewHandler(g, j)
	client.RegisterHandlers(e, h)

	return e, g, j
}
