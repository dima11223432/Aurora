package app

import (
	grpcApp "authService/internal/app/grpc"
	"authService/internal/services/auth"
	"authService/internal/storage/sqlite"
	"log/slog"
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
