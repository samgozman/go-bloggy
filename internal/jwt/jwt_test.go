package jwt

import (
	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_CreateTokenString(t *testing.T) {
	service := NewService("testKey")

	t.Run("OK", func(t *testing.T) {
		userID := "testUser1"
		expiresAt := time.Now().Add(time.Hour)

		token, err := service.CreateTokenString(userID, expiresAt)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Invalid date", func(t *testing.T) {
		userID := "testUser2"
		expiresAt := time.Now().Add(-time.Hour)

		token, err := service.CreateTokenString(userID, expiresAt)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func Test_ParseTokenString(t *testing.T) {
	service := NewService("testKey")

	t.Run("OK", func(t *testing.T) {
		userID := "testUser3"
		expiresAt := time.Now().Add(time.Hour)

		token, err := service.CreateTokenString(userID, expiresAt)
		assert.NoError(t, err)

		parsedUserID, err := service.ParseTokenString(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, parsedUserID)
	})

	t.Run("invalid signKey", func(t *testing.T) {
		token, err := service.CreateTokenString("testUser", time.Now().Add(time.Hour))
		assert.NoError(t, err)

		serviceInvalidKey := NewService("invalidKey")

		parsedUserID, err := serviceInvalidKey.ParseTokenString(token)
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)
	})

	t.Run("invalid token", func(t *testing.T) {
		parsedUserID, err := service.ParseTokenString("invalidToken")
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)

		parsedUserID, err = service.ParseTokenString("")
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)
	})

	t.Run("expired token", func(t *testing.T) {
		userID := "testUser"
		keyByte := []byte(service.signKey)

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

		parsedUserID, err := service.ParseTokenString(ss)
		assert.Error(t, err)
		assert.Empty(t, parsedUserID)
	})
}

func Test_NewService(t *testing.T) {
	service := NewService("testKey")
	assert.Equal(t, "testKey", service.signKey)
	assert.NotNil(t, service.CreateTokenString)
	assert.NotNil(t, service.ParseTokenString)
}
