package handler

import (
	"github.com/labstack/echo/v4"
	"testing"

	"github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/server/middlewares"
	captchaMock "github.com/samgozman/go-bloggy/mocks/captcha"
	mockGithub "github.com/samgozman/go-bloggy/mocks/github"
	jwtMock "github.com/samgozman/go-bloggy/mocks/jwt"
	mockMailer "github.com/samgozman/go-bloggy/mocks/mailer"
)

// registerHandlers creates a new echo instance and registers the handlers for testing.
func registerHandlers(t *testing.T, conn *db.Database, adminsIDs []string) (
	s *echo.Echo,
	githubService *mockGithub.MockServiceInterface,
	jwtService *jwtMock.MockServiceInterface,
	mailerService *mockMailer.MockServiceInterface,
	captchaService *captchaMock.MockClientInterface,
) {
	// Create mocks
	g := mockGithub.NewMockServiceInterface(t)
	j := jwtMock.NewMockServiceInterface(t)
	hc := captchaMock.NewMockClientInterface(t)
	ms := mockMailer.NewMockServiceInterface(t)

	// Create echo instance
	e := echo.New()
	cfg := ProvideConfig(&config.Config{AdminsExternalIDs: adminsIDs})
	h := ProvideHandler(cfg, g, j, conn, hc, ms)
	e.Use(middlewares.JWTAuth(j))

	api.RegisterHandlers(e, h)

	return e, g, j, ms, hc
}
