package usecase_test

import (
	"testing"

	"github.com/omelaymy/users/internal/auth"
	"github.com/omelaymy/users/internal/auth/repository"
	"github.com/omelaymy/users/internal/auth/usecase"
	"github.com/omelaymy/users/pkg/secure"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	password, _ := secure.HashPassword("password")
	users := map[string]*auth.User{
		"testuser": {
			Username: "testuser",
			Password: password,
			Admin:    false,
		},
	}
	repo := repository.NewFakeRepository(users)
	authUsecase := usecase.NewAuth(repo)

	assert.True(t, authUsecase.Authentication("testuser", "password"))
	assert.False(t, authUsecase.Authentication("testuser", "wrongpassword"))
	assert.False(t, authUsecase.Authentication("nonexistentuser", "password"))
}

func TestAdminAuthorization(t *testing.T) {
	password, _ := secure.HashPassword("password")
	users := map[string]*auth.User{
		"testadmin": {
			Username: "testadmin",
			Password: password,
			Admin:    true,
		},
		"testuser": {
			Username: "testuser",
			Password: password,
			Admin:    false,
		},
	}
	repo := repository.NewFakeRepository(users)
	authUsecase := usecase.NewAuth(repo)

	assert.True(t, authUsecase.AdminAuthorization("testadmin", "password"))
	assert.False(t, authUsecase.AdminAuthorization("testuser", "password"))
	assert.False(t, authUsecase.AdminAuthorization("testadmin", "wrongpassword"))
	assert.False(t, authUsecase.AdminAuthorization("nonexistentuser", "password"))
}
