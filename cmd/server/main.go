package main

import (
	"fmt"

	"go.uber.org/zap"

	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/initialize"
	"Go-CollabSpace/internal/server"
	"Go-CollabSpace/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}

	logger.InitLogger(cfg.Server.Mode)
	defer func(Log *zap.Logger) {
		_ = Log.Sync()
	}(logger.Log)

	logger.Log.Info("Starting server", zap.Int("Port", cfg.Server.Port))

	db, err := initialize.InitDB(cfg.Database)
	if err != nil {
		return
	}

	srv := server.NewServer(cfg, db)
	if err := srv.Run(); err != nil {
		logger.Log.Fatal("Failed", zap.Error(err))
	}
}
