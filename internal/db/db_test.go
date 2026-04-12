package db

import (
	"github.com/samgozman/go-bloggy/internal/db/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	t.Run("NewDatabase", func(t *testing.T) {
		conn := &gorm.DB{}
		m := NewModels(
			models.NewUserRepository(conn),
			models.NewPostRepository(conn),
			models.NewSubscribersRepository(conn),
		)
		got := NewDatabase(conn, m)
		assert.NotNil(t, got)
		assert.Equal(t, conn, got.conn)
		assert.Equal(t, m, got.models)
		assert.Equal(t, m, got.Models())
		assert.Equal(t, conn, got.GetConn())
		assert.NotNil(t, got.models.Users())
		assert.NotNil(t, got.models.Posts())
		assert.NotNil(t, got.models.Subscribers())
	})
}
