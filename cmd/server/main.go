package main

import (
	"github.com/kataras/hcaptcha"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapi "github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/mailer"
	"github.com/samgozman/go-bloggy/internal/middlewares"
)

func main() {
	cfg := config.NewConfigFromEnv()

	dnConn, err := db.InitDatabase(cfg.DSN)
	if err != nil {
		panic(err)
	}

	// TODO: Replace initialization with google/wire

	ghService := github.NewService(cfg.GithubClientID, cfg.GithubClientSecret)
	jwtService := jwt.NewService(cfg.JWTSecretKey)
	hcaptchaService := hcaptcha.New(cfg.HCaptchaSecret)
	mailerService := mailer.NewService(cfg.MailerJet.PublicKey, cfg.MailerJet.PrivateKey, &mailer.Options{
		FromEmail:                    cfg.MailerJet.FromEmail,
		FromName:                     cfg.MailerJet.FromName,
		ConfirmationTemplateID:       cfg.MailerJet.ConfirmationTemplateID,
		ConfirmationTemplateURLParam: cfg.MailerJet.ConfirmationTemplateURLParam,
		PostTemplateID:               cfg.MailerJet.PostTemplateID,
		PostTemplateURLParam:         cfg.MailerJet.PostTemplateURLParam,
		UnsubscribeURLParam:          cfg.MailerJet.UnsubscribeURLParam,
	})

	apiHandler := handler.NewHandler(
		ghService,
		jwtService,
		dnConn,
		hcaptchaService,
		mailerService,
		cfg.AdminsExternalIDs,
	)
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

	oapi.RegisterHandlers(server, apiHandler)

	server.Logger.Fatal(server.Start(":" + cfg.Port))
}
