package main

import (
	"context"
	"github.com/kovalyov-valentin/data-ops/internal/config"
	v1 "github.com/kovalyov-valentin/data-ops/internal/http/v1"
	"github.com/kovalyov-valentin/data-ops/internal/services"
	"github.com/kovalyov-valentin/data-ops/internal/storage"
	"github.com/kovalyov-valentin/data-ops/internal/storage/postgres"
	"github.com/kovalyov-valentin/data-ops/internal/storage/redis"
	http_server "github.com/kovalyov-valentin/data-ops/pkg/http-server"
	"github.com/sirupsen/logrus"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.CtxTimeout)
	defer cancel()

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))

	db, err := postgres.NewPostgresDB(cfg.PostgresDB)
	if err != nil {
		logrus.Fatalf("failed to connect postgres db: %s", err.Error())
	}
	defer db.Close()

	rdb, err := redis.NewRedisDB(ctx, cfg.Redis)
	if err != nil {
		logrus.Fatalf("failed to connect redis: %s", err.Error())
	}
	defer rdb.Close()

	// init clickhouse

	//click, err := clickhouse.NewClickhouseDB(ctx, cfg.Clickhouse)
	//if err != nil {
	//	logrus.Fatalf("failed to connect clickhouse: %s", err.Error())
	//}

	// init nats

	//nc, err := nats.Connect(cfg.Nats.Port)
	//if err != nil {
	//	panic(err)
	//}
	//defer func() { _ = nc.Drain() }()
	//njs, err := nc.JetStream()
	//if err != nil {
	//	panic(err)
	//}

	storage := storage.NewRepository(db, rdb)
	service := services.NewService(storage)
	handlers := v1.NewHandler(cfg, service)

	srv := http_server.Server{}

	serverErrors := make(chan error, 1)
	go func() {
		logrus.Printf("Start listen http service on %s at %s\n", cfg.Address, time.Now().Format(time.DateTime))
		err := srv.Run(cfg.HTTPServer, handlers.InitRoutes())
		if err != nil {
			logrus.Printf("shutting down the server: %s\n", cfg.Address)
		}
		serverErrors <- err
	}()

	osSignal := make(chan os.Signal, 1)

	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT)
	select {
	case err := <-serverErrors:
		logrus.Printf("error starting server: %v\n", err)
	case <-osSignal:
		logrus.Println("start shutdown...")
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Printf("graceful shutdown error: %v\n", err)
			os.Exit(1)
		}
	}
	logrus.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
