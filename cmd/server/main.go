package main

import (
	"server/config"
	"server/internal/logger"
	"server/internal/repository"

	"go.uber.org/zap"
)

func main() {
	config.GetConfig()
	logger.GetLogger()

	database, err := repository.NewConnection()
	if err != nil {
		logger.Log.Error("Error connecting to database", zap.Error(err))
	}

	database.Ping()
}
