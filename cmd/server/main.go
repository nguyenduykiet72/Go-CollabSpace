package main

import (
	"Go-CollabSpace/config"
	"Go-CollabSpace/internal/controller"
	"Go-CollabSpace/internal/initialize"
	"Go-CollabSpace/internal/repository"
	"Go-CollabSpace/internal/router"
	"Go-CollabSpace/internal/service"
	"Go-CollabSpace/pkg/logger"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading config: %v", err))
	}

	logger.InitLogger(cfg.Server.Mode)
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {

		}
	}(logger.Log)

	logger.Log.Info("Starting server", zap.Int("Port", cfg.Server.Port))

	db, err := initialize.InitDB(cfg.Database)
	if err != nil {
		return
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	r := router.NewRouter(userController)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Log.Info("Server is running on ", zap.String("address:", addr))

	if err := r.Run(addr); err != nil {
		logger.Log.Fatal("Server start failed", zap.Error(err))
	}
}
