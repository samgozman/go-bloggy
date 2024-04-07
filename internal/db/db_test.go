package db

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"

	modelsMock "github.com/samgozman/go-bloggy/mocks/db/models"
)

func TestNewDatabase(t *testing.T) {
	t.Run("NewDatabase", func(t *testing.T) {
		conn := &gorm.DB{}
		models := &Models{
			users:       modelsMock.NewMockUserRepositoryInterface(t),
			posts:       modelsMock.NewMockPostRepositoryInterface(t),
			subscribers: modelsMock.NewMockSubscriberRepositoryInterface(t),
		}
		got := NewDatabase(conn, models)
		assert.NotNil(t, got)
		assert.Equal(t, conn, got.conn)
		assert.Equal(t, models, got.models)
		assert.Equal(t, models, got.Models())
		assert.Equal(t, conn, got.GetConn())
		assert.NotNil(t, got.models.Users())
		assert.NotNil(t, got.models.Posts())
		assert.NotNil(t, got.models.Subscribers())
	})
}
