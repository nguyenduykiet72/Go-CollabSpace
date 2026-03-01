package main

import (
	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/initialize"
	"Go-CollabSpace/internal/server"
	"Go-CollabSpace/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}

	logger.InitLogger(cfg.Server.Mode)

	db, err := initialize.InitDB(cfg.Database)
	if err != nil {
		logger.Log.Fatal("Failed to initialize database", zap.Error(err))
	}

	s := server.NewServer(cfg, db)
	if err := s.Run(); err != nil {
		logger.Log.Fatal("Failed to run server", zap.Error(err))
	}
}
