package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"server/internal/app/config"
	"server/internal/pkg/logger"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/sync/errgroup"

	// -- only for debug
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

type DatabaseAdapter struct {
	DB *sql.DB
}

func New() (*DatabaseAdapter, error) {

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

	return &DatabaseAdapter{DB: db}, nil
}

func (db *DatabaseAdapter) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {

		err := db.RunMigrations()
		if err != nil {
			return fmt.Errorf("failed to setup database: %v", err)
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		err := db.DB.Close()
		if err != nil {
			return err
		}

		return nil
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	return nil
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
