package bank_card

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"

	"server/internal/app/config"
	domain "server/internal/app/domain/bank_card"
	"server/internal/pkg/encryption/aes"

	"github.com/DATA-DOG/go-sqlmock"
)

func init() {
	config.InitTestConfig()
}

func requireBankCardEncKey(t *testing.T) string {
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

func TestRepository_GetByID_OK_DecryptsPID(t *testing.T) {
	t.Parallel()

	key := requireBankCardEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	plainPID := "PID-123"
	enc, err := aes.EncryptAES([]byte(plainPID), []byte(key))
	if err != nil {
		t.Fatalf("EncryptAES error: %v", err)
	}

	const q = `
		SELECT id, user_id, bank_name, pid
		FROM bank_data
		WHERE id = $1`

	rows := sqlmock.NewRows([]string{"id", "user_id", "bank_name", "pid"}).
		AddRow(int64(10), int64(7), "maib", enc)

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(10)).
		WillReturnRows(rows)

	got, err := repo.GetByID(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected card, got nil")
	}
	if got.CardId != 10 || got.UserId != 7 {
		t.Fatalf("unexpected ids: %+v", got)
	}
	if got.Pid != plainPID {
		t.Fatalf("expected decrypted pid=%q, got %q", plainPID, got.Pid)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_GetByID_NoRows_ReturnsNotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	const q = `
		SELECT id, user_id, bank_name, pid
		FROM bank_data
		WHERE id = $1`

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByID(context.Background(), 999)
	if !errors.Is(err, domain.ErrBankCardNotFound) {
		t.Fatalf("expected ErrBankCardNotFound, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_Create_OK_ReturnsID(t *testing.T) {
	t.Parallel()

	_ = requireBankCardEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	const q = `
		INSERT INTO bank_data (user_id, bank_name, pid)
		VALUES ($1, $2, $3)
		RETURNING id`

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(7), "maib", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(123)))

	id, err := repo.Create(context.Background(), &domain.BankCard{
		UserId: 7,
		Bank:   "maib",
		Pid:    "PID-123",
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
func TestRepository_Update_OK(t *testing.T) {
	t.Parallel()

	_ = requireBankCardEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	const q = `
		UPDATE bank_data SET
		bank_name = $1, pid = $2
		WHERE id = $3`

	mock.ExpectExec(sqlRe(q)).
		WithArgs("maib", sqlmock.AnyArg(), int64(55)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), &domain.BankCard{
		CardId: 55,
		Bank:   "maib",
		Pid:    "PID-NEW",
	})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sql expectations: %v", err)
	}
}

func TestRepository_GetByUserID_OK_ListDecryptsPID(t *testing.T) {
	t.Parallel()

	key := requireBankCardEncKey(t)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New error: %v", err)
	}
	defer db.Close()

	repo := &Repository{db: db}

	enc1, err := aes.EncryptAES([]byte("PID1"), []byte(key))
	if err != nil {
		t.Fatalf("EncryptAES error: %v", err)
	}
	enc2, err := aes.EncryptAES([]byte("PID2"), []byte(key))
	if err != nil {
		t.Fatalf("EncryptAES error: %v", err)
	}

	const q = `
		SELECT id, user_id, bank_name, pid
		FROM bank_data
		WHERE user_id = $1`

	rows := sqlmock.NewRows([]string{"id", "user_id", "bank_name", "pid"}).
		AddRow(int64(1), int64(7), "maib", enc1).
		AddRow(int64(2), int64(7), "victoriabank", enc2)

	mock.ExpectQuery(sqlRe(q)).
		WithArgs(int64(7)).
		WillReturnRows(rows)

	list, err := repo.GetByUserID(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetByUserID error: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(list))
	}
	if list[0].Pid != "PID1" {
		t.Fatalf("expected decrypted PID1, got %q", list[0].Pid)
	}
	if list[1].Pid != "PID2" {
		t.Fatalf("expected decrypted PID2, got %q", list[1].Pid)
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
