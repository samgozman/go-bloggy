package main

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	oapi "github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/config"
)

// TODO: 1. Fix visibility of the providers init/new functions

func newServerApp(
	server *echo.Echo,
	handler oapi.ServerInterface,
) *serverApp {
	return &serverApp{
		Server:  server,
		Handler: handler,
	}
}

type serverApp struct {
	Server  *echo.Echo
	Handler oapi.ServerInterface
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.NewConfigFromEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := initApp(ctx, cfg)
	if err != nil {
		return err
	}

	oapi.RegisterHandlers(app.Server, app.Handler)

	if err := app.Server.Start(":" + cfg.Port); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
