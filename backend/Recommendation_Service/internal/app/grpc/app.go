package grpcApp

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	// apiv1 "recommendationService/api/gen/v1"
	"recommendationService/internal/domain/models"
	grpcUserDataProvider "recommendationService/internal/grpc/userDataProvider"
)

type ParsingChannelsProvider interface {
	GetAllParsingChannels(ctx context.Context) ([]string, error)
	AddNewParsingChannel(ctx context.Context, channel string) error
	DeleteParsingChannel(ctx context.Context, channel string) error
	GetParsingChannelsWithCategories(ctx context.Context) (map[string][]string, error)
	GetAllCategories(ctx context.Context) ([]string, error)
}

type userDataProvider interface {
	GetUserPriorityChannels(ctx context.Context, userID int64) ([]models.PriorityChannel, error)
}

type NewsDataProvider interface {
	GetRecommendatedPosts(ctx context.Context, userID int64, cursor *models.Cursor) ([]models.Post, *models.Cursor, error)
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, userDataProvider userDataProvider, parsingChannelsProvider ParsingChannelsProvider, newsDataProvider NewsDataProvider, port int) *App {
	gRPCServer := grpc.NewServer()
	grpcUserDataProvider.Register(gRPCServer, userDataProvider, parsingChannelsProvider, newsDataProvider)
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
