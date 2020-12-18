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

func (uc *SubscriptionUseCase) Create(subscription *models.Subscription) *errors.Error {
	dbSubscription, customErr := uc.GetByUserID(subscription.UserID)
	if customErr != nil {
		return customErr
	}
	if dbSubscription != nil {
		return errors.Get(consts.CodeSubscriptionAlreadyExist)
	}

	err := uc.rep.Insert(subscription)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *SubscriptionUseCase) GetByUserID(userID uint64) (*models.Subscription, *errors.Error) {
	dbSubscription, err := uc.rep.SelectByUserID(userID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	if dbSubscription.Expires.Before(time.Now()) {
		if dbSubscription.Active {
			customErr := uc.updateTerm(dbSubscription)
			if customErr != nil {
				return dbSubscription, customErr
			}
			return dbSubscription, nil
		} else {
			err := uc.rep.DeleteByID(dbSubscription.ID)
			if err != nil {
				return nil, errors.New(consts.CodeInternalError, err)
			}
			return nil, nil
		}
	}

	return dbSubscription, nil
}

func (uc *SubscriptionUseCase) updateTerm(dbSubscription *models.Subscription) *errors.Error {
	// TODO: здесь же "оплата"
	dbSubscription.Expires = time.Now().AddDate(0, 1, 0)
	err := uc.rep.Update(dbSubscription)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *SubscriptionUseCase) DeleteByUserID(userID uint64) *errors.Error {
	dbSubscription, customErr := uc.GetByUserID(userID)
	if customErr != nil {
		return customErr
	}
	if dbSubscription == nil {
		return errors.Get(consts.CodeSubscriptionDoesNotExist)
	}

	if dbSubscription.Expires.Before(time.Now()) {
		err := uc.rep.DeleteByID(dbSubscription.ID)
		if err != nil {
			return errors.New(consts.CodeInternalError, err)
		}
	}

	dbSubscription.Active = false

	err := uc.rep.Update(dbSubscription)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}
