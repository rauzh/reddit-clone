package user

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNoUser    = errors.New("no user found")
	ErrBadPass   = errors.New("invald password")
	ErrExistUser = errors.New("login is already taken")
)

type UserRepo struct {
	mu        *sync.RWMutex
	counterID uint64
	data      map[string]*User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		mu:        &sync.RWMutex{},
		data:      map[string]*User{},
		counterID: 0,
	}
}

func (repo *UserRepo) Authorize(username, password string) (*User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	u, ok := repo.data[username]
	if !ok {
		return nil, ErrNoUser
	}
	if u.password != password {
		return nil, ErrBadPass
	}
	return u, nil
}

func (repo *UserRepo) UserExist(username, userID string) bool {
	for _, user := range repo.data {
		if fmt.Sprint(user.ID) == userID && user.Username == username {
			return true
		}
	}

	return false
}

func (repo *UserRepo) Registration(username, password string) (*User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	_, ok := repo.data[username]
	if ok {
		return nil, ErrExistUser
	}

	u := &User{
		ID:       repo.counterID,
		Username: username,
		password: password,
	}

	repo.counterID++
	repo.data[username] = u
	return u, nil
}
