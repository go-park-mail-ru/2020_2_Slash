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
		"description", "short_description", "year", "images", "type"})
	rows.AddRow(content.ContentID, content.Name, content.OriginalName, content.Description,
		content.ShortDescription, content.Year, content.Images, content.Type)
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
