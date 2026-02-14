package app

import (
	grpcApp "API_Service/internal/app/grpc"
	apiKafka "API_Service/internal/broker/kafka"
	"API_Service/internal/cache"
	"API_Service/internal/config"
	apiService "API_Service/internal/services"
	"API_Service/internal/storage/postgres"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCApp *grpcApp.App
}

func New(log *logrus.Entry, cfg *config.Config) *App {
	storage, err := postgres.New(cfg.StoragePass)
	if err != nil {
		panic(err)
	}
	redisCache := cache.NewRedisCache(
		cfg.RedisConfig.Host,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
		cfg.RedisConfig.Port,
		time.Duration(1)*time.Hour,
	)

	err = redisCache.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	publisher := apiKafka.NewProducer([]string{"localhost:9092"}, "notification.sent")

	apiService := apiService.New(storage, publisher, redisCache)
	grpcApp := grpcApp.New(apiService, cfg.GRPC.Port, log, cfg.Auth.JwtSecret, cfg.Auth.PublicMethods)

	return &App{
		GRPCApp: grpcApp,
	}

}
