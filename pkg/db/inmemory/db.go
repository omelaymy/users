package inmemory

import (
	"sync"

	"github.com/google/uuid"
)

type InMemoryDatabase struct {
	idIndex       map[uuid.UUID]*User
	usernameIndex map[string]*User
	mu            *sync.RWMutex
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		idIndex:       make(map[uuid.UUID]*User),
		usernameIndex: make(map[string]*User),
		mu:            &sync.RWMutex{},
	}
}

func (db *InMemoryDatabase) GetUserByUsername(username string) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	user, ok := db.usernameIndex[username]
	if !ok {
		return User{}, NotFoundError
	}

	return *user, nil
}

func (db *InMemoryDatabase) InsertUser(user User) (uuid.UUID, error) {
	if !db.validateUniqueUsername(user.Username) {
		return uuid.UUID{}, AlreadyExistsError
	}

	id := uuid.New()

	db.mu.Lock()
	defer db.mu.Unlock()

	user.ID = id
	db.idIndex[id] = &user
	db.usernameIndex[user.Username] = &user

	return id, nil
}

func (db *InMemoryDatabase) GetUserById(id uuid.UUID) (User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.idIndex[id]
	if !ok {
		return User{}, NotFoundError
	}

	return User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Admin:    user.Admin,
	}, nil
}

func (db *InMemoryDatabase) GetUsers() []User {
	db.mu.RLock()
	defer db.mu.RUnlock()

	i := 0
	users := make([]User, len(db.idIndex))
	for id, user := range db.idIndex {
		users[i] = User{
			ID:       id,
			Email:    user.Email,
			Username: user.Username,
			Admin:    user.Admin,
		}
		i++
	}

	return users
}

func (db *InMemoryDatabase) UpdateUser(userUpdated User) error {
	db.mu.RLock()
	user, ok := db.idIndex[userUpdated.ID]
	if !ok {
		db.mu.RUnlock()
		return NotFoundError
	}
	db.mu.RUnlock()

	if user.Username != userUpdated.Username &&
		!db.validateUniqueUsername(userUpdated.Username) {
		return AlreadyExistsError
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.usernameIndex, user.Username)
	db.usernameIndex[userUpdated.Username] = user

	user.Username = userUpdated.Username
	user.Email = userUpdated.Email
	user.Admin = userUpdated.Admin
	user.Password = userUpdated.Password

	return nil
}

func (db *InMemoryDatabase) DeleteUser(id uuid.UUID) {
	db.mu.RLock()
	user, ok := db.idIndex[id]
	if !ok {
		db.mu.RUnlock()
		return
	}
	db.mu.RUnlock()

	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.idIndex, id)
	delete(db.usernameIndex, user.Username)
}

func (db *InMemoryDatabase) validateUniqueUsername(username string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	_, ok := db.usernameIndex[username]
	return !ok
}
