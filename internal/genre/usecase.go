package genre

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type GenreUsecase interface {
	Create(genre *models.Genre) *errors.Error
	UpdateByID(genreID uint64, newGenreData *models.Genre) (*models.Genre, *errors.Error)
	DeleteByID(genreID uint64) *errors.Error
	GetByID(genreID uint64) (*models.Genre, *errors.Error)
	GetByName(name string) (*models.Genre, *errors.Error)
	List() ([]*models.Genre, *errors.Error)
	ListByID(genresID []uint64) ([]*models.Genre, *errors.Error)
}
