package content

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type ContentUsecase interface {
	Create(content *models.Content) *errors.Error
	UpdateByID(contentID uint64, newContentData *models.Content) (*models.Content, *errors.Error)
	UpdatePosters(content *models.Content, newPostersDir string) *errors.Error
	DeleteByID(contentID uint64) *errors.Error
	GetByID(contentID uint64) (*models.Content, *errors.Error)
	GetFullByID(contentID uint64) (*models.Content, *errors.Error)
	FillContent(content *models.Content) *errors.Error
	GetCountriesByID(contentID uint64) ([]*models.Country, *errors.Error)
	GetGenresByID(contentID uint64) ([]*models.Genre, *errors.Error)
	GetActorsByID(contentID uint64) ([]*models.Actor, *errors.Error)
	GetDirectorsByID(contentID uint64) ([]*models.Director, *errors.Error)
}
