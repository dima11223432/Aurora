package app

import (
	"log/slog"
	grpcApp "recommendationService/internal/app/grpc"
	userDataProvider "recommendationService/internal/services/user_data_provider"
	"recommendationService/internal/storage/postgres"
	"recommendationService/internal/storage/redis"
	"time"
)

type App struct {
	GRPCapp *grpcApp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.New(storagePath)
	redis := redis.NewRedisController(
		"localhost:6379",
		"1111",
		0,
		1,
		tokenTTL,
	)
	if err != nil {
		panic(err)
	}
	userDataProviderService := userDataProvider.New(log,
		storage, redis, tokenTTL)

	grpcapp := grpcApp.New(log, userDataProviderService, userDataProviderService, grpcPort)
	return &App{
		GRPCapp: grpcapp,
	}

}
