package mocks

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockInsertSuccess(mock sqlmock.Sqlmock, favourite *models.Favourite) {
	mock.ExpectBegin()
	res := sqlmock.NewResult(0, 1)
	mock.ExpectExec(`INSERT INTO favourites`).
		WithArgs(favourite.UserID, favourite.ContentID, favourite.Created).
		WillReturnResult(res)
	mock.ExpectCommit()
}

func MockDeleteSuccess(mock sqlmock.Sqlmock, favourite *models.Favourite) {
	mock.ExpectBegin()
	res := sqlmock.NewResult(0, 1)
	mock.ExpectExec(`DELETE FROM favourites`).
		WithArgs(favourite.UserID, favourite.ContentID).
		WillReturnResult(res)
	mock.ExpectCommit()
}

func MockSelectFavouriteContentReturnRows(mock sqlmock.Sqlmock, userID uint64,
	resContent []*models.Content) {
	rows := sqlmock.NewRows([]string{"id", "name", "original_name", "description",
		"short_description", "year", "images", "type"})
	for _, favourite := range resContent {
		rows.AddRow(favourite.ContentID, "", "", "", "", 0, "", "")
	}
	mock.ExpectQuery(`SELECT`).
		WithArgs(userID).
		WillReturnRows(rows)
}

func MockSelectFavouriteContentReturnErrNoRows(mock sqlmock.Sqlmock, userID uint64) {
	mock.ExpectQuery(`SELECT`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)
}

func MockSelectReturnRows(mock sqlmock.Sqlmock, reqFavorite *models.Favourite,
	resFavorite *models.Favourite) {
	rows := sqlmock.NewRows([]string{"user_id", "content_id", "created"})
	rows.AddRow(resFavorite.UserID, resFavorite.ContentID, resFavorite.Created)
	mock.ExpectQuery(`SELECT`).
		WithArgs(reqFavorite.UserID, reqFavorite.ContentID).
		WillReturnRows(rows)
}

func MockSelectReturnErrNoRows(mock sqlmock.Sqlmock, reqFavourite *models.Favourite) {
	mock.ExpectQuery(`SELECT`).
		WithArgs(reqFavourite.UserID, reqFavourite.ContentID).
		WillReturnError(sql.ErrNoRows)
}
