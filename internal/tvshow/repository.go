package tvshow

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type TVShowRepository interface {
	Insert(tvshow *models.TVShow) error
	SelectByID(tvshowID uint64) (*models.TVShow, error)
	SelectFullByID(tvshowID uint64, curUserID uint64) (*models.TVShow, error)
	SelectByContentID(contentID uint64) (*models.TVShow, error)
	SelectWhereNameLike(name string, pgnt *models.Pagination,
		curUserID uint64) ([]*models.TVShow, error)
}
