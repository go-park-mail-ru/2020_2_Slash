package mocks

import (
	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockContentRepoDeleteReturnResultOk(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM content`).
		WithArgs(id).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
}

func MockContentRepoSelectByIDReturnRows(mock sqlmock.Sqlmock, id uint64, content *models.Content) {
	rows := sqlmock.NewRows([]string{"id", "name", "original_name",
		"description", "short_description", "year", "images", "type", "is_free"})
	rows.AddRow(content.ContentID, content.Name, content.OriginalName, content.Description,
		content.ShortDescription, content.Year, content.Images, content.Type, content.IsFree)
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnRows(rows)
}

func MockContentRepoSelectByIDReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnError(sql.ErrNoRows)
}

func MockContentRepoSelectCountriesReturnRows(mock sqlmock.Sqlmock, id uint64, countries []uint64) {
	rows := sqlmock.NewRows([]string{"id"})
	for _, country := range countries {
		rows.AddRow(country)
	}
	mock.ExpectQuery(`SELECT country_id`).WithArgs(id).WillReturnRows(rows)
}

func MockContentRepoSelectDirectorsReturnRows(mock sqlmock.Sqlmock, id uint64, directors []uint64) {
	rows := sqlmock.NewRows([]string{"id"})
	for _, director := range directors {
		rows.AddRow(director)
	}
	mock.ExpectQuery(`SELECT director_id`).WithArgs(id).WillReturnRows(rows)
}

func MockContentRepoSelectActorsReturnRows(mock sqlmock.Sqlmock, id uint64, actors []uint64) {
	rows := sqlmock.NewRows([]string{"id"})
	for _, actor := range actors {
		rows.AddRow(actor)
	}
	mock.ExpectQuery(`SELECT actor_id`).WithArgs(id).WillReturnRows(rows)
}

func MockContentRepoSelectGenresReturnRows(mock sqlmock.Sqlmock, id uint64, genres []uint64) {
	rows := sqlmock.NewRows([]string{"id"})
	for _, genre := range genres {
		rows.AddRow(genre)
	}
	mock.ExpectQuery(`SELECT genre_id`).WithArgs(id).WillReturnRows(rows)
}

func MockContentRepoInsertReturnResultOk(mock sqlmock.Sqlmock, content *models.Content, country *models.Country,
	genre *models.Genre, actor *models.Actor, director *models.Director) {
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(content.ContentID)
	mock.ExpectQuery(`INSERT INTO content`).
		WithArgs(content.Name, content.OriginalName, content.Description,
			content.ShortDescription, content.Year, content.Images, content.Type, content.IsFree).WillReturnRows(rows)

	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, country.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, genre.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, actor.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, director.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectCommit()
}

func MockContentRepoUpdateReturnResultOk(mock sqlmock.Sqlmock, content *models.Content, country *models.Country,
	genre *models.Genre, actor *models.Actor, director *models.Director) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE`).
		WithArgs(content.ContentID, content.Name, content.OriginalName, content.Description,
			content.ShortDescription, content.Year, content.Images, content.Type, content.IsFree).WillReturnResult(driver.ResultNoRows)

	mock.ExpectExec(`DELETE`).
		WithArgs(content.ContentID).WillReturnResult(driver.ResultNoRows)
	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, country.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectExec(`DELETE`).
		WithArgs(content.ContentID).WillReturnResult(driver.ResultNoRows)
	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, genre.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectExec(`DELETE`).
		WithArgs(content.ContentID).WillReturnResult(driver.ResultNoRows)
	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, actor.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectExec(`DELETE`).
		WithArgs(content.ContentID).WillReturnResult(driver.ResultNoRows)
	mock.ExpectPrepare(``).ExpectExec().WithArgs(content.ContentID, director.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(``).WithArgs().WillReturnResult(driver.ResultNoRows)

	mock.ExpectCommit()
}

func MockContentRepoUpdateImagesReturnResultOk(mock sqlmock.Sqlmock, content *models.Content) {
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE`).
		WithArgs(content.ContentID, content.Images).WillReturnResult(driver.ResultNoRows)

	mock.ExpectCommit()
}
