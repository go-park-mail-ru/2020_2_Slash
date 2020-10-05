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
	Avatar   string `json:"avatar"`
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

func (ur *UserRepo) Get(userID uint64) (*User, bool) {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	for _, user := range ur.data {
		if user.ID == userID {
			return user, true
		}
	}
	return nil, false
}

func (ur *UserRepo) Delete(userID uint64) (*User, error) {
	user, has := ur.Get(userID)
	if has {
		ur.mu.Lock()
		defer ur.mu.Unlock()
		delete(ur.data, user.Email)
		return user, nil
	}
	return nil, errors.New("There is no such user")
}

func (ur *UserRepo) UpdateEmail(userID uint64, email string) error {
	user, has := ur.Get(userID)
	if !has {
		return errors.New("There is no such user")
	}
	if user.Email == email {
		return nil
	}
	if !ur.IsUniqEmail(email) {
		return errors.New("User with this Email already exists")
	}
	ur.mu.Lock()
	defer ur.mu.Unlock()
	delete(ur.data, user.Email)
	user.Email = email
	ur.data[user.Email] = user
	return nil
}

func (ur *UserRepo) Exists(userID uint64) bool {
	_, has := ur.Get(userID)
	return has
}

func (ur *UserRepo) IsUniqEmail(email string) bool {
	ur.mu.Lock()
	defer ur.mu.Unlock()
	_, has := ur.data[email]
	return !has
}

func (ur *UserRepo) Register(user *User) error {
	if !ur.IsUniqEmail(user.Email) {
		return errors.New("User with this Email already exists")
	}
	ur.mu.Lock()
	ur.mu.Unlock()
	user.ID = ur.count
	ur.data[user.Email] = user
	ur.count++
	return nil
}
