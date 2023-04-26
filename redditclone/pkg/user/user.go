package user

type User struct {
	ID       uint64
	Username string
	password string
}

func NewUser(id uint64, username, password string) *User {
	return &User{
		ID:       id,
		Username: username,
		password: password,
	}
}
