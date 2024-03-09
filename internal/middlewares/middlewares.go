package middlewares

import (
	"github.com/labstack/echo/v4"
	"strings"
)

// JWTAuth is a middleware that checks for JWT token in the request and validates it.
// If the token is not present or invalid, it returns 401 Unauthorized.
// If the token is valid, it adds user ID to the request context.
//
// It skips the middleware for all GET requests.
func JWTAuth(jwtService jwtService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Skip for all GET requests
			if ctx.Request().Method == "GET" {
				return next(ctx)
			}

			// Skip for /login requests and /subscriptions requests
			if strings.HasPrefix(ctx.Request().URL.Path, "/login") ||
				strings.HasPrefix(ctx.Request().URL.Path, "/subscriptions") {
				return next(ctx)
			}

			// get token from header
			token := ctx.Request().Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			if token == "" {
				return ctx.JSON(401, ErrAuthHeaderRequired)
			}

			// parse token
			externalUserID, err := jwtService.ParseTokenString(token)
			if err != nil {
				return ctx.JSON(401, ErrInvalidToken)
			}

			// add user ID to context
			ctx.Set("externalUserID", externalUserID)

			return next(ctx)
		}
	}
}

type jwtService interface {
	ParseTokenString(tokenString string) (externalUserID string, err error)
}
