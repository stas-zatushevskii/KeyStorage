package bank_card_obj

import (
	"context"
	"fmt"
	domain "server/internal/app/domain/bank_card_obj"
)

type Repository interface {
	GetByUserID(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	GetByID(ctx context.Context, cardId int64) (*domain.BankCard, error)
	Create(ctx context.Context, card *domain.BankCard) (int64, error)
	Update(ctx context.Context, card *domain.BankCard) error
}

type BankCardObj struct {
	repo Repository
}

func New(repo Repository) *BankCardObj {
	return &BankCardObj{repo: repo}
}

func (b *BankCardObj) GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	account, err := b.repo.GetByID(ctx, cardId)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	return account, nil
}

func (b *BankCardObj) GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	list, err := b.repo.GetByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get card list: %w", err)
	}
	if len(list) == 0 {
		return nil, domain.ErrEmptyBankCardsList
	}
	return list, nil
}

func (b *BankCardObj) CreateNewBankCardObj(ctx context.Context, card *domain.BankCard) (int64, error) {

	if card.UserId == 0 {
		return 0, fmt.Errorf("user id is zero")
	}

	if card.Bank == "" {
		return 0, fmt.Errorf("service name is zero")
	}

	if card.Pid == "" {
		return 0, fmt.Errorf("pid is zero")
	}

	id, err := b.repo.Create(ctx, card)
	if err != nil {
		return 0, fmt.Errorf("failed to create card: %w", err)
	}

	return id, nil
}

func (b *BankCardObj) UpdateBankCard(ctx context.Context, card *domain.BankCard) error {

	if card.UserId == 0 {
		return fmt.Errorf("user id is zero")
	}

	if card.Bank == "" {
		return fmt.Errorf("service name is zero")
	}

	if card.Pid == "" {
		return fmt.Errorf("pid is zero")
	}

	err := b.repo.Update(ctx, card)
	if err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}

	return nil
}
