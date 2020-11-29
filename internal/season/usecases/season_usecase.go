package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
)

type SeasonUsecase struct {
	rep           season.SeasonRepository
	tvShowUseCase tvshow.TVShowUsecase
}

func NewSeasonUsecase(rep season.SeasonRepository,
	tvShowUseCase tvshow.TVShowUsecase) season.SeasonUsecase {
	return &SeasonUsecase{
		rep:           rep,
		tvShowUseCase: tvShowUseCase,
	}
}

func (uc *SeasonUsecase) Create(season *models.Season) *errors.Error {
	_, customErr := uc.tvShowUseCase.GetByID(season.TVShowID)
	if customErr != nil {
		return customErr
	}

	isConflicts, customErr := uc.isConflicts(season)
	if customErr != nil {
		return customErr
	}
	if isConflicts {
		return errors.Get(consts.CodeSeasonAlreadyExist)
	}

	err := uc.rep.Insert(season)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *SeasonUsecase) Change(season *models.Season) *errors.Error {
	_, customErr := uc.tvShowUseCase.GetByID(season.TVShowID)
	if customErr != nil {
		return customErr
	}

	seasonDB, customErr := uc.Get(season.ID)
	if customErr != nil {
		return customErr
	}
	if seasonDB.Number == season.Number &&
		seasonDB.TVShowID == season.TVShowID {
		return nil
	}

	isConflicts, customErr := uc.isConflicts(season)
	if customErr != nil {
		return customErr
	}
	if isConflicts {
		return errors.Get(consts.CodeSeasonAlreadyExist)
	}

	err := uc.rep.Update(season)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
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
	if _, customErr := uc.isExist(id); customErr != nil {
		return customErr
	}

	if err := uc.rep.Delete(id); err != nil {
		return errors.New(consts.CodeInternalError, err)
	}

	return nil
}
