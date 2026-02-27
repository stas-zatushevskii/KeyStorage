package text_obj

import (
	"context"
	"database/sql"
	"errors"
	domain "server/internal/app/domain/text_obj"
	"server/internal/pkg/logger"

	"go.uber.org/zap"
)

func (u *Repository) GetByUserID(ctx context.Context, userId int64) ([]*domain.Text, error) {
	query := `
		SELECT id, user_id, title, text
		FROM text_data
		WHERE user_id = $1`

	var texts []*domain.Text

	rows, err := u.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Log.Error("rows.Close() failed", zap.Error(err))
		}
	}()

	for rows.Next() {
		obj := new(Text)

		if err := rows.Scan(&obj.ID, &obj.UserID, &obj.Title, &obj.Text); err != nil {
			return nil, err
		}

		texts = append(texts, obj.ToDomain())
	}
	return texts, nil
}

func (u *Repository) GetByID(ctx context.Context, cardId int64) (*domain.Text, error) {
	query := `
		SELECT id, user_id, title, text
		FROM text_data
		WHERE id = $1`

	obj := new(Text)

	if err := u.db.QueryRowContext(ctx, query, cardId).Scan(&obj.ID, &obj.UserID, &obj.Title, &obj.Text); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTextInformationNotFound
		}

		return nil, err
	}

	return obj.ToDomain(), nil
}

func (u *Repository) Create(ctx context.Context, card *domain.Text) (int64, error) {
	query := `
		INSERT INTO text_data (user_id, title, text)
		VALUES ($1, $2, $3)
		RETURNING id`

	var id sql.NullInt64

	if err := u.db.QueryRowContext(ctx, query, card.UserId, card.Title, card.Text).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrFailedCreateText
		}
		return 0, err
	}

	if !id.Valid {
		return 0, domain.ErrFailedCreateText
	}

	return id.Int64, nil
}

func (u *Repository) Update(ctx context.Context, card *domain.Text) error {
	query := `
		UPDATE text_data SET
		title = $1, text = $2
		WHERE id = $3`

	if _, err := u.db.ExecContext(ctx, query, card.Title, card.Text, card.TextId); err != nil {
		return err
	}
	return nil
}
