package cacheApp

import (
	"context"
	"log/slog"
)

type Consumer interface {
	StartWorkerPull(ctx context.Context) error
	Consume(ctx context.Context) error
	Close()
}

type App struct {
	log      *slog.Logger
	consumer Consumer

	cancel context.CancelFunc
}

func New(log *slog.Logger, consumer Consumer) *App {
	return &App{
		log:      log,
		consumer: consumer,
	}
}

func (a *App) MustRun() {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

func (a *App) Run(ctx context.Context) error {
	const op = "grpcApp"

	a.log.With(slog.String("op", op))

	go a.consumer.StartWorkerPull(ctx)
	go a.consumer.Consume(ctx)
	return nil

}

func (a *App) Stop() {
	const op = "grpcApp"

	a.log.With(slog.String("op", op))

	if a.cancel != nil {
		a.cancel()
	}

}
