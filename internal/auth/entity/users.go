package entity

// UserAccount - db schema
type UserAccount struct {
	Username  string
	Password  string
	CreatedAt string
}

type RegisterUserRequest struct {
	Username string
	Password string
}

type LoginUserRequest struct {
	Username string
	Password string
}
