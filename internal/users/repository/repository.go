package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/omelaymy/users/internal/users"
	"github.com/omelaymy/users/pkg/db/inmemory"
	"github.com/rs/zerolog"
)

type UsersRepository struct {
	db  *inmemory.InMemoryDatabase
	log *zerolog.Logger
}

func NewUsersRepository(
	db *inmemory.InMemoryDatabase,
	log *zerolog.Logger,
) *UsersRepository {
	return &UsersRepository{
		db:  db,
		log: log,
	}
}

func (r *UsersRepository) CreateUser(user *users.User) (uuid.UUID, error) {
	id, err := r.db.InsertUser(
		inmemory.User{
			Email:    user.Email,
			Username: user.Username,
			Password: user.Password,
			Admin:    user.Admin,
		},
	)
	if err != nil {
		if errors.Is(err, inmemory.AlreadyExistsError) {
			return uuid.UUID{}, users.UserAlreadyExistsError
		}
		return uuid.UUID{}, users.UnknownError
	}

	return id, nil
}

func (r *UsersRepository) GetUserById(id uuid.UUID) (*users.User, error) {
	user, err := r.db.GetUserById(id)
	if err != nil {
		if errors.Is(err, inmemory.NotFoundError) {
			return nil, users.UserNotFoundError
		}
		return nil, users.UnknownError
	}

	return &users.User{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Admin:    user.Admin,
	}, nil
}

func (r *UsersRepository) GetUsers() []*users.User {
	return castUsersFromDB(r.db.GetUsers())
}

func (r *UsersRepository) UpdateUser(user *users.User) error {
	err := r.db.UpdateUser(
		inmemory.User{
			ID:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			Admin:    user.Admin,
		},
	)
	if err != nil {
		if errors.Is(err, inmemory.NotFoundError) {
			return users.UserNotFoundError
		}
		if errors.Is(err, inmemory.AlreadyExistsError) {
			return users.UserAlreadyExistsError
		}
		return users.UnknownError
	}

	return nil
}

func (r *UsersRepository) DeleteUser(id uuid.UUID) {
	r.db.DeleteUser(id)
}

func castUsersFromDB(inmemoryUsers []inmemory.User) []*users.User {
	res := make([]*users.User, len(inmemoryUsers))
	for i, user := range inmemoryUsers {
		res[i] = &users.User{
			Id:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Admin:    user.Admin,
		}
	}

	return res
}
