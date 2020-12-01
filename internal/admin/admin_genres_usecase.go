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

func (am *AdminMicroservice) CreateGenre(ctx context.Context, genre *Genre) (*Genre, error) {
	if err := am.checkGenreByName(genre.Name); err == nil {
		return &Genre{}, status.Error(codes.Code(consts.CodeGenreNameAlreadyExists), "")
	}

	modelGenre := GenreGRPCToModel(genre)
	if err := am.genresRep.Insert(modelGenre); err != nil {
		return &Genre{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	genre.ID = modelGenre.ID

	return genre, nil
}

func (am *AdminMicroservice) ChangeGenre(ctx context.Context, newGenreData *Genre) (*empty.Empty, error) {
	modelGenre, err := am.getGenreByID(newGenreData.GetID())
	if err != nil {
		return &empty.Empty{}, err
	}

	if newGenreData.Name == "" || modelGenre.Name == newGenreData.Name {
		return &empty.Empty{}, nil
	}

	if err := am.checkGenreByName(newGenreData.Name); err == nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeGenreNameAlreadyExists), "")
	}

	modelGenre.Name = newGenreData.Name
	if err := am.genresRep.Update(modelGenre); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteGenreByID(ctx context.Context, genreID *ID) (*empty.Empty, error) {
	if err := am.checkGenreByID(genreID.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeGenreDoesNotExist), "")
	}

	if err := am.genresRep.DeleteByID(genreID.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) getGenreByID(genreID uint64) (*models.Genre, error) {
	genre, err := am.genresRep.SelectByID(genreID)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeGenreDoesNotExist), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return genre, nil
}

func (am *AdminMicroservice) getGenreByName(name string) (*models.Genre, error) {
	genre, err := am.genresRep.SelectByName(name)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeGenreNameAlreadyExists), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return genre, nil
}

func (am *AdminMicroservice) checkGenreByID(genreID uint64) error {
	_, err := am.getGenreByID(genreID)
	return err
}

func (am *AdminMicroservice) checkGenreByName(name string) error {
	_, err := am.getGenreByName(name)
	return err
}
