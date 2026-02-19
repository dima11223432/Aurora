package grpc

import (
	grpcAuth "API_Service/internal/grpc"
	authinterceptor "API_Service/internal/grpc/AuthInterceptor"
	"API_Service/internal/services"
	"fmt"
	"log"
	"net"

	ssov1 "github.com/dima11223432/protos/gen/go/sso"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log  *logrus.Entry
	gRPC *grpc.Server
	port int
}

func New(port int, logger *logrus.Entry, jwtSecret string, publicRoutes []string) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			authinterceptor.AuthInterceptor(authinterceptor.AuthConfig{
				JwtSecret:    jwtSecret,
				PublicRoutes: publicRoutes,
			}),
		),
	)
	authConn, err := grpc.NewClient(":44044", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("cant connect to authService: %v", err)
	}
	authClient := ssov1.NewAuthClient(authConn)
	authService := services.NewAuthService(authClient)

	grpcAuth.RegisterGrpcServer(gRPCServer, api, authService)
	reflection.Register(gRPCServer)
	logger.Infof("gRPC server initialized on port %d", port)

	return &App{
		log:  logger,
		gRPC: gRPCServer,
		port: port,
	}
}

func (a *App) MustRun() {
	a.log.Infof("Starting gRPC server on port %d...", a.port)
	if err := a.Run(); err != nil {
		a.log.Fatalf("gRPC server failed to start: %v", err)
	}
}

func (a *App) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.log.Errorf("Failed to listen on port %d: %v", a.port, err)
		return err
	}

	a.log.Infof("gRPC server listening on %s", listener.Addr().String())

	if err := a.gRPC.Serve(listener); err != nil {
		a.log.Errorf("gRPC server stopped with error: %v", err)
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
