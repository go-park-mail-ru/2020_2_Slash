package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"strings"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type FavouritePgRepository struct {
	dbConn *sql.DB
}

func NewFavouritePgRepository(conn *sql.DB) favourite.FavouriteRepository {
	return &FavouritePgRepository{
		dbConn: conn,
	}
}

func (rep *FavouritePgRepository) Insert(favourite *models.Favourite) error {
	tx, err := rep.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO favourites(user_id, content_id, created)
		VALUES ($1, $2, $3)`, favourite.UserID, favourite.ContentID, favourite.Created)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (rep *FavouritePgRepository) Select(favourite *models.Favourite) error {
	dbFavourite := &models.Favourite{}

	row := rep.dbConn.QueryRow(`
		SELECT user_id, content_id, created
		FROM favourites
		WHERE user_id=$1 AND content_id=$2`,
		favourite.UserID, favourite.ContentID)
	err := row.Scan(&dbFavourite.UserID, &dbFavourite.ContentID, &dbFavourite.Created)

	return err
}

func (rep *FavouritePgRepository) SelectFavouriteMovies(userID uint64,
	limit uint64, offset uint64) ([]*models.Movie, error) {
	var values []interface{}
	selectQuery := `
		SELECT m.id, m.video, c.id, c.name, c.original_name, c.description,
		c.short_description, c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN movies as m ON m.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id
		WHERE f.user_id=$1
		ORDER BY created DESC`
	values = append(values, userID)

	var pgntQuery string
	if limit != 0 {
		pgntQuery = "LIMIT $2 OFFSET $3"
		values = append(values, limit, offset)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := rep.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favouriteMovies []*models.Movie

	for rows.Next() {
		movie := &models.Movie{}
		cnt := &models.Content{}

		err := rows.Scan(&movie.ID, &movie.Video, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)
		if err != nil {
			return nil, err
		}

		movie.Content = *cnt
		favouriteMovies = append(favouriteMovies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favouriteMovies, nil
}

func (rep *FavouritePgRepository) SelectFavouriteTVShows(userID uint64,
	limit uint64, offset uint64) ([]*models.TVShow, error) {
	var values []interface{}
	selectQuery := `
		SELECT t.id, t.seasons, c.id, c.name, c.original_name, c.description,
		c.short_description, c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN tv_shows as t ON t.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id
		WHERE f.user_id=$1
		ORDER BY created DESC`
	values = append(values, userID)

	var pgntQuery string
	if limit != 0 {
		pgntQuery = "LIMIT $2 OFFSET $3"
		values = append(values, limit, offset)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := rep.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favouriteTvShows []*models.TVShow

	for rows.Next() {
		tvshow := &models.TVShow{}
		cnt := &models.Content{}

		err := rows.Scan(&tvshow.ID, &tvshow.Seasons, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)
		if err != nil {
			return nil, err
		}

		tvshow.Content = *cnt
		favouriteTvShows = append(favouriteTvShows, tvshow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favouriteTvShows, nil
}

func (rep *FavouritePgRepository) Delete(favourite *models.Favourite) error {
	tx, err := rep.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM favourites
		WHERE user_id=$1 AND content_id=$2`,
		favourite.UserID, favourite.ContentID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logger.Error(rollbackErr.Error())
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
