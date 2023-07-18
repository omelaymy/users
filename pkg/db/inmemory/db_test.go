package inmemory_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/omelaymy/users/pkg/db/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByUsername(t *testing.T) {

	db := inmemory.NewInMemoryDatabase()

	testUser := inmemory.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Admin:    false,
		Password: "password",
	}

	_, _ = db.InsertUser(testUser)

	user, err := db.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, testUser.Admin, user.Admin)
	assert.Equal(t, testUser.Username, user.Username)
	assert.Equal(t, testUser.Password, user.Password)
	assert.Equal(t, testUser.Email, user.Email)

	_, err = db.GetUserByUsername("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, inmemory.NotFoundError, err)
}

func TestInsertUser(t *testing.T) {
	db := inmemory.NewInMemoryDatabase()

	// Test InsertUser
	testUser := inmemory.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Admin:    false,
	}

	id, err := db.InsertUser(testUser)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, id)

	user, err := db.GetUserById(id)
	assert.NoError(t, err)
	assert.Equal(t, testUser.Admin, user.Admin)
	assert.Equal(t, testUser.Username, user.Username)
	assert.Equal(t, testUser.Email, user.Email)
}

func TestGetUserById(t *testing.T) {
	// Create a new instance of InMemoryDatabase for testing
	db := inmemory.NewInMemoryDatabase()

	// Insert a test user
	testUser := inmemory.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Admin:    false,
		Password: "password",
	}
	id, _ := db.InsertUser(testUser)

	user, err := db.GetUserById(id)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.UUID{}, user.ID)
	assert.Equal(t, testUser.Email, user.Email)
	assert.Equal(t, testUser.Username, user.Username)
	assert.Equal(t, testUser.Admin, user.Admin)

	_, err = db.GetUserById(uuid.New())
	assert.Error(t, err)
	assert.Equal(t, inmemory.NotFoundError, err)
}

func TestGetUsers(t *testing.T) {
	db := inmemory.NewInMemoryDatabase()

	// Insert test users
	testUsers := []inmemory.User{
		{
			Username: "user1",
			Email:    "user1@example.com",
			Admin:    false,
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
			Admin:    true,
		},
	}

	for _, user := range testUsers {
		_, _ = db.InsertUser(user)
	}

	users := db.GetUsers()
	assert.Len(t, users, len(testUsers))

	for i, user := range users {
		assert.Equal(t, testUsers[i].Email, user.Email)
		assert.Equal(t, testUsers[i].Username, user.Username)
		assert.Equal(t, testUsers[i].Admin, user.Admin)
	}
}

func TestUpdateUser(t *testing.T) {
	db := inmemory.NewInMemoryDatabase()

	testUser := inmemory.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Admin:    false,
		Password: "password",
	}
	id, _ := db.InsertUser(testUser)

	testUser.Username = "updateduser"
	testUser.Email = "updateduser@example.com"

	testUser.ID = id
	err := db.UpdateUser(testUser)
	assert.NoError(t, err)

	updatedUser, err := db.GetUserById(testUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, testUser.Username, updatedUser.Username)
	assert.Equal(t, testUser.Email, updatedUser.Email)

	nonexistentUser := inmemory.User{
		Username: "nonexistent",
		Email:    "nonexistent@example.com",
		Admin:    false,
		Password: "password",
	}
	err = db.UpdateUser(nonexistentUser)
	assert.Error(t, err)
	assert.Equal(t, inmemory.NotFoundError, err)
}

func TestDeleteUser(t *testing.T) {
	db := inmemory.NewInMemoryDatabase()

	testUser := inmemory.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Admin:    false,
		Password: "password",
	}
	id, _ := db.InsertUser(testUser)

	db.DeleteUser(id)

	_, err := db.GetUserById(id)
	assert.Error(t, err)
	assert.Equal(t, inmemory.NotFoundError, err)

	nonexistentID := uuid.New()
	db.DeleteUser(nonexistentID)
}
