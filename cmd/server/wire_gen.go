// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"github.com/samgozman/go-bloggy/internal/captcha"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/handler"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/mailer"
	"github.com/samgozman/go-bloggy/internal/server"
)

import (
	_ "github.com/google/subcommands"
)

// Injectors from wire.go:

func initApp(ctx context.Context, cfg *config.Config) (*serverApp, error) {
	jwtSecretKey := jwt.ProvideJWTSecretKey(cfg)
	service := jwt.ProvideService(jwtSecretKey)
	echo := server.ProvideServer(service)
	handlerConfig := handler.ProvideConfig(cfg)
	githubConfig := github.ProvideConfig(cfg)
	githubService := github.ProvideService(githubConfig)
	dsn := db.ProvideDSN(cfg)
	gormDB, err := db.ProvideConnection(dsn)
	if err != nil {
		return nil, err
	}
	models := db.ProvideModels(gormDB)
	database, err := db.ProvideDatabase(gormDB, models)
	if err != nil {
		return nil, err
	}
	hCaptchaSecret := captcha.ProvideHCaptchaSecret(cfg)
	client := captcha.ProvideClient(hCaptchaSecret)
	mailerConfig := mailer.ProvideConfig(cfg)
	mailerService := mailer.ProvideService(mailerConfig)
	handlerHandler := handler.ProvideHandler(handlerConfig, githubService, service, database, client, mailerService)
	mainServerApp := newServerApp(echo, handlerHandler)
	return mainServerApp, nil
}
