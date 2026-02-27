package bank_card

import (
	"database/sql"
	domain "server/internal/app/domain/bank_card"
)

type Card struct {
	ID     sql.NullInt64
	UserId sql.NullInt64
	Bank   sql.NullString
	PID    []byte
}

func (c *Card) ToDomain() *domain.BankCard {
	return &domain.BankCard{
		CardId: c.ID.Int64,
		UserId: c.UserId.Int64,
		Bank:   c.Bank.String,
		Pid:    string(c.PID),
	}
}
