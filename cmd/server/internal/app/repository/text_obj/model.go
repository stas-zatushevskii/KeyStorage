package text_obj

import (
	"database/sql"
	domain "server/internal/app/domain/text_obj"
)

type Text struct {
	Text   sql.NullString `db:"text"`
	Title  sql.NullString `db:"title"`
	UserID sql.NullInt64  `db:"user_id"`
	ID     sql.NullInt64  `db:"id"`
}

func (t *Text) ToDomain() *domain.Text {
	return &domain.Text{
		Text:   t.Text.String,
		UserId: t.UserID.Int64,
		TextId: t.ID.Int64,
		Title:  t.Title.String,
	}
}
