package models

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SubscriptionDB struct {
	conn *gorm.DB
}

func NewSubscriptionDB(conn *gorm.DB) *SubscriptionDB {
	return &SubscriptionDB{
		conn: conn,
	}
}

type Subscription struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	CreatedAt time.Time
}

func (s *Subscription) Validate() error {
	if s.Email == "" {
		return ErrSubscriptionEmailRequired
	}
	return nil
}

func (s *Subscription) BeforeCreate(_ *gorm.DB) error {
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

func (db *SubscriptionDB) Create(ctx context.Context, s *Subscription) error {
	err := db.conn.WithContext(ctx).Create(s).Error
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateSubscription, mapGormError(err))
	}

	return nil
}

func (db *SubscriptionDB) GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error) {
	var s Subscription
	err := db.conn.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetSubscription, mapGormError(err))
	}

	return &s, nil
}

func (db *SubscriptionDB) GetEmails(ctx context.Context) ([]string, error) {
	var emails []string
	err := db.conn.WithContext(ctx).Model(&Subscription{}).Pluck("email", &emails).Error
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetSubscriptionEmails, mapGormError(err))
	}

	return emails, nil
}

func (db *SubscriptionDB) Delete(ctx context.Context, id uuid.UUID) error {
	res := db.conn.WithContext(ctx).Where("id = ?", id).Delete(&Subscription{})
	if res.Error != nil {
		return fmt.Errorf("%w: %w", ErrDeleteSubscription, mapGormError(res.Error))
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: %w", ErrNotFound, gorm.ErrRecordNotFound)
	}

	return nil
}
