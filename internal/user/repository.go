package user

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type UserRepository interface {
	Insert(user *models.User) error
	SelectByEmail(email string) (*models.User, error)
	SelectByID(userID uint64) (*models.User, error)
	Update(user *models.User) error
}
