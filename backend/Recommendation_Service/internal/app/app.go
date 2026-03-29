package app

import (
	"context"
	"log/slog"
	grpcApp "recommendationService/internal/app/grpc"
	"recommendationService/internal/config"
	userDataProvider "recommendationService/internal/services/user_data_provider"
	"recommendationService/internal/storage/postgres"
	"recommendationService/internal/storage/redis"
)

type App struct {
	GRPCapp *grpcApp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	storage, err := postgres.New(cfg.StoragePass)
	redis := redis.NewRedisController(
		cfg.Redis.Host,
		cfg.Redis.Password,
		0,
		1,
		cfg.TokenTTL,
	)
	redis.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	userDataProviderService := userDataProvider.New(log,
		storage, redis, cfg.TokenTTL)

	grpcapp := grpcApp.New(log, userDataProviderService, userDataProviderService, cfg.GRPC.Port)
	return &App{
		GRPCapp: grpcapp,
	}
}
