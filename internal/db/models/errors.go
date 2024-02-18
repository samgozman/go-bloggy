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
