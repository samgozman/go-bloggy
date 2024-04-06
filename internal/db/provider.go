package db

import (
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
)

// ProvideDSN provides the DSN from the config.
func ProvideDSN(cfg *config.Config) string {
	return cfg.DSN
}

// ProvideDatabase provides a new database connection.
func ProvideDatabase(dsn string) (*Database, error) {
	return InitDatabase(dsn)
}

// ProviderSet is a wire provider set that provides the database connection.
var ProviderSet = wire.NewSet(
	ProvideDSN,
	ProvideDatabase,
)
