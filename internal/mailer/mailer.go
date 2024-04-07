package mailer

import (
	"fmt"
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/samgozman/go-bloggy/internal/mailer/types"
)

type Service struct {
	client  types.MailjetInterface
	options *types.Options
}

func NewService(publicKey, privateKey string, options *types.Options) *Service {
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
			Subject:          "Please confirm your subscription",
			TemplateID:       s.options.ConfirmationTemplateID,
			TemplateLanguage: true,
			Variables: map[string]interface{}{
				"confirm_link":     s.options.ConfirmationTemplateURLParam + confirmationID,
				"unsubscribe_link": s.options.UnsubscribeURLParam + confirmationID,
			},
		},
	}

	_, err := s.client.SendMailV31(&mailjet.MessagesV31{Info: messagesInfo})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSendConfirmationMail, err)
	}

	return nil
}

func (s *Service) SendPostEmail(pe *types.PostEmailSend) error {
	messageFrom := mailjet.RecipientV31{
		Email: s.options.FromEmail,
		Name:  s.options.FromName,
	}
	subject := fmt.Sprintf("New post: %s", pe.Title)

	messagesInfo := make([]mailjet.InfoMessagesV31, len(pe.To))
	for i, sub := range pe.To {
		messagesInfo[i] = mailjet.InfoMessagesV31{
			From: &messageFrom,
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: sub.Email,
				},
			},
			Subject:          subject,
			TemplateID:       s.options.PostTemplateID,
			TemplateLanguage: true,
			Variables: map[string]interface{}{
				"email_title":      pe.Title,
				"email_paragraph":  pe.Description,
				"post_link":        s.options.PostTemplateURLParam + pe.Slug,
				"unsubscribe_link": s.options.UnsubscribeURLParam + sub.ID,
			},
		}
	}

	_, err := s.client.SendMailV31(&mailjet.MessagesV31{Info: messagesInfo})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSendPostMail, err)
	}

	return nil
}
