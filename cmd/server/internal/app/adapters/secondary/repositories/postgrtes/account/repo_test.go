package account_obj

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"

	"server/internal/app/config"
	domain "server/internal/app/domain/account"
	"server/internal/pkg/encryption/aes"

	"github.com/DATA-DOG/go-sqlmock"
)

func init() {
	config.InitTestConfig()
}

func requireAccountEncKey(t *testing.T) string {
	t.Helper()

	if config.App == nil {
		t.Skip("config.App is nil (config not initialized in tests)")
	}

	key := config.App.GetAccountObjEncryptionKey()
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		t.Skipf(
			"account encryption key is not initialized or invalid length (len=%d). Configure config.App.GetAccountObjEncryptionKey() for tests",
			len(key),
		)
	}
	return key
}

func TestRepository_GetByID_OK_DecryptsPassword(t *testing.T) {
	t.Parallel()

	key := requireAccountEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	plaintext := "my-pass"
	enc, err := aes.EncryptAES([]byte(plaintext), []byte(key))
	if err != nil {
		t.Fatalf("EncryptAES error: %v", err)
	}
	encStr := string(enc)

	const q = `
		SELECT id, service_name, username, user_id, password
		FROM account_data
		WHERE id = $1`

	rows := sqlmock.NewRows([]string{"id", "service_name", "username", "user_id", "password"}).
		AddRow(int64(10), "telegram", "stas", int64(7), encStr)

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(10)).
		WillReturnRows(rows)

	got, err := repo.GetByID(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected account, got nil")
	}
	if got.AccountId != 10 || got.UserId != 7 {
		t.Fatalf("unexpected ids: %+v", got)
	}
	if got.Password != plaintext {
		t.Fatalf("expected decrypted password=%q, got %q", plaintext, got.Password)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_GetByID_NoRows_ReturnsDomainNotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	const q = `
		SELECT id, service_name, username, user_id, password
		FROM account_data
		WHERE id = $1`

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByID(context.Background(), 999)
	if !errors.Is(err, domain.ErrAccountInformationNotFound) {
		t.Fatalf("expected ErrAccountInformationNotFound, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_Create_OK_ReturnsID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	const q = `
		INSERT INTO account_data (user_id, service_name, username, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(7), "telegram", "stas", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(123)))

	id, err := repo.Create(context.Background(), &domain.Account{
		UserId:      7,
		ServiceName: "telegram",
		UserName:    "stas",
		Password:    "my-pass",
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if id != 123 {
		t.Fatalf("expected id=123, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_Create_InvalidID_ReturnsDomainFailed(t *testing.T) {
	t.Parallel()

	_ = requireAccountEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	const q = `
		INSERT INTO account_data (user_id, service_name, username, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(7), "telegram", "stas", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(nil))

	_, err = repo.Create(context.Background(), &domain.Account{
		UserId:      7,
		ServiceName: "telegram",
		UserName:    "stas",
		Password:    "my-pass",
	})
	if !errors.Is(err, domain.ErrFaildeCreateAccountObject) {
		t.Fatalf("expected ErrFaildeCreateAccountObject, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_Update_OK(t *testing.T) {
	t.Parallel()

	_ = requireAccountEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	const q = `
		UPDATE account_data SET
		service_name = $1, username = $2, password = $3
		WHERE id = $4`

	mock.ExpectExec(sqlRe(q)).
		WithArgs("telegram", "stas", sqlmock.AnyArg(), int64(55)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), &domain.Account{
		AccountId:   55,
		ServiceName: "telegram",
		UserName:    "stas",
		Password:    "new-pass",
	})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_GetByUserID_OK_ListDecryptsPasswords(t *testing.T) {
	t.Parallel()

	key := requireAccountEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	enc1, err := aes.EncryptAES([]byte("p1"), []byte(key))
	if err != nil {
		t.Fatalf("EncryptAES error: %v", err)
	}
	enc2, err := aes.EncryptAES([]byte("p2"), []byte(key))
	if err != nil {
		t.Fatalf("EncryptAES error: %v", err)
	}

	const q = `
		SELECT id, service_name, username, user_id, password
		FROM account_data
		WHERE user_id = $1`

	rows := sqlmock.NewRows([]string{"id", "service_name", "username", "user_id", "password"}).
		AddRow(int64(1), "telegram", "u1", int64(7), string(enc1)).
		AddRow(int64(2), "shopify", "u2", int64(7), string(enc2))

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(7)).
		WillReturnRows(rows)

	list, err := repo.GetByUserID(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetByUserID error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(list))
	}
	if list[0].Password != "p1" {
		t.Fatalf("expected decrypted p1, got %q", list[0].Password)
	}
	if list[1].Password != "p2" {
		t.Fatalf("expected decrypted p2, got %q", list[1].Password)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
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
