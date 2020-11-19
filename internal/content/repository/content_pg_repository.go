package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/lib/pq"
)

type ContentPgRepository struct {
	dbConn *sql.DB
}

func NewContentPgRepository(conn *sql.DB) content.ContentRepository {
	return &ContentPgRepository{
		dbConn: conn,
	}
}

func (cr *ContentPgRepository) Insert(content *models.Content) error {
	tx, err := cr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		`INSERT INTO content(name, original_name, description, short_description,
		year, images, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		content.Name, content.OriginalName, content.Description,
		content.ShortDescription, content.Year, content.Images, content.Type)

	err = row.Scan(&content.ContentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert countries
	if err := InsertCountries(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	// Insert genres
	if err := InsertGenres(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	// Insert actors
	if err := InsertActors(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	// Insert directors
	if err := InsertDirectors(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *ContentPgRepository) Update(content *models.Content) error {
	tx, err := cr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE content
		SET name = $2, original_name = $3, description = $4,
		short_description = $5, year = $6, images = $7, type = $8
		WHERE id = $1;`,
		content.ContentID, content.Name, content.OriginalName, content.Description,
		content.ShortDescription, content.Year, content.Images, content.Type)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Update countries
	if err := UpdateCountries(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	// Update genres
	if err := UpdateGenres(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	// Update actors
	if err := UpdateActors(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	// Update directors
	if err := UpdateDirectors(tx, content); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *ContentPgRepository) UpdateImages(content *models.Content) error {
	tx, err := cr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE content
		SET images = $2
		WHERE id = $1;`,
		content.ContentID, content.Images)

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *ContentPgRepository) DeleteByID(contentID uint64) error {
	tx, err := cr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM content
		WHERE id=$1`,
		contentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (cr *ContentPgRepository) SelectByID(contentID uint64) (*models.Content, error) {
	content := &models.Content{}

	row := cr.dbConn.QueryRow(
		`SELECT id, name, original_name, description, short_description,
		year, images, type
		FROM content
		WHERE id=$1`,
		contentID)

	err := row.Scan(&content.ContentID, &content.Name, &content.OriginalName, &content.Description,
		&content.ShortDescription, &content.Year, &content.Images, &content.Type)

	if err != nil {
		return nil, err
	}
	return content, nil
}

func (cr *ContentPgRepository) SelectCountriesByID(contentID uint64) ([]uint64, error) {
	var countries []uint64

	rows, err := cr.dbConn.Query(
		`SELECT country_id
		FROM content_country
		WHERE content_id=$1`,
		contentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var country uint64
		err := rows.Scan(&country)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return countries, nil
}

func (cr *ContentPgRepository) SelectGenresByID(contentID uint64) ([]uint64, error) {
	var genres []uint64

	rows, err := cr.dbConn.Query(
		`SELECT genre_id
		FROM content_genre
		WHERE content_id=$1`,
		contentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre uint64
		err := rows.Scan(&genre)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return genres, nil
}

func (cr *ContentPgRepository) SelectActorsByID(contentID uint64) ([]uint64, error) {
	var actors []uint64

	rows, err := cr.dbConn.Query(
		`SELECT actor_id
		FROM content_actor
		WHERE content_id=$1`,
		contentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor uint64
		err := rows.Scan(&actor)
		if err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return actors, nil
}

func (cr *ContentPgRepository) SelectDirectorsByID(contentID uint64) ([]uint64, error) {
	var directors []uint64

	rows, err := cr.dbConn.Query(
		`SELECT director_id
		FROM content_director
		WHERE content_id=$1`,
		contentID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var director uint64
		err := rows.Scan(&director)
		if err != nil {
			return nil, err
		}
		directors = append(directors, director)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return directors, nil
}

func InsertCountries(tx *sql.Tx, content *models.Content) error {
	stmt, err := tx.Prepare(pq.CopyIn("content_country", "content_id", "country_id"))
	if err != nil {
		return err
	}

	for _, country := range content.Countries {
		_, err = stmt.Exec(content.ContentID, country.ID)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	return nil
}

func InsertGenres(tx *sql.Tx, content *models.Content) error {
	stmt, err := tx.Prepare(pq.CopyIn("content_genre", "content_id", "genre_id"))
	if err != nil {
		return err
	}

	for _, genre := range content.Genres {
		_, err = stmt.Exec(content.ContentID, genre.ID)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	return nil
}

func InsertActors(tx *sql.Tx, content *models.Content) error {
	stmt, err := tx.Prepare(pq.CopyIn("content_actor", "content_id", "actor_id"))
	if err != nil {
		return err
	}

	for _, actor := range content.Actors {
		_, err = stmt.Exec(content.ContentID, actor.ID)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	return nil
}

func InsertDirectors(tx *sql.Tx, content *models.Content) error {
	stmt, err := tx.Prepare(pq.CopyIn("content_director", "content_id", "director_id"))
	if err != nil {
		return err
	}

	for _, director := range content.Directors {
		_, err = stmt.Exec(content.ContentID, director.ID)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}
	if err = stmt.Close(); err != nil {
		return err
	}
	return nil
}

func DeleteCountries(tx *sql.Tx, contentID uint64) error {
	_, err := tx.Exec(
		`DELETE FROM content_country
		WHERE content_id=$1`,
		contentID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteGenres(tx *sql.Tx, contentID uint64) error {
	_, err := tx.Exec(
		`DELETE FROM content_genre
		WHERE content_id=$1`,
		contentID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteActors(tx *sql.Tx, contentID uint64) error {
	_, err := tx.Exec(
		`DELETE FROM content_actor
		WHERE content_id=$1`,
		contentID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDirectors(tx *sql.Tx, contentID uint64) error {
	_, err := tx.Exec(
		`DELETE FROM content_director
		WHERE content_id=$1`,
		contentID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCountries(tx *sql.Tx, content *models.Content) error {
	if err := DeleteCountries(tx, content.ContentID); err != nil {
		return err
	}
	if err := InsertCountries(tx, content); err != nil {
		return err
	}
	return nil
}

func UpdateGenres(tx *sql.Tx, content *models.Content) error {
	if err := DeleteGenres(tx, content.ContentID); err != nil {
		return err
	}
	if err := InsertGenres(tx, content); err != nil {
		return err
	}
	return nil
}

func UpdateActors(tx *sql.Tx, content *models.Content) error {
	if err := DeleteActors(tx, content.ContentID); err != nil {
		return err
	}
	if err := InsertActors(tx, content); err != nil {
		return err
	}
	return nil
}

func UpdateDirectors(tx *sql.Tx, content *models.Content) error {
	if err := DeleteDirectors(tx, content.ContentID); err != nil {
		return err
	}
	if err := InsertDirectors(tx, content); err != nil {
		return err
	}
	return nil
}
