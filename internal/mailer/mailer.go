package mailer

import (
	"fmt"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type Service struct {
	client  Mailer
	options *Options
}

type Mailer interface {
	SendMailV31(data *mailjet.MessagesV31, options ...mailjet.RequestOptions) (*mailjet.ResultsV31, error)
}

type Options struct {
	FromEmail                    string
	FromName                     string
	ConfirmationTemplateID       int
	ConfirmationTemplateURLParam string
	PostTemplateID               int
	PostTemplateURLParam         string
	UnsubscribeURLParam          string
}

type PostEmailSend struct {
	To           []string
	SubscriberID string
	Title        string
	Description  string
	PostSlug     string
}

func NewService(publicKey, privateKey string, options *Options) *Service {
	return &Service{
		client:  mailjet.NewMailjetClient(publicKey, privateKey),
		options: options,
	}
}

func (s *Service) SendConfirmationEmail(to, confirmationID string) error {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: s.options.FromEmail,
				Name:  s.options.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: to,
				},
			},
			Subject:    "Please confirm your subscription",
			TemplateID: s.options.ConfirmationTemplateID,
			Variables: map[string]interface{}{
				"confirm_link":     fmt.Sprintf("%s?token=%s", s.options.ConfirmationTemplateURLParam, confirmationID),
				"unsubscribe_link": fmt.Sprintf("%s?token=%s", s.options.UnsubscribeURLParam, confirmationID),
			},
		},
	}

	_, err := s.client.SendMailV31(&mailjet.MessagesV31{Info: messagesInfo})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSendConfirmationMail, err)
	}

	return nil
}

func (s *Service) SendPostEmail(pe *PostEmailSend) error {
	r := make([]mailjet.RecipientV31, len(pe.To))
	for i, email := range pe.To {
		r[i] = mailjet.RecipientV31{
			Email: email,
		}
	}
	recipients := mailjet.RecipientsV31(r)

	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: s.options.FromEmail,
				Name:  s.options.FromName,
			},
			To:         &recipients,
			Subject:    fmt.Sprintf("New post: %s", pe.Title),
			TemplateID: s.options.PostTemplateID,
			Variables: map[string]interface{}{
				"email_title":      pe.Title,
				"email_paragraph":  pe.Description,
				"post_link":        fmt.Sprintf("%s/%s", s.options.PostTemplateURLParam, pe.PostSlug),
				"unsubscribe_link": fmt.Sprintf("%s?token=%s", s.options.UnsubscribeURLParam, pe.SubscriberID),
			},
		},
	}

	_, err := s.client.SendMailV31(&mailjet.MessagesV31{Info: messagesInfo})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSendPostMail, err)
	}

	return nil
}
