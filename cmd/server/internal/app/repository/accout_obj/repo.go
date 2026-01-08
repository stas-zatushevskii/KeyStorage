package account_obj

import (
	"context"
	"database/sql"
	"errors"
	domain "server/internal/app/domain/account_obj"
)

func (u *Repository) GetByUserID(ctx context.Context, userId int64) ([]*domain.Account, error) {
	query := `
		SELECT id, service_name, username, user_id, password
		FROM account_data
		WHERE user_id = $1`

	var accounts []*domain.Account

	rows, err := u.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		obj := new(Account)

		if err := rows.Scan(&obj.ID, &obj.ServiceName, &obj.UserName, &obj.UserId, &obj.Password); err != nil {
			return nil, err
		}

		if obj.UserId.Valid && obj.ID.Valid {
			accounts = append(accounts, obj.ToDomain())
		}
	}
	return accounts, nil
}

func (u *Repository) GetByID(ctx context.Context, accountId int64) (*domain.Account, error) {
	query := `
		SELECT id, service_name, username, user_id, password
		FROM account_data
		WHERE id = $1`

	obj := new(Account)

	if err := u.db.QueryRowContext(ctx, query, accountId).Scan(&obj.ID, &obj.ServiceName, &obj.UserName, &obj.UserId, &obj.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountInformationNotFound
		}
		return nil, err
	}

	return obj.ToDomain(), nil
}

func (u *Repository) Create(ctx context.Context, account *domain.Account) (int64, error) {
	query := `
		INSERT INTO account_data (user_id, service_name, username, user_id, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	if _, err := u.db.ExecContext(ctx, query, account.UserId, account.ServiceName, account.UserName, account.UserId, account.Password); err != nil {
		return 0, err
	}
	return account.AccountId, nil
}

func (u *Repository) Update(ctx context.Context, account *domain.Account) error {
	query := `
		UPDATE account_data SET
		service_name = $1, username = $2, password = $3
		WHERE id = $4`

	if _, err := u.db.ExecContext(ctx, query, account.ServiceName, account.UserName, account.Password, account.AccountId); err != nil {
		return err
	}
	return nil
}
