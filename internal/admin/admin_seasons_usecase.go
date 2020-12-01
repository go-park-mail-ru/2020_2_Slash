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

func (am *AdminMicroservice) CreateSeason(ctx context.Context, season *Season) (*Season, error) {
	_, err := am.GetTVShowByID(season.GetTVShowID())
	if err != nil {
		return &Season{}, err
	}

	modelSeason := SeasonGRPCToModel(season)
	isConflicts, err := am.IsSeasonConflicts(modelSeason)
	if err != nil {
		return &Season{}, err
	}
	if isConflicts {
		return &Season{}, status.Error(codes.Code(consts.CodeSeasonAlreadyExist), err.Error())
	}

	err = am.seasonsRep.Insert(modelSeason)
	if err != nil {
		return &Season{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	season.ID = modelSeason.ID

	return season, nil
}

func (am *AdminMicroservice) ChangeSeason(ctx context.Context, season *Season) (*empty.Empty, error) {
	_, err := am.GetTVShowByID(season.TVShowID)
	if err != nil {
		return &empty.Empty{}, err
	}

	seasonDB, err := am.GetSeason(season.ID)
	if err != nil {
		return &empty.Empty{}, err
	}
	if int64(seasonDB.Number) == season.Number &&
		seasonDB.TVShowID == season.TVShowID {
		return &empty.Empty{}, nil
	}

	modelSeason := SeasonGRPCToModel(season)

	isConflicts, err := am.IsSeasonConflicts(modelSeason)
	if err != nil {
		return &empty.Empty{}, err
	}
	if isConflicts {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeSeasonAlreadyExist), "")
	}

	err = am.seasonsRep.Update(modelSeason)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) DeleteSeasonsByID(ctx context.Context, id *ID) (*empty.Empty, error) {
	if _, err := am.isSeasonExist(id.GetID()); err != nil {
		return &empty.Empty{}, err
	}

	if err := am.seasonsRep.Delete(id.GetID()); err != nil {
		return &empty.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}

	return &empty.Empty{}, nil
}

func (am *AdminMicroservice) GetTVShowByID(tvshowID uint64) (*models.TVShow, error) {
	tvshow, err := am.tvshowsRep.SelectByID(tvshowID)
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeTVShowDoesNotExist), err.Error())
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return tvshow, nil
}

func (am *AdminMicroservice) IsSeasonConflicts(season *models.Season) (bool, error) {
	_, err := am.seasonsRep.Select(season)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, status.Error(codes.Code(consts.CodeTVShowDoesNotExist), err.Error())
	}
	return true, nil
}

func (am *AdminMicroservice) GetSeason(id uint64) (*models.Season, error) {
	season, err := am.seasonsRep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.Code(consts.CodeSeasonDoesNotExist), err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return season, nil
}

func (am *AdminMicroservice) isSeasonExist(id uint64) (bool, error) {
	_, err := am.GetSeason(id)
	if err != nil {
		return false, err
	}
	return true, nil
}
