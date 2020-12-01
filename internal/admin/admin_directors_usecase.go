package admin

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (am *AdminMicroservice) CreateDirector(ctx context.Context, director *Director) (*Director, error) {
	modelDirector := DirectorGRPCToModel(director)
	err := am.directorsRep.Insert(modelDirector)
	if err != nil {
		return &Director{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	director.ID = modelDirector.ID

	return director, nil
}

func (am *AdminMicroservice) ChangeDirector(ctx context.Context, newDirector *Director) (*empty.Empty, error) {
	if _, err := am.getDirector(newDirector.GetID()); err != nil {
		return &empty.Empty{}, err
	}

	if err := am.directorsRep.Update(DirectorGRPCToModel(newDirector)); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteDirectorByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	if _, err := am.getDirector(id.GetID()); err != nil {
		return &empty.Empty{}, err
	}

	if err := am.directorsRep.DeleteById(id.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) getDirector(id uint64) (*models.Director, error) {
	dbDirector, err := am.directorsRep.SelectById(id)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.Code(consts.CodeDirectorDoesNotExist), "")
	} else if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return dbDirector, nil
}
