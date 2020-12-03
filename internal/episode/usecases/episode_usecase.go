package usecases

import (
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
)

type EpisodeUsecase struct {
	rep           episode.EpisodeRepository
	seasonUseCase season.SeasonUsecase
}

func NewEpisodeUsecase(rep episode.EpisodeRepository, seasonUseCase season.SeasonUsecase) episode.EpisodeUsecase {
	return &EpisodeUsecase{
		rep:           rep,
		seasonUseCase: seasonUseCase,
	}
}

func (uc *EpisodeUsecase) Create(episode *models.Episode) *errors.Error {
	_, customErr := uc.seasonUseCase.Get(episode.SeasonID)
	if customErr != nil {
		return customErr
	}

	isConflict, customErr := uc.isConflict(episode)
	if customErr != nil {
		return customErr
	}
	if isConflict {
		return errors.Get(consts.CodeEpisodeAlreadyExist)
	}

	err := uc.rep.Insert(episode)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *EpisodeUsecase) Change(episode *models.Episode) *errors.Error {
	_, customErr := uc.seasonUseCase.Get(episode.SeasonID)
	if customErr != nil {
		return customErr
	}

	isExist, customErr := uc.isExist(episode)
	if customErr != nil {
		return customErr
	}
	if !isExist {
		return errors.Get(consts.CodeEpisodeDoesNotExist)
	}

	err := uc.rep.Update(episode)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *EpisodeUsecase) isConflict(episode *models.Episode) (bool, *errors.Error) {
	_, err := uc.rep.SelectByNumberAndSeason(episode.Number, episode.SeasonID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.New(consts.CodeInternalError, err)
	}
	return true, nil
}

func (uc *EpisodeUsecase) isExist(episode *models.Episode) (bool, *errors.Error) {
	_, err := uc.rep.SelectByID(episode.ID)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, errors.New(consts.CodeInternalError, err)
	}
	return true, nil
}

func (uc *EpisodeUsecase) GetByID(id uint64) (*models.Episode, *errors.Error) {
	dbEpisode, err := uc.rep.SelectByID(id)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeEpisodeDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return dbEpisode, nil
}

func (uc *EpisodeUsecase) DeleteByID(id uint64) *errors.Error {
	_, customErr := uc.GetByID(id)
	if customErr != nil {
		return customErr
	}

	err := uc.rep.DeleteByID(id)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *EpisodeUsecase) GetContentByEID(eid uint64) (*models.Content, *errors.Error) {
	content, err := uc.rep.SelectContentByID(eid)
	if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return content, nil
}

func (uc *EpisodeUsecase) GetSeasonNumber(eid uint64) (int, *errors.Error) {
	seasonNumber, err := uc.rep.SelectSeasonNumberByID(eid)
	if err != nil {
		return 0, errors.New(consts.CodeInternalError, err)
	}
	return seasonNumber, nil
}

func (uc *EpisodeUsecase) UpdatePoster(episode *models.Episode, newPosterPath string) *errors.Error {
	prevPosterPath := episode.Poster
	if newPosterPath == prevPosterPath {
		return nil
	}

	episode.Poster = newPosterPath
	if err := uc.rep.UpdatePoster(episode); err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *EpisodeUsecase) UpdateVideo(episode *models.Episode, newVideoPath string) *errors.Error {
	prevVideoPath := episode.Video
	if newVideoPath == prevVideoPath {
		return nil
	}

	episode.Video = newVideoPath
	if err := uc.rep.UpdateVideo(episode); err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}
