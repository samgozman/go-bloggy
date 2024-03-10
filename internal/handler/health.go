package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/api"
	"net/http"
)

// GetHealth returns health status of the service.
func (h *Handler) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, api.HealthCheckResponse{
		Status: "OK",
	})
}
