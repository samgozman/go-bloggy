package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	jwtMock "github.com/samgozman/go-bloggy/mocks/jwt"
)

func Test_JWTAuth(t *testing.T) {
	mockJwtService := jwtMock.NewMockServiceInterface(t)
	middleware := JWTAuth(mockJwtService)

	t.Run("valid token", func(t *testing.T) {
		mockJwtService.On("ParseTokenString", "validToken").Return("SuperUserID", nil)

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)
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
		mockJwtService.On("ParseTokenString", "invalidToken").Return("", errors.New("unauthorized"))

		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)
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
		req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)

		_ = middleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})(ctx)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, fmt.Sprintf("\"%s\"\n", ErrAuthHeaderRequired), rec.Body.String())
	})

	t.Run("skip allowed paths", func(t *testing.T) {
		t.Run("GET request", func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := echo.New().NewContext(req, rec)

			_ = middleware(func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})(ctx)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "test", rec.Body.String())
		})

		testCases := []string{"/login", "/subscribers"}
		for _, path := range testCases {
			t.Run(path, func(t *testing.T) {
				req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, path, nil)
				rec := httptest.NewRecorder()
				ctx := echo.New().NewContext(req, rec)

				_ = middleware(func(c echo.Context) error {
					return c.String(http.StatusOK, "test")
				})(ctx)

				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, "test", rec.Body.String())
			})
		}
	})
}
