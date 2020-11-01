package mocks

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
)

func MockActorRepoInsertReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery(`INSERT INTO actors`).
		WithArgs(name).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockActorRepoInsertReturnErrNoRows(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO actors`).
		WithArgs(name).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()
}

func MockActorRepoUpdateReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE actors`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func MockActorRepoUpdateReturnResultZero(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE actors`).
		WithArgs(id, name).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func MockActorRepoDeleteReturnResultOk(mock sqlmock.Sqlmock, id uint64, name string) {
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM actors`).
		WithArgs(id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockActorRepoSelectReturnRows(mock sqlmock.Sqlmock, id uint64, name string) {
	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(id, name)
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnRows(rows)
}

func MockActorRepoSelectReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnError(sql.ErrNoRows)
}
