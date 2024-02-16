package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/pkg/client"
	"net/http"
)

// Handler for the service API endpoints.
type Handler struct {
	githubService githubService
}

// NewHandler creates a new Handler.
func NewHandler(githubService githubService) *Handler {
	return &Handler{
		githubService: githubService,
	}
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

	token, err := s.githubService.ExchangeCodeForToken(ctx.Request().Context(), req.Code)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errExchangeCode,
			Message: "Error while exchanging GitHub code for token",
		})
	}

	// TODO: Generate JWT token
	// TODO: Save data to DB (or update if exists)

	return ctx.JSON(http.StatusOK, ctx.JSON(http.StatusOK, client.JWTToken{
		Token: token,
	}))
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

// githubService is an interface for the github.Service.
type githubService interface {
	ExchangeCodeForToken(ctx context.Context, code string) (string, error)
	GetUserInfo(ctx context.Context, token string) (*github.UserInfo, error)
}
