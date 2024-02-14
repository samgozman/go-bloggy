package jwt

import (
	"fmt"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"time"
)

// Claims is a custom JWT claims type. It embeds the standard JWT claims and adds a custom field.
type Claims struct {
	UserID string `json:"userId"`
	jwtgo.RegisteredClaims
}

// CreateTokenString creates a JWT token string with the given sign key and expiration time.
func CreateTokenString(signKey, userID string, expiresAt time.Time) (string, error) {
	if expiresAt.Before(time.Now()) {
		return "", fmt.Errorf("expiresAt must be in the future")
	}

	keyByte := []byte(signKey)

	claims := Claims{
		userID,
		jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(expiresAt),
			IssuedAt:  jwtgo.NewNumericDate(time.Now()),
			NotBefore: jwtgo.NewNumericDate(time.Now()),
			Issuer:    "go-bloggy",
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	ss, err := token.SignedString(keyByte)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return ss, nil
}

// ParseTokenString parses a JWT token string and returns User ID.
func ParseTokenString(signKey, tokenString string) (string, error) {
	keyByte := []byte(signKey)

	token, err := jwtgo.ParseWithClaims(tokenString, &Claims{}, func(token *jwtgo.Token) (interface{}, error) {
		return keyByte, nil
	})
	if err != nil {
		return "", fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
