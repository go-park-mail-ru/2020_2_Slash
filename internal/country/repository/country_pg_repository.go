package repository

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
)

type CountryPgRepository struct {
	dbConn *sql.DB
}

func NewCountryPgRepository(conn *sql.DB) country.CountryRepository {
	return &CountryPgRepository{
		dbConn: conn,
	}
}

func (cr *CountryPgRepository) Insert(country *models.Country) error {
	tx, err := cr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		`INSERT INTO countries(name)
		VALUES ($1)
		RETURNING id`,
		country.Name)

	err = row.Scan(&country.ID)
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

func (cr *CountryPgRepository) Update(country *models.Country) error {
	tx, err := cr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE countries
		SET name = $2
		WHERE id = $1;`,
		country.ID, country.Name)
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

func (cr *CountryPgRepository) DeleteByID(countryID uint64) error {
	tx, err := cr.dbConn.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`DELETE FROM countries
		WHERE id=$1`,
		countryID)
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

func (cr *CountryPgRepository) SelectByName(countryName string) (*models.Country, error) {
	country := &models.Country{}

	row := cr.dbConn.QueryRow(
		`SELECT id, name
		FROM countries
		WHERE name=$1`,
		countryName)

	if err := row.Scan(&country.ID, &country.Name); err != nil {
		return nil, err
	}
	return country, nil
}

func (cr *CountryPgRepository) SelectByID(countryID uint64) (*models.Country, error) {
	country := &models.Country{}

	row := cr.dbConn.QueryRow(
		`SELECT id, name
		FROM countries
		WHERE id=$1`,
		countryID)

	if err := row.Scan(&country.ID, &country.Name); err != nil {
		return nil, err
	}
	return country, nil
}

func (cr *CountryPgRepository) SelectAll() ([]*models.Country, error) {
	rows, err := cr.dbConn.Query(
		`SELECT id, name
		FROM countries`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []*models.Country
	for rows.Next() {
		country := &models.Country{}
		err := rows.Scan(&country.ID, &country.Name)
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
