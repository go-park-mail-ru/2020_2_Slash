package subscription

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type SubscriptionRepository interface {
	Insert(subscription *models.Subscription) error
	Update(subscription *models.Subscription) error
	SelectByUserID(userID uint64) (*models.Subscription, error)
	DeleteByID(id uint64) error
}
