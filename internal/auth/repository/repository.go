package repository

import (
	"errors"

	"github.com/omelaymy/users/internal/auth"
	"github.com/omelaymy/users/pkg/db/inmemory"
	"github.com/rs/zerolog"
)

type AuthRepository struct {
	db  *inmemory.InMemoryDatabase
	log *zerolog.Logger
}

func NewAuthRepository(
	db *inmemory.InMemoryDatabase,
	log *zerolog.Logger,
) *AuthRepository {
	return &AuthRepository{
		db:  db,
		log: log,
	}
}

func (r *AuthRepository) GetUserByUsername(username string) (*auth.User, error) {
	user, err := r.db.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, inmemory.NotFoundError) {
			return nil, auth.UserNotFoundError
		}
		r.log.Err(err).Msg("failed to get user")
		return nil, auth.UnknownError
	}

	return &auth.User{
		Username: user.Username,
		Password: user.Password,
		Admin:    user.Admin,
	}, nil
}
