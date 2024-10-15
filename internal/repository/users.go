package repository

import (
	"errors"
	"sync"
)

type User struct {
	Username string
	Password string
	Balance  float64
}

type Repository struct {
	users map[string]*User
	mu    sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		users: make(map[string]*User),
	}
}

func (r *Repository) CreateUser(username, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[username]; exists {
		return errors.New("user already exists")
	}

	r.users[username] = &User{
		Username: username,
		Password: password,
		Balance:  100, // Начальный баланс
	}
	return nil
}

func (r *Repository) Authenticate(username, password string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[username]
	if !exists || user.Password != password {
		return "", errors.New("invalid username or password")
	}

	return "token-for-user-" + username, nil // Пример токена
}

func (r *Repository) GetUsers() ([]*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var users []*User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}
