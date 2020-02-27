package models

import (
	"fmt"
	"sync"
)

type User struct {
	Login    string `json:"login"` // unique field
	Password string `json:"password"`
}

type UsersStore struct {
	users map[string]User
	mu    *sync.Mutex
}

func (store *UsersStore) StoreUser(user User) error {
	defer store.mu.Unlock()
	store.mu.Lock()

	if _, isUnique := store.users[user.Login]; isUnique {
		return fmt.Errorf("user already exists")
	}
	store.users[user.Login] = user
	return nil
}

func (store *UsersStore) GetUser(userLogin string) (User, error) {
	user, inStore := store.users[userLogin]
	if !inStore {
		return user, fmt.Errorf("user not exits")
	}
	return user, nil
}
