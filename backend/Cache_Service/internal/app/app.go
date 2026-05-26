// Package app initializes and wires together Cache Service components.
package app

import (
	"CacheService/internal/app/cacheApp"
	"CacheService/internal/config"
	"context"
	"log/slog"
)

// Consumer defines the interface for Kafka message consumption.
type Consumer interface {
	StartWorkerPull(ctx context.Context) error
	Consume(ctx context.Context) error
	Close()
}

// CacheServiceApp defines the interface for the cache service application lifecycle.
type CacheServiceApp interface {
	MustRun()
	Run(ctx context.Context) error
	Stop()
}

// App is the top-level application container for the Cache Service.
type App struct {
	CacheServiceApp CacheServiceApp
}

// New creates a new App instance, initializing the cache service app.
func New(log *slog.Logger, cfg *config.Config) *App {

	app := cacheApp.New(log, cfg)
	return &App{
		CacheServiceApp: app,
	}

}
