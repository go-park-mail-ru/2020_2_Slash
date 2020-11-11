package usecases

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"golang.org/x/crypto/bcrypt"
	"os"
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

func (uu *UserUsecase) GetByID(userID uint64) (*models.User, *errors.Error) {
	user, err := uu.userRepo.SelectByID(userID)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeUserDoesNotExist)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return user, nil
}

func (uu *UserUsecase) UpdateProfile(userID uint64, newUserData *models.User) (*models.User, *errors.Error) {
	user, err := uu.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Update email
	if newUserData.Email != "" && user.Email != newUserData.Email {
		if err := uu.checkByEmail(newUserData.Email); err == nil {
			return nil, errors.Get(CodeEmailAlreadyExists)
		}
		user.Email = newUserData.Email
	}

	// Update nickname
	if newUserData.Nickname != "" {
		user.Nickname = newUserData.Nickname
	}

	// Update password
	if newUserData.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUserData.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New(CodeInternalError, err)
		}
		user.Password = string(hashedPassword)
	}

	if err := uu.userRepo.Update(user); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return user, nil
}

func (uu *UserUsecase) UpdateAvatar(userID uint64, newAvatar string) (*models.User, *errors.Error) {
	user, customErr := uu.GetByID(userID)
	if customErr != nil {
		return nil, customErr
	}

	// Update user avatar
	prevAvatar := user.Avatar
	user.Avatar = newAvatar
	if err := uu.userRepo.Update(user); err != nil {
		return nil, errors.New(CodeInternalError, err)
	}

	// Delete prev avatar image
	if prevAvatar != "" {
		if err := os.Remove("." + prevAvatar); err != nil {
			return nil, errors.New(CodeInternalError, err)
		}
	}
	return user, nil
}

func (uu *UserUsecase) checkByEmail(email string) *errors.Error {
	_, err := uu.GetByEmail(email)
	return err
}

func (uu *UserUsecase) CheckPassword(user *models.User, password string) *errors.Error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),
		[]byte(password)); err != nil {
		return errors.Get(CodeWrongPassword)
	}
	return nil
}

func (uu *UserUsecase) IsAdmin(userID uint64) (bool, *errors.Error) {
	user, err := uu.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.Role == Admin, nil
}
