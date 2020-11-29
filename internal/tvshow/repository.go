package tvshow

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type TVShowRepository interface {
	Insert(tvshow *models.TVShow) error
	SelectByID(tvshowID uint64) (*models.TVShow, error)
	SelectShortByID(tvshowID uint64) (*models.TVShow, error)
	SelectFullByID(tvshowID uint64, curUserID uint64) (*models.TVShow, error)
	SelectByContentID(contentID uint64) (*models.TVShow, error)
	SelectWhereNameLike(name string, pgnt *models.Pagination,
		curUserID uint64) ([]*models.TVShow, error)
	SelectByParams(params *models.ContentFilter, pgnt *models.Pagination,
		curUserID uint64) ([]*models.TVShow, error)
	SelectLatest(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, error)
	SelectByRating(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, error)
}
