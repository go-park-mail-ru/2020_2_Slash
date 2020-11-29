package repository

import (
	"context"
	"database/sql"
	"strings"

	queryBuilder "github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/query_builder"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
)

type TVShowPgRepository struct {
	dbConn *sql.DB
}

func NewTVShowPgRepository(conn *sql.DB) tvshow.TVShowRepository {
	return &TVShowPgRepository{
		dbConn: conn,
	}
}

func (tr *TVShowPgRepository) Insert(tvshow *models.TVShow) error {
	tx, err := tr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		`INSERT INTO tv_shows(content_id)
		VALUES ($1)
		RETURNING id`,
		tvshow.ContentID)

	err = row.Scan(&tvshow.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (tr *TVShowPgRepository) SelectByID(tvshowID uint64) (*models.TVShow, error) {
	tvshow := &models.TVShow{}

	row := tr.dbConn.QueryRow(
		`SELECT id, seasons, content_id
		FROM tv_shows
		WHERE id=$1`,
		tvshowID)

	if err := row.Scan(&tvshow.ID, &tvshow.Seasons, &tvshow.ContentID); err != nil {
		return nil, err
	}
	return tvshow, nil
}

func (tr *TVShowPgRepository) SelectShortByID(tvshowID uint64) (*models.TVShow, error) {
	tvshow := &models.TVShow{}
	cnt := &models.Content{}

	row := tr.dbConn.QueryRow(
		`SELECT tv.id, c.name
		FROM content AS c
		JOIN tv_shows as tv ON tv.content_id=c.id AND tv.id=$1`,
		tvshowID)

	err := row.Scan(&tvshow.ID, &cnt.Name)
	if err != nil {
		return nil, err
	}

	tvshow.Content = *cnt
	return tvshow, nil
}

func (tr *TVShowPgRepository) SelectFullByID(tvshowID uint64, curUserID uint64) (*models.TVShow, error) {
	tvshow := &models.TVShow{}
	cnt := &models.Content{}

	row := tr.dbConn.QueryRow(
		`SELECT tv.id, tv.seasons, c.id, c.name, c.original_name, c.description, c.short_description,
		c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN tv_shows as tv ON tv.content_id=c.id AND tv.id=$1
		LEFT OUTER JOIN rates as r ON r.user_id=$2 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$2 AND f.content_id=c.id`,
		tvshowID, curUserID)

	err := row.Scan(&tvshow.ID, &tvshow.Seasons, &cnt.ContentID, &cnt.Name,
		&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
		&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)

	if err != nil {
		return nil, err
	}

	tvshow.Content = *cnt
	return tvshow, nil
}

func (tr *TVShowPgRepository) SelectByContentID(contentID uint64) (*models.TVShow, error) {
	tvshow := &models.TVShow{}
	row := tr.dbConn.QueryRow(
		`SELECT id, seasons, content_id
		FROM tv_shows
		WHERE content_id=$1`,
		contentID)

	if err := row.Scan(&tvshow.ID, &tvshow.Seasons, &tvshow.ContentID); err != nil {
		return nil, err
	}
	return tvshow, nil
}

func (tr *TVShowPgRepository) SelectWhereNameLike(name string,
	pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, error) {
	var values []interface{}

	selectQuery := `
		SELECT tv.id, tv.seasons, c.id, c.name, c.original_name,
		c.description, c.short_description, c.rating,
		c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN tv_shows as tv ON tv.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id
		WHERE c.name ILIKE $2 OR c.original_name ILIKE $2
		ORDER BY c.year DESC`
	values = append(values, curUserID)
	searchName := "%" + name + "%"
	values = append(values, searchName)

	var pgntQuery string
	if pgnt.Count != 0 {
		pgntQuery = "LIMIT $3 OFFSET $4"
		values = append(values, pgnt.Count, pgnt.From)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := tr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tvshows []*models.TVShow
	for rows.Next() {
		tvshow := &models.TVShow{}
		cnt := &models.Content{}

		err := rows.Scan(&tvshow.ID, &tvshow.Seasons, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription, &cnt.Rating,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)
		if err != nil {
			return nil, err
		}

		tvshow.Content = *cnt
		tvshows = append(tvshows, tvshow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tvshows, nil
}

func (tr *TVShowPgRepository) SelectByParams(params *models.ContentFilter,
	pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, error) {

	selectQuery := `
		SELECT tv.id, tv.seasons, c.id, c.name, c.original_name, c.description,
		c.short_description, c.rating, c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content as c`

	var values []interface{}

	joinTVShowQuery := `
		JOIN tv_shows as tv ON tv.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id`
	values = append(values, curUserID)

	var pgntQuery string
	if pgnt.Count != 0 {
		pgntQuery = "ORDER BY tv.id LIMIT $2 OFFSET $3"
		values = append(values, pgnt.Count, pgnt.From)
	}

	filtersJoinQuery, values := queryBuilder.GetContentJoinFiltersByParams(values, params)
	filtersWhereQuery, values := queryBuilder.GetContentWhereQueryByParams(values, params)

	resultQuery := strings.Join([]string{
		selectQuery,
		filtersJoinQuery,
		joinTVShowQuery,
		filtersWhereQuery,
		pgntQuery,
	}, " ")

	rows, err := tr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tvshows []*models.TVShow
	for rows.Next() {
		tvshow := &models.TVShow{}
		cnt := &models.Content{}

		err := rows.Scan(&tvshow.ID, &tvshow.Seasons, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription, &cnt.Rating,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)
		if err != nil {
			return nil, err
		}

		tvshow.Content = *cnt
		tvshows = append(tvshows, tvshow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tvshows, nil
}

func (tr *TVShowPgRepository) SelectLatest(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, error) {
	var values []interface{}

	selectQuery := `
		SELECT tv.id, tv.seasons, c.id, c.name, c.original_name,
		c.description, c.short_description, c.rating,
		c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN tv_shows as tv ON tv.content_id=c.id
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

	rows, err := tr.dbConn.Query(resultQuery, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tvshows []*models.TVShow
	for rows.Next() {
		tvshow := &models.TVShow{}
		cnt := &models.Content{}

		err := rows.Scan(&tvshow.ID, &tvshow.Seasons, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription, &cnt.Rating,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)
		if err != nil {
			return nil, err
		}

		tvshow.Content = *cnt
		tvshows = append(tvshows, tvshow)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tvshows, nil
}

func (tr *TVShowPgRepository) SelectByRating(pgnt *models.Pagination, curUserID uint64) ([]*models.TVShow, error) {
	var values []interface{}

	selectQuery := `
		SELECT tv.id, tv.seasons, c.id, c.name, c.original_name,
		c.description, c.short_description,
		c.rating, c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN tv_shows as tv ON tv.content_id=c.id
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

	rows, err := tr.dbConn.Query(resultQuery, values...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tvshows []*models.TVShow
	for rows.Next() {
		tvshow := &models.TVShow{}
		cnt := &models.Content{}

		err := rows.Scan(&tvshow.ID, &tvshow.Seasons, &cnt.ContentID, &cnt.Name,
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
			&cnt.Rating, &cnt.Year, &cnt.Images, &cnt.Type,
			&cnt.IsLiked, &cnt.IsFavourite)
		if err != nil {
			return nil, err
		}
		tvshow.Content = *cnt
		tvshows = append(tvshows, tvshow)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tvshows, nil
}
