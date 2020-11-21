package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testSeason = &models.Season{
	ID:             0,
	Number:         3,
	EpisodesNumber: 3,
	TVShowID:       1,
	Episodes:       nil,
}

var testEpisodes = []*models.Episode{
	&models.Episode{
		ID:          1,
		Name:        "Рикбег из Рикшенка",
		Number:      1,
		Video:       "/videos/rickandmorty_22/3/1",
		Description: "Саммер решает спасти Рика из тюрьмы.",
		Poster:      "/images/rickandmorty_22/3/1",
		SeasonID:    3,
	},
	&models.Episode{
		ID:          2,
		Name:        "Рикман с камнем",
		Number:      2,
		Video:       "/videos/rickandmorty_22/3/2",
		Description: "Рик, Морти и Саммер охотятся за новым источником энергии в постакалиптической версии Земли.",
		Poster:      "/images/rickandmorty_22/3/2",
		SeasonID:    3,
	},
	&models.Episode{
		ID:          2,
		Name:        "Огурчик Рик",
		Number:      3,
		Video:       "/videos/rickandmorty_22/3/3",
		Description: "Рик превращает себя в огурчик.",
		Poster:      "/images/rickandmorty_22/3/3",
		SeasonID:    3,
	},
}


func BuildMockAndRepo() (sqlmock.Sqlmock, season.SeasonRepository, error){
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	seasonPgRep := NewSeasonPgRepository(db)
	return mock, seasonPgRep, nil
}

func TestSeasonPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectInsertSuccess(mock, testSeason)

	err = seasonPgRep.Insert(testSeason)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectUpdateSuccess(mock, testSeason)

	err = seasonPgRep.Update(testSeason)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_Delete_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectDeleteSuccess(mock, testSeason)

	err = seasonPgRep.Delete(testSeason.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectByIDReturnRows(mock, testSeason)

	season, err := seasonPgRep.SelectByID(testSeason.ID)
	assert.Equal(t, testSeason, season)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_SelectByID_NoRows(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectByIDReturnErrNoRows(mock, testSeason)

	season, err := seasonPgRep.SelectByID(testSeason.ID)
	assert.Nil(t, season)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_Select_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectReturnRows(mock, testSeason)

	season, err := seasonPgRep.Select(testSeason)
	assert.Equal(t, testSeason, season)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_Select_NoRows(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectReturnErrNoRows(mock, testSeason)

	season, err := seasonPgRep.Select(testSeason)
	assert.Nil(t, season)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_SelectEpisodes_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectEpisodesReturnRows(mock, testSeason, testEpisodes)

	episodes, err := seasonPgRep.SelectEpisodes(testSeason.ID)
	assert.Equal(t, testEpisodes, episodes)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSeasonPgRepository_SelectEpisodes_NoRows(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectEpisodesReturnErrNoRows(mock, testSeason)

	episodes, err := seasonPgRep.SelectEpisodes(testSeason.ID)
	assert.Nil(t, episodes)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
