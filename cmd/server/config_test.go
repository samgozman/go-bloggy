package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromEnv(t *testing.T) {
	t.Setenv("GITHUB_CLIENT_ID", "test_id")
	t.Setenv("GITHUB_CLIENT_SECRET", "test_secret")
	t.Setenv("JWT_SECRET_KEY", "test_jwt")

	config := NewConfigFromEnv()

	assert.Equal(t, "test_id", config.GithubClientID)
	assert.Equal(t, "test_secret", config.GithubClientSecret)
	assert.Equal(t, "test_jwt", config.JWTSecretKey)
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
