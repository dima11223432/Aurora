package grpc

import (
	grpcAuth "API_Service/internal/grpc"
	authinterceptor "API_Service/internal/grpc/AuthInterceptor"
	"API_Service/internal/services"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"log/slog"
	"net"
	"strconv"

	"API_Service/internal/config"

	ssov1 "github.com/dima11223432/Aurora_SSO_Protos/api/gen/v1"
	recv1 "github.com/dima11223432/recommendationService_protos/api/gen/v1"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
}

func New(port int, logger *slog.Logger, jwtSecret string, publicRoutes []string, ssoConfig config.ServiceConfig, recsConfig config.ServiceConfig) *App {
	authInterceptor := authinterceptor.NewAuthInterceptor(authinterceptor.AuthConfig{JwtSecret: jwtSecret, PublicRoutes: publicRoutes})

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			authInterceptor.SetAuthInterceptor(),
		),
	)

	authAddr := fmt.Sprintf("%s:%d", ssoConfig.SSO.Host, ssoConfig.SSO.Port)
	recsAddr := fmt.Sprintf("%s:%d", recsConfig.RECS.Host, recsConfig.RECS.Port)

	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Fatalf("cant connect to authService: %v", err)
	}
	authClient := ssov1.NewAuthServiceClient(authConn)
	authService := services.NewAuthService(logger, authClient, authInterceptor)

	recsConn, err := grpc.NewClient(recsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("cant connect to recommendationService: %v", slog.Any("err", err))
	}
	recsClient := recv1.NewRecommendationServiceClient(recsConn)
	recommendationService := services.NewRecommendationService(recsClient, logger, authInterceptor)

	grpcAuth.RegisterGrpcServer(gRPCServer, authService, recommendationService)
	reflection.Register(gRPCServer)
	logger.Info("gRPC server initialized", slog.String("gRPC_Port", strconv.Itoa(port)))

	grpc_prometheus.Register(gRPCServer)

	return &App{
		log:  logger,
		gRPC: gRPCServer,
		port: port,
	}
}

func (a *App) MustRun() {
	a.log.Info("Starting gRPC server", slog.Int("port", a.port))
	if err := a.Run(); err != nil {
		a.log.Error("gRPC server failed to start", slog.Any("error", err))
	}
}

func (a *App) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.log.Error("Failed to listen on port", slog.Int("port", a.port), slog.Any("error", err))
		return err
	}

	a.log.Info("gRPC server listening", slog.String("addr", listener.Addr().String()))

	if err := a.gRPC.Serve(listener); err != nil {
		a.log.Error("gRPC server stopped with error", slog.Any("error", err))
		return err
	}

	a.log.Info("gRPC server stopped gracefully")
	return nil
}

func (a *App) Close() {
	a.log.Info("Shutting down gRPC server...")
	a.gRPC.GracefulStop()
	a.log.Info("gRPC server stopped")
}
