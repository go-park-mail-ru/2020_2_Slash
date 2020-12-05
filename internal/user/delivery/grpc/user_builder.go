package grpc

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type UserBuilder struct{}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (um *UserBuilder) CreateRegularUser() *User {
	return &User{
		ID:       1,
		Nickname: "Just a user",
		Email:    "user@mail.ru",
		Password: "123456",
		Avatar:   "",
		Role:     "user",
	}
}

func (um *UserBuilder) CreateRegularUserModel() *models.User {
	return &models.User{
		ID:       1,
		Nickname: "Just a user",
		Email:    "user@mail.ru",
		Password: "123456",
		Avatar:   "",
		Role:     "user",
	}
}
func (um *UserBuilder) CreateUserWithoutRoleModel() *models.User {
	return &models.User{
		ID:       1,
		Nickname: "Just a user",
		Email:    "user@mail.ru",
		Password: "123456",
		Avatar:   "",
		Role:     "",
	}
}

func (um *UserBuilder) CreateAdmin() *User {
	return &User{
		ID:       2,
		Nickname: "admin",
		Email:    "admin@mail.ru",
		Password: "123456",
		Avatar:   "",
		Role:     "admin",
	}
}

func (um *UserBuilder) CreateUserWithEmptyNickname() *User {
	return &User{
		ID:       3,
		Nickname: "",
		Email:    "nick@mail.ru",
		Password: "123456",
		Avatar:   "",
		Role:     "user",
	}
}

func (um *UserBuilder) CreateNinthUserFromDB() *User {
	return &User{
		ID:       9,
		Nickname: "oleg",
		Email:    "oleg@mail.ru",
		Password: "123456",
		Avatar:   "/avatars/userid_9_ae0d514d-dd82-4943-b31f-9a6a03007d2f.png",
		Role:     "admin",
	}
}

func (um *UserBuilder) CreateNinthUserFromDBModel() *models.User {
	return &models.User{
		ID:       9,
		Nickname: "oleg",
		Email:    "oleg@mail.ru",
		Password: "123456",
		Avatar:   "/avatars/userid_9_ae0d514d-dd82-4943-b31f-9a6a03007d2f.png",
		Role:     "admin",
	}
}

func (um *UserBuilder) CreateUserWithNotExistedID() *User {
	existedUser := um.CreateNinthUserFromDB()
	existedUser.ID = 9321 // there is no such id in database
	return existedUser
}

func (um *UserBuilder) CreateUserWithConflictEmail() *User {
	existedUser := um.CreateNinthUserFromDB()
	// users in db has 9 and 10 ids
	existedUser.ID = 10
	return existedUser
}

func (um *UserBuilder) CreateUserWithConflictEmailModel() *models.User {
	existedUser := um.CreateUserWithConflictEmail()
	return GrpcUserToModel(existedUser)
}
