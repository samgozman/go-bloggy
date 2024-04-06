package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// AuthMethod is the method of authentication used by the user.
type AuthMethod string

const (
	GitHubAuthMethod AuthMethod = "github"
)

// UserRepository is the database for the user data.
type UserRepository struct {
	conn *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

// User is the model for the user data.
type User struct {
	ID         int        `json:"id" gorm:"primaryKey;autoIncrement"`
	ExternalID string     `json:"external_id" gorm:"uniqueIndex"` // ExternalID is the ID of the user in the AuthMethod
	Login      string     `json:"login"`
	AuthMethod AuthMethod `json:"auth_method"` // AuthMethod is the method of authentication used by the user
	Posts      []Post     `json:"posts" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Login == "" {
		return ErrUserLoginRequired
	}
	if u.ExternalID == "" {
		return ErrUserExternalIDRequired
	}
	if u.AuthMethod == "" {
		return ErrUserAuthMethodRequired
	}
	return nil
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	err := u.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	// Because BeforeCreate can act as a hook for both Create and Update in Upsert operations
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
	return nil
}

// UserRepositoryInterface is the interface for the UserRepository.
type UserRepositoryInterface interface {
	Upsert(ctx context.Context, user *User) error
	GetByExternalID(ctx context.Context, externalID string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
}

// Upsert inserts or updates the User data.
func (db *UserRepository) Upsert(ctx context.Context, user *User) error {
	err := db.conn.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"login"}),
		}).
		Create(user).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCreateUser, mapGormError(err))
	}

	return nil
}

// GetByExternalID returns the User data by the User.ExternalID.
func (db *UserRepository) GetByExternalID(ctx context.Context, externalID string) (*User, error) {
	var user User
	err := db.conn.WithContext(ctx).Where("external_id = ?", externalID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w. external_id=%s", ErrFailedToGetUser, mapGormError(err), externalID)
	}

	return &user, nil
}

// GetByID returns the User data by the User.ID.
func (db *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := db.conn.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w. id=%v", ErrFailedToGetUser, mapGormError(err), id)
	}

	return &user, nil
}
