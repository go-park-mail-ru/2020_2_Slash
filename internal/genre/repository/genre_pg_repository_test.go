package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenrePgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoInsertReturnRows(mock, 1, genre.Name)
	err = genrePgRep.Insert(genre)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_Insert_GenreAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoInsertReturnRows(mock, 1, genre.Name)
	err = genrePgRep.Insert(genre)
	assert.NoError(t, err)

	mocks.MockGenreRepoInsertReturnErrNoUniq(mock, 2, genre.Name)
	err = genrePgRep.Insert(genre)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoUpdateReturnResultOk(mock, genre.ID, genre.Name)
	err = genrePgRep.Update(genre)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_Update_NoGenreWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoUpdateReturnResultZero(mock, genre.ID, genre.Name)
	err = genrePgRep.Update(genre)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_DeleteByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoDeleteReturnResultOk(mock, genre.ID, genre.Name)
	err = genrePgRep.DeleteByID(genre.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoSelectByIDReturnRows(mock, genre.ID, genre.Name)
	dbGenre, err := genrePgRep.SelectByID(genre.ID)
	assert.Equal(t, genre, dbGenre)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_SelectById_NoGenreWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoSelectByIDReturnErrNoRows(mock, genre.ID)
	dbGenre, err := genrePgRep.SelectByID(genre.ID)
	assert.Equal(t, dbGenre, (*models.Genre)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_SelectByName_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoSelectByNameReturnRows(mock, genre.ID, genre.Name)
	dbGenre, err := genrePgRep.SelectByName(genre.Name)
	assert.Equal(t, genre, dbGenre)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_SelectById_NoGenreWithThisName(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genre := &models.Genre{
		ID:   1,
		Name: "USA",
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoSelectByNameReturnErrNoRows(mock, genre.Name)
	dbGenre, err := genrePgRep.SelectByName(genre.Name)
	assert.Equal(t, dbGenre, (*models.Genre)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenrePgRepository_SelectAll_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	genres := []*models.Genre{
		&models.Genre{
			ID:   1,
			Name: "USA",
		},
		&models.Genre{
			ID:   2,
			Name: "GB",
		},
	}

	genrePgRep := NewGenrePgRepository(db)

	mocks.MockGenreRepoSelectAllReturnRows(mock, genres)
	dbGenre, err := genrePgRep.SelectAll()
	assert.Equal(t, genres, dbGenre)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
