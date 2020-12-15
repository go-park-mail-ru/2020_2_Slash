package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type EpisodeRepository struct {
	db *sql.DB
}

func NewEpisodeRepository(db *sql.DB) episode.EpisodeRepository {
	return &EpisodeRepository{db: db}
}

func (rep *EpisodeRepository) Insert(episode *models.Episode) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(`
		INSERT INTO episodes(number, name, video, description, poster, season_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		episode.Number, episode.Name, episode.Video,
		episode.Description, episode.Poster, episode.SeasonID).Scan(&episode.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *EpisodeRepository) SelectByID(id uint64) (*models.Episode, error) {
	dbEpisode := &models.Episode{}

	row := rep.db.QueryRow(`
		SELECT id, number, name, video, description, poster, season_id
		FROM episodes
		WHERE id=$1`, id)
	err := row.Scan(&dbEpisode.ID, &dbEpisode.Number, &dbEpisode.Name, &dbEpisode.Video,
		&dbEpisode.Description, &dbEpisode.Poster, &dbEpisode.SeasonID)
	if err != nil {
		return nil, err
	}

	return dbEpisode, nil
}

func (rep *EpisodeRepository) SelectByNumberAndSeason(number int,
	seasonID uint64) (*models.Episode, error) {
	dbEpisode := &models.Episode{}

	row := rep.db.QueryRow(`
		SELECT id, number, name, video, description, poster, season_id
		FROM episodes
		WHERE number=$1 AND season_id=$2`, number, seasonID)
	err := row.Scan(&dbEpisode.ID, &dbEpisode.Number, &dbEpisode.Name, &dbEpisode.Video,
		&dbEpisode.Description, &dbEpisode.Poster, &dbEpisode.SeasonID)
	if err != nil {
		return nil, err
	}

	return dbEpisode, nil
}

func (rep *EpisodeRepository) Update(newEpisode *models.Episode) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE episodes
		SET number = $1,
		    name = $2,
		    video = $3,
		    description = $4,
		    poster = $5,
		    season_id = $6
		WHERE id=$7`, newEpisode.Number, newEpisode.Name, newEpisode.Video,
		newEpisode.Description, newEpisode.Poster, newEpisode.SeasonID, newEpisode.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *EpisodeRepository) DeleteByID(id uint64) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM episodes
		WHERE id=$1`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *EpisodeRepository) UpdatePoster(episode *models.Episode) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE episodes
		SET poster=$1
		WHERE id=$2`, episode.Poster, episode.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *EpisodeRepository) UpdateVideo(episode *models.Episode) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE episodes
		SET video=$1
		WHERE id=$2`, episode.Video, episode.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *EpisodeRepository) SelectContentByID(id uint64) (*models.Content, error) {
	content := &models.Content{}
	row := rep.db.QueryRow(`
		SELECT content_id, c.name, c.original_name, c.description, c.short_description,
		       c.rating, c.year, c.images, c.type
		FROM episodes
		LEFT OUTER JOIN seasons ON season_id=seasons.id
		LEFT OUTER JOIN tv_shows on tv_show_id=tv_shows.id
		LEFT OUTER JOIN content c on tv_shows.content_id = c.id
		WHERE episodes.id=$1`, id)
	err := row.Scan(&content.ContentID, &content.Name, &content.OriginalName,
		&content.Description, &content.ShortDescription,
		&content.Rating, &content.Year, &content.Images, &content.Type)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (rep *EpisodeRepository) SelectSeasonNumberByID(id uint64) (int, error) {
	seasonNumber := 0
	row := rep.db.QueryRow(`
		SELECT s.number
		FROM episodes e
		LEFT OUTER JOIN seasons s ON e.season_id=s.id
		WHERE e.id=$1`, id)
	err := row.Scan(&seasonNumber)
	if err != nil {
		return 0, err
	}
	return seasonNumber, nil
}
