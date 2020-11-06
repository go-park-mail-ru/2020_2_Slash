package rating

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type RatingUsecase interface {
	Create(rating *models.Rating) *errors.Error
	Change(rating *models.Rating) *errors.Error
	GetByUserIDContentID(userID uint64, contentID uint64) (*models.Rating, *errors.Error)
	GetContentRating(contentID uint64) (int, *errors.Error)
	Delete(rating *models.Rating) *errors.Error
}
