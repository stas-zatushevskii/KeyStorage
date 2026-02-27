package file

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"server/internal/app/config"
	"strings"
	"testing"
	"time"

	domain "server/internal/app/domain/file_obj"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5/pgconn"
)

func init() {
	config.InitTestConfig()
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("nil file -> error", func(t *testing.T) {
		t.Parallel()

		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		r := &Repository{db: db}
		_, err = r.Create(context.Background(), nil)
		if err == nil || !strings.Contains(err.Error(), "file is nil") {
			t.Fatalf("expected 'file is nil' error, got: %v", err)
		}
	})

	t.Run("ok -> sets ID and CreatedAt", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		r := &Repository{db: db}

		now := time.Now().UTC().Truncate(time.Microsecond)

		f := &domain.File{
			UserID: 7,
			Title:  "t",
			Storage: domain.StorageRef{
				BucketName: "bucket",
				ObjectKey:  "key",
			},
			SizeBytes:   10,
			ContentType: "text/plain",
			ETag:        "etag",
		}

		const q = `
			INSERT INTO file_data (
				user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, created_at
		`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(
				int64(7),
				"t",
				"bucket",
				"key",
				int64(10),
				"text/plain",
				"etag",
			).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(int64(123), now))

		id, err := r.Create(context.Background(), f)
		if err != nil {
			t.Fatalf("Create error: %v", err)
		}
		if id != 123 {
			t.Fatalf("expected id=123, got %d", id)
		}
		if f.ID != 123 {
			t.Fatalf("expected f.ID=123, got %d", f.ID)
		}
		if f.CreatedAt.IsZero() {
			t.Fatalf("expected CreatedAt to be set")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("unique violation uq_file_object -> wrapped message", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		r := &Repository{db: db}

		f := &domain.File{
			UserID: 7,
			Storage: domain.StorageRef{
				BucketName: "b",
				ObjectKey:  "k",
			},
			SizeBytes:   1,
			ContentType: "x",
		}

		pgErr := &pgconn.PgError{Code: "23505", ConstraintName: "uq_file_object"}

		const q = `
			INSERT INTO file_data (
				user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, created_at
		`

		mock.ExpectQuery(sqlRe(q)).
			WillReturnError(pgErr)

		_, err = r.Create(context.Background(), f)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.As(err, &pgErr) {
			t.Fatalf("expected pg error in chain, got: %v", err)
		}
		if !strings.Contains(err.Error(), "file already exists in storage") {
			t.Fatalf("expected unique violation message, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("fk violation fk_file_data_user -> wrapped message", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		r := &Repository{db: db}

		f := &domain.File{
			UserID: 999,
			Storage: domain.StorageRef{
				BucketName: "b",
				ObjectKey:  "k",
			},
			SizeBytes:   1,
			ContentType: "x",
		}

		pgErr := &pgconn.PgError{Code: "23503", ConstraintName: "fk_file_data_user"}

		const q = `
			INSERT INTO file_data (
				user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, created_at
		`

		mock.ExpectQuery(sqlRe(q)).
			WillReturnError(pgErr)

		_, err = r.Create(context.Background(), f)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "user not found") {
			t.Fatalf("expected fk violation message, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("other error -> wrapped insert file_data", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		r := &Repository{db: db}

		f := &domain.File{
			UserID: 7,
			Storage: domain.StorageRef{
				BucketName: "b",
				ObjectKey:  "k",
			},
			SizeBytes:   1,
			ContentType: "x",
		}

		dbErr := errors.New("db down")

		const q = `
			INSERT INTO file_data (
				user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, created_at
		`

		mock.ExpectQuery(sqlRe(q)).WillReturnError(dbErr)

		_, err = r.Create(context.Background(), f)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped dbErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "insert file_data") {
			t.Fatalf("expected context, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})
}

func TestRepository_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("invalid id -> ErrInvalidFileID", func(t *testing.T) {
		t.Parallel()

		db, _, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}
		_, err := r.GetByID(context.Background(), 0)
		if !errors.Is(err, domain.ErrInvalidFileID) {
			t.Fatalf("expected ErrInvalidFileID, got: %v", err)
		}
	})

	t.Run("no rows -> ErrFileNotFound", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		const q = `
			SELECT
				id, user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag,
				created_at
			FROM file_data
			WHERE id = $1
		`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(10)).
			WillReturnError(sql.ErrNoRows)

		_, err := r.GetByID(context.Background(), 10)
		if !errors.Is(err, domain.ErrFileNotFound) {
			t.Fatalf("expected ErrFileNotFound, got: %v", err)
		}
	})

	t.Run("other error -> wrapped select file_data by id", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		dbErr := errors.New("db down")

		const q = `
			SELECT
				id, user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag,
				created_at
			FROM file_data
			WHERE id = $1
		`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(10)).
			WillReturnError(dbErr)

		_, err := r.GetByID(context.Background(), 10)
		if err == nil {
			t.Fatalf("expected error")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped dbErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "select file_data by id=10") {
			t.Fatalf("expected context, got: %v", err)
		}
	})

	t.Run("ok -> returns file", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		now := time.Now().UTC()

		const q = `
			SELECT
				id, user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag,
				created_at
			FROM file_data
			WHERE id = $1
		`

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "title",
			"bucket_name", "object_key",
			"size_bytes", "content_type", "etag",
			"created_at",
		}).AddRow(
			int64(1), int64(7), "title",
			"b", "k",
			int64(100), "text/plain", "etag",
			now,
		)

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(1)).
			WillReturnRows(rows)

		f, err := r.GetByID(context.Background(), 1)
		if err != nil {
			t.Fatalf("GetByID error: %v", err)
		}
		if f.ID != 1 || f.UserID != 7 {
			t.Fatalf("unexpected file: %+v", f)
		}
		if f.Storage.BucketName != "b" || f.Storage.ObjectKey != "k" {
			t.Fatalf("unexpected storage: %+v", f.Storage)
		}
	})
}

