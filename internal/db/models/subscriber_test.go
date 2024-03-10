package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscriptionDB(t *testing.T) {
	db, err := NewTestDB("file::memory:")
	assert.NoError(t, err)

	subscriptionDB := NewSubscribersDB(db)

	t.Run("Create", func(t *testing.T) {
		t.Run("create a new subscription", func(t *testing.T) {
			subscription := &Subscriber{
				Email: genEmail(),
			}

			err := subscriptionDB.Create(context.Background(), subscription)
			assert.NoError(t, err)
			assert.NotEmpty(t, subscription.ID)
			assert.NotZero(t, subscription.CreatedAt)
		})

		t.Run("return error if email is not unique", func(t *testing.T) {
			subscription := &Subscriber{
				Email: genEmail(),
			}

			err := subscriptionDB.Create(context.Background(), subscription)
			assert.NoError(t, err)

			err = subscriptionDB.Create(context.Background(), subscription)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrDuplicate)
		})

		t.Run("return error if email is empty", func(t *testing.T) {
			subscription := &Subscriber{
				Email: "",
			}

			err := subscriptionDB.Create(context.Background(), subscription)
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrValidationFailed)
			assert.ErrorIs(t, err, ErrSubscriptionEmailRequired)
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		subscription := &Subscriber{
			Email: genEmail(),
		}

		err := subscriptionDB.Create(context.Background(), subscription)
		assert.NoError(t, err)

		t.Run("should get the subscription", func(t *testing.T) {
			retrievedSubscription, err := subscriptionDB.GetByID(context.Background(), subscription.ID.String())
			assert.NoError(t, err)
			assert.Equal(t, subscription.ID, retrievedSubscription.ID)
		})

		t.Run("should return error if not found", func(t *testing.T) {
			_, err := subscriptionDB.GetByID(context.Background(), uuid.New().String())
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		subscription := &Subscriber{
			Email: genEmail(),
		}

		err := subscriptionDB.Create(context.Background(), subscription)
		assert.NoError(t, err)

		t.Run("should delete the subscription", func(t *testing.T) {
			err := subscriptionDB.Delete(context.Background(), subscription.ID.String())
			assert.NoError(t, err)

			_, err = subscriptionDB.GetByID(context.Background(), subscription.ID.String())
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrNotFound)
		})

		t.Run("should return error if not found", func(t *testing.T) {
			err := subscriptionDB.Delete(context.Background(), uuid.New().String())
			assert.Error(t, err)
			assert.ErrorIs(t, err, ErrNotFound)
		})
	})

	t.Run("GetEmails", func(t *testing.T) {
		t.Run("should return a list of emails", func(t *testing.T) {
			// Create a few subscriptions
			for i := 0; i < 3; i++ {
				subscription := &Subscriber{
					Email: genEmail(),
				}
				err := subscriptionDB.Create(context.Background(), subscription)
				assert.NoError(t, err)
			}

			emails, err := subscriptionDB.GetEmails(context.Background())
			assert.NoError(t, err)
			assert.NotEmpty(t, emails)
			// check that all emails are unique
			uniqueEmails := make(map[string]struct{})
			for _, email := range emails {
				uniqueEmails[email] = struct{}{}
			}
			assert.Len(t, uniqueEmails, len(emails))
		})
	})
}

func genEmail() string {
	return uuid.New().String() + "@example.com"
}
