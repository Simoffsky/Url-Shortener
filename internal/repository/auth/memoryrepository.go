package auth

import (
	"sync"
	"url-shorter/internal/models"
)

type UserRepositoryMemory struct {
	mx    *sync.RWMutex
	users map[string]models.User
}

func NewAuthRepositoryMemory() *UserRepositoryMemory {
	return &UserRepositoryMemory{
		mx:    &sync.RWMutex{},
		users: make(map[string]models.User),
	}
}

func (r *UserRepositoryMemory) Exists(login string) (bool, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	_, ok := r.users[login]
	return ok, nil
}

func (r *UserRepositoryMemory) CreateUser(user *models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	exists, _ := r.Exists(user.Login)
	if exists {
		return models.ErrUserExists
	}

	r.users[user.Login] = *user
	return nil
}

func (r *UserRepositoryMemory) GetUser(login string) (*models.User, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	user, ok := r.users[login]
	if !ok {
		return nil, models.ErrUserNotFound
	}

	return &user, nil
}
