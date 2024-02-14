package jwt

import (
	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_CreateTokenString(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		signKey := "testKey1"
		userID := "testUser1"
		expiresAt := time.Now().Add(time.Hour)

		token, err := CreateTokenString(signKey, userID, expiresAt)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid date", func(t *testing.T) {
		signKey := "testKey2"
		userID := "testUser2"
		expiresAt := time.Now().Add(-time.Hour)

		token, err := CreateTokenString(signKey, userID, expiresAt)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func Test_ParseTokenString(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		signKey := "testKey3"
		userID := "testUser3"
		expiresAt := time.Now().Add(time.Hour)

		token, err := CreateTokenString(signKey, userID, expiresAt)
		assert.NoError(t, err)

		parsedUserID, err := ParseTokenString(signKey, token)
		assert.NoError(t, err)
		assert.Equal(t, userID, parsedUserID)
	})

	t.Run("invalid signKey", func(t *testing.T) {
		token, err := CreateTokenString("testKey", "testUser", time.Now().Add(time.Hour))
		assert.NoError(t, err)

		parsedUserID, err := ParseTokenString("invalidKey", token)
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)
	})

	t.Run("invalid token", func(t *testing.T) {
		parsedUserID, err := ParseTokenString("testKey", "invalidToken")
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)

		parsedUserID, err = ParseTokenString("testKey", "")
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)
	})

	t.Run("expired token", func(t *testing.T) {
		signKey := "testKey"
		userID := "testUser"
		keyByte := []byte(signKey)

		claims := Claims{
			userID,
			jwtgo.RegisteredClaims{
				ExpiresAt: jwtgo.NewNumericDate(time.Now().Add(-time.Hour)),
				IssuedAt:  jwtgo.NewNumericDate(time.Now()),
				NotBefore: jwtgo.NewNumericDate(time.Now()),
				Issuer:    "go-bloggy",
			},
		}

		token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
		ss, err := token.SignedString(keyByte)
		assert.NoError(t, err)

		parsedUserID, err := ParseTokenString(signKey, ss)
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)
	})
}
