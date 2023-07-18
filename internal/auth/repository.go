package auth

type Repository interface {
	GetUserByUsername(username string) (*User, error)
}
