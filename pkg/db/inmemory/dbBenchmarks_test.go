package inmemory_test

import (
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/omelaymy/users/pkg/db/inmemory"
	"github.com/stretchr/testify/assert"
)

func BenchmarkInsertUser(b *testing.B) {
	db := inmemory.NewInMemoryDatabase()
	users := generateTestUsers(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := db.InsertUser(users[i])
		assert.NoError(b, err)
	}
}

func BenchmarkConcurrentInsertUser(b *testing.B) {
	db := inmemory.NewInMemoryDatabase()
	users := generateTestUsers(b.N)

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(user inmemory.User) {
			_, err := db.InsertUser(user)
			assert.NoError(b, err)
			wg.Done()
		}(users[i])
	}

	wg.Wait()
}

func BenchmarkGetUserByUsername(b *testing.B) {
	db := inmemory.NewInMemoryDatabase()
	users := generateTestUsers(b.N)

	for i := 0; i < b.N; i++ {
		_, err := db.InsertUser(users[i])
		assert.NoError(b, err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := db.GetUserByUsername(users[i].Username)
		assert.NoError(b, err)
	}
}

func BenchmarkConcurrentGetUserByUsername(b *testing.B) {
	db := inmemory.NewInMemoryDatabase()
	users := generateTestUsers(b.N)

	for i := 0; i < b.N; i++ {
		_, err := db.InsertUser(users[i])
		assert.NoError(b, err)
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(username string) {
			_, err := db.GetUserByUsername(username)
			assert.NoError(b, err)
			wg.Done()
		}(users[i].Username)
	}

	wg.Wait()
}

func generateTestUsers(count int) []inmemory.User {
	users := make([]inmemory.User, count)
	for i := 0; i < count; i++ {
		user := inmemory.User{
			Username: uuid.New().String(),
			Email:    uuid.New().String() + "@example.com",
			Admin:    i%2 == 0,
			Password: uuid.New().String(),
		}
		users[i] = user
	}
	return users
}
