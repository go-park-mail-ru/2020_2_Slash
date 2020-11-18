package episode

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type EpisodeUsecase interface {
	Create(episode *models.Episode) *errors.Error
	Change(episode *models.Episode) *errors.Error
	GetByID(id uint64) (*models.Episode, *errors.Error)
	DeleteByID(id uint64) *errors.Error
	GetContentByEID(eid uint64) (*models.Content, *errors.Error)
	GetSeasonNumber(eid uint64) (int, *errors.Error)
	UpdatePoster(episode *models.Episode, posters string) *errors.Error
	UpdateVideo(episode *models.Episode, video string) *errors.Error
}
