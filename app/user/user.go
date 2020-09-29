package user

import (
	"errors"
	"sync"
)

type User struct {
	ID       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UserRepo struct {
	data  map[string]*User
	count uint64
	mu    *sync.Mutex
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data:  make(map[string]*User),
		count: 0,
		mu:    &sync.Mutex{},
	}
}

func (repo *UserRepo) Get(nickname string) (*User, bool) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	user, ok := repo.data[nickname]
	return user, ok
}

func (repo *UserRepo) exists(user *User) bool {
	_, ok := repo.data[user.Nickname]
	return ok
}

func (repo *UserRepo) Register(user *User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.exists(user) {
		return errors.New("User with this nickname already exists")
	}
	user.ID = repo.count
	repo.data[user.Nickname] = user
	repo.count++
	return nil
}
