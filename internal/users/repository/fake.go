package repository

import (
	"github.com/google/uuid"
	"github.com/omelaymy/users/internal/users"
)

type FakeRepository struct {
	users map[uuid.UUID]*users.User
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		users: make(map[uuid.UUID]*users.User),
	}
}

func (f *FakeRepository) CreateUser(user *users.User) (uuid.UUID, error) {
	for _, u := range f.users {
		if u.Username == user.Username {
			return uuid.UUID{}, users.UserAlreadyExistsError
		}
	}

	user.Id = uuid.New()
	f.users[user.Id] = user

	return user.Id, nil
}

func (f *FakeRepository) GetUserById(id uuid.UUID) (*users.User, error) {
	user, ok := f.users[id]
	if !ok {
		return nil, users.UserNotFoundError
	}

	return user, nil
}

func (f *FakeRepository) GetUsers() []*users.User {
	users := make([]*users.User, 0, len(f.users))
	for _, user := range f.users {
		users = append(users, user)
	}

	return users
}

func (f *FakeRepository) UpdateUser(user *users.User) error {
	for _, u := range f.users {
		if u.Id != user.Id && u.Username == user.Username {
			return users.UserAlreadyExistsError
		}
	}

	_, ok := f.users[user.Id]
	if !ok {
		return users.UserNotFoundError
	}

	f.users[user.Id] = user
	return nil
}

func (f *FakeRepository) DeleteUser(id uuid.UUID) {
	delete(f.users, id)
}
