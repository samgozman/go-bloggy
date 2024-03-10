package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SubscribersDB struct {
	conn *gorm.DB
}

func NewSubscribersDB(conn *gorm.DB) *SubscribersDB {
	return &SubscribersDB{
		conn: conn,
	}
}

type Subscriber struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Email       string    `json:"email" gorm:"uniqueIndex"`
	IsConfirmed bool      `json:"is_confirmed"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Subscriber) Validate() error {
	if s.Email == "" {
		return ErrSubscriptionEmailRequired
	}
	return nil
}

func (s *Subscriber) BeforeCreate(_ *gorm.DB) error {
	err := s.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}

	s.CreatedAt = time.Now()
	return nil
}

func (db *SubscribersDB) Create(ctx context.Context, s *Subscriber) error {
	err := db.conn.WithContext(ctx).Create(s).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateSubscription, mapGormError(err))
	}

	return nil
}

func (db *SubscribersDB) GetByID(ctx context.Context, id string) (*Subscriber, error) {
	var s Subscriber
	err := db.conn.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetSubscription, mapGormError(err))
	}

	return &s, nil
}

func (db *SubscribersDB) GetConfirmedEmails(ctx context.Context) ([]string, error) {
	var emails []string
	err := db.conn.WithContext(ctx).
		Model(&Subscriber{}).
		Where("is_confirmed = ?", true).
		Pluck("email", &emails).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetSubscriptionEmails, mapGormError(err))
	}

	return emails, nil
}

func (db *SubscribersDB) Delete(ctx context.Context, id string) error {
	res := db.conn.WithContext(ctx).Where("id = ?", id).Delete(&Subscriber{})
	if res.Error != nil {
		return fmt.Errorf("%w: %w", ErrDeleteSubscription, mapGormError(res.Error))
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: %w", ErrNotFound, gorm.ErrRecordNotFound)
	}

	return nil
}
