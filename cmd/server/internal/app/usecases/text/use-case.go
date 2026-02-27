package text

import (
	"context"

	domain "server/internal/app/domain/text"
)

type Repository interface {
	GetByUserID(ctx context.Context, userId int64) ([]*domain.Text, error)
	GetByID(ctx context.Context, textId int64) (*domain.Text, error)
	Create(ctx context.Context, text *domain.Text) (int64, error)
	Update(ctx context.Context, text *domain.Text) error
}

type Text struct {
	repo Repository
}

func New(repo Repository) *Text {
	return &Text{repo: repo}
}

func (b *Text) GetText(ctx context.Context, textId int64) (*domain.Text, error) {
	if textId <= 0 {
		return nil, domain.ErrInvalidTextID
	}

	item, err := b.repo.GetByID(ctx, textId)
	if err != nil {
		return nil, domain.ErrTextNotFound
	}

	return item, nil
}

func (b *Text) GetTextList(ctx context.Context, userId int64) ([]*domain.Text, error) {
	if userId <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	list, err := b.repo.GetByUserID(ctx, userId)
	if err != nil {
		return nil, domain.ErrTextNotFound
	}

	if len(list) == 0 {
		return nil, domain.ErrEmptyTextsList
	}

	return list, nil
}

func (b *Text) CreateNewText(ctx context.Context, text *domain.Text) (int64, error) {
	if text == nil {
		return 0, domain.ErrFailedCreateText
	}

	if text.UserId <= 0 {
		return 0, domain.ErrInvalidUserID
	}

	if text.Title == "" {
		return 0, domain.ErrEmptyTitle
	}

	if text.Text == "" {
		return 0, domain.ErrEmptyText
	}

	id, err := b.repo.Create(ctx, text)
	if err != nil {
		return 0, domain.ErrFailedCreateText
	}

	return id, nil
}

func (b *Text) UpdateText(ctx context.Context, text *domain.Text) error {
	if text == nil {
		return domain.ErrFailedUpdateText
	}

	if text.UserId <= 0 {
		return domain.ErrInvalidUserID
	}

	if text.Title == "" {
		return domain.ErrEmptyTitle
	}

	if text.Text == "" {
		return domain.ErrEmptyText
	}

	if err := b.repo.Update(ctx, text); err != nil {
		return domain.ErrFailedUpdateText
	}

	return nil
}
