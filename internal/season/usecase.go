package season

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type SeasonUsecase interface {
	Create(season *models.Season) *errors.Error
	Change(season *models.Season) *errors.Error
	Get(id uint64) (*models.Season, *errors.Error)
	GetEpisodes(id uint64) ([]*models.Episode, *errors.Error)
	Delete(id uint64) *errors.Error
}
