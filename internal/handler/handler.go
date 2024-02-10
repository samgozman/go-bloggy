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
