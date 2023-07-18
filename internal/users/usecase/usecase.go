package usecase

import (
	"github.com/google/uuid"
	"github.com/omelaymy/users/config"
	"github.com/omelaymy/users/internal/users"
	"github.com/omelaymy/users/pkg/secure"
)

type Users struct {
	cfg        *config.Config
	repository users.Repository
}

func NewUsers(
	cfg *config.Config,
	repository users.Repository,
) *Users {
	return &Users{
		cfg:        cfg,
		repository: repository,
	}
}

func (u *Users) CreateUser(user *users.User) (uuid.UUID, error) {
	hashedPassword, err := secure.HashPassword(user.Password)
	if err != nil {
		return uuid.UUID{}, users.UnknownError
	}
	user.Password = hashedPassword

	return u.repository.CreateUser(user)
}

func (u *Users) GetUser(id uuid.UUID) (*users.User, error) {
	return u.repository.GetUserById(id)
}

func (u *Users) GetUsers() []*users.User {
	return u.repository.GetUsers()
}

func (u *Users) UpdateUser(user *users.User) error {
	hashedPassword, err := secure.HashPassword(user.Password)
	if err != nil {
		return users.UnknownError
	}
	user.Password = hashedPassword

	return u.repository.UpdateUser(user)
}

func (u *Users) DeleteUser(id uuid.UUID) {
	u.repository.DeleteUser(id)
}
