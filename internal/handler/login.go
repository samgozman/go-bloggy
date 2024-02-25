package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/samgozman/go-bloggy/pkg/client"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// PostLoginGithubAuthorize handles the request to authorize with GitHub.
func (h *Handler) PostLoginGithubAuthorize(ctx echo.Context) error {
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

	token, err := h.githubService.ExchangeCodeForToken(ctx.Request().Context(), req.Code)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errExchangeCode,
			Message: "Error while exchanging GitHub code for token",
		})
	}

	user, err := h.githubService.GetUserInfo(ctx.Request().Context(), token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errGetUserInfo,
			Message: "Error while getting user info from GitHub",
		})
	}

	// Check if user is an admin
	if !slices.Contains(h.adminsExternalIDs, strconv.Itoa(user.ID)) {
		return ctx.JSON(http.StatusForbidden, client.RequestError{
			Code:    errForbidden,
			Message: "User is not an admin",
		})
	}

	// TODO: Store JWT expiration time in config
	jwtToken, err := h.jwtService.CreateTokenString(strconv.Itoa(user.ID), time.Now().Add(time.Minute))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, client.RequestError{
			Code:    errCreateToken,
			Message: "Error while creating JWT token",
		})
	}

	err = h.db.Models.Users.Upsert(ctx.Request().Context(), &models.User{
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
func (h *Handler) PostLoginRefresh(ctx echo.Context) error {
	token := ctx.Request().Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return ctx.JSON(http.StatusUnauthorized, client.RequestError{
			Code:    errForbidden,
			Message: "Authorization header is required",
		})
	}

	userID, err := h.jwtService.ParseTokenString(token)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, client.RequestError{
			Code:    errForbidden,
			Message: "Invalid token",
		})
	}

	// TODO: Store JWT expiration time in config
	newToken, err := h.jwtService.CreateTokenString(userID, time.Now().Add(time.Minute))
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
