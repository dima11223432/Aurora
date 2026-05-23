// Package grpcApp sets up the gRPC server for the SSO service.
package grpcApp

import (
	grpcauth "authService/internal/grpc/auth"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// App represents the gRPC server for the SSO service.
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates a new gRPC server App with the given auth service.
func New(log *slog.Logger, authService grpcauth.Auth, port int) *App {
	gRPCServer := grpc.NewServer()
	grpcauth.Register(gRPCServer, authService)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun starts the gRPC server and panics on error.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run starts the gRPC server on the configured port.
func (a *App) Run() error {
	const op = "grpcApp.Run"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("starting gRPC server")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("failed to listen:%s, %w", op, err)
	}
	log.Info("gRPC server listening", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	return nil
}

// Stop gracefully stops the gRPC server.
func (a *App) Stop() {
	const op = "grpcApp"

	a.log.With(slog.String("op", op),
		slog.Int("port", a.port),
	)
	a.gRPCServer.GracefulStop()
}
