package file_obj

import (
	"context"
	"database/sql"
	"io"
)

type Storage interface {
	PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64, contentType string) (etag string, err error)

	DeleteObject(ctx context.Context, bucket, key string) error
}

type Repository struct {
	db      *sql.DB
	storage Storage
}

func New(db *sql.DB, storage Storage) *Repository {
	return &Repository{db: db, storage: storage}
}
