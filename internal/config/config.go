package config

import (
	"os"
	"strconv"
	"strings"
)

type DSN string
type JWTSecretKey string

type Config struct {
	GithubClientID     string       // GithubClientID is the client ID for GitHub OAuth.
	GithubClientSecret string       // GithubClientSecret is the secret key for GitHub OAuth.
	JWTSecretKey       JWTSecretKey // JWTSecretKey is the secret key for JWT token creation.
	Port               string       // Port for server to listen on.
	DSN                DSN          // DSN - Database Source Name for Postgres.
	AdminsExternalIDs  []string     // AdminsExternalIDs list of admins allowed to auth, separated by comma.
	HCaptchaSecret     string       // HCaptchaSecret is the secret key for HCaptcha verification.
	MailerJet          MailerConfig
}

type MailerConfig struct {
	PublicKey                    string // PublicKey is the public key for Mailjet API.
	PrivateKey                   string // PrivateKey is the private key for Mailjet API.
	FromEmail                    string // FromEmail is the email address to send emails from.
	FromName                     string // FromName is the name to send emails from.
	ConfirmationTemplateID       int    // ConfirmationTemplateID is the ID of the Mailjet template for confirmation emails.
	ConfirmationTemplateURLParam string // ConfirmationTemplateURLParam e.g. "https://example.com/confirm?token="
	PostTemplateID               int    // PostTemplateID is the ID of the Mailjet template for post-emails.
	PostTemplateURLParam         string // PostTemplateURLParam e.g. "https://example.com/post/" to append the posts slug.
	UnsubscribeURLParam          string // UnsubscribeURLParam e.g. "https://example.com/unsubscribe?id="
}

// NewConfigFromEnv creates a new Config.
func NewConfigFromEnv() *Config {
	admins := getEnvOrPanic("ADMINS_EXTERNAL_IDS")
	var adminsList []string
	if admins != "" {
		// admins is a comma-separated list of external IDs
		adminsList = strings.Split(admins, ",")
	}

	confirmationTemplateID, _ := strconv.Atoi(getEnvOrPanic("MAILJET_CONFIRMATION_TEMPLATE_ID"))
	postTemplateID, _ := strconv.Atoi(getEnvOrPanic("MAILJET_POST_TEMPLATE_ID"))

	return &Config{
		GithubClientID:     getEnvOrPanic("GITHUB_CLIENT_ID"),
		GithubClientSecret: getEnvOrPanic("GITHUB_CLIENT_SECRET"),
		JWTSecretKey:       JWTSecretKey(getEnvOrPanic("JWT_SECRET_KEY")),
		Port:               getEnvOrPanic("PORT"),
		DSN:                DSN(getEnvOrPanic("DSN")),
		AdminsExternalIDs:  adminsList,
		HCaptchaSecret:     getEnvOrPanic("HCAPTCHA_SECRET"),
		MailerJet: MailerConfig{
			PublicKey:                    getEnvOrPanic("MAILJET_PUBLIC_KEY"),
			PrivateKey:                   getEnvOrPanic("MAILJET_PRIVATE_KEY"),
			FromEmail:                    getEnvOrPanic("MAILJET_MAIL_FROM"),
			FromName:                     getEnvOrPanic("MAILJET_MAIL_FROM_NAME"),
			ConfirmationTemplateID:       confirmationTemplateID,
			ConfirmationTemplateURLParam: getEnvOrPanic("MAILJET_CONFIRMATION_TEMPLATE_URL_PARAM"),
			PostTemplateID:               postTemplateID,
			PostTemplateURLParam:         getEnvOrPanic("MAILJET_POST_TEMPLATE_URL_PARAM"),
			UnsubscribeURLParam:          getEnvOrPanic("MAILJET_UNSUBSCRIBE_URL_PARAM"),
		},
	}
}

// getEnvOrPanic returns the value of the environment variable or panics if it is not set.
func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("missing env variable " + key)
	}
	return value
}
