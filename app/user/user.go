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

func (r *UserRepo) Get(email string) (*User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, has := r.data[email]
	return user, has
}

func (r *UserRepo) exists(user *User) bool {
	_, has := r.data[user.Email]
	return has
}

func (r *UserRepo) Register(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.exists(user) {
		return errors.New("User with this Email already exists")
	}
	user.ID = r.count
	r.data[user.Email] = user
	r.count++
	return nil
}
