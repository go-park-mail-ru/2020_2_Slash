package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var content_inst *models.Content = &models.Content{
	ContentID: 1,
}

var movie_inst *models.Movie = &models.Movie{
	ID:      1,
	Video:   "movie_inst.mp4",
	Content: *content_inst,
}

func TestMoviePgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoInsertReturnRows(mock, 1, movie_inst)
	err = moviePgRep.Insert(movie_inst)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_Insert_ContentAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoInsertReturnRows(mock, 1, movie_inst)
	err = moviePgRep.Insert(movie_inst)
	assert.NoError(t, err)

	mocks.MockMovieRepoInsertReturnErrNoUniq(mock, 2, movie_inst)
	err = moviePgRep.Insert(movie_inst)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoUpdateReturnResultOk(mock, movie_inst.ID, movie_inst)
	err = moviePgRep.Update(movie_inst)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_Update_NoMovieWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoUpdateReturnResultZero(mock, movie_inst.ID, movie_inst)
	err = moviePgRep.Update(movie_inst)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_DeleteByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoDeleteReturnResultOk(mock, movie_inst.ID, movie_inst)
	err = moviePgRep.DeleteByID(movie_inst.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoSelectByIDReturnRows(mock, movie_inst.ID, movie_inst)
	dbMovie, err := moviePgRep.SelectByID(movie_inst.ID)
	assert.Equal(t, movie_inst, dbMovie)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectById_NoMovieWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoSelectByIDReturnErrNoRows(mock, movie_inst.ID)
	dbMovie, err := moviePgRep.SelectByID(movie_inst.ID)
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectByName_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoSelectByContentIDReturnRows(mock, movie_inst.ID, movie_inst)
	dbMovie, err := moviePgRep.SelectByContentID(movie_inst.ContentID)
	assert.Equal(t, movie_inst, dbMovie)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectById_NoMovieWithThisContentID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoSelectByContentIDReturnErrNoRows(mock, movie_inst)
	dbMovie, err := moviePgRep.SelectByContentID(movie_inst.ContentID)
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
