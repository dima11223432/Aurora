package app

import (
	"log/slog"
	grpcApp "recommendationService/internal/app/grpc"
	userDataProvider "recommendationService/internal/services/user_data_provider"
	"recommendationService/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCapp *grpcApp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}
	userDataProviderService := userDataProvider.New(log, storage, tokenTTL)

	grpcapp := grpcApp.New(log, userDataProviderService, grpcPort)
	return &App{
		GRPCapp: grpcapp,
	}

}
