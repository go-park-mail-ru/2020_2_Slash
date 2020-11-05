package favourite

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type FavouriteRepository interface {
	Insert(favourite *models.Favourite) error
	Select(favourite *models.Favourite) error
	SelectFavouriteContent(userID uint64) ([]*models.Content, error)
	Delete(favourite *models.Favourite) error
}
