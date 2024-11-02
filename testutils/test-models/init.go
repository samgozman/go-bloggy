package testmodels

import (
	"fmt"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	testdb "github.com/samgozman/go-bloggy/testutils/test-db"
)

// InitDatabaseWithModelsTest creates a new database connection & migrates the schema for testing.
func InitDatabaseWithModelsTest() (*db.Database, error) {
	gormDB, err := testdb.InitDatabaseTest()
	if err != nil {
		return nil, fmt.Errorf("error init test db: %w", err)
	}

	err = gormDB.AutoMigrate(&models.User{}, &models.Post{}, &models.Subscriber{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	m := db.NewModels(
		models.NewUserRepository(gormDB),
		models.NewPostRepository(gormDB),
		models.NewSubscribersRepository(gormDB),
	)

	return db.NewDatabase(gormDB, m), nil
}
