package handler

import (
	"github.com/google/wire"
	oapi "github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/captcha"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/jwt"
	mailer "github.com/samgozman/go-bloggy/internal/mailer/types"
)

type Config struct {
	AdminsExternalIDs config.AdminsExternalIDs
}

// Handler for the service API endpoints.
type Handler struct {
	githubService     github.ServiceInterface
	jwtService        jwt.ServiceInterface
	hcaptchaService   captcha.ClientInterface
	db                *db.Database
	mailerService     mailer.ServiceInterface
	adminsExternalIDs []string
}

func ProvideConfig(cfg *config.Config) *Config {
	return &Config{
		AdminsExternalIDs: cfg.AdminsExternalIDs,
	}
}

// ProvideHandler is a wire provider function that creates a new Handler.
func ProvideHandler(
	cfg *Config,
	g github.ServiceInterface,
	j jwt.ServiceInterface,
	db *db.Database,
	h captcha.ClientInterface,
	ms mailer.ServiceInterface,
) *Handler {
	return &Handler{
		githubService:     g,
		jwtService:        j,
		db:                db,
		hcaptchaService:   h,
		mailerService:     ms,
		adminsExternalIDs: cfg.AdminsExternalIDs,
	}
}

// ProviderSet is a wire provider set that includes all the providers for the handler package.
var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideHandler,
	wire.Bind(new(oapi.ServerInterface), new(*Handler)),
)
