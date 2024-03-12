package main

import (
	"github.com/kataras/hcaptcha"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapi "github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/mailer"
	"github.com/samgozman/go-bloggy/internal/middlewares"
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
	mailerService := mailer.NewService(config.MailerJet.PublicKey, config.MailerJet.PrivateKey, &mailer.Options{
		FromEmail:                    config.MailerJet.FromEmail,
		FromName:                     config.MailerJet.FromName,
		ConfirmationTemplateID:       config.MailerJet.ConfirmationTemplateID,
		ConfirmationTemplateURLParam: config.MailerJet.ConfirmationTemplateURLParam,
		PostTemplateID:               config.MailerJet.PostTemplateID,
		PostTemplateURLParam:         config.MailerJet.PostTemplateURLParam,
		UnsubscribeURLParam:          config.MailerJet.UnsubscribeURLParam,
	})

	apiHandler := handler.NewHandler(
		ghService,
		jwtService,
		dnConn,
		hcaptchaService,
		mailerService,
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
