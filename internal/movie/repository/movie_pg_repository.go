package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
)

type MoviePgRepository struct {
	dbConn *sql.DB
}

func NewMoviePgRepository(conn *sql.DB) movie.MovieRepository {
	return &MoviePgRepository{
		dbConn: conn,
	}
}

func (mr *MoviePgRepository) Insert(movie *models.Movie) error {
	tx, err := mr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		`INSERT INTO movies(video, content_id)
		VALUES ($1, $2)
		RETURNING id`,
		movie.Video, movie.ContentID)

	err = row.Scan(&movie.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (mr *MoviePgRepository) Update(movie *models.Movie) error {
	tx, err := mr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE movies
		SET video = $2, content_id = $3
		WHERE id = $1;`,
		movie.ID, movie.Video, movie.ContentID)
	if err != nil {
		tx.Rollback()
		return nil
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (mr *MoviePgRepository) DeleteByID(movieID uint64) error {
	tx, err := mr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM movies
		WHERE id=$1`,
		movieID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (mr *MoviePgRepository) SelectByID(movieID uint64) (*models.Movie, error) {
	movie := &models.Movie{}
	row := mr.dbConn.QueryRow(
		`SELECT id, video, content_id
		FROM movies
		WHERE id=$1`,
		movieID)

	if err := row.Scan(&movie.ID, &movie.Video, &movie.ContentID); err != nil {
		return nil, err
	}
	return movie, nil
}

func (mr *MoviePgRepository) SelectByContentID(contentID uint64) (*models.Movie, error) {
	movie := &models.Movie{}
	row := mr.dbConn.QueryRow(
		`SELECT id, video, content_id
		FROM movies
		WHERE content_id=$1`,
		contentID)

	if err := row.Scan(&movie.ID, &movie.Video, &movie.ContentID); err != nil {
		return nil, err
	}
	return movie, nil
}

func (mr *MoviePgRepository) SelectByGenre(genreID uint64) ([]*models.Movie, error) {
	rows, err := mr.dbConn.Query(
		`SELECT m.id, m.video, c.id, c.name, c.original_name, c.description, c.short_description,
		c.year, c.images, c.type
		FROM content AS c
		JOIN content_genre AS cg ON c.id=cg.content_id AND cg.genre_id=$1
		LEFT OUTER JOIN movies as m ON m.content_id=c.id`,
		genreID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
	for rows.Next() {
		movie := &models.Movie{}
		cnt := &models.Content{}

		err := rows.Scan(&movie.ID, &movie.Video, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
			&cnt.Year, &cnt.Images, &cnt.Type)
		if err != nil {
			return nil, err
		}

		movie.Content = *cnt
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}
