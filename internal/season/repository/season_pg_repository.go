package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
)

type SeasonPgRepository struct {
	db *sql.DB
}

func NewSeasonPgRepository(db *sql.DB) season.SeasonRepository {
	return &SeasonPgRepository{db: db}
}

func (rep *SeasonPgRepository) Insert(season *models.Season) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(`
		INSERT INTO seasons(number, episodes, tv_show_id) 
		VALUES ($1, $2, $3)
		RETURNING id`, season.Number, season.EpisodesNumber, season.TVShowID)
	err = row.Scan(&season.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *SeasonPgRepository) Update(season *models.Season) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE seasons
		SET number=$1,
		    episodes=$2,
		    tv_show_id=$3
		WHERE id=$4`, season.Number, season.EpisodesNumber,
		season.TVShowID, season.ID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *SeasonPgRepository) SelectByID(id uint64) (*models.Season, error) {
	season := &models.Season{}
	row := rep.db.QueryRow(`
		SELECT id, number, episodes, tv_show_id
		FROM seasons
		WHERE id=$1`, id)
	err := row.Scan(&season.ID, &season.Number, &season.EpisodesNumber, &season.TVShowID)
	if err != nil {
		return nil, err
	}
	return season, nil
}

func (rep *SeasonPgRepository) Delete(id uint64) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE
		FROM seasons
		WHERE id=$1`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *SeasonPgRepository) Select(season *models.Season) (*models.Season, error) {
	dbSeason := &models.Season{}
	row := rep.db.QueryRow(`
		SELECT id, number, episodes, tv_show_id
		FROM seasons
		WHERE number=$1 AND tv_show_id=$2`, season.Number, season.TVShowID)
	err := row.Scan(&dbSeason.ID, &dbSeason.Number, &dbSeason.EpisodesNumber, &dbSeason.TVShowID)
	if err != nil {
		return nil, err
	}
	return dbSeason, nil
}

func (rep *SeasonPgRepository) SelectEpisodes(id uint64) ([]*models.Episode, error) {
	rows, err := rep.db.Query(`
		SELECT id, number, name, video, description, poster, season_id
		FROM episodes
		WHERE season_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []*models.Episode

	for rows.Next() {
		episode := &models.Episode{}

		err := rows.Scan(&episode.ID, &episode.Number, &episode.Name,
			&episode.Video, &episode.Description,
			&episode.Poster, &episode.SeasonID)
		if err != nil {
			return nil, err
		}

		episodes = append(episodes, episode)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return episodes, nil
}
