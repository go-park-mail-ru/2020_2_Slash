package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	"github.com/jinzhu/copier"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
)

type SeasonUsecase struct {
	rep              season.SeasonRepository
	tvShowUseCase    tvshow.TVShowUsecase
	adminPanelClient admin.AdminPanelClient
}

func NewSeasonUsecase(rep season.SeasonRepository,
	tvShowUseCase tvshow.TVShowUsecase, client admin.AdminPanelClient) season.SeasonUsecase {
	return &SeasonUsecase{
		rep:              rep,
		tvShowUseCase:    tvShowUseCase,
		adminPanelClient: client,
	}
}

func (uc *SeasonUsecase) Create(season *models.Season) *errors.Error {
	grcpSeason, err := uc.adminPanelClient.CreateSeason(context.Background(),
		admin.SeasonModelToGRPC(season))
	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	if err := copier.Copy(season, admin.SeasonGRPCToModel(grcpSeason)); err != nil {
		return errors.New(consts.CodeInternalError, err)
	}

	return nil
}

func (uc *SeasonUsecase) Change(newSeason *models.Season) *errors.Error {
	_, err := uc.adminPanelClient.ChangeSeason(context.Background(),
		admin.SeasonModelToGRPC(newSeason))

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (uc *SeasonUsecase) isExist(id uint64) (bool, *errors.Error) {
	_, err := uc.Get(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (uc *SeasonUsecase) Get(id uint64) (*models.Season, *errors.Error) {
	season, err := uc.rep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeSeasonDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return season, nil
}

func (uc *SeasonUsecase) GetEpisodes(id uint64) ([]*models.Episode, *errors.Error) {
	episodes, err := uc.rep.SelectEpisodes(id)
	if err == sql.ErrNoRows {
		return []*models.Episode{}, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	if episodes == nil {
		return []*models.Episode{}, nil
	}
	return episodes, nil
}

func (uc *SeasonUsecase) isConflicts(season *models.Season) (bool, *errors.Error) {
	_, err := uc.rep.Select(season)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.New(consts.CodeInternalError, err)
	}
	return true, nil
}

func (uc *SeasonUsecase) Delete(id uint64) *errors.Error {
	_, err := uc.adminPanelClient.DeleteSeasonsByID(context.Background(),
		&admin.ID{ID: id})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (uc *SeasonUsecase) ListByTVShow(tvshowID uint64) ([]*models.Season, *errors.Error) {
	seasons, err := uc.rep.SelectByTVShow(tvshowID)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	if len(seasons) == 0 {
		return []*models.Season{}, nil
	}
	return seasons, nil
}
