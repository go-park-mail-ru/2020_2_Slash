package user

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type UserUsecase interface {
	Create(user *models.User) *errors.Error
	GetByEmail(email string) (*models.User, *errors.Error)
	GetByID(userID uint64) (*models.User, *errors.Error)
	UpdateProfile(userID uint64, newUserData *models.User) (*models.User, *errors.Error)
	UpdateAvatar(userID uint64, newAvatar string) (*models.User, *errors.Error)
	CheckPassword(user *models.User, password string) *errors.Error
	IsAdmin(userID uint64) (bool, *errors.Error)
}
