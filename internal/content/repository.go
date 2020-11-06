package content

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type ContentRepository interface {
	Insert(content *models.Content) error
	Update(content *models.Content) error
	UpdateImages(content *models.Content) error
	DeleteByID(contentID uint64) error
	SelectByID(contentID uint64) (*models.Content, error)
	SelectCountriesByID(contentID uint64) ([]uint64, error)
	SelectGenresByID(contentID uint64) ([]uint64, error)
	SelectActorsByID(contentID uint64) ([]uint64, error)
	SelectDirectorsByID(contentID uint64) ([]uint64, error)
}
