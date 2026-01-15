package text_obj

import (
	"context"
	"fmt"
	domain "server/internal/app/domain/text_obj"
)

type Repository interface {
	GetByUserID(ctx context.Context, userId int64) ([]*domain.Text, error)
	GetByID(ctx context.Context, textId int64) (*domain.Text, error)
	Create(ctx context.Context, text *domain.Text) (int64, error)
	Update(ctx context.Context, text *domain.Text) error
}

type TextObj struct {
	repo Repository
}

func New(repo Repository) *TextObj {
	return &TextObj{repo: repo}
}

func (b *TextObj) GetText(ctx context.Context, textId int64) (*domain.Text, error) {
	account, err := b.repo.GetByID(ctx, textId)
	if err != nil {
		return nil, fmt.Errorf("failed to get text: %w", err)
	}
	return account, nil
}

func (b *TextObj) GetTextList(ctx context.Context, userId int64) ([]*domain.Text, error) {
	list, err := b.repo.GetByUserID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get text list: %w", err)
	}
	if len(list) == 0 {
		return nil, domain.ErrEmptyTextsList
	}
	return list, nil
}

func (b *TextObj) CreateNewTextObj(ctx context.Context, text *domain.Text) (int64, error) {

	if text.UserId == 0 {
		return 0, fmt.Errorf("user id is zero")
	}

	if text.Title == "" {
		return 0, fmt.Errorf("title is zero")
	}

	if text.Text == "" {
		return 0, fmt.Errorf("text is zero")
	}

	id, err := b.repo.Create(ctx, text)
	if err != nil {
		return 0, fmt.Errorf("failed to create text: %w", err)
	}

	return id, nil
}

func (b *TextObj) UpdateText(ctx context.Context, text *domain.Text) error {

	if text.UserId == 0 {
		return fmt.Errorf("user id is zero")
	}

	if text.Title == "" {
		return fmt.Errorf("title is zero")
	}

	if text.Text == "" {
		return fmt.Errorf("text is zero")
	}

	err := b.repo.Update(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to create text: %w", err)
	}

	return nil
}
