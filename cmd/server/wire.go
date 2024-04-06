//go:build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
)

func initApp(ctx context.Context, cfg *config.Config) (*tempApp, error) {
	wire.Build(
		db.ProviderSet,

		newTempApp,
	)

	return &tempApp{}, nil
}
