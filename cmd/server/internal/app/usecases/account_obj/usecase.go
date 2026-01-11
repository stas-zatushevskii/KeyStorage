package account_obj

import (
	"context"
	"fmt"
	domain "server/internal/app/domain/account_obj"
)

type Repository interface {
	GetByUserID(ctx context.Context, userId int64) ([]*domain.Account, error)
	GetByID(ctx context.Context, accountId int64) (*domain.Account, error)
	Create(ctx context.Context, account *domain.Account) (int64, error)
	Update(ctx context.Context, account *domain.Account) error
}

type AccountObj struct {
	repo Repository
}

func New(repo Repository) *AccountObj {
	return &AccountObj{repo: repo}
}

func (a *AccountObj) GetAccountsList(ctx context.Context, userId int64) ([]*domain.Account, error) {
	list, err := a.repo.GetByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts list: %w", err)
	}
	if len(list) == 0 {
		return nil, domain.ErrEmptyAccountsList
	}
	return list, nil
}

func (a *AccountObj) GetAccount(ctx context.Context, accountId int64) (*domain.Account, error) {
	account, err := a.repo.GetByID(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return account, nil
}

func (a *AccountObj) CreateNewAccountObj(ctx context.Context, account *domain.Account) (int64, error) {

	if account.UserId == 0 {
		return 0, fmt.Errorf("user id is zero")
	}

	if account.ServiceName == "" {
		return 0, fmt.Errorf("service name is zero")
	}

	id, err := a.repo.Create(ctx, account)
	if err != nil {
		return 0, fmt.Errorf("failed to create account: %w", err)
	}

	return id, nil
}

func (a *AccountObj) UpdateAccount(ctx context.Context, account *domain.Account) error {

	if account.AccountId == 0 {
		return fmt.Errorf("account id is zero")
	}

	if account.ServiceName == "" {
		return fmt.Errorf("service name cant be empty")
	}

	err := a.repo.Update(ctx, account)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}
