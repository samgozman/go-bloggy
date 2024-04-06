package main

import (
	"context"
	"github.com/labstack/echo/v4"
	oapi "github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/config"
)

// TODO: 1. Fix visibility of the providers init/new functions
// TODO: 2. Add service mocks generation by interface
// TODO: 3. Add generate command to the Makefile
// TODO: 4. Fix tests

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
	cfg := config.NewConfigFromEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := initApp(ctx, cfg)
	if err != nil {
		panic(err)
	}

	oapi.RegisterHandlers(app.Server, app.Handler)

	app.Server.Logger.Fatal(app.Server.Start(":" + cfg.Port))
}
