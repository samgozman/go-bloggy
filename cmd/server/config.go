package main

import "os"

type Config struct {
	GithubClientID     string
	GithubClientSecret string
	JWTSecretKey       string
}

// NewConfigFromEnv creates a new Config.
func NewConfigFromEnv() *Config {
	return &Config{
		GithubClientID:     getEnvOrPanic("GITHUB_CLIENT_ID"),
		GithubClientSecret: getEnvOrPanic("GITHUB_CLIENT_SECRET"),
		JWTSecretKey:       getEnvOrPanic("JWT_SECRET_KEY"),
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
