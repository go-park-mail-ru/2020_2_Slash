package episode

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type EpisodeRepository interface {
	Insert(episode *models.Episode) error
	Update(newEpisode *models.Episode) error
	SelectByID(id uint64) (*models.Episode, error)
	SelectByNumberAndSeason(number int, seasonID uint64) (*models.Episode, error)
	SelectContentByID(id uint64) (*models.Content, error)
	SelectSeasonNumberByID(id uint64) (int, error)
	DeleteByID(id uint64) error
	UpdatePoster(episode *models.Episode) error
	UpdateVideo(episode *models.Episode) error
}
