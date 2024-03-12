package handler

import (
	"context"
	"github.com/kataras/hcaptcha"
	"github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/mailer"
	"github.com/samgozman/go-bloggy/internal/middlewares"
	"github.com/stretchr/testify/mock"
	"time"

	"github.com/labstack/echo/v4"
)

type MockGithubService struct {
	mock.Mock
}

func (m *MockGithubService) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

func (m *MockGithubService) GetUserInfo(ctx context.Context, token string) (*github.UserInfo, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*github.UserInfo), args.Error(1) //nolint:wrapcheck
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) CreateTokenString(userID string, expiresAt time.Time) (string, error) {
	args := m.Called(userID, expiresAt)
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

func (m *MockJWTService) ParseTokenString(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

type MockHCaptchaService struct {
	mock.Mock
}

func (m *MockHCaptchaService) VerifyToken(tkn string) (response hcaptcha.Response) {
	args := m.Called(tkn)
	return args.Get(0).(hcaptcha.Response)
}

type MockMailerService struct {
	mock.Mock
}

func (m *MockMailerService) SendConfirmationEmail(to, confirmationID string) error {
	args := m.Called(to, confirmationID)
	return args.Error(0) //nolint:wrapcheck
}

func (m *MockMailerService) SendPostEmail(pe *mailer.PostEmailSend) error {
	args := m.Called(pe)
	return args.Error(0) //nolint:wrapcheck
}

// registerHandlers creates a new echo instance and registers the handlers for testing.
func registerHandlers(conn *db.Database, adminsIDs []string) (s *echo.Echo, githubService *MockGithubService, jwtService *MockJWTService) {
	// Create mocks
	g := new(MockGithubService)
	j := new(MockJWTService)
	hc := new(MockHCaptchaService)
	ms := new(MockMailerService)

	// Create echo instance
	e := echo.New()
	h := NewHandler(g, j, conn, hc, ms, adminsIDs)
	e.Use(middlewares.JWTAuth(j))

	api.RegisterHandlers(e, h)

	return e, g, j
}
