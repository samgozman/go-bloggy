package jwt

import "errors"

var (
	ErrExpiresAtMustBeInTheFuture = errors.New("expiresAt must be in the future")
	ErrErrorSigningToken          = errors.New("error signing token")
	ErrErrorParsingToken          = errors.New("error parsing token")
	ErrInvalidToken               = errors.New("invalid token")
)
