package main

import (
	"server/internal/app"
	"server/internal/app/config"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	config.GetConfig()
	logger.GetLogger()

	application, err := app.New()
	if err != nil {
		logger.Log.Error("Failed to create Application", zap.Error(err))
		return
	}

	err = application.Start()
	if err != nil {
		logger.Log.Error("Failed to start application", zap.Error(err))
		return
	}
}
