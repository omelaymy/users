package usecase

import (
	"github.com/omelaymy/users/internal/auth"
	"github.com/omelaymy/users/pkg/secure"
)

type Auth struct {
	repository auth.Repository
}

func NewAuth(
	repository auth.Repository,
) *Auth {
	return &Auth{
		repository: repository,
	}
}

func (a *Auth) Authentication(username, password string) bool {
	user, err := a.repository.GetUserByUsername(username)
	if err != nil {
		return false
	}

	if err := secure.ComparePasswords(user.Password, password); err != nil {
		return false
	}

	return true
}

func (a *Auth) AdminAuthorization(username, password string) bool {
	user, err := a.repository.GetUserByUsername(username)
	if err != nil {
		return false
	}

	if err := secure.ComparePasswords(user.Password, password); err != nil {
		return false
	}

	return user.Admin
}
