package account_obj

import (
	"database/sql"
	domain "server/internal/app/domain/account"
)

type Account struct {
	ID          sql.NullInt64
	UserId      sql.NullInt64
	ServiceName sql.NullString
	UserName    sql.NullString
	Password    sql.NullString
}

func (u *Account) ToDomain() *domain.Account {
	return &domain.Account{
		AccountId:   u.ID.Int64,
		UserId:      u.UserId.Int64,
		ServiceName: u.ServiceName.String,
		UserName:    u.UserName.String,
		Password:    u.Password.String,
	}
}
