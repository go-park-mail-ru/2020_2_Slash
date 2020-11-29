package favourite

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type FavouriteRepository interface {
	Insert(favourite *models.Favourite) error
	Select(favourite *models.Favourite) error
	SelectFavouriteMovies(userID uint64, limit uint64, offset uint64) ([]*models.Movie, error)
	SelectFavouriteTVShows(userID uint64, limit uint64, offset uint64) ([]*models.TVShow, error)
	Delete(favourite *models.Favourite) error
}
