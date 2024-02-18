package main

import "os"

type Config struct {
	GithubClientID     string // GithubClientID is the client ID for GitHub OAuth.
	GithubClientSecret string // GithubClientSecret is the secret key for GitHub OAuth.
	JWTSecretKey       string // JWTSecretKey is the secret key for JWT token creation.
	Port               string // Port for server to listen on.
	DSN                string // DSN - Database Source Name. For sqlite, it's the file path.
}

// NewConfigFromEnv creates a new Config.
func NewConfigFromEnv() *Config {
	return &Config{
		GithubClientID:     getEnvOrPanic("GITHUB_CLIENT_ID"),
		GithubClientSecret: getEnvOrPanic("GITHUB_CLIENT_SECRET"),
		JWTSecretKey:       getEnvOrPanic("JWT_SECRET_KEY"),
		Port:               getEnvOrPanic("PORT"),
		DSN:                getEnvOrPanic("DSN"),
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
