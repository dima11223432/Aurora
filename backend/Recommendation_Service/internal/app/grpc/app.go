package grpcApp

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"recommendationService/internal/domain/models"
	grpcUserDataProvider "recommendationService/internal/grpc/userDataProvider"
)

type userDataProvider interface {
	GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, userDataProvider userDataProvider, port int) *App {
	gRPCServer := grpc.NewServer()
	grpcUserDataProvider.Register(gRPCServer, userDataProvider)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

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

func (a *App) Stop() {
	const op = "grpcApp"

	a.log.With(slog.String("op", op),
		slog.Int("port", a.port),
	)
	a.gRPCServer.GracefulStop()
}
