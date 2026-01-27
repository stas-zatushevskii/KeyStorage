package file_obj

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "server/internal/app/domain/file_obj"

	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Create(ctx context.Context, f *domain.File) (int64, error) {
	if f == nil {
		return 0, fmt.Errorf("file is nil")
	}

	query := `
		INSERT INTO file_data (
			user_id, title,
			bucket_name, object_key,
			size_bytes, content_type, etag
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at
	`

	var (
		id        int64
		createdAt = f.CreatedAt
	)

	err := r.db.QueryRowContext(ctx, query,
		f.UserID,
		nullIfEmpty(f.Title),
		f.Storage.BucketName,
		f.Storage.ObjectKey,
		f.SizeBytes,
		nullIfEmpty(f.ContentType),
		nullIfEmpty(f.ETag),
	).Scan(&id, &createdAt)

	if err != nil {
		if isUniqueViolation(err, "uq_file_object") {
			return 0, fmt.Errorf(
				"file already exists in storage (bucket=%s key=%s): %w",
				f.Storage.BucketName, f.Storage.ObjectKey, err,
			)
		}
		if isFKViolation(err, "fk_file_data_user") {
			return 0, fmt.Errorf("user not found (user_id=%d): %w", f.UserID, err)
		}
		return 0, fmt.Errorf("insert file_data: %w", err)
	}

	f.ID = id
	f.CreatedAt = createdAt

	return f.ID, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*domain.File, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidFileID
	}

	const q = `
		SELECT
			id, user_id, title,
			bucket_name, object_key,
			size_bytes, content_type, etag,
			created_at
		FROM file_data
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, q, id)

	f, err := scanFile(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrFileNotFound
		}
		return nil, fmt.Errorf("select file_data by id=%d: %w", id, err)
	}

	return f, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID int64) ([]*domain.File, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	query := `
		SELECT
			id, user_id, title,
			bucket_name, object_key,
			size_bytes, content_type, etag,
			created_at
		FROM file_data
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list file_data by user_id=%d: %w", userID, err)
	}
	defer rows.Close()

	var out []*domain.File
	for rows.Next() {
		f, err := scanFile(rows)
		if err != nil {
			return nil, fmt.Errorf("scan file_data row: %w", err)
		}
		out = append(out, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return out, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidFileID
	}

	query := `DELETE FROM file_data WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, int64(id))
	if err != nil {
		return fmt.Errorf("delete file_data id=%d: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete file_data id=%d: rows affected: %w", id, err)
	}

	if affected == 0 {
		return domain.ErrFileNotFound
	}

	return nil
}

// help func

type scanner interface {
	Scan(dest ...any) error
}

func scanFile(s scanner) (*domain.File, error) {
	var (
		id         int64
		userID     int64
		title      sql.NullString
		bucketName string
		objectKey  string
		sizeBytes  int64
		ct         sql.NullString
		etag       sql.NullString
		createdAt  sql.NullTime
	)

	err := s.Scan(
		&id, &userID, &title,
		&bucketName, &objectKey,
		&sizeBytes, &ct, &etag,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	ref, err := domain.NewStorageRef(bucketName, objectKey)
	if err != nil {
		return nil, err
	}

	f := &domain.File{
		ID:          id,
		UserID:      userID,
		Title:       nullStringToString(title),
		Storage:     ref,
		SizeBytes:   sizeBytes,
		ContentType: nullStringToString(ct),
		ETag:        nullStringToString(etag),
	}

	if createdAt.Valid {
		f.CreatedAt = createdAt.Time
	}

	return f, nil
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func nullStringToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23505 = unique_violation
		return pgErr.Code == "23505" && pgErr.ConstraintName == constraint
	}
	return false
}

func isFKViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23503 = foreign_key_violatio
		return pgErr.Code == "23503" && pgErr.ConstraintName == constraint
	}
	return false
}
