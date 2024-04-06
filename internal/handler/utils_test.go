package handler

import (
	"fmt"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// initDatabaseTest creates a new database connection & migrates the schema for testing.
// Using SQLite in a memory database only to simplify testing.
func initDatabaseTest() (*db.Database, error) {
	gormDB, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign key constraint enforcement in SQLite
	sqliteDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = sqliteDB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign key constraints: %w", err)
	}

	err = gormDB.AutoMigrate(&models.User{}, &models.Post{}, &models.Subscriber{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	m := db.NewModels(
		models.NewUserRepository(gormDB),
		models.NewPostRepository(gormDB),
		models.NewSubscribersRepository(gormDB),
	)

	return db.NewDatabase(gormDB, m), nil
}
