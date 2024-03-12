package handler

import (
	"context"
	"github.com/kataras/hcaptcha"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/mailer"
	"time"
)

// Handler for the service API endpoints.
type Handler struct {
	githubService     githubService
	jwtService        jwtService
	hcaptchaService   hcaptchaService
	db                *db.Database
	mailerService     mailerService
	adminsExternalIDs []string
}

// NewHandler creates a new Handler.
func NewHandler(
	g githubService,
	j jwtService,
	db *db.Database,
	h hcaptchaService,
	ms mailerService,
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

// githubService is an interface for the github.Service.
type githubService interface {
	ExchangeCodeForToken(ctx context.Context, code string) (string, error)
	GetUserInfo(ctx context.Context, token string) (*github.UserInfo, error)
}

type jwtService interface {
	CreateTokenString(userID string, expiresAt time.Time) (jwtToken string, err error)
	ParseTokenString(tokenString string) (externalUserID string, err error)
}

type hcaptchaService interface {
	VerifyToken(tkn string) (response hcaptcha.Response)
}

type mailerService interface {
	SendConfirmationEmail(to, confirmationID string) error
	SendPostEmail(pe *mailer.PostEmailSend) error
}
