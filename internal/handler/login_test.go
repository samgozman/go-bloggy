package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"
)

func Test_PostLoginGithubAuthorize(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:?cache=shared")
	if errDB != nil {
		t.Fatal(errDB)
	}

	t.Run("OK", func(t *testing.T) {
		adminExternalID := rand.Int()
		e, mockGithubService, mockJwtService := registerHandlers(conn, []string{strconv.Itoa(adminExternalID)})

		ghUser := &github.UserInfo{
			ID:    adminExternalID,
			Login: uuid.New().String(),
		}
		ghUserID := strconv.Itoa(ghUser.ID)

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(ghUser, nil)

		mockJwtService.
			On("CreateTokenString", ghUserID, mock.Anything).
			Return("someToken", nil)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())
		mockGithubService.AssertExpectations(t)
		mockJwtService.AssertExpectations(t)

		var body client.JWTToken
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.NotEmpty(t, body.Token)
		assert.Equal(t, "someToken", body.Token)

		// Check if the user was created in the database
		dbUser, err := conn.Models.Users.GetByExternalID(context.Background(), ghUserID)
		assert.NoError(t, err)
		assert.Equal(t, ghUserID, dbUser.ExternalID)
		assert.Equal(t, ghUser.Login, dbUser.Login)
	})

	t.Run("should work for existing user", func(t *testing.T) {
		adminExternalID := rand.Int()
		e, mockGithubService, mockJwtService := registerHandlers(conn, []string{strconv.Itoa(adminExternalID)})

		ghUser := &github.UserInfo{
			ID:    adminExternalID,
			Login: uuid.New().String(),
		}
		ghUserID := strconv.Itoa(ghUser.ID)

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(ghUser, nil)

		mockJwtService.
			On("CreateTokenString", ghUserID, mock.Anything).
			Return("someToken", nil)

		// Create user in the database
		err := conn.Models.Users.Upsert(context.Background(), &models.User{
			ExternalID: ghUserID,
			Login:      ghUser.Login,
			AuthMethod: models.GitHubAuthMethod,
		})
		assert.NoError(t, err)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())
		mockGithubService.AssertExpectations(t)
		mockJwtService.AssertExpectations(t)

		var body client.JWTToken
		err = res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.NotEmpty(t, body.Token)
		assert.Equal(t, "someToken", body.Token)
	})

	t.Run("ValidationError", func(t *testing.T) {
		adminExternalID := rand.Int()
		e, _, _ := registerHandlers(conn, []string{strconv.Itoa(adminExternalID)})

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "", // empty code
		})

		res := testutil.NewRequest().
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
		adminExternalID := rand.Int()
		e, mockGithubService, _ := registerHandlers(conn, []string{strconv.Itoa(adminExternalID)})

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("", assert.AnError)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusInternalServerError, res.Code())
		mockGithubService.AssertExpectations(t)

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errExchangeCode, body.Code)
		assert.Equal(t, "Error while exchanging GitHub code for token", body.Message)
	})

	t.Run("GitHub GetUserInfo error", func(t *testing.T) {
		adminExternalID := rand.Int()
		e, mockGithubService, _ := registerHandlers(conn, []string{strconv.Itoa(adminExternalID)})

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(&github.UserInfo{}, assert.AnError)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusInternalServerError, res.Code())
		mockGithubService.AssertExpectations(t)

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errGetUserInfo, body.Code)
		assert.Equal(t, "Error while getting user info from GitHub", body.Message)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		adminExternalID := rand.Int()
		e, _, _ := registerHandlers(conn, []string{strconv.Itoa(adminExternalID)})

		res := testutil.NewRequest().
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
		adminExternalID := rand.Int()
		e, mockGithubService, mockJwtService := registerHandlers(conn, []string{strconv.Itoa(adminExternalID)})

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(&github.UserInfo{
				ID:    adminExternalID,
				Login: "testUser",
			}, nil)

		mockJwtService.
			On("CreateTokenString", strconv.Itoa(adminExternalID), mock.Anything).
			Return("", assert.AnError)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusInternalServerError, res.Code())
		mockGithubService.AssertExpectations(t)
		mockJwtService.AssertExpectations(t)

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errCreateToken, body.Code)
		assert.Equal(t, "Error while creating JWT token", body.Message)
	})

	t.Run("Auth forbidden for non-admin", func(t *testing.T) {
		// Fake admin id
		e, mockGithubService, _ := registerHandlers(conn, []string{"000000"})

		mockGithubService.
			On("ExchangeCodeForToken", mock.Anything, "123").
			Return("someToken", nil)

		mockGithubService.
			On("GetUserInfo", mock.Anything, "someToken").
			Return(&github.UserInfo{
				ID:    rand.Int(),
				Login: "testUser",
			}, nil)

		rb, _ := json.Marshal(client.GitHubAuthRequestBody{
			Code: "123",
		})

		res := testutil.NewRequest().
			Post("/login/github/authorize").
			WithHeader("Content-Type", "application/json").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusForbidden, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errForbidden, body.Code)
		assert.Equal(t, "User is not an admin", body.Message)
	})
}

func Test_PostLoginRefresh(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(nil, nil)

		mockJwtService.
			On("ParseTokenString", "token").
			Return("123", nil)

		mockJwtService.
			On("CreateTokenString", "123", mock.Anything).
			Return("someToken", nil)

		res := testutil.NewRequest().
			Post("/login/refresh").
			WithJWSAuth("token").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())
		mockJwtService.AssertExpectations(t)

		var body client.JWTToken
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)
		assert.Equal(t, "someToken", body.Token)
	})

	t.Run("Forbidden Authorization header is required", func(t *testing.T) {
		e, _, _ := registerHandlers(nil, nil)

		res := testutil.NewRequest().Post("/login/refresh").GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusUnauthorized, res.Code())

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errForbidden, body.Code)
		assert.Equal(t, "Authorization header is required", body.Message)
	})

	t.Run("JWTService ParseTokenString error", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(nil, nil)

		mockJwtService.
			On("ParseTokenString", "token").
			Return("", assert.AnError)

		res := testutil.NewRequest().
			Post("/login/refresh").
			WithJWSAuth("token").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusUnauthorized, res.Code())
		mockJwtService.AssertExpectations(t)

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errForbidden, body.Code)
		assert.Equal(t, "Invalid token", body.Message)
	})

	t.Run("CreateTokenString error", func(t *testing.T) {
		e, _, mockJwtService := registerHandlers(nil, nil)

		mockJwtService.
			On("ParseTokenString", "token").
			Return("123", nil)

		mockJwtService.
			On("CreateTokenString", "123", mock.Anything).
			Return("", assert.AnError)

		res := testutil.NewRequest().
			Post("/login/refresh").
			WithJWSAuth("token").
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusInternalServerError, res.Code())
		mockJwtService.AssertExpectations(t)

		var body client.RequestError
		err := res.UnmarshalBodyToObject(&body)
		assert.NoError(t, err)

		assert.Equal(t, errCreateToken, body.Code)
		assert.Equal(t, "Error while creating JWT token", body.Message)
	})
}
