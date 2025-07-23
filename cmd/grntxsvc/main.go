package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/qwond/grntx/api/v1"
	"github.com/qwond/grntx/internal/service"
	"github.com/qwond/grntx/pkg/grinex"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

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

	// Construct rates service
	rateSvc := service.New(grinex.New("https://grinex.io"))

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
