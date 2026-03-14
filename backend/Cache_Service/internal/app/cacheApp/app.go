package cacheApp

import (
	"CacheService/internal/brokers/kafka"
	"CacheService/internal/config"
	analyseDataProvider "CacheService/internal/services/AnalyseDataProvider"
	"CacheService/internal/storage/redis"
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

func New(log *slog.Logger, cfg *config.Config) *App {

	redisCache := redis.NewRedisController(
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
		cfg.RedisConfig.Port,
		cfg.TokenTTL,
	)

	analyseDataProvider := analyseDataProvider.NewRedisService(log, redisCache, cfg.TokenTTL)

	consumer := kafka.NewConsumer(
		log,
		[]string{"localhost:9092"},
		"news_data",
		"news_consumer_group",
		10,
		analyseDataProvider,
	)
	log.Info("Consumer created")
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
