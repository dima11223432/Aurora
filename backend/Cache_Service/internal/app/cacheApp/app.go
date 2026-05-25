// Package cacheApp implements the main Cache Service application that consumes
// analysed data from Kafka and stores it in Redis.
package cacheApp

import (
	"CacheService/internal/brokers/kafka"
	"CacheService/internal/config"
	analyseDataProvider "CacheService/internal/services/AnalyseDataProvider"
	"CacheService/internal/storage/redis"
	"context"
	"log/slog"
)

// Consumer defines the interface for Kafka message consumption lifecycle.
type Consumer interface {
	StartWorkerPull(ctx context.Context) error
	Consume(ctx context.Context) error
	Close()
}

// App represents the Cache Service application with Kafka consumer lifecycle.
type App struct {
	log      *slog.Logger
	consumer Consumer

	cancel context.CancelFunc
}

// New creates a new App instance, initializing Redis and Kafka consumer.
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
		[]string{cfg.KafkaConfig.Host},
		cfg.KafkaConfig.Topic,
		cfg.KafkaConfig.ConsumerGroup,
		int32(cfg.KafkaConfig.MaxPulls),
		analyseDataProvider,
	)
	log.Info("Consumer created")
	return &App{
		log:      log,
		consumer: consumer,
	}
}

// MustRun starts the cache service and panics on error.
func (a *App) MustRun() {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	if err := a.Run(ctx); err != nil {
		panic(err)
	}
}

// Run starts the Kafka consumer and worker pull in background goroutines.
func (a *App) Run(ctx context.Context) error {
	const op = "grpcApp"

	a.log.With(slog.String("op", op))

	go a.consumer.StartWorkerPull(ctx)
	go a.consumer.Consume(ctx)
	return nil

}

// Stop cancels the context, gracefully shutting down the consumer.
func (a *App) Stop() {
	const op = "grpcApp"

	a.log.With(slog.String("op", op))

	if a.cancel != nil {
		a.cancel()
	}

}
