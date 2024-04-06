package db

import (
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
)

// ProvideDSN provides the DSN from the config.
func ProvideDSN(cfg *config.Config) config.DSN {
	return cfg.DSN
}

// ProvideDatabase provides a new database connection.
func ProvideDatabase(dsn config.DSN) (*Database, error) {
	return InitDatabase(string(dsn))
}

// ProviderSet is a wire provider set that provides the database connection.
var ProviderSet = wire.NewSet(
	ProvideDSN,
	ProvideDatabase,
)
