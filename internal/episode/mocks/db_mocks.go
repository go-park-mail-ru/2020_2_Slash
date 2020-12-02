package mocks

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func ExpectInsertSuccess(mock sqlmock.Sqlmock, episode *models.Episode) {
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(episode.ID)
	mock.
		ExpectQuery(`INSERT INTO episodes`).
		WithArgs(episode.Number, episode.Name, episode.Video,
			episode.Description, episode.Poster, episode.SeasonID).
		WillReturnRows(rows)
	mock.ExpectCommit()
}

func ExpectUpdateSuccess(mock sqlmock.Sqlmock, newEpisode *models.Episode) {
	mock.ExpectBegin()
	mock.
		ExpectExec(`UPDATE episodes`).
		WithArgs(newEpisode.Number, newEpisode.Name, newEpisode.Video,
			newEpisode.Description, newEpisode.Poster, newEpisode.SeasonID, newEpisode.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func ExpectSelectByIDReturnRows(mock sqlmock.Sqlmock, episode *models.Episode) {
	rows := sqlmock.NewRows([]string{"id", "number", "name", "video",
		"description", "poster", "season_id"})
	rows.AddRow(episode.ID, episode.Number, episode.Name, episode.Video,
		episode.Description, episode.Poster, episode.SeasonID)
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(episode.ID).
		WillReturnRows(rows)
}

func ExpectSelectByIDReturnErrNoRows(mock sqlmock.Sqlmock, episode *models.Episode) {
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(episode.ID).
		WillReturnError(sql.ErrNoRows)
}

func ExpectSelectByNumberAndSeason(mock sqlmock.Sqlmock, episode *models.Episode) {
	rows := sqlmock.NewRows([]string{"id", "number", "name", "video",
		"description", "poster", "season_id"})
	rows.AddRow(episode.ID, episode.Number, episode.Name, episode.Video,
		episode.Description, episode.Poster, episode.SeasonID)
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(episode.Number, episode.SeasonID).
		WillReturnRows(rows)
}

func ExpectSelectByNumberAndSeasonReturnErrNoRows(mock sqlmock.Sqlmock, episode *models.Episode) {
	mock.
		ExpectQuery(`SELECT`).
		WithArgs(episode.Number, episode.SeasonID).
		WillReturnError(sql.ErrNoRows)
}

func ExpectDeleteSuccess(mock sqlmock.Sqlmock, episode *models.Episode) {
	mock.ExpectBegin()
	mock.
		ExpectExec(`DELETE FROM episodes`).
		WithArgs(episode.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}
