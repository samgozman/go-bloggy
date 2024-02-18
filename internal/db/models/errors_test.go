package models

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func Test_mapGormError(t *testing.T) {
	t.Run("should map gorm errors to application errors if possible", func(t *testing.T) {
		err := mapGormError(gorm.ErrRecordNotFound)
		assert.ErrorIs(t, err, ErrNotFound)

		err = mapGormError(gorm.ErrDuplicatedKey)
		assert.ErrorIs(t, err, ErrDuplicate)

		err = mapGormError(gorm.ErrInvalidValue)
		assert.ErrorIs(t, err, ErrValidationFailed)
	})

	t.Run("should return original error if it cannot be mapped", func(t *testing.T) {
		err := mapGormError(gorm.ErrInvalidField)
		assert.ErrorIs(t, err, gorm.ErrInvalidField)
	})
}
