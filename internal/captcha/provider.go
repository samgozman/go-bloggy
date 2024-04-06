package captcha

import (
	"github.com/google/wire"
	"github.com/kataras/hcaptcha"
	"github.com/samgozman/go-bloggy/internal/config"
)

// ProvideHCaptchaSecret is a Wire provider function that returns HCaptcha secret from the config.
func ProvideHCaptchaSecret(cfg *config.Config) config.HCaptchaSecret {
	return cfg.HCaptchaSecret
}

// ProvideClient is a Wire provider function that creates a new HCaptcha client.
func ProvideClient(secret config.HCaptchaSecret) *Client {
	return hcaptcha.New(string(secret))
}

// ProviderSet is a wire provider set for HCaptcha.
var ProviderSet = wire.NewSet(
	ProvideHCaptchaSecret,
	ProvideClient,
)
