package mocks

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func ExpectInsertSuccess(mock sqlmock.Sqlmock, season *models.Season) {
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(season.ID)
	mock.
		ExpectQuery(`INSERT INTO seasons`).
		WithArgs(season.Number, season.EpisodesNumber, season.TVShowID).
		WillReturnRows(rows)
	mock.ExpectCommit()
}

func ExpectUpdateSuccess(mock sqlmock.Sqlmock, season *models.Season) {
	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE seasons`).
		WithArgs(season.Number, season.EpisodesNumber, season.TVShowID, season.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func ExpectSelectByIDReturnRows(mock sqlmock.Sqlmock, season *models.Season) {
	rows := sqlmock.NewRows([]string{"id", "number", "episodes", "tv_show_id"})
	rows.AddRow(season.ID, season.Number, season.EpisodesNumber, season.TVShowID)
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(season.ID).
		WillReturnRows(rows)
}

func ExpectSelectByIDReturnErrNoRows(mock sqlmock.Sqlmock, season *models.Season) {
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(season.ID).
		WillReturnError(sql.ErrNoRows)
}

func ExpectSelectReturnRows(mock sqlmock.Sqlmock, season *models.Season) {
	rows := sqlmock.NewRows([]string{"id", "number", "episodes", "tv_show_id"})
	rows.AddRow(season.ID, season.Number, season.EpisodesNumber, season.TVShowID)
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(season.Number, season.TVShowID).
		WillReturnRows(rows)
}

func ExpectSelectReturnErrNoRows(mock sqlmock.Sqlmock, season *models.Season) {
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(season.Number, season.TVShowID).
		WillReturnError(sql.ErrNoRows)
}

func ExpectSelectEpisodesReturnRows(mock sqlmock.Sqlmock, season *models.Season,
	returnEpisodes []*models.Episode) {
	rows := sqlmock.NewRows([]string{"id", "number", "name",
		"video", "description", "poster", "season_id"})
	for _, episode := range returnEpisodes {
		rows.AddRow(episode.ID, episode.Number, episode.Name, episode.Video,
			episode.Description, episode.Poster, episode.SeasonID)
	}
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(season.ID).
		WillReturnRows(rows)
}

func ExpectSelectEpisodesReturnErrNoRows(mock sqlmock.Sqlmock, season *models.Season) {
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(season.ID).
		WillReturnError(sql.ErrNoRows)
}

func ExpectDeleteSuccess(mock sqlmock.Sqlmock, season *models.Season) {
	mock.ExpectBegin()
	mock.
		ExpectExec(`DELETE FROM seasons`).
		WithArgs(season.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}
