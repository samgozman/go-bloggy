package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/pkg/client"
	"net/http"
)

// Handler for the service API endpoints.
type Handler struct{}

// NewHandler creates a new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

// GetHealth returns health status of the service.
func (s *Handler) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, client.HealthCheckResponse{
		Status: "OK",
	})
}

// PostLoginGithubAuthorize handles the request to authorize with GitHub.
func (s *Handler) PostLoginGithubAuthorize(ctx echo.Context) error {
	var req client.GitHubAuthRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, client.RequestError{
			Code:    errRequestBodyBinding,
			Message: "Error binding request body",
		})
	}

	if req.Code == "" {
		return ctx.JSON(http.StatusBadRequest, client.RequestError{
			Code:    errBodyValidation,
			Message: "Code field is required",
		})
	}

	// TODO: request to github
	// TODO: Generate JWT token
	// TODO: Save data to DB (or update if exists)

	return ctx.JSON(http.StatusOK, client.PostLoginGithubAuthorizeResponse{
		JSON200: &client.JWTToken{
			Token: "",
		},
	})
}

// PostLoginRefresh handles the request to refresh the JWT token.
func (s *Handler) PostLoginRefresh(ctx echo.Context) error {
	token := ctx.Request().Header.Get("Authorization")
	if token == "" {
		return ctx.JSON(http.StatusUnauthorized, client.RequestError{
			Code:    errForbidden,
			Message: "Authorization header is required",
		})
	}

	// TODO: parse token and get data
	// TODO: Check data in DB
	// TODO: Check if GitHub token from DB is still valid (if possible)
	// TODO: Generate new JWT token

	return ctx.JSON(http.StatusOK, client.JWTToken{
		Token: "",
	})
}
