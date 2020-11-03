package mocks

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockGenreRepoInsertReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery(`INSERT INTO genres`).
		WithArgs(name).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockGenreRepoInsertReturnErrNoUniq(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO genres`).
		WithArgs(name).
		WillReturnError(errors.New("No UNIQUE"))
	mock.ExpectRollback()
}

func MockGenreRepoUpdateReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE genres`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func MockGenreRepoUpdateReturnResultZero(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE genres`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func MockGenreRepoDeleteReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM genres`).
		WithArgs(id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockGenreRepoSelectByIDReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(id, name)
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnRows(rows)
}

func MockGenreRepoSelectByIDReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnError(sql.ErrNoRows)
}

func MockGenreRepoSelectByNameReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(id, name)
	mock.ExpectQuery(`SELECT`).WithArgs(name).WillReturnRows(rows)
}

func MockGenreRepoSelectByNameReturnErrNoRows(mock sqlmock.Sqlmock, name string) {
	mock.ExpectQuery(`SELECT`).WithArgs(name).WillReturnError(sql.ErrNoRows)
}

func MockGenreRepoSelectAllReturnRows(mock sqlmock.Sqlmock, genres []*models.Genre) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, genre := range genres {
		rows.AddRow(genre.ID, genre.Name)
	}
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
}
