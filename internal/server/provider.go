package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/server/middlewares"
)

// ProvideServer is a provider for the echo server.
func ProvideServer(jwtService jwt.ServiceInterface) *echo.Echo {
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: Pass allowed origins from cfg
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))
	server.Use(middlewares.JWTAuth(jwtService))

	return server
}
