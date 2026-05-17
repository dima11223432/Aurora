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
	errorhandler "API_Service/internal/grpc/ErrorHandler"

	v1 "API_Service/api/gen/v1"
	app "API_Service/internal/app"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
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
	log := setupLogger(cfg.Env)
	log.Info("starting app", slog.String("env", cfg.Env))
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
		runtime.WithErrorHandler(errorhandler.CallHanlder),
	)
	err := v1.RegisterApiServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%d", cfg.GRPC.Port),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Error("failed to register gateway", slog.String("error", err.Error()))
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

	go func() {
		metrixMux := http.NewServeMux()
		metrixMux.Handle("/metrics", promhttp.Handler())
		log.Info("Starting Prometheus metrics on :2112")

		handleWithCors := c.Handler(metrixMux)

		if err := http.ListenAndServe(":2112", handleWithCors); err != nil {
			log.Error("metrics server failed", slog.String("error", err.Error()))
		}
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
