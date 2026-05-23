// Package app initializes and wires together SSO application components.
package app

import (
	grpcApp "authService/internal/app/grpc"
	"authService/internal/services/auth"
	"authService/internal/storage/postgres"
	"log/slog"
	"time"
)

// App is the top-level application container for the SSO service.
type App struct {
	GRPCapp *grpcApp.App
}

// New creates a new App instance, initializing storage, auth service, and gRPC server.
func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: implement
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcapp := grpcApp.New(log, authService, grpcPort)
	return &App{
		GRPCapp: grpcapp,
	}

}
