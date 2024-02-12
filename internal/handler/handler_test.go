package handler

import (
	"encoding/json"
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/pkg/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
)

func Test_GetHealth(t *testing.T) {
	e := registerHandlers()

	res := testutil.NewRequest().Get("/health").GoWithHTTPHandler(t, e)

	var body client.HealthCheckResponse
	err := res.UnmarshalBodyToObject(&body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.Code())
	assert.Equal(t, "OK", body.Status)
}

func Test_PostLoginGithubAuthorize_OK(t *testing.T) {
	e := registerHandlers()

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

	assert.Equal(t, "", body.Token) // TODO: check token
}

func Test_PostLoginGithubAuthorize_ValidationError(t *testing.T) {
	e := registerHandlers()

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
}

func Test_PostLoginRefresh_OK(t *testing.T) {
	e := registerHandlers()

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
}

func Test_PostLoginRefresh_Forbidden(t *testing.T) {
	e := registerHandlers()

	res := testutil.NewRequest().Post("/login/refresh").GoWithHTTPHandler(t, e)

	assert.Equal(t, http.StatusUnauthorized, res.Code())

	var body client.RequestError
	err := res.UnmarshalBodyToObject(&body)
	assert.NoError(t, err)

	assert.Equal(t, errForbidden, body.Code)
	assert.Equal(t, "Authorization header is required", body.Message)
}

func registerHandlers() *echo.Echo {
	e := echo.New()
	h := NewHandler()
	client.RegisterHandlers(e, h)

	return e
}
