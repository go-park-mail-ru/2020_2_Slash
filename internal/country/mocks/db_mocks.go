package mocks

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockCountryRepoInsertReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery(`INSERT INTO countries`).
		WithArgs(name).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockCountryRepoInsertReturnErrNoUniq(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO countries`).
		WithArgs(name).
		WillReturnError(errors.New("No UNIQUE"))
	mock.ExpectRollback()
}

func MockCountryRepoUpdateReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE countries`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func MockCountryRepoUpdateReturnResultZero(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE countries`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func MockCountryRepoDeleteReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM countries`).
		WithArgs(id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockCountryRepoSelectByIDReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(id, name)
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnRows(rows)
}

func MockCountryRepoSelectByIDReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnError(sql.ErrNoRows)
}

func MockCountryRepoSelectByNameReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(id, name)
	mock.ExpectQuery(`SELECT`).WithArgs(name).WillReturnRows(rows)
}

func MockCountryRepoSelectByNameReturnErrNoRows(mock sqlmock.Sqlmock, name string) {
	mock.ExpectQuery(`SELECT`).WithArgs(name).WillReturnError(sql.ErrNoRows)
}

func MockCountryRepoSelectAllReturnRows(mock sqlmock.Sqlmock, countries []*models.Country) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, country := range countries {
		rows.AddRow(country.ID, country.Name)
	}
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
}
