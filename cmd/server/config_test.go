package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	t.Setenv("MAILJET_PUBLIC_KEY", "test_public_key")
	t.Setenv("MAILJET_PRIVATE_KEY", "test_private_key")
	t.Setenv("MAILJET_MAIL_FROM", "test_mail_from")
	t.Setenv("MAILJET_MAIL_FROM_NAME", "test_mail_from_name")
	t.Setenv("MAILJET_CONFIRMATION_TEMPLATE_ID", "1")
	t.Setenv("MAILJET_CONFIRMATION_TEMPLATE_URL_PARAM", "test_confirmation_template_url_param")
	t.Setenv("MAILJET_POST_TEMPLATE_ID", "2")
	t.Setenv("MAILJET_POST_TEMPLATE_URL_PARAM", "test_post_template_url_param")
	t.Setenv("MAILJET_UNSUBSCRIBE_URL_PARAM", "test_unsubscribe_url_param")

	config := NewConfigFromEnv()

	assert.Equal(t, "test_id", config.GithubClientID)
	assert.Equal(t, "test_secret", config.GithubClientSecret)
	assert.Equal(t, "test_jwt", config.JWTSecretKey)
	assert.Equal(t, "3000", config.Port)
	assert.Equal(t, "test_dsn", config.DSN)
	assert.Equal(t, []string{"test_admin1", "test_admin2"}, config.AdminsExternalIDs)
	assert.Equal(t, "test_env", config.Environment)
	assert.Equal(t, "test_h", config.HCaptchaSecret)
	assert.Equal(t, "test_public_key", config.MailerJet.PublicKey)
	assert.Equal(t, "test_private_key", config.MailerJet.PrivateKey)
	assert.Equal(t, "test_mail_from", config.MailerJet.FromEmail)
	assert.Equal(t, "test_mail_from_name", config.MailerJet.FromName)
	assert.Equal(t, 1, config.MailerJet.ConfirmationTemplateID)
	assert.Equal(t, "test_confirmation_template_url_param", config.MailerJet.ConfirmationTemplateURLParam)
	assert.Equal(t, 2, config.MailerJet.PostTemplateID)
	assert.Equal(t, "test_post_template_url_param", config.MailerJet.PostTemplateURLParam)
	assert.Equal(t, "test_unsubscribe_url_param", config.MailerJet.UnsubscribeURLParam)
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
