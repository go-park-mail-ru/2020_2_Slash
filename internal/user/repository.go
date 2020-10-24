package user

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type UserRepository interface {
	Insert(user *models.User) error
	SelectByEmail(email string) (*models.User, error)
}
