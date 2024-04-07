package db

import (
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestProvideDSN(t *testing.T) {
	t.Run("ProvideDSN", func(t *testing.T) {
		cfg := &config.Config{
			DSN: "test",
		}
		got := ProvideDSN(cfg)
		assert.Equal(t, cfg.DSN, got)
	})
}

func TestProvideModels(t *testing.T) {
	t.Run("ProvideModels", func(t *testing.T) {
		conn := &gorm.DB{}
		got := ProvideModels(conn)
		assert.NotNil(t, got)
		assert.NotNil(t, got.Users())
		assert.NotNil(t, got.Posts())
		assert.NotNil(t, got.Subscribers())
	})
}

func TestProvideDatabase(t *testing.T) {
	t.Run("ProvideDatabase", func(t *testing.T) {
		conn := &gorm.DB{}
		models := &Models{}
		got, err := ProvideDatabase(conn, models)
		assert.Nil(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, conn, got.GetConn())
		assert.Equal(t, models, got.models)
	})
}
