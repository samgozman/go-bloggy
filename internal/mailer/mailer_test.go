package mailer

import (
	"errors"
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/samgozman/go-bloggy/internal/mailer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	mockMailer "github.com/samgozman/go-bloggy/mocks/mailer"
)

func TestService_SendConfirmationEmail(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockClient := mockMailer.NewMockMailjetInterface(t)
		s := NewService("", "", &types.Options{})
		s.client = mockClient

		mockClient.On("SendMailV31", mock.Anything).Return(&mailjet.ResultsV31{}, nil)

		err := s.SendConfirmationEmail("test@example.com", "123")
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := mockMailer.NewMockMailjetInterface(t)
		s := NewService("", "", &types.Options{})
		s.client = mockClient

		mockClient.On("SendMailV31", mock.Anything).Return(&mailjet.ResultsV31{}, errors.New("error"))

		err := s.SendConfirmationEmail("test@example.com", "123")
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrSendConfirmationMail)
		mockClient.AssertExpectations(t)
	})
}

func TestService_SendPostEmail(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockClient := mockMailer.NewMockMailjetInterface(t)
		s := NewService("", "", &types.Options{})
		s.client = mockClient

		mockClient.On("SendMailV31", mock.Anything).Return(&mailjet.ResultsV31{}, nil)

		err := s.SendPostEmail(&types.PostEmailSend{
			To: []*types.Subscriber{
				{
					ID:    "123",
					Email: "some@example.com",
				},
			},
			Title:       "Test Title",
			Description: "Test Description",
			Slug:        "test-slug",
		})
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := mockMailer.NewMockMailjetInterface(t)
		s := NewService("", "", &types.Options{})
		s.client = mockClient

		mockClient.On("SendMailV31", mock.Anything).Return(&mailjet.ResultsV31{}, errors.New("error"))

		err := s.SendPostEmail(&types.PostEmailSend{
			To: []*types.Subscriber{
				{
					ID:    "123",
					Email: "some@example.com",
				},
			},
			Title:       "Test Title",
			Description: "Test Description",
			Slug:        "test-slug",
		})
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrSendPostMail)
		mockClient.AssertExpectations(t)
	})
}
