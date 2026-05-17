package app

import (
	grpcApp "API_Service/internal/app/grpc"
	"API_Service/internal/cache"
	"API_Service/internal/config"
	"context"
	"log/slog"
	"time"
)

type App struct {
	GRPCApp *grpcApp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	redisCache := cache.NewRedisCache(
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
		cfg.RedisConfig.Port,
		time.Duration(1)*time.Hour,
	)

	err := redisCache.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	grpcApp := grpcApp.New(cfg.GRPC.Port, log, cfg.Auth.JwtSecret, cfg.Auth.PublicMethods, cfg.Services.SSO, cfg.Services.RECS)

	return &App{
		GRPCApp: grpcApp,
	}

}
