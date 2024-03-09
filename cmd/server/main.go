package main

import (
	"github.com/kataras/hcaptcha"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/middlewares"
	oapi "github.com/samgozman/go-bloggy/pkg/server"
)

func main() {
	config := NewConfigFromEnv()
	dnConn, err := db.InitDatabase(config.DSN)
	if err != nil {
		panic(err)
	}

	// TODO: Replace initialization with google/wire

	ghService := github.NewService(config.GithubClientID, config.GithubClientSecret)
	jwtService := jwt.NewService(config.JWTSecretKey)
	hcaptchaService := hcaptcha.New(config.HCaptchaSecret)

	apiHandler := handler.NewHandler(
		ghService,
		jwtService,
		dnConn,
		hcaptchaService,
		config.AdminsExternalIDs,
	)
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: Pass allowed origins from config
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))
	server.Use(middlewares.JWTAuth(jwtService))

	oapi.RegisterHandlers(server, apiHandler)

	server.Logger.Fatal(server.Start(":" + config.Port))
}
