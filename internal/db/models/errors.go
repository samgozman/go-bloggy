package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

var (
	ErrValidationFailed = errors.New("ERR_VALIDATION_FAILED")
	ErrNotFound         = errors.New("ERR_NOT_FOUND") // ErrNotFound is returned if item is not found
	ErrDuplicate        = errors.New("ERR_DUPLICATE") // ErrDuplicate is returned if item already exists

	ErrUserLoginRequired      = errors.New("ERR_USER_LOGIN_REQUIRED")
	ErrUserExternalIDRequired = errors.New("ERR_USER_EXTERNAL_ID_REQUIRED")
	ErrUserAuthMethodRequired = errors.New("ERR_USER_AUTH_METHOD_REQUIRED")
	ErrFailedToCreateUser     = errors.New("ERR_FAILED_TO_CREATE_USER")
	ErrFailedToGetUser        = errors.New("ERR_FAILED_TO_GET_USER")

	ErrPostURLRequired         = errors.New("ERR_POST_URL_REQUIRED")
	ErrPostTitleRequired       = errors.New("ERR_POST_TITLE_REQUIRED")
	ErrPostDescriptionRequired = errors.New("ERR_POST_DESCRIPTION_REQUIRED")
	ErrPostContentRequired     = errors.New("ERR_POST_CONTENT_REQUIRED")
	ErrPostWrongKeywordsString = errors.New("ERR_POST_WRONG_KEYWORDS_STRING")
	ErrPostInvalidSlug         = errors.New("ERR_POST_INVALID_SLUG")
	ErrPostUserIDRequired      = errors.New("ERR_POST_USER_ID_REQUIRED")

	ErrCreateSubscription        = errors.New("ERR_CREATE_SUBSCRIPTION")
	ErrSubscriptionEmailRequired = errors.New("ERR_SUBSCRIPTION_EMAIL_REQUIRED")
	ErrGetSubscription           = errors.New("ERR_GET_SUBSCRIPTION")
	ErrGetSubscriptionEmails     = errors.New("ERR_GET_SUBSCRIPTION_EMAILS")
	ErrDeleteSubscription        = errors.New("ERR_DELETE_SUBSCRIPTION")
	ErrUpdateSubscription        = errors.New("ERR_UPDATE_SUBSCRIPTION")
)

// mapGormError maps gorm errors to application errors if possible.
func mapGormError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fmt.Errorf("%w: %w", ErrNotFound, err)

	case errors.Is(err, gorm.ErrDuplicatedKey),
		errors.Is(err, gorm.ErrForeignKeyViolated),
		strings.Contains(err.Error(), "UNIQUE constraint failed"):
		return fmt.Errorf("%w: %w", ErrDuplicate, err)

	case errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, gorm.ErrInvalidValueOfLength),
		errors.Is(err, gorm.ErrInvalidData),
		errors.Is(err, gorm.ErrInvalidField),
		errors.Is(err, gorm.ErrEmptySlice):
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}

	return err
}
