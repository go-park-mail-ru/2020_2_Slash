package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/subscription"
	"time"
)

type SubscriptionUseCase struct {
	rep subscription.SubscriptionRepository
}

func NewSubscriptionUseCase(
	rep subscription.SubscriptionRepository) subscription.SubscriptionUseCase {
	return &SubscriptionUseCase{rep: rep}
}

func (uc *SubscriptionUseCase) Create(subscription *models.Subscription) (*models.Subscription, *errors.Error) {
	dbSubscription, customErr := uc.GetByUserID(subscription.UserID)
	if customErr != nil {
		return nil, customErr
	}

	if dbSubscription != nil {
		dbSubscription.IsPaid = true
		dbSubscription.IsCanceled = false
		if !isExpired(dbSubscription) {
			oldExpires := dbSubscription.Expires
			dbSubscription.Expires = oldExpires.AddDate(0, 1, 0)
		} else {
			dbSubscription.Expires = subscription.Expires
		}
		err := uc.rep.Update(dbSubscription)
		if err != nil {
			return nil, errors.New(consts.CodeInternalError, err)
		}
		return dbSubscription, nil
	}

	err := uc.rep.Insert(subscription)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return subscription, nil
}

func (uc *SubscriptionUseCase) RecoverSubscriptionByUserID(userID uint64) (*models.Subscription, *errors.Error) {
	dbSubscription, err := uc.rep.SelectByUserID(userID)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	if dbSubscription == nil {
		return nil, errors.Get(consts.CodeSubscriptionDoesNotExist)
	}

	if dbSubscription.IsCanceled {
		dbSubscription.IsCanceled = false
	}

	err = uc.rep.Update(dbSubscription)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return dbSubscription, nil
}

func isExpired(subscription *models.Subscription) bool {
	return subscription.Expires.Before(time.Now())
}

func (uc *SubscriptionUseCase) GetByUserID(userID uint64) (*models.Subscription, *errors.Error) {
	dbSubscription, err := uc.rep.SelectByUserID(userID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	if isExpired(dbSubscription) {
		if dbSubscription.IsCanceled {
			err := uc.rep.DeleteByID(dbSubscription.ID)
			if err != nil {
				return nil, errors.New(consts.CodeInternalError, err)
			}
			return nil, nil
		} else {
			dbSubscription.IsPaid = false
			err := uc.rep.Update(dbSubscription)
			if err != nil {
				return nil, errors.New(consts.CodeInternalError, err)
			}
			return dbSubscription, nil
		}
	}

	return dbSubscription, nil
}

func (uc *SubscriptionUseCase) DeleteByUserID(userID uint64) (*models.Subscription, *errors.Error) {
	dbSubscription, err := uc.rep.SelectByUserID(userID)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	if dbSubscription == nil {
		return nil, errors.Get(consts.CodeSubscriptionDoesNotExist)
	}

	if isExpired(dbSubscription) {
		err := uc.rep.DeleteByID(dbSubscription.ID)
		if err != nil {
			return nil, errors.New(consts.CodeInternalError, err)
		}
	}

	dbSubscription.IsCanceled = true

	err = uc.rep.Update(dbSubscription)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return dbSubscription, nil
}
