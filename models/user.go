package models

import (
	"regexp"
	"strings"
	"sync"
)

type User struct {
	Login    string `json:"login"` // unique field
	Email    string `json:"email"` // unique field
	Password string `json:"password;omitempty"` // password in UsersStore == some_hash(password)
}

const (
	_ = iota
	StatusUserNotExist
	StatusUserLoginNotUnique
	StatusUserEmailNotUnique
	StatusUserLoginInvalid
	StatusUserEmailInvalid
)

var EmailValidator = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type UsersStore struct {
	users map[string]User
	mu    *sync.Mutex
}

func (store *UsersStore) ValidateUser(user User) CommonError {
	if _, isUniqueLogin := store.users[user.Login]; !isUniqueLogin {
		return NewModelError(`user with same login already exists`, StatusUserLoginNotUnique)
	}
	if _, isUniqueEmail := store.users[user.Email]; !isUniqueEmail {
		return NewModelError(`user with same email already exists`, StatusUserEmailNotUnique)
	}
	if len(strings.Fields(user.Login)) != 1 {
		return NewModelError(`user login invalid`, StatusUserLoginInvalid)
	}
	if EmailValidator.MatchString(user.Email) == false {
		return NewModelError(`user email invalid`, StatusUserEmailInvalid)
	}
	return nil
}

func (store *UsersStore) StoreUser(user User) CommonError {
	defer store.mu.Unlock()
	store.mu.Lock()

	if err := store.ValidateUser(user); err != nil {
		return err
	}
	store.users[user.Login] = user
	return nil
}

func (store *UsersStore) GetUser(userLogin string) (User, CommonError) {
	user, inStore := store.users[userLogin]
	if !inStore {
		return user, NewModelError(`user not exits`, StatusUserNotExist)
	}
	return user, nil
}
