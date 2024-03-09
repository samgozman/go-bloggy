package handler

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/samgozman/go-bloggy/pkg/server"
	"net/http"
	"os"
	"regexp"
)

func (h *Handler) PostSubscriptions(ctx echo.Context) error {
	var req server.SubscriptionRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, server.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	// validate email
	if !isValidEmail(req.Email) {
		return ctx.JSON(http.StatusBadRequest, server.RequestError{
			Code:    errValidationEmail,
			Message: "Invalid email",
		})
	}

	// TODO: Get ENVIRONMENT from config
	if os.Getenv("ENVIRONMENT") == "production" {
		if hr := h.hcaptchaService.VerifyToken(req.Captcha); !hr.Success {
			return ctx.JSON(http.StatusBadRequest, server.RequestError{
				Code:    errValidationCaptcha,
				Message: "Invalid captcha",
			})
		}
	}

	subscription := models.Subscription{
		Email: req.Email,
	}

	if err := h.db.Models.Subscriptions.Create(ctx.Request().Context(), &subscription); err != nil {
		// Note: we shouldn't tell duplicate error to the user for security reasons
		if !errors.Is(err, models.ErrDuplicate) {
			return ctx.JSON(http.StatusInternalServerError, server.RequestError{
				Code:    errCreateSubscription,
				Message: "Error creating subscription",
			})
		}
	}

	return ctx.NoContent(http.StatusCreated)
}

func (h *Handler) DeleteSubscriptions(ctx echo.Context) error {
	var req server.UnsubscribeRequest
	if err := ctx.Bind(&req); err != nil {
		var errorMessage string
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			errorMessage = fmt.Sprintf("%v", echoErr.Message)
		}

		return ctx.JSON(http.StatusBadRequest, server.RequestError{
			Code:    errRequestBodyBinding,
			Message: fmt.Sprintf("Error binding request body: %v", errorMessage),
		})
	}

	subscriptionID, err := h.db.Models.Subscriptions.GetByID(ctx.Request().Context(), req.SubscriptionId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, server.RequestError{
			Code:    errGetSubscription,
			Message: "Subscription is not found or error getting subscription by ID",
		})
	}

	if err := h.db.Models.Subscriptions.Delete(ctx.Request().Context(), subscriptionID.ID.String()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, server.RequestError{
			Code:    errDeleteSubscription,
			Message: "Error deleting subscription",
		})
	}

	// TODO: Log req.Reason for unsubscribing in Sentry

	return ctx.NoContent(http.StatusNoContent)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,24}$`)
	return re.MatchString(email)
}
