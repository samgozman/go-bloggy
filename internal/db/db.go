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
	users       models.UserRepositoryInterface
	posts       models.PostRepositoryInterface
	subscribers models.SubscriberRepositoryInterface
}

// NewModels creates a new Models instance.
func NewModels(users models.UserRepositoryInterface,
	posts models.PostRepositoryInterface,
	subscribers models.SubscriberRepositoryInterface,
) *Models {
	return &Models{
		users:       users,
		posts:       posts,
		subscribers: subscribers,
	}
}

// Users returns the models.UserRepository.
func (m *Models) Users() models.UserRepositoryInterface {
	return m.users
}

// Posts returns the models.PostRepository.
func (m *Models) Posts() models.PostRepositoryInterface {
	return m.posts
}

// Subscribers returns the models.SubscriberRepository.
func (m *Models) Subscribers() models.SubscriberRepositoryInterface {
	return m.subscribers
}

type ModelsInterface interface {
	Users() models.UserRepositoryInterface
	Posts() models.PostRepositoryInterface
	Subscribers() models.SubscriberRepositoryInterface
}

// Database is the database connection.
type Database struct {
	conn   *gorm.DB
	models ModelsInterface
}

// NewDatabase creates a new Database instance.
func NewDatabase(conn *gorm.DB, models ModelsInterface) *Database {
	return &Database{
		conn:   conn,
		models: models,
	}
}

// Models returns the database models.
func (d *Database) Models() ModelsInterface {
	return d.models
}

// GetConn returns the database connection.
func (d *Database) GetConn() *gorm.DB {
	return d.conn
}

// DatabaseInterface is the interface for the database.
type DatabaseInterface interface {
	GetConn() *gorm.DB
	Models() ModelsInterface
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