func TestRepository_ListByUserID(t *testing.T) {
	t.Parallel()

	t.Run("invalid user -> ErrInvalidUserID", func(t *testing.T) {
		t.Parallel()

		db, _, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}
		_, err := r.ListByUserID(context.Background(), 0)
		if !errors.Is(err, domain.ErrInvalidUserID) {
			t.Fatalf("expected ErrInvalidUserID, got: %v", err)
		}
	})

	t.Run("query error -> wrapped", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}
		dbErr := errors.New("db down")

		const q = `
			SELECT
				id, user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag,
				created_at
			FROM file_data
			WHERE user_id = $1
			ORDER BY created_at DESC, id DESC
		`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7)).
			WillReturnError(dbErr)

		_, err := r.ListByUserID(context.Background(), 7)
		if err == nil {
			t.Fatalf("expected error")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped dbErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "list file_data by user_id=7") {
			t.Fatalf("expected context, got: %v", err)
		}
	})

	t.Run("scan error -> wrapped", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		const q = `
			SELECT
				id, user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag,
				created_at
			FROM file_data
			WHERE user_id = $1
			ORDER BY created_at DESC, id DESC
		`

		// missing required columns -> Scan will error
		rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(1))

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7)).
			WillReturnRows(rows)

		_, err := r.ListByUserID(context.Background(), 7)
		if err == nil {
			t.Fatalf("expected error")
		}
		if !strings.Contains(err.Error(), "scan file_data row") {
			t.Fatalf("expected scan context, got: %v", err)
		}
	})

	t.Run("rows.Err -> wrapped", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		const q = `
			SELECT
				id, user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag,
				created_at
			FROM file_data
			WHERE user_id = $1
			ORDER BY created_at DESC, id DESC
		`

		rowErr := errors.New("rows fail")

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "title",
			"bucket_name", "object_key",
			"size_bytes", "content_type", "etag",
			"created_at",
		}).AddRow(
			int64(1), int64(7), "t",
			"b", "k",
			int64(1), "ct", "etag",
			time.Now(),
		).RowError(0, rowErr)

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7)).
			WillReturnRows(rows)

		_, err := r.ListByUserID(context.Background(), 7)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "rows err:") {
			t.Fatalf("expected rows err context, got: %v", err)
		}
	})

	t.Run("ok -> returns list", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		now := time.Now().UTC()

		const q = `
			SELECT
				id, user_id, title,
				bucket_name, object_key,
				size_bytes, content_type, etag,
				created_at
			FROM file_data
			WHERE user_id = $1
			ORDER BY created_at DESC, id DESC
		`

		rows := sqlmock.NewRows([]string{
			"id", "user_id", "title",
			"bucket_name", "object_key",
			"size_bytes", "content_type", "etag",
			"created_at",
		}).
			AddRow(int64(2), int64(7), "t2", "b", "k2", int64(2), "ct", "e2", now).
			AddRow(int64(1), int64(7), "t1", "b", "k1", int64(1), "ct", "e1", now.Add(-time.Minute))

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7)).
			WillReturnRows(rows)

		list, err := r.ListByUserID(context.Background(), 7)
		if err != nil {
			t.Fatalf("ListByUserID error: %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2, got %d", len(list))
		}
		if list[0].ID != 2 || list[1].ID != 1 {
			t.Fatalf("unexpected order/ids: %+v %+v", list[0], list[1])
		}
	})
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	t.Run("invalid id -> ErrInvalidFileID", func(t *testing.T) {
		t.Parallel()

		db, _, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}
		err := r.Delete(context.Background(), 0)
		if !errors.Is(err, domain.ErrInvalidFileID) {
			t.Fatalf("expected ErrInvalidFileID, got: %v", err)
		}
	})

	t.Run("exec error -> wrapped", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}
		dbErr := errors.New("db down")

		mock.ExpectExec(sqlRe(`DELETE FROM file_data WHERE id = $1`)).
			WithArgs(int64(10)).
			WillReturnError(dbErr)

		err := r.Delete(context.Background(), 10)
		if err == nil {
			t.Fatalf("expected error")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped dbErr, got: %v", err)
		}
		if !strings.Contains(err.Error(), "delete file_data id=10") {
			t.Fatalf("expected context, got: %v", err)
		}
	})

	t.Run("rows affected error -> wrapped", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		// sqlmock умеет имитировать ошибку RowsAffected через ResultError
		mock.ExpectExec(sqlRe(`DELETE FROM file_data WHERE id = $1`)).
			WithArgs(int64(10)).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected fail")))

		err := r.Delete(context.Background(), 10)
		if err == nil {
			t.Fatalf("expected error")
		}
		if !strings.Contains(err.Error(), "rows affected") {
			t.Fatalf("expected rows affected context, got: %v", err)
		}
	})

	t.Run("affected=0 -> ErrFileNotFound", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		mock.ExpectExec(sqlRe(`DELETE FROM file_data WHERE id = $1`)).
			WithArgs(int64(10)).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := r.Delete(context.Background(), 10)
		if !errors.Is(err, domain.ErrFileNotFound) {
			t.Fatalf("expected ErrFileNotFound, got: %v", err)
		}
	})

	t.Run("ok -> nil", func(t *testing.T) {
		t.Parallel()

		db, mock, _ := sqlmock.New()
		defer db.Close()

		r := &Repository{db: db}

		mock.ExpectExec(sqlRe(`DELETE FROM file_data WHERE id = $1`)).
			WithArgs(int64(10)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := r.Delete(context.Background(), 10)
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}
	})
}

func sqlRe(q string) string {
	// 1) normalize whitespace (как выглядит actual sql в ошибке sqlmock)
	s := strings.TrimSpace(q)
	s = strings.Join(strings.Fields(s), " ")

	// 2) escape ALL regexp metacharacters: $, (, ), etc.
	s = regexp.QuoteMeta(s)

	// 3) make spaces flexible
	s = strings.ReplaceAll(s, `\ `, `\s+`)
	return s
}
