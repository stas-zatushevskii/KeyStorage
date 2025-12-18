package main

import (
	"server/config"
	"server/internal/logger"
	"server/internal/repository/db"

	"go.uber.org/zap"
)

func main() {
	config.GetConfig()
	logger.GetLogger()

	database, err := db.NewConnection()
	if err != nil {
		logger.Log.Error("Error connecting to database", zap.Error(err))
	}

	database.Ping()
}
