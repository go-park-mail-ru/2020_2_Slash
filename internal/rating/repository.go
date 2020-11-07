package rating

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type RatingRepository interface {
	Insert(rating *models.Rating) error
	SelectByUserIDContentID(userID uint64, contentID uint64) (*models.Rating, error)
	SelectRatesCount(contentID uint64) (int, error)
	SelectLikesCount(contentID uint64) (int, error)
	Update(rating *models.Rating) error
	Delete(rating *models.Rating) error
}
