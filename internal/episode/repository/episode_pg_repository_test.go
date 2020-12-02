package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testEpisode = &models.Episode{
	ID:          2,
	Name:        "Огурчик Рик",
	Number:      3,
	Video:       "/videos/rickandmorty_22/3/3",
	Description: "Рик превращает себя в огурчик.",
	Poster:      "/images/rickandmorty_22/3/3",
	SeasonID:    3,
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
	testEpisode,
}

func BuildMockAndRepo() (sqlmock.Sqlmock, episode.EpisodeRepository, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	seasonPgRep := NewEpisodeRepository(db)
	return mock, seasonPgRep, nil
}

func TestEpisodePgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectInsertSuccess(mock, testEpisode)

	err = seasonPgRep.Insert(testEpisode)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEpisodePgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectUpdateSuccess(mock, testEpisode)

	err = seasonPgRep.Update(testEpisode)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEpisodePgRepository_Delete_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectDeleteSuccess(mock, testEpisode)

	err = seasonPgRep.DeleteByID(testEpisode.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEpisodePgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectByIDReturnRows(mock, testEpisode)

	episode, err := seasonPgRep.SelectByID(testEpisode.ID)
	assert.Equal(t, testEpisode, episode)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEpisodePgRepository_SelectByID_NoRows(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectByIDReturnErrNoRows(mock, testEpisode)

	episode, err := seasonPgRep.SelectByID(testEpisode.ID)
	assert.Nil(t, episode)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestEpisodePgRepository_SelectByNumberAndSeason_OK(t *testing.T) {
	t.Parallel()
	mock, seasonPgRep, err := BuildMockAndRepo()
	if err != nil {
		t.Fatal(err)
	}

	mocks.ExpectSelectByNumberAndSeason(mock, testEpisode)

	episode, err := seasonPgRep.SelectByNumberAndSeason(testEpisode.Number, testEpisode.SeasonID)
	assert.Equal(t, testEpisode, episode)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
