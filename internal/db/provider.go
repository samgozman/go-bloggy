package db

import (
	"fmt"
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db/models"
)

// ProvideDSN provides the DSN from the config.
func ProvideDSN(cfg *config.Config) string {
	return cfg.DSN
}

// ProvideDatabase provides a new database connection.
func ProvideDatabase(dsn string) (*Database, error) {
	db, err := connectToPG(dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToConnectDatabase, err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Subscriber{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToMigrateDatabase, err)
	}

	return &Database{
		Conn: db,
		Models: &Models{
			Users:       models.NewUserDB(db),
			Posts:       models.NewPostDB(db),
			Subscribers: models.NewSubscribersDB(db),
		},
	}, nil
}

// ProviderSet is a wire provider set that provides the database connection.
var ProviderSet = wire.NewSet(
	ProvideDSN,
	ProvideDatabase,
)
