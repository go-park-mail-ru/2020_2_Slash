package search

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type SearchUsecase interface {
	Search(curUserID uint64, query string, pagination *models.Pagination) (
		[]*models.Movie, []*models.Actor, *errors.Error)
}
