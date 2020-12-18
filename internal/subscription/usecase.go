package subscription

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type SubscriptionUseCase interface {
	Create(subscription *models.Subscription) *errors.Error
	GetByUserID(userID uint64) (*models.Subscription, *errors.Error)
	DeleteByUserID(userID uint64) *errors.Error
}
