package db

import (
	"database/sql"
	"fmt"
	"os"
	"server/internal/app/config"
	"server/internal/pkg/logger"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	// -- only for debug
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func NewConnection() (*sql.DB, error) {

	dsn := config.App.GetDSN()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("[Db] failed to set connection with database: %w", err)
	}
	db.SetConnMaxLifetime(config.App.GetConnMaxLifetime() * time.Second)
	db.SetMaxIdleConns(config.App.GetMaxIdleConns())
	db.SetMaxOpenConns(config.App.GetMaxOpenConns())

	logger.Log.Info(fmt.Sprintf("[Db] set connection to database: %v", dsn))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[Db] database ping failed: %w", err)
	}

	if config.App.GetDebugMode() {
		db = debugModeConnection(dsn, db)
	}

	return db, nil
}

func debugModeConnection(dsn string, db *sql.DB) *sql.DB {
	loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))

	db = sqldblogger.OpenDriver(
		dsn,
		db.Driver(),
		loggerAdapter,
		sqldblogger.WithLogArguments(false),
		sqldblogger.WithMinimumLevel(sqldblogger.LevelDebug),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithLogArguments(true),
	)

	return db
}
