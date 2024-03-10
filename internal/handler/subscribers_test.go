package handler

import (
	"context"
	"encoding/json"
	"github.com/oapi-codegen/testutil"
	"github.com/samgozman/go-bloggy/internal/api"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_PostSubscribers(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
	if errDB != nil {
		t.Fatal(errDB)
	}

	t.Run("Created", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		rb, _ := json.Marshal(api.CreateSubscriberRequest{
			Email:   "some@email.com",
			Captcha: "some-captcha",
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Post("/subscribers").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusCreated, res.Code())

		// Check that the subscription was created
		emails, err := conn.Models.Subscribers.GetConfirmedEmails(context.Background())
		assert.NoError(t, err)
		assert.Contains(t, emails, "some@email.com")
	})

	t.Run("BadRequest", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		rb, _ := json.Marshal(api.CreateSubscriberRequest{
			Email:   "invalid-email",
			Captcha: "some-captcha",
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Post("/subscribers").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())
	})
}

func Test_DeleteSubscribers(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
	if errDB != nil {
		t.Fatal(errDB)
	}

	t.Run("NoContent", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		sub := models.Subscriber{
			Email: "some@email.space",
		}

		// Create a subscription
		err := conn.Models.Subscribers.Create(context.Background(), &sub)
		assert.NoError(t, err)

		rb, _ := json.Marshal(api.UnsubscribeRequest{
			SubscriptionId: sub.ID.String(),
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Delete("/subscribers").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusNoContent, res.Code())

		// Check that the subscription was deleted
		_, err = conn.Models.Subscribers.GetByID(context.Background(), sub.ID.String())
		assert.Error(t, err)
		assert.ErrorIs(t, err, models.ErrNotFound)
	})

	t.Run("StatusBadRequest ", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		rb, _ := json.Marshal(api.UnsubscribeRequest{
			SubscriptionId: "f87c5cc0-ec7b-41eb-8d23-0abe0938efd2",
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Delete("/subscribers").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())

		var errRes api.RequestError
		err := res.UnmarshalBodyToObject(&errRes)
		assert.NoError(t, err)
		assert.Equal(t, errRes.Code, errGetSubscription)
		assert.Equal(t, errRes.Message, "Subscriber is not found or error getting subscription by ID")
	})
}

func Test_PostSubscribersConfirm(t *testing.T) {
	conn, errDB := db.InitDatabase("file::memory:")
	if errDB != nil {
		t.Fatal(errDB)
	}

	t.Run("OK - NoContent", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		sub := models.Subscriber{
			Email:       "some@email.space",
			IsConfirmed: false,
		}

		// Create a subscription
		err := conn.Models.Subscribers.Create(context.Background(), &sub)
		assert.NoError(t, err)

		rb, _ := json.Marshal(api.ConfirmSubscriberRequest{
			Token: sub.ID.String(),
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Post("/subscribers/confirm").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusOK, res.Code())

		// Check that the subscription was confirmed
		retrievedSubscription, err := conn.Models.Subscribers.GetByID(context.Background(), sub.ID.String())
		assert.NoError(t, err)
		assert.True(t, retrievedSubscription.IsConfirmed)
	})

	t.Run("StatusBadRequest - not found", func(t *testing.T) {
		e, _, _ := registerHandlers(conn, nil)

		rb, _ := json.Marshal(api.ConfirmSubscriberRequest{
			Token: "ce247e1d-a371-42fc-b36b-26b566c0096c",
		})

		res := testutil.NewRequest().
			WithHeader("Content-Type", "application/json").
			Post("/subscribers/confirm").
			WithBody(rb).
			GoWithHTTPHandler(t, e)

		assert.Equal(t, http.StatusBadRequest, res.Code())
	})
}
