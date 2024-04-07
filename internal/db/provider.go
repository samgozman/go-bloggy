package db

import (
	"fmt"
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"gorm.io/gorm"
)

// ProvideDSN provides the DSN from the config.
func ProvideDSN(cfg *config.Config) config.DSN {
	return cfg.DSN
}

// ProvideConnection provides a new database connection.
func ProvideConnection(dsn config.DSN) (*gorm.DB, error) {
	conn, err := connectToPG(string(dsn))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToConnectDatabase, err)
	}

	// Migrate the schema
	// TODO: add external migrator, do not use AutoMigrate in production
	err = conn.AutoMigrate(&models.User{}, &models.Post{}, &models.Subscriber{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToMigrateDatabase, err)
	}

	return conn, nil
}

// ProvideModels provides the models.
func ProvideModels(conn *gorm.DB) *Models {
	return NewModels(
		models.NewUserRepository(conn),
		models.NewPostRepository(conn),
		models.NewSubscribersRepository(conn),
	)
}

// ProvideDatabase provides a new database connection.
func ProvideDatabase(conn *gorm.DB, models ModelsInterface) (*Database, error) {
	return NewDatabase(conn, models), nil
}

// ProviderSet is a wire provider set that provides the database connection.
var ProviderSet = wire.NewSet( //nolint:gochecknoglobals // required by Wire
	ProvideDSN,
	ProvideConnection,
	ProvideModels,
	ProvideDatabase,
	wire.Bind(new(DatabaseInterface), new(*Database)),
	wire.Bind(new(ModelsInterface), new(*Models)),
)
