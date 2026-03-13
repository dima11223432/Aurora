package app

import (
	"CacheService/internal/app/grpc"
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

func New(log *slog.Logger, consumer Consumer) *App {

	app := cacheApp.New(log, consumer)
	return &App{
		CacheServiceApp: app,
	}

}
