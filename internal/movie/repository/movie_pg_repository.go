package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"

	queryBuilder "github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/query_builder"
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

func (mr *MoviePgRepository) SelectFullByID(movieID uint64, curUserID uint64) (*models.Movie, error) {
	movie := &models.Movie{}
	cnt := &models.Content{}

	row := mr.dbConn.QueryRow(
		`SELECT m.id, m.video, c.id, c.name, c.original_name, c.description, c.short_description,
		c.year, c.images, c.type, c.is_free, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN movies as m ON m.content_id=c.id AND m.id=$1
		LEFT OUTER JOIN rates as r ON r.user_id=$2 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$2 AND f.content_id=c.id`,
		movieID, curUserID)

	err := row.Scan(&movie.ID, &movie.Video, &cnt.ContentID, &cnt.Name,
		&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
		&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsFree, &cnt.IsLiked, &cnt.IsFavourite)

	if err != nil {
		return nil, err
	}

	movie.Content = *cnt
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

func (mr *MoviePgRepository) SelectByParams(params *models.ContentFilter,
	pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, error) {

	selectQuery := `
		SELECT m.id, m.video, c.id, c.name, c.original_name, c.description,
		c.short_description, c.rating, c.year, c.images, c.type, c.is_free, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content as c`

	var values []interface{}

	joinUserQuery := `
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id`
	values = append(values, curUserID)

	pgntQuery := "ORDER BY m.id"
	if pgnt.Count != 0 {
		pgntQuery += " LIMIT $2 OFFSET $3"
		values = append(values, pgnt.Count, pgnt.From)
	}

	joinMovieQuery := "JOIN movies as m ON m.content_id=c.id"
	if params.IsFree != nil {
		ind := len(values) + 1
		joinMovieQuery = fmt.Sprintf("%s AND c.is_free=$%d", joinMovieQuery, ind)
		values = append(values, params.IsFree)
	}

	filtersJoinQuery, values := queryBuilder.GetContentJoinFiltersByParams(values, params)

	resultQuery := strings.Join([]string{
		selectQuery,
		joinMovieQuery,
		filtersJoinQuery,
		joinUserQuery,
		pgntQuery,
	}, " ")

	rows, err := mr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
	for rows.Next() {
		movie := &models.Movie{}
		cnt := &models.Content{}

		err := rows.Scan(&movie.ID, &movie.Video, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription, &cnt.Rating,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsFree, &cnt.IsLiked, &cnt.IsFavourite)
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

func (mr *MoviePgRepository) SelectLatest(pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, error) {
	var values []interface{}

	selectQuery := `
		SELECT m.id, m.video, c.id, c.name, c.original_name,
		c.description, c.short_description, c.rating,
		c.year, c.images, c.type, c.is_free, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN movies as m ON m.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id
		ORDER BY c.year DESC`
	values = append(values, curUserID)

	var pgntQuery string
	if pgnt.Count != 0 {
		pgntQuery = "LIMIT $2 OFFSET $3"
		values = append(values, pgnt.Count, pgnt.From)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := mr.dbConn.Query(resultQuery, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
	for rows.Next() {
		movie := &models.Movie{}
		cnt := &models.Content{}

		err := rows.Scan(&movie.ID, &movie.Video, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription, &cnt.Rating,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsFree, &cnt.IsLiked, &cnt.IsFavourite)
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

func (mr *MoviePgRepository) SelectByRating(pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, error) {
	var values []interface{}

	selectQuery := `
		SELECT m.id, m.video, c.id, c.name, c.original_name,
		c.description, c.short_description,
		c.rating, c.year, c.images, c.type, c.is_free, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN movies as m ON m.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id
		ORDER BY c.rating DESC`
	values = append(values, curUserID)

	var pgntQuery string
	if pgnt.Count != 0 {
		pgntQuery = "LIMIT $2 OFFSET $3"
		values = append(values, pgnt.Count, pgnt.From)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := mr.dbConn.Query(resultQuery, values...)

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
			&cnt.Rating, &cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsFree,
			&cnt.IsLiked, &cnt.IsFavourite)
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

func (mr *MoviePgRepository) SelectWhereNameLike(curUserID uint64, name string, limit, offset uint64) ([]*models.Movie, error) {
	selectQuery := `
		SELECT m.id, m.video, c.id, c.name, c.original_name,
		c.description, c.short_description,
		c.rating, c.year, c.images, c.type, c.is_free, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN movies as m ON m.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id
		WHERE c.name ILIKE $2 OR c.original_name ILIKE $2
		ORDER BY m.id`

	var values []interface{}
	values = append(values, curUserID)
	searchName := "%" + name + "%"
	values = append(values, searchName)

	var pgntQuery string
	if limit != 0 {
		pgntQuery = "LIMIT $3 OFFSET $4"
		values = append(values, limit, offset)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := mr.dbConn.Query(resultQuery, values...)
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
			&cnt.Rating, &cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsFree,
			&cnt.IsLiked, &cnt.IsFavourite)
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
