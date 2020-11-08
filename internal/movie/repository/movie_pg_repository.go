package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"strconv"
	"strings"
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

func (mr *MoviePgRepository) SelectFullByID(movieID uint64, curUserID uint64) (*models.Movie, error) {
	movie := &models.Movie{}
	cnt := &models.Content{}

	row := mr.dbConn.QueryRow(
		`SELECT m.id, m.video, c.id, c.name, c.original_name, c.description, c.short_description,
		c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content AS c
		JOIN movies as m ON m.content_id=c.id AND m.id=$1
		LEFT OUTER JOIN rates as r ON r.user_id=$2 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$2 AND f.content_id=c.id`,
		movieID, curUserID)

	err := row.Scan(&movie.ID, &movie.Video, &cnt.ContentID, &cnt.Name,
		&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
		&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)

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

func buildJoinContentFilter(entity string, valInd int) string {
	entityTable := "content_" + entity                  // content_genre
	entityID := entityTable + "." + entity + "_id"      // content.genre_id
	entityContentID := entityTable + "." + "content_id" // content_genre.content_id

	// JOIN content_genre ON c.id=cg.content_id AND cg.genre_id=$1
	filter := "JOIN " + entityTable + " ON " + "c.id=" +
		entityContentID + " AND " + entityID + "=$" + strconv.Itoa(valInd)

	return filter
}

func getJoinFiltersByParams(values []interface{}, params *models.ContentFilter) (string, []interface{}) {
	var filters []string

	if params.Genre != 0 {
		filter := buildJoinContentFilter("genre", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Genre)
	}

	if params.Country != 0 {
		filter := buildJoinContentFilter("country", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Country)
	}

	if params.Actor != 0 {
		filter := buildJoinContentFilter("actor", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Actor)
	}

	if params.Director != 0 {
		filter := buildJoinContentFilter("director", len(values)+1)
		filters = append(filters, filter)
		values = append(values, params.Director)
	}

	filtersQuery := strings.Join(filters, " ")
	return filtersQuery, values
}

func getWhereQueryByParams(values []interface{}, params *models.ContentFilter) (string, []interface{}) {
	var filters []string

	if params.Year != 0 {
		ind := len(values) + 1
		filter := `WHERE c.year=$` + strconv.Itoa(ind)
		filters = append(filters, filter)
		values = append(values, params.Year)
	}

	filtersQuery := strings.Join(filters, " ")
	return filtersQuery, values
}

func (mr *MoviePgRepository) SelectByParams(params *models.ContentFilter,
	pgnt *models.Pagination, curUserID uint64) ([]*models.Movie, error) {

	selectQuery := `SELECT m.id, m.video, c.id, c.name, c.original_name, c.description,
		c.short_description, c.year, c.images, c.type, r.likes,
		CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
		FROM content as c`

	var values []interface{}

	joinMovieQuery := `LEFT OUTER JOIN movies as m ON m.content_id=c.id
		LEFT OUTER JOIN rates as r ON r.user_id=$1 AND r.content_id=c.id
		LEFT OUTER JOIN favourites as f ON f.user_id=$1 AND f.content_id=c.id`
	values = append(values, curUserID)

	var pgntQuery string
	if pgnt.Count != 0 {
		pgntQuery = "ORDER BY m.id LIMIT $2 OFFSET $3"
		values = append(values, pgnt.Count, pgnt.From)
	}

	filtersJoinQuery, values := getJoinFiltersByParams(values, params)
	filtersWhereQuery, values := getWhereQueryByParams(values, params)

	resultQuery := strings.Join([]string{
		selectQuery,
		filtersJoinQuery,
		joinMovieQuery,
		filtersWhereQuery,
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
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)
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

	selectQuery := `SELECT m.id, m.video, c.id, c.name, c.original_name, c.description, c.short_description,
	c.year, c.images, c.type, r.likes,
	CASE WHEN f.content_id IS NULL THEN false ELSE true END AS is_favourite
	FROM content AS c
	LEFT OUTER JOIN movies as m ON m.content_id=c.id
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
			&cnt.OriginalName, &cnt.Description, &cnt.ShortDescription,
			&cnt.Year, &cnt.Images, &cnt.Type, &cnt.IsLiked, &cnt.IsFavourite)
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
