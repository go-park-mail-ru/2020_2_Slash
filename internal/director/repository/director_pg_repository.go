package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type DirectorPgRepository struct {
	dbConn *sql.DB
}

func NewDirectorPgRepository(conn *sql.DB) director.DirectorRepository {
	return &DirectorPgRepository{
		dbConn: conn,
	}
}

func (dr *DirectorPgRepository) Insert(director *models.Director) error {
	tx, err := dr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		`INSERT INTO directors(name)
		VALUES ($1)
		RETURNING id`,
		director.Name).Scan(&director.ID)
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

func (dr *DirectorPgRepository) Update(director *models.Director) error {
	tx, err := dr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE directors
		SET name = $2
		WHERE id = $1`,
		director.ID, director.Name)
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

func (dr *DirectorPgRepository) DeleteById(id uint64) error {
	tx, err := dr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM directors
		WHERE id = $1`, id)
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

func (dr *DirectorPgRepository) SelectById(id uint64) (*models.Director, error) {
	dbDirector := &models.Director{}

	row := dr.dbConn.QueryRow(
		`SELECT id, name
		FROM directors WHERE id=$1`, id)

	err := row.Scan(&dbDirector.ID, &dbDirector.Name)
	if err != nil {
		return nil, err
	}
	return dbDirector, nil
}

func (dr *DirectorPgRepository) SelectAll(pgnt *models.Pagination) ([]*models.Director, error) {
	var values []interface{}

	selectQuery := `
		SELECT id, name
		FROM directors
		ORDER BY name`

	var pgntQuery string
	if pgnt.Count != 0 {
		pgntQuery = "LIMIT $1 OFFSET $2"
		values = append(values, pgnt.Count, pgnt.From)
	}

	resultQuery := strings.Join([]string{
		selectQuery,
		pgntQuery,
	}, " ")

	rows, err := dr.dbConn.Query(resultQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var directors []*models.Director
	for rows.Next() {
		director := &models.Director{}
		err := rows.Scan(&director.ID, &director.Name)
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
