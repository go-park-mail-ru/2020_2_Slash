package favourite

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type FavouriteUsecase interface {
	Create(favourite *models.Favourite) *errors.Error
	GetUserFavouriteMovies(userID uint64, pagination *models.Pagination) ([]*models.Movie, *errors.Error)
	Delete(favourite *models.Favourite) *errors.Error
}
