package genre

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type GenreRepository interface {
	Insert(genre *models.Genre) error
	Update(genre *models.Genre) error
	DeleteByID(genreID uint64) error
	SelectByID(genreID uint64) (*models.Genre, error)
	SelectByName(name string) (*models.Genre, error)
	SelectAll() ([]*models.Genre, error)
}
