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

func Test_GetHealth(t *testing.T) {
	mockGithubService := new(MockGithubService)
	e := registerHandlers(mockGithubService)

	res := testutil.NewRequest().Get("/health").GoWithHTTPHandler(t, e)

	var body client.HealthCheckResponse
	err := res.UnmarshalBodyToObject(&body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.Code())
	assert.Equal(t, "OK", body.Status)
}

func Test_PostLoginGithubAuthorize(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockGithubService := new(MockGithubService)
		e := registerHandlers(mockGithubService)

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
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

		// TODO: Check for JWT token
		assert.Equal(t, "someToken", body.Token)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockGithubService := new(MockGithubService)
		e := registerHandlers(mockGithubService)

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

	t.Run("GitHub error", func(t *testing.T) {
		mockGithubService := new(MockGithubService)
		e := registerHandlers(mockGithubService)

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
}

func Test_PostLoginRefresh(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockGithubService := new(MockGithubService)
		e := registerHandlers(mockGithubService)

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
		mockGithubService := new(MockGithubService)
		e := registerHandlers(mockGithubService)

		res := testutil.NewRequest().Post("/login/refresh").GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusUnauthorized, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errForbidden, body.Code)
		assert.Equal(t, "Authorization header is required", body.Message)
	})
}

func registerHandlers(githubService githubService) *echo.Echo {
	e := echo.New()
	h := NewHandler(githubService)
	client.RegisterHandlers(e, h)

	return e
}
