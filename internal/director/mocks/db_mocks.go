package mocks

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
)

func MockDirectorRepoInsertReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery(`INSERT INTO directors`).
		WithArgs(name).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockDirectorRepoInsertReturnErrNoRows(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO directors`).
		WithArgs(name).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
}

func MockDirectorRepoUpdateReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE directors`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func MockDirectorRepoUpdateReturnResultZero(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE directors`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func MockDirectorRepoDeleteReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM directors`).
		WithArgs(id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockDirectorRepoSelectReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(id, name)
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnRows(rows)
}

func MockDirectorRepoSelectReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnError(sql.ErrNoRows)
}
