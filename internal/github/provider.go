package github

import (
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
)

// Config is a struct that holds the configuration for the GitHub service.
type Config struct {
	ClientID     string
	ClientSecret string
}

// ProvideConfig is a Wire provider function that creates a Config.
func ProvideConfig(cfg *config.Config) *Config {
	return &Config{
		ClientID:     cfg.GithubClientID,
		ClientSecret: cfg.GithubClientSecret,
	}
}

// ProvideService is a Wire provider function that creates a Service.
func ProvideService(cfg *Config) *Service {
	return NewService(cfg.ClientID, cfg.ClientSecret)
}

// ProviderSet is a Wire provider set that includes all the providers from the github package.
var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideService,
	wire.Bind(new(ServiceInterface), new(*Service)),
)
