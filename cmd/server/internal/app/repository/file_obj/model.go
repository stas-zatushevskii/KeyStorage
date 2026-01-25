package file_obj

import (
	"database/sql"
	domain "server/internal/app/domain/file_obj"
)

type Card struct {
	ID     sql.NullInt64
	UserId sql.NullInt64
	Bank   sql.NullString
	PID    []byte
}

func (c *Card) ToDomain() *domain.File {
	return nil
}
