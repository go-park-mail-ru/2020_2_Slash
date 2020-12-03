package mocks

import (
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

func MockTVShowRepoInsertReturnRows(mock sqlmock.Sqlmock, id uint64, tvshow *models.TVShow) {
	mock.ExpectBegin()
	insertAnswer := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery(`INSERT INTO tv_shows`).
		WithArgs(tvshow.ContentID).
		WillReturnRows(insertAnswer)
	mock.ExpectCommit()
}

func MockTVShowRepoInsertReturnErrNoUniq(mock sqlmock.Sqlmock, id uint64, tvshow *models.TVShow) {
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO tv_shows`).
		WithArgs(tvshow.ContentID).
		WillReturnError(errors.New("No UNIQUE"))
	mock.ExpectRollback()
}

func MockTVShowRepoSelectByIDReturnRows(mock sqlmock.Sqlmock, id uint64, tvshow *models.TVShow) {
	rows := sqlmock.NewRows([]string{"id", "seasons", "content_id"})
	rows.AddRow(id, tvshow.Seasons, tvshow.ContentID)
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnRows(rows)
}

func MockTVShowRepoSelectByIDReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).WithArgs(id).WillReturnError(sql.ErrNoRows)
}

func MockTVShowRepoSelectFullByIDReturnRows(mock sqlmock.Sqlmock, id uint64, curUserID uint64,
	tvshow *models.TVShow) {

	rows := sqlmock.NewRows([]string{"tv.id", "tv.seasons", "c.id", "c.name",
		"c.original_name", "c.description", "c.short_description",
		"c.year", "c.images", "c.type", "r.likes", "is_favourite"})
	rows.AddRow(tvshow.ID, tvshow.Seasons, tvshow.ContentID, tvshow.Name,
		tvshow.OriginalName, tvshow.Description, tvshow.ShortDescription,
		tvshow.Year, tvshow.Images, tvshow.Type, tvshow.IsLiked, tvshow.IsFavourite)
	mock.ExpectQuery(`SELECT tv.id, tv.seasons, c.id, c.name`).WithArgs(id, curUserID).WillReturnRows(rows)
}

func MockTVShowRepoSelectByContentIDReturnRows(mock sqlmock.Sqlmock, id uint64, tvshow *models.TVShow) {
	rows := sqlmock.NewRows([]string{"id", "seasons", "content_id"})
	rows.AddRow(id, tvshow.Seasons, tvshow.ContentID)
	mock.ExpectQuery(`SELECT`).WithArgs(tvshow.ContentID).WillReturnRows(rows)
}

func MockTVShowRepoSelectByContentIDReturnErrNoRows(mock sqlmock.Sqlmock, tvshow *models.TVShow) {
	mock.ExpectQuery(`SELECT`).WithArgs(tvshow.ContentID).WillReturnError(sql.ErrNoRows)
}

func MockTVShowRepoSelectByGenreReturnRows(mock sqlmock.Sqlmock, genreID uint64, tv_shows []*models.TVShow) {
	rows := sqlmock.NewRows([]string{"id", "seasons", "id", "name", "original_name",
		"description", "short_description", "year", "images", "type"})
	for _, tvshow := range tv_shows {
		rows.AddRow(tvshow.ID, tvshow.Seasons, tvshow.ContentID, tvshow.Name, tvshow.OriginalName, tvshow.Description,
			tvshow.ShortDescription, tvshow.Year, tvshow.Images, tvshow.Type)
	}
	mock.ExpectQuery(`SELECT tv.id, tv.seasons, c.id`).WithArgs(genreID).WillReturnRows(rows)
}

func MockTVShowRepoSelectWhereNameLikeReturnRows(mock sqlmock.Sqlmock, pgnt *models.Pagination, curUserID uint64,
	tvshows []*models.TVShow, name string) {

	rows := sqlmock.NewRows([]string{"tv.id", "tv.seasons", "c.id", "c.name",
		"c.original_name", "c.description", "c.short_description", "c.rating",
		"c.year", "c.images", "c.type", "r.likes", "is_favourite"})
	for _, tvshow := range tvshows {
		rows.AddRow(tvshow.ID, tvshow.Seasons, tvshow.ContentID, tvshow.Name,
			tvshow.OriginalName, tvshow.Description, tvshow.ShortDescription, tvshow.Rating,
			tvshow.Year, tvshow.Images, tvshow.Type, tvshow.IsLiked, tvshow.IsFavourite)
	}
	query := `
		SELECT tv.id, tv.seasons, c.id, c.name`

	searchName := "%" + name + "%"
	mock.ExpectQuery(query).WithArgs(curUserID, searchName, pgnt.Count, pgnt.From).WillReturnRows(rows)
}

func MockTVShowRepoSelectByParamsReturnRows(mock sqlmock.Sqlmock, params *models.ContentFilter,
	pgnt *models.Pagination, curUserID uint64, tv_shows []*models.TVShow) {

	rows := sqlmock.NewRows([]string{"tv.id", "tv.seasons", "c.id", "c.name",
		"c.original_name", "c.description", "c.short_description", "c.rating",
		"c.year", "c.images", "c.type", "r.likes", "is_favourite"})
	for _, tvshow := range tv_shows {
		rows.AddRow(tvshow.ID, tvshow.Seasons, tvshow.ContentID, tvshow.Name,
			tvshow.OriginalName, tvshow.Description, tvshow.ShortDescription, tvshow.Rating,
			tvshow.Year, tvshow.Images, tvshow.Type, tvshow.IsLiked, tvshow.IsFavourite)
	}
	query := `
		SELECT tv.id, tv.seasons, c.id, c.name`

	mock.ExpectQuery(query).WithArgs(curUserID, pgnt.Count, pgnt.From, params.Genre,
		params.Country, params.Actor, params.Director, params.Year).WillReturnRows(rows)
}

func MockTVShowRepoSelectLatestReturnRows(mock sqlmock.Sqlmock, pgnt *models.Pagination, curUserID uint64,
	tv_shows []*models.TVShow) {

	rows := sqlmock.NewRows([]string{"tv.id", "tv.seasons", "c.id", "c.name",
		"c.original_name", "c.description", "c.short_description", "c.rating",
		"c.year", "c.images", "c.type", "r.likes", "is_favourite"})
	for _, tvshow := range tv_shows {
		rows.AddRow(tvshow.ID, tvshow.Seasons, tvshow.ContentID, tvshow.Name,
			tvshow.OriginalName, tvshow.Description, tvshow.ShortDescription, tvshow.Rating,
			tvshow.Year, tvshow.Images, tvshow.Type, tvshow.IsLiked, tvshow.IsFavourite)
	}
	query := `
		SELECT tv.id, tv.seasons, c.id, c.name`

	mock.ExpectQuery(query).WithArgs(curUserID, pgnt.Count, pgnt.From).WillReturnRows(rows)
}

func MockTVShowRepoSelectByRatingReturnRows(mock sqlmock.Sqlmock, pgnt *models.Pagination, curUserID uint64,
	tv_shows []*models.TVShow) {

	rows := sqlmock.NewRows([]string{"tv.id", "tv.seasons", "c.id", "c.name",
		"c.original_name", "c.description", "c.short_description", "c.rating",
		"c.year", "c.images", "c.type", "r.likes", "is_favourite"})
	for _, tvshow := range tv_shows {
		rows.AddRow(tvshow.ID, tvshow.Seasons, tvshow.ContentID, tvshow.Name,
			tvshow.OriginalName, tvshow.Description, tvshow.ShortDescription, tvshow.Rating,
			tvshow.Year, tvshow.Images, tvshow.Type, tvshow.IsLiked, tvshow.IsFavourite)
	}
	query := `
		SELECT tv.id, tv.seasons, c.id, c.name`

	mock.ExpectQuery(query).WithArgs(curUserID, pgnt.Count, pgnt.From).WillReturnRows(rows)
}
