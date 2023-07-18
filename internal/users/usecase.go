package users

import "github.com/google/uuid"

type Usecase interface {
	CreateUser(user *User) (uuid.UUID, error)
	GetUser(id uuid.UUID) (*User, error)
	GetUsers() []*User
	UpdateUser(user *User) error
	DeleteUser(id uuid.UUID)
}
