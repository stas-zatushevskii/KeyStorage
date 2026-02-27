package user

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	domain "server/internal/app/domain/user"
	"server/internal/pkg/token"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
)

func mustMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db, mock
}

func sqlRe(q string) string {
	s := strings.TrimSpace(q)
	s = strings.Join(strings.Fields(s), " ") // нормализовали любые пробелы в один пробел

	// экранируем спецсимволы regexp
	s = regexp.QuoteMeta(s)

	// пробел в regexp заменяем на \s+
	s = strings.ReplaceAll(s, `\ `, `\s+`)

	// разрешаем чтобы sqlmock матчился независимо от leading/trailing пробелов
	// и переносов/табов в фактической строке
	return `(?s)` + s
}

func TestRepository_GetById(t *testing.T) {
	ctx := context.Background()

	t.Run("no_rows -> ErrUserNotFound", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			SELECT id, username, password_hash
			FROM users
			WHERE id = $1
		`)

		mock.ExpectQuery(q).
			WithArgs(int64(10)).
			WillReturnError(sql.ErrNoRows)

		u, err := repo.GetById(ctx, 10)
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got: %v", err)
		}
		if u != nil {
			t.Fatalf("expected nil user, got: %+v", u)
		}
	})

	t.Run("other db error -> returned", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		dbErr := errors.New("db down")

		q := sqlRe(`
			SELECT id, username, password_hash
			FROM users
			WHERE id = $1
		`)

		mock.ExpectQuery(q).
			WithArgs(int64(10)).
			WillReturnError(dbErr)

		_, err := repo.GetById(ctx, 10)
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected dbErr, got: %v", err)
		}
	})

}

func TestRepository_GetByUsername(t *testing.T) {
	ctx := context.Background()

	t.Run("no_rows -> ErrUserNotFound", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			SELECT id, username, password_hash
			FROM users
			WHERE username = $1
		`)

		mock.ExpectQuery(q).
			WithArgs("nope").
			WillReturnError(sql.ErrNoRows)

		u, err := repo.GetByUsername(ctx, "nope")
		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got: %v", err)
		}
		if u != nil {
			t.Fatalf("expected nil user, got: %+v", u)
		}
	})

	t.Run("other db error -> returned", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		dbErr := errors.New("db down")

		q := sqlRe(`
			SELECT id, username, password_hash
			FROM users
			WHERE username = $1
		`)

		mock.ExpectQuery(q).
			WithArgs("john").
			WillReturnError(dbErr)

		_, err := repo.GetByUsername(ctx, "john")
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected dbErr, got: %v", err)
		}
	})

	t.Run("ok -> returns user", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			SELECT id, username, password_hash
			FROM users
			WHERE username = $1
		`)

		rows := sqlmock.NewRows([]string{"id", "username", "password_hash"}).
			AddRow(int64(7), "john", "hash123")

		mock.ExpectQuery(q).
			WithArgs("john").
			WillReturnRows(rows)

		u, err := repo.GetByUsername(ctx, "john")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if u == nil || u.ID != 7 || u.Username != "john" || u.Password != "hash123" {
			t.Fatalf("unexpected user: %+v", u)
		}
	})
}

func TestRepository_CreateNewUser(t *testing.T) {
	ctx := context.Background()

	t.Run("unique_violation -> ErrUsernameAlreadyExists", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			INSERT INTO users (username, password_hash)
			VALUES ($1, $2)
			RETURNING id
		`)

		pgErr := &pgconn.PgError{Code: "23505"}

		mock.ExpectQuery(q).
			WithArgs("john", "hash").
			WillReturnError(pgErr)

		id, err := repo.CreateNewUser(ctx, &domain.User{Username: "john", Password: "hash"})
		if !errors.Is(err, domain.ErrUsernameAlreadyExists) {
			t.Fatalf("expected ErrUsernameAlreadyExists, got: %v", err)
		}
		if id != 0 {
			t.Fatalf("expected id=0, got %d", id)
		}
	})

	t.Run("other db error -> returned", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		dbErr := errors.New("insert failed")

		q := sqlRe(`
			INSERT INTO users (username, password_hash)
			VALUES ($1, $2)
			RETURNING id
		`)

		mock.ExpectQuery(q).
			WithArgs("john", "hash").
			WillReturnError(dbErr)

		_, err := repo.CreateNewUser(ctx, &domain.User{Username: "john", Password: "hash"})
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected dbErr, got: %v", err)
		}
	})

	t.Run("ok -> returns id", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			INSERT INTO users (username, password_hash)
			VALUES ($1, $2)
			RETURNING id
		`)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(int64(123))

		mock.ExpectQuery(q).
			WithArgs("john", "hash").
			WillReturnRows(rows)

		id, err := repo.CreateNewUser(ctx, &domain.User{Username: "john", Password: "hash"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if id != 123 {
			t.Fatalf("expected id=123, got %d", id)
		}
	})
}

func TestRepository_AddTokens(t *testing.T) {
	ctx := context.Background()

	t.Run("exec error -> returned", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		dbErr := errors.New("insert tokens failed")

		q := sqlRe(`
			INSERT INTO user_tokens (user_id, refresh_token, refresh_token_expires_at, revoked_at)
			VALUES ($1, $2, $3, $4)
		`)

		tk := token.NewTokens(7)
		tk.RefreshToken = "rt"
		tk.RefreshTokenExpAt = time.Now().Add(time.Hour)

		mock.ExpectExec(q).
			WithArgs(int64(7), "rt", tk.RefreshTokenExpAt, nil).
			WillReturnError(dbErr)

		err := repo.AddTokens(ctx, 7, tk)
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected dbErr, got: %v", err)
		}
	})

	t.Run("ok -> nil", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			INSERT INTO user_tokens (user_id, refresh_token, refresh_token_expires_at, revoked_at)
			VALUES ($1, $2, $3, $4)
		`)

		tk := token.NewTokens(7)
		tk.RefreshToken = "rt"
		tk.RefreshTokenExpAt = time.Now().Add(time.Hour)

		mock.ExpectExec(q).
			WithArgs(int64(7), "rt", tk.RefreshTokenExpAt, nil).
			WillReturnResult(sqlmock.NewResult(0, 1))

		if err := repo.AddTokens(ctx, 7, tk); err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})
}

