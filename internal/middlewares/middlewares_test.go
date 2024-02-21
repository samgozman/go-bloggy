package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockJwtService struct {
	mock.Mock
}

func (m *mockJwtService) ParseTokenString(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0), args.Error(1)
}

func Test_JWTAuth(t *testing.T) {
	mockJwtService := new(mockJwtService)
	middleware := JWTAuth(mockJwtService)

	t.Run("valid token", func(t *testing.T) {
		mockJwtService.On("ParseTokenString", "validToken").Return("SuperUserID", nil)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer validToken")
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test", rec.Body.String())
		assert.Equal(t, "SuperUserID", ctx.Get("userID"))
	})

	t.Run("invalid token", func(t *testing.T) {
		mockJwtService.On("ParseTokenString", "invalidToken").Return("", echo.ErrUnauthorized)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		req.Header.Set("Authorization", "Bearer invalidToken")
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
