//go:build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package main

import (
	"context"
	_ "github.com/google/subcommands" //nolint:goimports // required by Wire
	"github.com/google/wire"

	"github.com/samgozman/go-bloggy/internal/captcha"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/mailer"
	"github.com/samgozman/go-bloggy/internal/server"
)

func initApp(ctx context.Context, cfg *config.Config) (*serverApp, error) {
	wire.Build(
		db.ProviderSet,
		github.ProviderSet,
		jwt.ProviderSet,
		captcha.ProviderSet,
		mailer.ProviderSet,
		server.ProviderSet,
		handler.ProviderSet,

		newServerApp,
	)

	return &serverApp{}, nil
}
