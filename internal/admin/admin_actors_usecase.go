package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (am *AdminMicroservice) CreateActor(ctx context.Context, actor *Actor) (*Actor, error) {
	modelActor := ActorGRPCToModel(actor)
	err := am.actorsRep.Insert(modelActor)
	if err != nil {
		return &Actor{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	actor.ID = modelActor.ID

	return actor, nil
}

func (am *AdminMicroservice) ChangeActor(ctx context.Context, newActor *Actor) (*empty.Empty, error) {
	if _, err := am.getActor(newActor.GetID()); err != nil {
		return &empty.Empty{}, err
	}

	if err := am.actorsRep.Update(ActorGRPCToModel(newActor)); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteActorByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	if _, err := am.getActor(id.GetID()); err != nil {
		return &empty.Empty{}, err
	}

	if err := am.actorsRep.DeleteById(id.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) getActor(id uint64) (*models.Actor, error) {
	dbActor, err := am.actorsRep.SelectById(id)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.Code(consts.CodeActorDoesNotExist), "")
	} else if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return dbActor, nil
}
