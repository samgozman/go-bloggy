package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/pkg/client"
	"net/http"
)

// GetHealth returns health status of the service.
func (h *Handler) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, client.HealthCheckResponse{
		Status: "OK",
	})
}
