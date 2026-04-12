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
	users       *models.UserRepository
	posts       *models.PostRepository
	subscribers *models.SubscribersRepository
}

// NewModels creates a new Models instance.
func NewModels(users *models.UserRepository,
	posts *models.PostRepository,
	subscribers *models.SubscribersRepository,
) *Models {
	return &Models{
		users:       users,
		posts:       posts,
		subscribers: subscribers,
	}
}

// Users returns the models.UserRepository.
func (m *Models) Users() *models.UserRepository {
	return m.users
}

// Posts returns the models.PostRepository.
func (m *Models) Posts() *models.PostRepository {
	return m.posts
}

// Subscribers returns the models.SubscribersRepository.
func (m *Models) Subscribers() *models.SubscribersRepository {
	return m.subscribers
}

// Database is the database connection.
type Database struct {
	conn   *gorm.DB
	models *Models
}

// NewDatabase creates a new Database instance.
func NewDatabase(conn *gorm.DB, models *Models) *Database {
	return &Database{
		conn:   conn,
		models: models,
	}
}

// Models returns the database models.
func (d *Database) Models() *Models {
	return d.models
}

// GetConn returns the database connection.
func (d *Database) GetConn() *gorm.DB {
	return d.conn
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
