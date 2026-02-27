package text_obj

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"server/internal/app/config"
	"strings"
	"testing"

	domain "server/internal/app/domain/text_obj"

	"github.com/DATA-DOG/go-sqlmock"
)

func init() {
	config.InitTestConfig()
}

func TestRepository_GetByUserID(t *testing.T) {
	t.Parallel()

	t.Run("query error -> returns error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			SELECT id, user_id, title, text
			FROM text_data
			WHERE user_id = $1`

		dbErr := errors.New("db down")
		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7)).
			WillReturnError(dbErr)

		_, err = repo.GetByUserID(context.Background(), 7)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped/returned dbErr, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("scan error -> returns error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			SELECT id, user_id, title, text
			FROM text_data
			WHERE user_id = $1`

		// Неправильные колонки => Scan упадёт
		rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(1))

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7)).
			WillReturnRows(rows)

		_, err = repo.GetByUserID(context.Background(), 7)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("ok -> returns list", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			SELECT id, user_id, title, text
			FROM text_data
			WHERE user_id = $1`

		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "text"}).
			AddRow(int64(1), int64(7), "t1", "body1").
			AddRow(int64(2), int64(7), "t2", "body2")

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7)).
			WillReturnRows(rows)

		list, err := repo.GetByUserID(context.Background(), 7)
		if err != nil {
			t.Fatalf("GetByUserID error: %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2 items, got %d", len(list))
		}
		if list[0].TextId != 1 || list[0].UserId != 7 || list[0].Title != "t1" || list[0].Text != "body1" {
			t.Fatalf("unexpected item[0]: %+v", list[0])
		}
		if list[1].TextId != 2 || list[1].UserId != 7 || list[1].Title != "t2" || list[1].Text != "body2" {
			t.Fatalf("unexpected item[1]: %+v", list[1])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})
}

func TestRepository_GetByID(t *testing.T) {
	t.Parallel()

	t.Run("no rows -> ErrTextInformationNotFound", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			SELECT id, user_id, title, text
			FROM text_data
			WHERE id = $1`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(99)).
			WillReturnError(sql.ErrNoRows)

		_, err = repo.GetByID(context.Background(), 99)
		if !errors.Is(err, domain.ErrTextInformationNotFound) {
			t.Fatalf("expected ErrTextInformationNotFound, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("other error -> returned", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			SELECT id, user_id, title, text
			FROM text_data
			WHERE id = $1`

		dbErr := errors.New("db down")
		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(99)).
			WillReturnError(dbErr)

		_, err = repo.GetByID(context.Background(), 99)
		if err == nil {
			t.Fatalf("expected error")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected returned dbErr, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("ok -> returns item", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			SELECT id, user_id, title, text
			FROM text_data
			WHERE id = $1`

		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "text"}).
			AddRow(int64(5), int64(7), "hello", "world")

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(5)).
			WillReturnRows(rows)

		item, err := repo.GetByID(context.Background(), 5)
		if err != nil {
			t.Fatalf("GetByID error: %v", err)
		}
		if item.TextId != 5 || item.UserId != 7 || item.Title != "hello" || item.Text != "world" {
			t.Fatalf("unexpected item: %+v", item)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	t.Run("ok -> returns id", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			INSERT INTO text_data (user_id, title, text)
			VALUES ($1, $2, $3)
			RETURNING id`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7), "t", "body").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(123)))

		id, err := repo.Create(context.Background(), &domain.Text{
			UserId: 7,
			Title:  "t",
			Text:   "body",
		})
		if err != nil {
			t.Fatalf("Create error: %v", err)
		}
		if id != 123 {
			t.Fatalf("expected id=123, got %d", id)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("no rows -> ErrFailedCreateText", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			INSERT INTO text_data (user_id, title, text)
			VALUES ($1, $2, $3)
			RETURNING id`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7), "t", "body").
			WillReturnError(sql.ErrNoRows)

		_, err = repo.Create(context.Background(), &domain.Text{UserId: 7, Title: "t", Text: "body"})
		if !errors.Is(err, domain.ErrFailedCreateText) {
			t.Fatalf("expected ErrFailedCreateText, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("null id -> ErrFailedCreateText", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			INSERT INTO text_data (user_id, title, text)
			VALUES ($1, $2, $3)
			RETURNING id`

		mock.ExpectQuery(sqlRe(q)).
			WithArgs(int64(7), "t", "body").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(nil))

		_, err = repo.Create(context.Background(), &domain.Text{UserId: 7, Title: "t", Text: "body"})
		if !errors.Is(err, domain.ErrFailedCreateText) {
			t.Fatalf("expected ErrFailedCreateText, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	t.Run("exec error -> returns error", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			UPDATE text_data SET
			title = $1, text = $2
			WHERE id = $3`

		dbErr := errors.New("db down")
		mock.ExpectExec(sqlRe(q)).
			WithArgs("t", "body", int64(9)).
			WillReturnError(dbErr)

		err = repo.Update(context.Background(), &domain.Text{TextId: 9, Title: "t", Text: "body"})
		if err == nil {
			t.Fatalf("expected error")
		}
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected returned dbErr, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
		}
	})

	t.Run("ok -> nil", func(t *testing.T) {
		t.Parallel()

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock.New: %v", err)
		}
		defer db.Close()

		repo := &Repository{db: db}

		const q = `
			UPDATE text_data SET
			title = $1, text = $2
			WHERE id = $3`

		mock.ExpectExec(sqlRe(q)).
			WithArgs("t", "body", int64(9)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = repo.Update(context.Background(), &domain.Text{TextId: 9, Title: "t", Text: "body"})
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("expectations: %v", err)
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
