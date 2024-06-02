package users

import (
	"sync"

	"application-design/internal/domain/users"
)

type InMemoryRepo struct {
	users map[string]users.User
	mu    sync.RWMutex
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		users: make(map[string]users.User),
	}
}

func (r *InMemoryRepo) SaveUser(user users.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.Email()] = user
	return nil
}

func (r *InMemoryRepo) GetUser(email string) (*users.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[email]
	if !ok {
		return nil, users.ErrUserNotFound
	}
	return &user, nil
}
