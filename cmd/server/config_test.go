package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromEnv(t *testing.T) {
	t.Setenv("GITHUB_CLIENT_ID", "test_id")
	t.Setenv("GITHUB_CLIENT_SECRET", "test_secret")
	t.Setenv("JWT_SECRET_KEY", "test_jwt")
	t.Setenv("PORT", "3000")
	t.Setenv("DSN", "test_dsn")
	t.Setenv("ADMINS_EXTERNAL_IDS", "test_admin1,test_admin2")
	t.Setenv("ENVIRONMENT", "test_env")
	t.Setenv("HCAPTCHA_SECRET", "test_h")

	config := NewConfigFromEnv()

	assert.Equal(t, "test_id", config.GithubClientID)
	assert.Equal(t, "test_secret", config.GithubClientSecret)
	assert.Equal(t, "test_jwt", config.JWTSecretKey)
	assert.Equal(t, "3000", config.Port)
	assert.Equal(t, "test_dsn", config.DSN)
	assert.Equal(t, []string{"test_admin1", "test_admin2"}, config.AdminsExternalIDs)
	assert.Equal(t, "test_env", config.Environment)
	assert.Equal(t, "test_h", config.HCaptchaSecret)
}

func TestGetEnvOrPanic(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		t.Setenv("TEST_ENV", "test_value")

		value := getEnvOrPanic("TEST_ENV")

		assert.Equal(t, "test_value", value)
	})

	t.Run("Panic", func(t *testing.T) {
		assert.Panics(t, func() { getEnvOrPanic("NON_EXISTING_ENV") })
	})
}
