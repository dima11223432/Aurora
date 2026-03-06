package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"API_Service/internal/config"

	v1 "API_Service/api/gen/v1"
	app "API_Service/internal/app"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {

	cfg := config.MustLoad()
	log := logrus.New().WithField("service", "api")
	log.Info(cfg)

	app := app.New(log, cfg)

	go app.GRPCApp.Run()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.ToLower(key) == "authorization" {
				return key, true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
	)
	err := v1.RegisterApiServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%d", cfg.GRPC.Port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatal(err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.Auth.Cors_urls,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	go func() {
		log.Info("HTTP gateway on :8081")
		http.ListenAndServe(":8081", c.Handler(mux))
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	app.GRPCApp.Close()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:

		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:

		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
