package app

import (
	"log/slog"
	grpcApp "recommendationService/internal/app/grpc"
	"recommendationService/internal/services/auth"
	"recommendationService/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCapp *grpcApp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: implement
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcapp := grpcApp.New(log, authService, grpcPort)
	return &App{
		GRPCapp: grpcapp,
	}

}
