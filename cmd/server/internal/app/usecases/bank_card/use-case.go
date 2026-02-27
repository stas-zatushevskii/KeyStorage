package bank_card

import (
	"context"

	domain "server/internal/app/domain/bank_card"
)

type Repository interface {
	GetByUserID(ctx context.Context, userId int64) ([]*domain.BankCard, error)
	GetByID(ctx context.Context, cardId int64) (*domain.BankCard, error)
	Create(ctx context.Context, card *domain.BankCard) (int64, error)
	Update(ctx context.Context, card *domain.BankCard) error
}

type BankCard struct {
	repo Repository
}

func New(repo Repository) *BankCard {
	return &BankCard{repo: repo}
}

func (b *BankCard) GetBankCard(ctx context.Context, cardId int64) (*domain.BankCard, error) {
	if cardId <= 0 {
		return nil, domain.ErrInvalidCardID
	}

	card, err := b.repo.GetByID(ctx, cardId)
	if err != nil {
		return nil, domain.ErrBankCardNotFound
	}

	return card, nil
}

func (b *BankCard) GetBankCardList(ctx context.Context, userId int64) ([]*domain.BankCard, error) {
	if userId <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	list, err := b.repo.GetByUserID(ctx, userId)
	if err != nil {
		return nil, domain.ErrBankCardNotFound
	}

	if len(list) == 0 {
		return nil, domain.ErrBankCardNotFound
	}

	return list, nil
}

func (b *BankCard) CreateNewBankCard(ctx context.Context, card *domain.BankCard) (int64, error) {
	if card == nil {
		return 0, domain.ErrFaildeCreateBankCardObject
	}

	if card.UserId <= 0 {
		return 0, domain.ErrInvalidUserID
	}

	if card.Bank == "" {
		return 0, domain.ErrEmptyBankName
	}

	if card.Pid == "" {
		return 0, domain.ErrEmptyPID
	}

	id, err := b.repo.Create(ctx, card)
	if err != nil {
		return 0, domain.ErrFaildeCreateBankCardObject
	}

	return id, nil
}

func (b *BankCard) UpdateBankCard(ctx context.Context, card *domain.BankCard) error {
	if card == nil {
		return domain.ErrFailedUpdateBankCard
	}

	if card.UserId <= 0 {
		return domain.ErrInvalidUserID
	}

	if card.Bank == "" {
		return domain.ErrEmptyBankName
	}

	if card.Pid == "" {
		return domain.ErrEmptyPID
	}

	if err := b.repo.Update(ctx, card); err != nil {
		return domain.ErrFailedUpdateBankCard
	}

	return nil
}
