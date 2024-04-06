package jwt

import (
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
)

// ProvideJWTSecretKey is a Wire provider function that returns JWT secret key from the config.
func ProvideJWTSecretKey(cfg *config.Config) config.JWTSecretKey {
	return cfg.JWTSecretKey
}

// ProvideService is a Wire provider function that creates a new JWT service.
func ProvideService(jwtSecretKey config.JWTSecretKey) *Service {
	return NewService(string(jwtSecretKey))
}

// ProviderSet is a wire provider set for JWT.
var ProviderSet = wire.NewSet(
	ProvideJWTSecretKey,
	ProvideService,
)
