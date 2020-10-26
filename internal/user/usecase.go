package user

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type UserUsecase interface {
	Create(user *models.User) *errors.Error
	GetByEmail(email string) (*models.User, *errors.Error)
}
