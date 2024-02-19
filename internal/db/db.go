package db

import (
	"fmt"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Models is a collection of all models in the database.
type Models struct {
	Users *models.UserDB
}

// Database is the database connection.
type Database struct {
	conn   *gorm.DB
	Models *Models
}

// InitDatabase creates a new database connection & migrates the schema.
func InitDatabase(dsn string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToConnectDatabase, err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToMigrateDatabase, err)
	}

	return &Database{
		conn: db,
		Models: &Models{
			Users: models.NewUserDB(db),
		},
	}, nil
}