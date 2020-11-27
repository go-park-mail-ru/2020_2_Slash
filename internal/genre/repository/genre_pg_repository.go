package repository

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type GenrePgRepository struct {
	dbConn *sql.DB
}

func NewGenrePgRepository(conn *sql.DB) genre.GenreRepository {
	return &GenrePgRepository{
		dbConn: conn,
	}
}

func (gr *GenrePgRepository) Insert(genre *models.Genre) error {
	tx, err := gr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		`INSERT INTO genres(name)
		VALUES ($1)
		RETURNING id`,
		genre.Name)

	err = row.Scan(&genre.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (gr *GenrePgRepository) Update(genre *models.Genre) error {
	tx, err := gr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE genres
		SET name = $2
		WHERE id = $1;`,
		genre.ID, genre.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (gr *GenrePgRepository) DeleteByID(genreID uint64) error {
	tx, err := gr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM genres
		WHERE id=$1`,
		genreID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (gr *GenrePgRepository) SelectByName(genreName string) (*models.Genre, error) {
	genre := &models.Genre{}

	row := gr.dbConn.QueryRow(
		`SELECT id, name
		FROM genres
		WHERE name=$1`,
		genreName)

	if err := row.Scan(&genre.ID, &genre.Name); err != nil {
		return nil, err
	}
	return genre, nil
}

func (gr *GenrePgRepository) SelectByID(genreID uint64) (*models.Genre, error) {
	genre := &models.Genre{}

	row := gr.dbConn.QueryRow(
		`SELECT id, name
		FROM genres
		WHERE id=$1`,
		genreID)

	if err := row.Scan(&genre.ID, &genre.Name); err != nil {
		return nil, err
	}
	return genre, nil
}

func (gr *GenrePgRepository) SelectAll() ([]*models.Genre, error) {
	rows, err := gr.dbConn.Query(
		`SELECT id, name
		FROM genres`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		genre := &models.Genre{}
		err := rows.Scan(&genre.ID, &genre.Name)
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
