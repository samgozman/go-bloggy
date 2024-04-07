package middlewares

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtMock "github.com/samgozman/go-bloggy/mocks/jwt"
)

func Test_JWTAuth(t *testing.T) {
	mockJwtService := jwtMock.NewMockServiceInterface(t)
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
		assert.Equal(t, "SuperUserID", ctx.Get("externalUserID"))
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
		assert.Equal(t, fmt.Sprintf("\"%s\"\n", ErrInvalidToken), rec.Body.String())
	})

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, fmt.Sprintf("\"%s\"\n", ErrAuthHeaderRequired), rec.Body.String())
	})
}
