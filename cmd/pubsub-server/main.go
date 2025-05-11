package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"awesomeProject3/internal/domain/repository"
	"awesomeProject3/internal/pubsub/delivery/grpc"
	"awesomeProject3/internal/usecase/publish"
	"awesomeProject3/internal/usecase/subscribe"
	"awesomeProject3/pkg/config"
	"awesomeProject3/pkg/logger"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		cfg = config.Default()
	}

	// Initialize logger
	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic(err)
	}

	// Create repository
	eventRepo := repository.NewInMemoryRepository()

	// Create use cases
	publishUC := publish.New(eventRepo, log)
	subscribeUC := subscribe.New(eventRepo, log)

	// Create gRPC handler with use cases
	handler := grpc.NewHandler(log, publishUC, subscribeUC)

	// Create and start gRPC server
	server := grpc.NewServer(handler, log, cfg.Server.Port)

	// Handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.WithError(err).Fatal("failed to start server")
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown
	server.Stop()
	if err := eventRepo.Close(context.Background()); err != nil {
		log.WithError(err).Error("failed to close repository")
	}
}
