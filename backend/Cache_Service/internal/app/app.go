package app

import (
	"CacheService/internal/app/cacheApp"
	"CacheService/internal/config"
	"context"
	"log/slog"
)

type Consumer interface {
	StartWorkerPull(ctx context.Context) error
	Consume(ctx context.Context) error
	Close()
}

type CacheServiceApp interface {
	MustRun()
	Run(ctx context.Context) error
	Stop()
}

type App struct {
	CacheServiceApp CacheServiceApp
}

func New(log *slog.Logger, cfg *config.Config) *App {

	app := cacheApp.New(log, cfg)
	return &App{
		CacheServiceApp: app,
	}

}
