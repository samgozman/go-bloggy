package main

import (
	"github.com/labstack/echo/v4"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/pkg/client"
)

func main() {
	config := NewConfigFromEnv()
	dnConn, err := db.InitDatabase(config.DSN)
	if err != nil {
		panic(err)
	}

	ghService := github.NewService(config.GithubClientID, config.GithubClientSecret)
	jwtService := jwt.NewService(config.JWTSecretKey)

	apiHandler := handler.NewHandler(ghService, jwtService, dnConn, config.AdminsExternalIDs)
	server := echo.New()
	client.RegisterHandlers(server, apiHandler)

	server.Logger.Fatal(server.Start(":" + config.Port))
}
