package main

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/pkg/client"
)

func main() {
	e := echo.New()

	// TODO: add Service clientID and clientSecret here from environment variables
	g := github.NewService("clientID", "clientSecret")
	j := jwt.NewService("testKey")
	h := handler.NewHandler(g, j)

	client.RegisterHandlers(e, h)

	e.Logger.Fatal(e.Start(":80"))
}
