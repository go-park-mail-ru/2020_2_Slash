package director

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type DirectorRepository interface {
	Insert(director *models.Director) error
	Update(director *models.Director) error
	DeleteById(id uint64) error
	SelectById(id uint64) (*models.Director, error)
}
