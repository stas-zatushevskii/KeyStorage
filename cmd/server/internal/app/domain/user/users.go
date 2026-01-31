package user

type User struct {
	ID       int64
	Username string
	Password string
}

func NewUser() *User {
	return &User{}
}
