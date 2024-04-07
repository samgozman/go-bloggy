package types

import "github.com/mailjet/mailjet-apiv3-go/v4"

type MailjetInterface interface {
	SendMailV31(data *mailjet.MessagesV31, options ...mailjet.RequestOptions) (*mailjet.ResultsV31, error)
}

type ServiceInterface interface {
	SendConfirmationEmail(to, confirmationID string) error
	SendPostEmail(pe *PostEmailSend) error
}

type PostEmailSend struct {
	To          []*Subscriber
	Title       string
	Description string
	Slug        string
}

type Subscriber struct {
	ID    string
	Email string
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
