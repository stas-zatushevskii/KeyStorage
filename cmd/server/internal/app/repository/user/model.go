package user

import "database/sql"

type User struct {
	ID       sql.NullInt64
	Username sql.NullString
	Password sql.NullString
}
