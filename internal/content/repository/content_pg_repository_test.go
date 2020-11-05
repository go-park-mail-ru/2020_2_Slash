package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var countries = []*models.Country{
	&models.Country{
		ID:   1,
		Name: "США",
	},
}

var genres = []*models.Genre{
	&models.Genre{
		Name: "Мультфильм",
	},
	&models.Genre{
		Name: "Комедия",
	},
}

var actors = []*models.Actor{
	&models.Actor{
		Name: "Майк Майерс",
	},
	&models.Actor{
		Name: "Эдди Мёрфи",
	},
}

var directors = []*models.Director{
	&models.Director{
		Name: "Эндрю Адамсон",
	},
	&models.Director{
		Name: "Вики Дженсон",
	},
}

var content_inst *models.Content = &models.Content{
	Name:             "Шрек",
	OriginalName:     "Shrek",
	Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
	ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
	Year:             2001,
	Countries:        countries,
	Genres:           genres,
	Actors:           actors,
	Directors:        directors,
	Type:             "movie",
}

func TestContentPgRepository_DeleteByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	contentPgRep := NewContentPgRepository(db)

	mocks.MockContentRepoDeleteReturnResultOk(mock, content_inst.ContentID)
	err = contentPgRep.DeleteByID(content_inst.ContentID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContentPgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var content_inst *models.Content = &models.Content{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		Type:             "movie",
	}

	contentPgRep := NewContentPgRepository(db)

	mocks.MockContentRepoSelectByIDReturnRows(mock, content_inst.ContentID, content_inst)
	dbContent, err := contentPgRep.SelectByID(content_inst.ContentID)
	assert.Equal(t, content_inst, dbContent)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContentPgRepository_SelectById_NoContentWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	contentPgRep := NewContentPgRepository(db)

	mocks.MockContentRepoSelectByIDReturnErrNoRows(mock, content_inst.ContentID)
	dbContent, err := contentPgRep.SelectByID(content_inst.ContentID)
	assert.Equal(t, dbContent, (*models.Content)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContentPgRepository_SelectCountries_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	countriesID := []uint64{1}

	contentPgRep := NewContentPgRepository(db)

	mocks.MockContentRepoSelectCountriesReturnRows(mock, content_inst.ContentID, countriesID)
	dbCountires, err := contentPgRep.SelectCountriesByID(content_inst.ContentID)
	assert.Equal(t, countriesID, dbCountires)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContentPgRepository_SelectDirectors_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	directorsID := []uint64{1, 2}

	contentPgRep := NewContentPgRepository(db)

	mocks.MockContentRepoSelectDirectorsReturnRows(mock, content_inst.ContentID, directorsID)
	dbDirectors, err := contentPgRep.SelectDirectorsByID(content_inst.ContentID)
	assert.Equal(t, directorsID, dbDirectors)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContentPgRepository_SelectActors_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	actorsID := []uint64{1, 2}

	contentPgRep := NewContentPgRepository(db)

	mocks.MockContentRepoSelectActorsReturnRows(mock, content_inst.ContentID, actorsID)
	dbActors, err := contentPgRep.SelectActorsByID(content_inst.ContentID)
	assert.Equal(t, actorsID, dbActors)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestContentPgRepository_SelectGenres_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genresID := []uint64{1, 2}

	contentPgRep := NewContentPgRepository(db)

	mocks.MockContentRepoSelectGenresReturnRows(mock, content_inst.ContentID, genresID)
	dbGenres, err := contentPgRep.SelectGenresByID(content_inst.ContentID)
	assert.Equal(t, genresID, dbGenres)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