func TestRepository_UpdateTokens(t *testing.T) {
	ctx := context.Background()

	t.Run("exec error -> returned", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		dbErr := errors.New("update tokens failed")

		q := sqlRe(`
			UPDATE user_tokens
			SET 
				refresh_token = $2,
				refresh_token_expires_at = $3,
				revoked_at = $4
			WHERE user_id = $1
		`)

		tk := token.NewTokens(7)
		tk.RefreshToken = "rt2"
		tk.RefreshTokenExpAt = time.Now().Add(2 * time.Hour)

		mock.ExpectExec(q).
			WithArgs(int64(7), "rt2", tk.RefreshTokenExpAt, nil).
			WillReturnError(dbErr)

		err := repo.UpdateTokens(ctx, 7, tk)
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected dbErr, got: %v", err)
		}
	})

	t.Run("ok -> nil", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			UPDATE user_tokens
			SET 
				refresh_token = $2,
				refresh_token_expires_at = $3,
				revoked_at = $4
			WHERE user_id = $1
		`)

		tk := token.NewTokens(7)
		tk.RefreshToken = "rt2"
		tk.RefreshTokenExpAt = time.Now().Add(2 * time.Hour)

		mock.ExpectExec(q).
			WithArgs(int64(7), "rt2", tk.RefreshTokenExpAt, nil).
			WillReturnResult(sqlmock.NewResult(0, 1))

		if err := repo.UpdateTokens(ctx, 7, tk); err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
	})
}

func TestRepository_GetTokens(t *testing.T) {
	ctx := context.Background()

	t.Run("no_rows -> ErrRefreshTokenNotFound", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		q := sqlRe(`
			SELECT refresh_token, refresh_token_expires_at, revoked_at
			FROM user_tokens
			WHERE user_id = $1
		`)

		mock.ExpectQuery(q).
			WithArgs(int64(7)).
			WillReturnError(sql.ErrNoRows)

		tk, err := repo.GetTokens(ctx, 7)
		if !errors.Is(err, domain.ErrRefreshTokenNotFound) {
			t.Fatalf("expected ErrRefreshTokenNotFound, got: %v", err)
		}
		if tk != nil {
			t.Fatalf("expected nil tokens, got: %+v", tk)
		}
	})

	t.Run("other db error -> returned", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		dbErr := errors.New("select failed")

		q := sqlRe(`
			SELECT refresh_token, refresh_token_expires_at, revoked_at
			FROM user_tokens
			WHERE user_id = $1
		`)

		mock.ExpectQuery(q).
			WithArgs(int64(7)).
			WillReturnError(dbErr)

		_, err := repo.GetTokens(ctx, 7)
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected dbErr, got: %v", err)
		}
	})

	t.Run("ok not revoked -> Revoked=false", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		exp := time.Now().Add(time.Hour)

		q := sqlRe(`
			SELECT refresh_token, refresh_token_expires_at, revoked_at
			FROM user_tokens
			WHERE user_id = $1
		`)

		rows := sqlmock.NewRows([]string{"refresh_token", "refresh_token_expires_at", "revoked_at"}).
			AddRow("hashed-rt", exp, sql.NullInt64{Valid: false})

		mock.ExpectQuery(q).
			WithArgs(int64(7)).
			WillReturnRows(rows)

		tk, err := repo.GetTokens(ctx, 7)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if tk == nil || tk.RefreshToken != "hashed-rt" || tk.RefreshTokenExpAt.IsZero() {
			t.Fatalf("unexpected tokens: %+v", tk)
		}
		if tk.Revoked {
			t.Fatalf("expected Revoked=false, got true")
		}
	})

	t.Run("ok revoked -> Revoked=true", func(t *testing.T) {
		db, mock := mustMockDB(t)
		repo := &Repository{db: db}

		exp := time.Now().Add(time.Hour)

		q := sqlRe(`
			SELECT refresh_token, refresh_token_expires_at, revoked_at
			FROM user_tokens
			WHERE user_id = $1
		`)

		rows := sqlmock.NewRows([]string{"refresh_token", "refresh_token_expires_at", "revoked_at"}).
			AddRow("hashed-rt", exp, sql.NullInt64{Int64: 123, Valid: true})

		mock.ExpectQuery(q).
			WithArgs(int64(7)).
			WillReturnRows(rows)

		tk, err := repo.GetTokens(ctx, 7)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if tk == nil {
			t.Fatalf("expected tokens, got nil")
		}
		if !tk.Revoked {
			t.Fatalf("expected Revoked=true, got false")
		}
	})
}
