package season

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type SeasonRepository interface {
	Insert(season *models.Season) error
	Update(season *models.Season) error
	SelectByID(id uint64) (*models.Season, error)
	Select(season *models.Season) (*models.Season, error)
	SelectEpisodes(id uint64) ([]*models.Episode, error)
	Delete(id uint64) error
	SelectByTVShow(tvshowID uint64) ([]*models.Season, error)
}
