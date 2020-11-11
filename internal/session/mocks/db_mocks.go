package mocks

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockInsertReturnRows(mock sqlmock.Sqlmock, session *models.Session) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(session.ID)
	mock.
		ExpectQuery(`INSERT INTO sessions`).
		WithArgs(session.Value, session.ExpiresAt, session.UserID).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockDeleteReturnResultOk(mock sqlmock.Sqlmock, sessionValue string) {
	mock.ExpectBegin()
	mock.
		ExpectExec(`DELETE FROM sessions`).
		WithArgs(sessionValue).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockSelectReturnRows(mock sqlmock.Sqlmock, session *models.Session) {
	rows := sqlmock.NewRows([]string{"id", "value", "expires", "user_id"})
	rows.AddRow(session.ID, session.Value, session.ExpiresAt, session.UserID)
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(session.Value).
		WillReturnRows(rows)
}

func MockSelectReturnErrNoRows(mock sqlmock.Sqlmock, sessionValue string) {
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(sessionValue).
		WillReturnError(sql.ErrNoRows)
}
