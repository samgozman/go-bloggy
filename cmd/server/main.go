package main

import (
	"github.com/kataras/hcaptcha"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapi "github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/captcha"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/mailer"
	"github.com/samgozman/go-bloggy/internal/middlewares"
)

// TODO: 1. Replace all structs with interfaces
// TODO: 2. Get rid of the tempApp struct
// TODO: 3. Fix visibility of the providers init/new functions
// TODO: 4. Add service mocks generation by interface
// TODO: 5. Add generate command to the Makefile

func newTempApp(
	database db.DatabaseInterface,
	gh github.ServiceInterface,
	jwt jwt.ServiceInterface,
	cap captcha.ClientInterface,
	mail mailer.ServiceInterface,
) *tempApp {
	return &tempApp{
		Database:      database,
		GithubService: gh,
		JWTService:    jwt,
		Captcha:       cap,
		Mailer:        mail,
	}
}

type tempApp struct {
	Database      db.DatabaseInterface
	GithubService github.ServiceInterface
	JWTService    jwt.ServiceInterface
	Captcha       captcha.ClientInterface
	Mailer        mailer.ServiceInterface
}

func main() {
	cfg := config.NewConfigFromEnv()

	dnConn, err := db.InitDatabase(string(cfg.DSN))
	if err != nil {
		panic(err)
	}

	// TODO: Replace initialization with google/wire

	ghService := github.NewService(cfg.GithubClientID, cfg.GithubClientSecret)
	jwtService := jwt.NewService(string(cfg.JWTSecretKey))
	hcaptchaService := hcaptcha.New(string(cfg.HCaptchaSecret))
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
