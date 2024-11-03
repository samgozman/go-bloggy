package server

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/server/middlewares"
)

type Config struct {
	SentryDSN string
}

func ProvideConfig(cfg *config.Config) *Config {
	return &Config{
		SentryDSN: cfg.SentryDSN,
	}
}

// ProvideServer is a provider for the echo server.
func ProvideServer(cfg *Config, jwtService jwt.ServiceInterface) *echo.Echo {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

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
	server.Use(middleware.Recover())

	// Add the Sentry middleware
	server.Use(sentryecho.New(sentryecho.Options{}))

	return server
}

var ProviderSet = wire.NewSet( //nolint:gochecknoglobals // required by Wire
	ProvideConfig,
	ProvideServer,
)
