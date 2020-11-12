package mocks

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockInsertReturnRows(mock sqlmock.Sqlmock, rating *models.Rating) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"user_id", "content_id", "likes"})
	insertAnswer.AddRow(rating.UserID, rating.ContentID, rating.Likes)
	mock.
		ExpectExec(`INSERT INTO rates`).
		WithArgs(rating.UserID, rating.ContentID, rating.Likes).
		WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectCommit()
}

func MockUpdateReturnResultOK(mock sqlmock.Sqlmock, rating *models.Rating) {
	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE rates`).
		WithArgs(rating.Likes, rating.ContentID, rating.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func MockDeleteReturnResultOK(mock sqlmock.Sqlmock, rating *models.Rating) {
	mock.ExpectBegin()
	mock.
		ExpectExec(`DELETE FROM rates`).
		WithArgs(rating.UserID, rating.ContentID).
		WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockSelectReturnRows(mock sqlmock.Sqlmock, rating *models.Rating) {
	rows := sqlmock.NewRows([]string{"user_id", "content_id", "likes"})
	rows.AddRow(rating.UserID, rating.ContentID, rating.Likes)
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(rating.UserID, rating.ContentID).
		WillReturnRows(rows)
}

func MockSelectReturnErrNoRows(mock sqlmock.Sqlmock, rating *models.Rating) {
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(rating.UserID, rating.ContentID).
		WillReturnError(sql.ErrNoRows)
}

func MockSelectCountReturnRows(mock sqlmock.Sqlmock, contentID uint64, count int) {
	rows := sqlmock.NewRows([]string{"count"})
	rows.AddRow(count)
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(contentID).
		WillReturnRows(rows)
}
