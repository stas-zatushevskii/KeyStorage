package account_obj

import (
	"context"
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
	if userId <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	list, err := a.repo.GetByUserID(ctx, userId)
	if err != nil {
		return nil, domain.ErrAccountNotFound
	}

	if len(list) == 0 {
		return nil, domain.ErrEmptyAccountsList
	}

	return list, nil
}

func (a *AccountObj) GetAccount(ctx context.Context, accountId int64) (*domain.Account, error) {
	if accountId <= 0 {
		return nil, domain.ErrInvalidAccountID
	}

	account, err := a.repo.GetByID(ctx, accountId)
	if err != nil {
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}

func (a *AccountObj) CreateNewAccountObj(ctx context.Context, account *domain.Account) (int64, error) {
	if account.UserId <= 0 {
		return 0, domain.ErrInvalidUserID
	}

	if account.ServiceName == "" {
		return 0, domain.ErrEmptyServiceName
	}

	id, err := a.repo.Create(ctx, account)
	if err != nil {
		return 0, domain.ErrFailedCreateAccount
	}

	return id, nil
}

func (a *AccountObj) UpdateAccount(ctx context.Context, account *domain.Account) error {
	if account.AccountId <= 0 {
		return domain.ErrInvalidAccountID
	}

	if account.ServiceName == "" {
		return domain.ErrEmptyServiceName
	}

	if err := a.repo.Update(ctx, account); err != nil {
		return domain.ErrFailedUpdateAccount
	}

	return nil
}
