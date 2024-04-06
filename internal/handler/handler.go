package handler

import (
	"github.com/samgozman/go-bloggy/internal/captcha"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/mailer"
)

// Handler for the service API endpoints.
type Handler struct {
	githubService     github.ServiceInterface
	jwtService        jwt.ServiceInterface
	hcaptchaService   captcha.ClientInterface
	db                *db.Database
	mailerService     mailer.ServiceInterface
	adminsExternalIDs []string
}

// NewHandler creates a new Handler.
func NewHandler(
	g github.ServiceInterface,
	j jwt.ServiceInterface,
	db *db.Database,
	h captcha.ClientInterface,
	ms mailer.ServiceInterface,
	adminsExternalIDs []string,
) *Handler {
	return &Handler{
		githubService:     g,
		jwtService:        j,
		db:                db,
		hcaptchaService:   h,
		mailerService:     ms,
		adminsExternalIDs: adminsExternalIDs,
	}
}
