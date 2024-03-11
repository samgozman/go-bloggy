package mailer

import "errors"

var (
	ErrSendConfirmationMail = errors.New("error sending confirmation mail")
	ErrSendPostMail         = errors.New("error sending post mail")
)
