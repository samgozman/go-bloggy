package captcha

import (
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvideClient(t *testing.T) {
	t.Run("ProvideClient", func(t *testing.T) {
		secret := config.HCaptchaSecret("test")
		got := ProvideClient(secret)
		assert.NotNil(t, got)
	})
}

func TestProvideHCaptchaSecret(t *testing.T) {
	t.Run("ProvideHCaptchaSecret", func(t *testing.T) {
		cfg := &config.Config{
			HCaptchaSecret: "test",
		}
		got := ProvideHCaptchaSecret(cfg)
		assert.Equal(t, cfg.HCaptchaSecret, got)
	})
}
