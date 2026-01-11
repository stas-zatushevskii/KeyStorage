package bank_card_obj

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"server/internal/app/config"
	domain "server/internal/app/domain/bank_card_obj"
	"server/internal/pkg/encryption/aes"
)

func (u *Repository) GetByUserID(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	query := `
		SELECT id, user_id, bank_name, pid
		FROM bank_data
		WHERE user_id = $1`

	var cards []*domain.BankCard

	rows, err := u.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		obj := new(Card)

		if err := rows.Scan(&obj.ID, &obj.UserId, &obj.Bank, &obj.PID); err != nil {
			return nil, err
		}

		if obj.PID.Valid {
			decrypted, err := aes.DecryptAES([]byte(obj.PID.String), []byte(config.App.GetBankCardObjEncryptionKey()))
			if err != nil {
				return nil, fmt.Errorf("failed decrypt bank card: %w", err)
			}
			obj.PID.String = string(decrypted)
		}

		cards = append(cards, obj.ToDomain())
	}
	return cards, nil
}

func (u *Repository) GetByID(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	query := `
		SELECT id, user_id, bank_name, pid
		FROM bank_data
		WHERE id = $1`

	obj := new(Card)

	if err := u.db.QueryRowContext(ctx, query, cardId).Scan(&obj.ID, &obj.UserId, &obj.Bank, &obj.PID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrBankCardInformationNotFound
		}

		decrypted, err := aes.DecryptAES([]byte(obj.PID.String), []byte(config.App.GetBankCardObjEncryptionKey()))
		if err != nil {
			return nil, fmt.Errorf("failed decrypt bank card: %w", err)
		}
		obj.PID.String = string(decrypted)
		return nil, err
	}

	return obj.ToDomain(), nil
}

func (u *Repository) Create(ctx context.Context, card *domain.BankCard) (int64, error) {
	query := `
		INSERT INTO bank_data (user_id, bank_name, pid)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	encryptedPassword, err := aes.EncryptAES([]byte(card.Pid), []byte(config.App.GetBankCardObjEncryptionKey()))
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt PID: %w", err)
	}

	card.Pid = string(encryptedPassword)

	if _, err := u.db.ExecContext(ctx, query, card.UserId, card.Bank, card.Pid); err != nil {
		return 0, err
	}

	var id sql.NullInt64

	if err := u.db.QueryRowContext(ctx, query).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrFaildeCreateBankCardObject
		}
		return 0, err
	}

	if !id.Valid {
		return 0, domain.ErrFaildeCreateBankCardObject
	}

	return id.Int64, nil
}

func (u *Repository) Update(ctx context.Context, card *domain.BankCard) error {
	query := `
		UPDATE bank_data SET
		bank_name = $1, pid = $2
		WHERE id = $3`

	encryptedPassword, err := aes.EncryptAES([]byte(card.Bank), []byte(config.App.GetBankCardObjEncryptionKey()))
	if err != nil {
		return fmt.Errorf("failed to encrypt PID: %w", err)
	}

	card.Pid = string(encryptedPassword)

	if _, err := u.db.ExecContext(ctx, query, card.Bank, card.Pid, card.Pid, card.CardId); err != nil {
		return err
	}
	return nil
}
