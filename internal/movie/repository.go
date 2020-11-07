package movie

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type MovieRepository interface {
	Insert(movie *models.Movie) error
	Update(movie *models.Movie) error
	DeleteByID(movieID uint64) error
	SelectByID(movieID uint64) (*models.Movie, error)
	SelectByContentID(contentID uint64) (*models.Movie, error)
	SelectByParams(params *models.ContentFilter) ([]*models.Movie, error)
}
