package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/jinzhu/copier"
)

type ActorUseCase struct {
	actorRepo        actor.ActorRepository
	adminPanelClient admin.AdminPanelClient
}

func NewActorUseCase(repo actor.ActorRepository,
	client admin.AdminPanelClient) actor.ActorUseCase {
	return &ActorUseCase{
		actorRepo:        repo,
		adminPanelClient: client,
	}
}

func (au *ActorUseCase) Create(actor *models.Actor) *errors.Error {
	grpcActor, err := au.adminPanelClient.CreateActor(context.Background(),
		admin.ActorModelToGRPC(actor))
	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	if err := copier.Copy(actor, admin.ActorGRPCToModel(grpcActor)); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (au *ActorUseCase) Change(newActor *models.Actor) *errors.Error {
	_, err := au.adminPanelClient.ChangeActor(context.Background(),
		admin.ActorModelToGRPC(newActor))

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (au *ActorUseCase) DeleteById(id uint64) *errors.Error {
	_, err := au.adminPanelClient.DeleteActorByID(context.Background(),
		&admin.ID{ID: id})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
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
