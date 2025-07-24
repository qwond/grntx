package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/qwond/grntx/api/v1"
	"github.com/qwond/grntx/database"
	"github.com/qwond/grntx/internal/grinex"
	"github.com/qwond/grntx/internal/repository"
	"github.com/qwond/grntx/internal/service"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	DSN       string `env:"DSN, required"`
	GrinexURL string `env:"GRINEX_URL, default=https://grinex.io"`
}

func main() {
	// init logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cant initialize logger: %v", err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("service starting")
	ctx := context.Background()

	// Load config from environment

	var cfg Config
	err = envconfig.Process(ctx, &cfg)
	if err != nil {
		logger.Fatal("cannot configure service", zap.Error(err))
	}

	// Migrate database
	err = database.MigrateDB(cfg.DSN)
	if err != nil {
		logger.Fatal("Failed to migrate database:", zap.Error(err))
	}

	// Create repository
	repo, err := repository.New(cfg.DSN)
	if err != nil {
		logger.Fatal("cannot create repository", zap.Error(err))
	}

	// Construct rates service
	rateSvc, err := service.New(logger, grinex.New(cfg.GrinexURL), repo)
	if err != nil {
		logger.Fatal("cannot prepare service", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	v1.RegisterRatesServiceServer(grpcServer, rateSvc)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	// Start server
	lstnr, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		logger.Info("server listening", zap.Any("Addr", lstnr.Addr()))
		errChan <- grpcServer.Serve(lstnr)
	}()

	select {
	case err := <-errChan:
		logger.Fatal("cannot serve", zap.Error(err))

	case sig := <-sigChan:
		logger.Info("shutdown start")
		defer logger.Info("shutdown complete", zap.Any("signal", sig))
		if err = lstnr.Close(); err != nil {
			logger.Error("cannot close listener", zap.Error(err))
		}
	}
}
