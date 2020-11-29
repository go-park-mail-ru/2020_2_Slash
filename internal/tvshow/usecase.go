package tvshow

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type TVShowUsecase interface {
	Create(tvshow *models.TVShow) *errors.Error
	GetByID(tvshowID uint64) (*models.TVShow, *errors.Error)
	GetShortByID(tvshowID uint64) (*models.TVShow, *errors.Error)
	GetFullByID(tvshowID uint64, curUserID uint64) (*models.TVShow, *errors.Error)
	GetByContentID(contentID uint64) (*models.TVShow, *errors.Error)
	ListByParams(params *models.ContentFilter, pgnt *models.Pagination,
		curUserID uint64) ([]*models.TVShow, *errors.Error)
	ListLatest(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, *errors.Error)
	ListByRating(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, *errors.Error)
}
