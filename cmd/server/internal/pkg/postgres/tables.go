package db

import (
	"github.com/pressly/goose/v3"
)

func (db *DatabaseAdapter) RunMigrations() error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db.DB, "./migrations")
}
