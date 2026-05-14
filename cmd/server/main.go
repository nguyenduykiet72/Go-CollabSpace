package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/initialize"
	"Go-CollabSpace/internal/server"
	"Go-CollabSpace/pkg/logger"
)

// Build metadata, injected at compile time via -ldflags:
//
//	go build -ldflags "-X main.version=v1.2.3 -X main.commit=abc123 -X main.buildDate=2026-01-01"
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}

	logger.InitLogger(cfg.Server.Mode)
	logger.Log.Info("Starting Go-CollabSpace",
		zap.String("version", version),
		zap.String("commit", commit),
		zap.String("build_date", buildDate),
		zap.String("mode", cfg.Server.Mode),
	)

	db, err := initialize.InitDB(cfg.Database)
	if err != nil {
		logger.Log.Fatal("Failed to initialize database", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := server.NewServer(cfg, db)
	if err := s.Run(ctx); err != nil {
		logger.Log.Fatal("Failed to run server", zap.Error(err))
	}
}
