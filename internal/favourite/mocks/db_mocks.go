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

func MockSelectFavouriteMoviesReturnRows(mock sqlmock.Sqlmock, userID uint64,
	resMovies []*models.Movie, limit uint64, offset uint64) {
	rows := sqlmock.NewRows([]string{"id", "video", "content_id", "name", "original_name",
		"description", "short_description", "year", "images", "type",
		"likes", "is_favourite"})
	for _, movie := range resMovies {
		rows.AddRow(movie.ID, movie.Video, movie.ContentID, movie.Name,
			movie.OriginalName, movie.Description, movie.ShortDescription,
			movie.Year, movie.Images, movie.Type, movie.IsLiked, movie.IsFavourite)
	}
	mock.ExpectQuery(`SELECT`).
		WithArgs(userID, limit, offset).
		WillReturnRows(rows)
}

func MockSelectFavouriteContentReturnErrNoRows(mock sqlmock.Sqlmock, userID uint64,
	limit uint64, offset uint64) {
	mock.ExpectQuery(`SELECT`).
		WithArgs(userID, limit, offset).
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
