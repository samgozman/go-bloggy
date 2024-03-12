package mailer

import (
	"errors"
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockMailjetClient struct {
	mock.Mock
}

func (m *MockMailjetClient) SendMailV31(data *mailjet.MessagesV31, _options ...mailjet.RequestOptions) (*mailjet.ResultsV31, error) {
	args := m.Called(data)
	return args.Get(0).(*mailjet.ResultsV31), args.Error(1) //nolint:wrapcheck
}

func TestService_SendConfirmationEmail(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		mockClient := new(MockMailjetClient)
		s := NewService("", "", &Options{})
		s.client = mockClient

		mockClient.On("SendMailV31", mock.Anything).Return(&mailjet.ResultsV31{}, nil)

		err := s.SendConfirmationEmail("test@example.com", "123")
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := new(MockMailjetClient)
		s := NewService("", "", &Options{})
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
		mockClient := new(MockMailjetClient)
		s := NewService("", "", &Options{})
		s.client = mockClient

		mockClient.On("SendMailV31", mock.Anything).Return(&mailjet.ResultsV31{}, nil)

		err := s.SendPostEmail(&PostEmailSend{
			To: []Subscriber{
				{
					ID:    "123",
					Email: "some@example.com",
				},
			},
			Title:       "Test Title",
			Description: "Test Description",
			PostSlug:    "test-slug",
		})
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := new(MockMailjetClient)
		s := NewService("", "", &Options{})
		s.client = mockClient

		mockClient.On("SendMailV31", mock.Anything).Return(&mailjet.ResultsV31{}, errors.New("error"))

		err := s.SendPostEmail(&PostEmailSend{
			To: []Subscriber{
				{
					ID:    "123",
					Email: "some@example.com",
				},
			},
			Title:       "Test Title",
			Description: "Test Description",
			PostSlug:    "test-slug",
		})
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrSendPostMail)
		mockClient.AssertExpectations(t)
	})
}
