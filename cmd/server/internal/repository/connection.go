package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	// -- only for debug
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

type Config interface {
	GetDSN() string
	GetMaxIdleConns() int
	GetMaxOpenConns() int
	GetConnMaxLifetime() time.Duration
	GetDebugMode() bool
}

func NewConnection(cfg Config) (*sql.DB, error) {

	var dsn string // TODO GetDSN function in config that crates DSN from config [DB] data

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// logger.Info("[Db] Failed to set connection with database") fixme
		return nil, fmt.Errorf("[Db] failed to set connection with database: %v", err)
	}
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime() * time.Second)
	db.SetMaxIdleConns(cfg.GetMaxIdleConns())
	db.SetMaxOpenConns(cfg.GetMaxOpenConns())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("[Db] database ping failed: %s", err)
	}

	if cfg.GetDebugMode() {
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
	)

	return db
}
