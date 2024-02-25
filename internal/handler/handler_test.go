package handler

import (
	"context"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/middlewares"
	"github.com/samgozman/go-bloggy/pkg/server"
	"github.com/stretchr/testify/mock"
	"time"

	"github.com/labstack/echo/v4"
)

type MockGithubService struct {
	mock.Mock
}

func (m *MockGithubService) ExchangeCodeForToken(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)          //nolint:typecheck
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

func (m *MockGithubService) GetUserInfo(ctx context.Context, token string) (*github.UserInfo, error) {
	args := m.Called(ctx, token)                         //nolint:typecheck
	return args.Get(0).(*github.UserInfo), args.Error(1) //nolint:wrapcheck
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) CreateTokenString(userID string, expiresAt time.Time) (string, error) {
	args := m.Called(userID, expiresAt)  //nolint:typecheck
	return args.String(0), args.Error(1) //nolint:wrapcheck
}

func (m *MockJWTService) ParseTokenString(token string) (string, error) {
	args := m.Called(token) //nolint:typecheck
	return args.String(0), args.Error(1)
}

// registerHandlers creates a new echo instance and registers the handlers for testing.
func registerHandlers(conn *db.Database, adminsIDs []string) (s *echo.Echo, githubService *MockGithubService, jwtService *MockJWTService) {
	// Create mocks
	g := new(MockGithubService)
	j := new(MockJWTService)

	// Create echo instance
	e := echo.New()
	h := NewHandler(g, j, conn, adminsIDs)
	e.Use(middlewares.JWTAuth(j))

	server.RegisterHandlers(e, h)

	return e, g, j
}
