package usecase_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/omelaymy/users/config"
	"github.com/omelaymy/users/internal/users"
	"github.com/omelaymy/users/internal/users/repository"
	"github.com/omelaymy/users/internal/users/usecase"
	"github.com/omelaymy/users/pkg/secure"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	repo := repository.NewFakeRepository()

	cfg := &config.Config{}
	usersUsecase := usecase.NewUsers(cfg, repo)

	user := &users.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Admin:    true,
	}

	id, err := usersUsecase.CreateUser(user)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, id)

	hashedPassword, _ := secure.HashPassword(user.Password)
	assert.NotEqual(t, user.Password, hashedPassword)

	createdUser, err := repo.GetUserById(id)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, createdUser.Username)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Admin, createdUser.Admin)
}

func TestGetUser(t *testing.T) {
	repo := repository.NewFakeRepository()

	cfg := &config.Config{}
	usersUsecase := usecase.NewUsers(cfg, repo)

	user := &users.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Admin:    true,
	}

	id, _ := usersUsecase.CreateUser(user)

	retrievedUser, err := usersUsecase.GetUser(id)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, retrievedUser.Username)
	assert.Equal(t, user.Email, retrievedUser.Email)
	assert.Equal(t, user.Admin, retrievedUser.Admin)
}

func TestGetUsers(t *testing.T) {
	repo := repository.NewFakeRepository()
	cfg := &config.Config{}
	usersUsecase := usecase.NewUsers(cfg, repo)
	usersData := []*users.User{
		{
			Username: "user1",
			Password: "password1",
			Email:    "user1@example.com",
			Admin:    true,
		},
		{
			Username: "user2",
			Password: "password2",
			Email:    "user2@example.com",
			Admin:    false,
		},
		{
			Username: "user3",
			Password: "password3",
			Email:    "user3@example.com",
			Admin:    true,
		},
	}

	for _, userData := range usersData {
		usersUsecase.CreateUser(userData)
	}

	allUsers := usersUsecase.GetUsers()
	assert.Len(t, allUsers, len(usersData))
	for _, userData := range usersData {
		found := false
		for _, user := range allUsers {
			if user.Username == userData.Username {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected user not found: "+userData.Username)
	}
}

func TestUpdateUser(t *testing.T) {
	repo := repository.NewFakeRepository()
	cfg := &config.Config{}
	usersUsecase := usecase.NewUsers(cfg, repo)
	user := &users.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Admin:    true,
	}

	id, _ := usersUsecase.CreateUser(user)
	user.Email = "updated@example.com"
	user.Admin = false
	err := usersUsecase.UpdateUser(user)
	assert.NoError(t, err)

	updatedUser, err := repo.GetUserById(id)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, updatedUser.Email)
	assert.Equal(t, user.Admin, updatedUser.Admin)
}

func TestDeleteUser(t *testing.T) {
	repo := repository.NewFakeRepository()

	cfg := &config.Config{}
	usersUsecase := usecase.NewUsers(cfg, repo)

	user := &users.User{
		Username: "testuser",
		Password: "password",
		Email:    "test@example.com",
		Admin:    true,
	}

	id, _ := usersUsecase.CreateUser(user)

	usersUsecase.DeleteUser(id)

	_, err := repo.GetUserById(id)
	assert.Equal(t, users.UserNotFoundError, err)
}
