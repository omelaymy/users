package users

import "github.com/google/uuid"

type Repository interface {
	CreateUser(user *User) (uuid.UUID, error)
	GetUserById(id uuid.UUID) (*User, error)
	GetUsers() []*User
	UpdateUser(user *User) error
	DeleteUser(id uuid.UUID)
}
