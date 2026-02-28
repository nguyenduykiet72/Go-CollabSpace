package main

import (
	"Go-CollabSpace/config"
	"Go-CollabSpace/pkg/logger"
	"fmt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("error loading config: %v", err))
	}

	logger.InitLogger(cfg.Server.Mode)
}
