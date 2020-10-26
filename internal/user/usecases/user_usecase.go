package usecases

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRepo: repo,
	}
}

func (uu *UserUsecase) Create(user *models.User) *errors.Error {
	if err := uu.checkByEmail(user.Email); err == nil {
		return errors.Get(CodeEmailAlreadyExists)
	}

	// If nickname wasn't sent, set nickname as email before @
	if user.Nickname == "" {
		user.Nickname = strings.Split(user.Email, "@")[0]
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(CodeInternalError, err)
	}
	user.Password = string(hashedPassword)

	if err := uu.userRepo.Insert(user); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (uu *UserUsecase) GetByEmail(email string) (*models.User, *errors.Error) {
	user, err := uu.userRepo.SelectByEmail(email)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeEmailDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return user, nil
}

func (uu *UserUsecase) checkByEmail(email string) *errors.Error {
	_, err := uu.GetByEmail(email)
	return err
}
