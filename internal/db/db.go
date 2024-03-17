package db

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

// Models is a collection of all models in the database.
type Models struct {
	Users       *models.UserDB
	Posts       *models.PostDB
	Subscribers *models.SubscribersDB
}

// Database is the database connection.
type Database struct {
	Conn   *gorm.DB
	Models *Models
}

// connectToPG connects to the Postgres database and retries until it is ready.
func connectToPG(dsn string) (*gorm.DB, error) {
	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval = 10 * time.Second
	bf.MaxInterval = 25 * time.Second
	bf.MaxElapsedTime = 90 * time.Second

	db, err := backoff.RetryWithData[*gorm.DB](func() (*gorm.DB, error) {
		conn, err := gorm.Open(postgres.New(postgres.Config{
			DSN: dsn,
		}))
		if err != nil {
			slog.Info("[connectToPG] Postgres not yet ready...")
			return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
		}
		slog.Info("[connectToPG] Connected to Postgres!")
		return conn, nil
	}, bf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	return db, nil
}

// InitDatabase creates a new database connection & migrates the schema.
func InitDatabase(dsn string) (*Database, error) {
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
