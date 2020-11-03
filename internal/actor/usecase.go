package actor

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type ActorUseCase interface {
	Create(actor *models.Actor) *errors.Error
	Get(id uint64) (*models.Actor, *errors.Error)
	Change(newActor *models.Actor) *errors.Error
	DeleteById(id uint64) *errors.Error
	ListByID(actorsID []uint64) ([]*models.Actor, *errors.Error)
}
