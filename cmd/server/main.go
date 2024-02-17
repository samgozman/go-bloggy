package main

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/pkg/client"
)

func main() {
	c := NewConfigFromEnv()
	g := github.NewService(c.GithubClientID, c.GithubClientSecret)
	j := jwt.NewService(c.JWTSecretKey)
	h := handler.NewHandler(g, j)
	e := echo.New()
	client.RegisterHandlers(e, h)
	e.Logger.Fatal(e.Start(":" + c.Port))
}
