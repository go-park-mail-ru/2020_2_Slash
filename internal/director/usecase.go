package director

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type DirectorUseCase interface {
	Create(director *models.Director) *errors.Error
	Get(id uint64) (*models.Director, *errors.Error)
	Change(newDirector *models.Director) *errors.Error
	DeleteById(id uint64) *errors.Error
	ListByID(directorsID []uint64) ([]*models.Director, *errors.Error)
	List(pgnt *models.Pagination) ([]*models.Director, *errors.Error)
}
