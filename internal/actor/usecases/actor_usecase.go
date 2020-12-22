package usecases

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type ActorUseCase struct {
	actorRepo actor.ActorRepository
}

func NewActorUseCase(repo actor.ActorRepository) actor.ActorUseCase {
	return &ActorUseCase{
		actorRepo: repo,
	}
}

func (au *ActorUseCase) Create(actor *models.Actor) *errors.Error {
	err := au.actorRepo.Insert(actor)
	if err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (au *ActorUseCase) Get(id uint64) (*models.Actor, *errors.Error) {
	dbActor, err := au.actorRepo.SelectById(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(CodeActorDoesNotExist)
	} else if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	return dbActor, nil
}

func (au *ActorUseCase) Change(newActor *models.Actor) *errors.Error {
	if _, customErr := au.Get(newActor.ID); customErr != nil {
		return customErr
	}

	if err := au.actorRepo.Update(newActor); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (au *ActorUseCase) DeleteById(id uint64) *errors.Error {
	if _, customErr := au.Get(id); customErr != nil {
		return customErr
	}

	if err := au.actorRepo.DeleteById(id); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (au *ActorUseCase) ListByID(actorsID []uint64) ([]*models.Actor, *errors.Error) {
	var actors []*models.Actor
	for _, actorID := range actorsID {
		actor, err := au.Get(actorID)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	return actors, nil
}

func (au *ActorUseCase) List(pgnt *models.Pagination) ([]*models.Actor, *errors.Error) {
	actors, err := au.actorRepo.SelectAll(pgnt)
	if err != nil {
		return nil, errors.New(CodeInternalError, err)
	}
	if len(actors) == 0 {
		return []*models.Actor{}, nil
	}
	return actors, nil
}
