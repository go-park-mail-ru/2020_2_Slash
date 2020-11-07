package mocks

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockMovieRepoInsertReturnRows(mock sqlmock.Sqlmock, id uint64, movie *models.Movie) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery(`INSERT INTO movies`).
		WithArgs(movie.Video, movie.ContentID).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockMovieRepoInsertReturnErrNoUniq(mock sqlmock.Sqlmock, id uint64, movie *models.Movie) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO movies`).
		WithArgs(movie.Video, movie.ContentID).
		WillReturnError(errors.New("No UNIQUE"))
	mock.ExpectRollback()
}

func MockMovieRepoUpdateReturnResultOk(mock sqlmock.Sqlmock, id uint64, movie *models.Movie) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE movies`).
		WithArgs(id, movie.Video, movie.ContentID).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	mock.ExpectCommit()
}

func MockMovieRepoUpdateReturnResultZero(mock sqlmock.Sqlmock, id uint64, movie *models.Movie) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE movies`).
		WithArgs(id, movie.Video, movie.ContentID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func MockMovieRepoDeleteReturnResultOk(mock sqlmock.Sqlmock, id uint64, movie *models.Movie) {
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM movies`).
		WithArgs(id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockMovieRepoSelectByIDReturnRows(mock sqlmock.Sqlmock, id uint64, movie *models.Movie) {
	rows := sqlmock.NewRows([]string{"id", "video", "content_id"})
	rows.AddRow(id, movie.Video, movie.ContentID)
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnRows(rows)
}

func MockMovieRepoSelectByIDReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnError(sql.ErrNoRows)
}

func MockMovieRepoSelectByContentIDReturnRows(mock sqlmock.Sqlmock, id uint64, movie *models.Movie) {
	rows := sqlmock.NewRows([]string{"id", "video", "content_id"})
	rows.AddRow(id, movie.Video, movie.ContentID)
	mock.ExpectQuery(`SELECT`).WithArgs(movie.ContentID).WillReturnRows(rows)
}

func MockMovieRepoSelectByContentIDReturnErrNoRows(mock sqlmock.Sqlmock, movie *models.Movie) {
	mock.ExpectQuery(`SELECT`).WithArgs(movie.ContentID).WillReturnError(sql.ErrNoRows)
}

func MockMovieRepoSelectByGenreReturnRows(mock sqlmock.Sqlmock, genreID uint64, movies []*models.Movie) {
	rows := sqlmock.NewRows([]string{"id", "video", "id", "name", "original_name",
		"description", "short_description", "year", "images", "type"})
	for _, movie := range movies {
		rows.AddRow(movie.ID, movie.Video, movie.ContentID, movie.Name, movie.OriginalName, movie.Description,
			movie.ShortDescription, movie.Year, movie.Images, movie.Type)
	}
	mock.ExpectQuery(`SELECT m.id, m.video, c.id`).WithArgs(genreID).WillReturnRows(rows)
}
