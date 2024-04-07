package mailer

import (
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/mailer/types"
)

// Config is a struct that holds all the configuration for mailer.
type Config struct {
	PublicKey  string
	PrivateKey string
	Options    *types.Options
}

// ProvideConfig is a wire provider function for mailer.Config.
func ProvideConfig(cfg *config.Config) *Config {
	return &Config{
		PublicKey:  cfg.MailerJet.PublicKey,
		PrivateKey: cfg.MailerJet.PrivateKey,
		Options: &types.Options{
			FromEmail:                    cfg.MailerJet.FromEmail,
			FromName:                     cfg.MailerJet.FromName,
			ConfirmationTemplateID:       cfg.MailerJet.ConfirmationTemplateID,
			ConfirmationTemplateURLParam: cfg.MailerJet.ConfirmationTemplateURLParam,
			PostTemplateID:               cfg.MailerJet.PostTemplateID,
			PostTemplateURLParam:         cfg.MailerJet.PostTemplateURLParam,
			UnsubscribeURLParam:          cfg.MailerJet.UnsubscribeURLParam,
		},
	}
}

// ProvideService is a wire provider function for mailer.Service.
func ProvideService(cfg *Config) *Service {
	return NewService(cfg.PublicKey, cfg.PrivateKey, cfg.Options)
}

// ProviderSet is a wire.ProviderSet for mailer package.
var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideService,
	wire.Bind(new(types.ServiceInterface), new(*Service)),
)
