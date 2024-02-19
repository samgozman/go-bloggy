package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/pkg/client"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Handler for the service API endpoints.
type Handler struct {
	githubService     githubService
	jwtService        jwtService
	db                *db.Database
	adminsExternalIDs []string
}

// NewHandler creates a new Handler.
func NewHandler(
	g githubService,
	j jwtService,
	dbConn *db.Database,
	adminsExternalIDs []string,
) *Handler {
	return &Handler{
		githubService:     g,
		jwtService:        j,
		db:                dbConn,
		adminsExternalIDs: adminsExternalIDs,
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

	user, err := s.githubService.GetUserInfo(ctx.Request().Context(), token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errGetUserInfo,
			Message: "Error while getting user info from GitHub",
		})
	}

	// Check if user is an admin
	if !slices.Contains(s.adminsExternalIDs, strconv.Itoa(user.ID)) {
		return ctx.JSON(http.StatusForbidden, client.RequestError{
			Code:    errForbidden,
			Message: "User is not an admin",
		})
	}

	// TODO: Store JWT expiration time in config
	jwtToken, err := s.jwtService.CreateTokenString(strconv.Itoa(user.ID), time.Now().Add(time.Minute))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errCreateToken,
			Message: "Error while creating JWT token",
		})
	}

	err = s.db.Models.Users.Upsert(ctx.Request().Context(), &models.User{
		ExternalID: strconv.Itoa(user.ID),
		Login:      user.Login,
		AuthMethod: models.GitHubAuthMethod,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errCreateUser,
			Message: "Error while creating user",
		})
	}

	return ctx.JSON(http.StatusOK, ctx.JSON(http.StatusOK, client.JWTToken{
		Token: jwtToken,
	}))
}

// PostLoginRefresh handles the request to refresh the JWT token.
func (s *Handler) PostLoginRefresh(ctx echo.Context) error {
	token := ctx.Request().Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return ctx.JSON(http.StatusUnauthorized, client.RequestError{
			Code:    errForbidden,
			Message: "Authorization header is required",
		})
	}

	userID, err := s.jwtService.ParseTokenString(token)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, client.RequestError{
			Code:    errForbidden,
			Message: "Invalid token",
		})
	}

	// TODO: Store JWT expiration time in config
	newToken, err := s.jwtService.CreateTokenString(userID, time.Now().Add(time.Minute))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errCreateToken,
			Message: "Error while creating JWT token",
		})
	}

	return ctx.JSON(http.StatusOK, client.JWTToken{
		Token: newToken,
	})
}

// githubService is an interface for the github.Service.
type githubService interface {
	ExchangeCodeForToken(ctx context.Context, code string) (string, error)
	GetUserInfo(ctx context.Context, token string) (*github.UserInfo, error)
}

type jwtService interface {
	CreateTokenString(userID string, expiresAt time.Time) (jwtToken string, err error)
	ParseTokenString(tokenString string) (userID string, err error)
}
