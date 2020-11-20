package movie

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type MovieUsecase interface {
	Create(movie *models.Movie) *errors.Error
	UpdateVideo(movie *models.Movie, newVideoPath string) *errors.Error
	DeleteByID(movieID uint64) *errors.Error
	GetByID(movieID uint64) (*models.Movie, *errors.Error)
	GetFullByID(movieID uint64, curUserID uint64) (*models.Movie, *errors.Error)
	GetByContentID(contentID uint64) (*models.Movie, *errors.Error)
	ListByParams(params *models.ContentFilter, pgnt *models.Pagination,
		curUserID uint64) ([]*models.Movie, *errors.Error)
	ListLatest(pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, *errors.Error)
	ListByRating(pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, *errors.Error)
}
