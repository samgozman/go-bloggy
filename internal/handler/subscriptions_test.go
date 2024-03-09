package handler

import (
	"context"
	"encoding/json"
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/samgozman/go-bloggy/pkg/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_PostSubscriptions(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
	if errDB != nil {
		t.Fatal(errDB)
	}

	t.Run("Created", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		rb, _ := json.Marshal(server.SubscriptionRequest{
			Email:   "some@email.com",
			Captcha: "some-captcha", // TODO: Add captcha support
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Post("/subscriptions").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusCreated, res.Code())

		// Check that the subscription was created
		emails, err := conn.Models.Subscriptions.GetEmails(context.Background())
		assert.NoError(t, err)
		assert.Contains(t, emails, "some@email.com")
	})

	t.Run("BadRequest", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		rb, _ := json.Marshal(server.SubscriptionRequest{
			Email:   "invalid-email",
			Captcha: "some-captcha", // TODO: Add captcha support
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Post("/subscriptions").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())
	})
}

func Test_DeleteSubscriptions(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
	if errDB != nil {
		t.Fatal(errDB)
	}

	t.Run("NoContent", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		sub := models.Subscription{
			Email: "some@email.space",
		}

		// Create a subscription
		err := conn.Models.Subscriptions.Create(context.Background(), &sub)
		assert.NoError(t, err)

		rb, _ := json.Marshal(server.UnsubscribeRequest{
			SubscriptionId: sub.ID.String(),
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Delete("/subscriptions").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusNoContent, res.Code())

		// Check that the subscription was deleted
		_, err = conn.Models.Subscriptions.GetByID(context.Background(), sub.ID.String())
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrNotFound)
	})

	t.Run("StatusBadRequest ", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		rb, _ := json.Marshal(server.UnsubscribeRequest{
			SubscriptionId: "f87c5cc0-ec7b-41eb-8d23-0abe0938efd2",
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Delete("/subscriptions").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var errRes server.RequestError
		err := res.UnmarshalBodyToObject(&errRes)
		assert.NoError(t, err)
		assert.Equal(t, errRes.Code, errGetSubscription)
		assert.Equal(t, errRes.Message, "Subscription is not found or error getting subscription by ID")
	})
}
