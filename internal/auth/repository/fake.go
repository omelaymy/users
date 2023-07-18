package repository

import "github.com/omelaymy/users/internal/auth"

type FakeRepository struct {
	users map[string]*auth.User
}

func NewFakeRepository(
	users map[string]*auth.User,
) *FakeRepository {
	return &FakeRepository{
		users: users,
	}
}

func (f *FakeRepository) GetUserByUsername(username string) (*auth.User, error) {
	user, ok := f.users[username]
	if !ok {
		return nil, auth.UserNotFoundError
	}
	return user, nil
}
