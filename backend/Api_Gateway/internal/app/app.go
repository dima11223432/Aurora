// Package app initializes and wires together the API Gateway application components,
// including Redis cache, gRPC server, and service clients.
package app

import (
	grpcApp "API_Service/internal/app/grpc"
	"API_Service/internal/cache"
	"API_Service/internal/config"
	"context"
	"log/slog"
	"time"
)

// App represents the top-level application container for the API Gateway.
type App struct {
	GRPCApp *grpcApp.App
}

// New creates a new App instance, initializing Redis cache and gRPC server.
func New(log *slog.Logger, cfg *config.Config) *App {
	redisCache := cache.NewRedisCache(
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
		cfg.RedisConfig.Port,
		time.Duration(1)*time.Hour,
	)

	log.Info(
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
		cfg.RedisConfig.Port,
	)

	err := redisCache.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	grpcApp := grpcApp.New(cfg.GRPC.Port,
		log,
		cfg.Auth.JwtSecret,
		cfg.Auth.PublicMethods,
		cfg.Services,
		cfg.Services,
	)

	return &App{
		GRPCApp: grpcApp,
	}

}
