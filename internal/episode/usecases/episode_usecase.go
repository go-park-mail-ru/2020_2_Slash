package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/jinzhu/copier"
)

type EpisodeUsecase struct {
	rep              episode.EpisodeRepository
	seasonUseCase    season.SeasonUsecase
	adminPanelClient admin.AdminPanelClient
}

func NewEpisodeUsecase(rep episode.EpisodeRepository,
	seasonUseCase season.SeasonUsecase, client admin.AdminPanelClient) episode.EpisodeUsecase {
	return &EpisodeUsecase{
		rep:              rep,
		seasonUseCase:    seasonUseCase,
		adminPanelClient: client,
	}
}

func (uc *EpisodeUsecase) Create(episode *models.Episode) *errors.Error {
	grcpEpisode, err := uc.adminPanelClient.CreateEpisode(context.Background(),
		admin.EpisodeModelToGRPC(episode))
	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	if err := copier.Copy(episode, admin.EpisodeGRPCToModel(grcpEpisode)); err != nil {
		return errors.New(consts.CodeInternalError, err)
	}

	return nil
}

func (uc *EpisodeUsecase) Change(episode *models.Episode) *errors.Error {
	_, err := uc.adminPanelClient.ChangeEpisode(context.Background(),
		admin.EpisodeModelToGRPC(episode))

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (uc *EpisodeUsecase) UpdatePoster(episode *models.Episode, newPosterPath string) *errors.Error {
	_, err := uc.adminPanelClient.UpdatePoster(context.Background(),
		&admin.EpisodePostersDir{
			Episode:    admin.EpisodeModelToGRPC(episode),
			PostersDir: newPosterPath,
		})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (uc *EpisodeUsecase) UpdateVideo(episode *models.Episode, newVideoPath string) *errors.Error {
	_, err := uc.adminPanelClient.UpdateVideo(context.Background(),
		&admin.EpisodeVideo{
			Episode: admin.EpisodeModelToGRPC(episode),
			Video:   newVideoPath,
		})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
}

func (uc *EpisodeUsecase) DeleteByID(id uint64) *errors.Error {
	_, err := uc.adminPanelClient.DeleteEpisodeByID(context.Background(),
		&admin.ID{ID: id})

	if err != nil {
		customErr := errors.GetCustomErr(err)
		return customErr
	}

	return nil
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
