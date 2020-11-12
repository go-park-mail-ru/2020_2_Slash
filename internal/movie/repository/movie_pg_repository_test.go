package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var contentInst *models.Content = &models.Content{
	ContentID: 1,
}

var movieInst *models.Movie = &models.Movie{
	ID:      1,
	Video:   "movie.mp4",
	Content: *contentInst,
}

func TestMoviePgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoInsertReturnRows(mock, 1, movieInst)
	err = moviePgRep.Insert(movieInst)
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

	mocks.MockMovieRepoInsertReturnRows(mock, 1, movieInst)
	err = moviePgRep.Insert(movieInst)
	assert.NoError(t, err)

	mocks.MockMovieRepoInsertReturnErrNoUniq(mock, 2, movieInst)
	err = moviePgRep.Insert(movieInst)
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

	var movieInst *models.Movie = &models.Movie{
		ID:      1,
		Video:   "movie.mp4",
		Content: *contentInst,
	}

	mocks.MockMovieRepoUpdateReturnResultOk(mock, movieInst.ID, movieInst)
	err = moviePgRep.Update(movieInst)
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

	mocks.MockMovieRepoUpdateReturnResultZero(mock, movieInst.ID, movieInst)
	err = moviePgRep.Update(movieInst)
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

	mocks.MockMovieRepoDeleteReturnResultOk(mock, movieInst.ID, movieInst)
	err = moviePgRep.DeleteByID(movieInst.ID)
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

	mocks.MockMovieRepoSelectByIDReturnRows(mock, movieInst.ID, movieInst)
	dbMovie, err := moviePgRep.SelectByID(movieInst.ID)
	assert.Equal(t, movieInst, dbMovie)
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

	mocks.MockMovieRepoSelectByIDReturnErrNoRows(mock, movieInst.ID)
	dbMovie, err := moviePgRep.SelectByID(movieInst.ID)
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectFullByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	var userID uint64 = 1
	mocks.MockMovieRepoSelectFullByIDReturnRows(mock, movieInst.ID, userID, movieInst)
	dbMovie, err := moviePgRep.SelectFullByID(movieInst.ID, userID)
	assert.Equal(t, movieInst, dbMovie)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectByContentID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	mocks.MockMovieRepoSelectByContentIDReturnRows(mock, movieInst.ID, movieInst)
	dbMovie, err := moviePgRep.SelectByContentID(movieInst.ContentID)
	assert.Equal(t, movieInst, dbMovie)
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

	mocks.MockMovieRepoSelectByContentIDReturnErrNoRows(mock, movieInst)
	dbMovie, err := moviePgRep.SelectByContentID(movieInst.ContentID)
	assert.Equal(t, dbMovie, (*models.Movie)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectByParams_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	var contentInst *models.Content = &models.Content{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		Countries:        nil,
		Genres:           nil,
		Actors:           nil,
		Directors:        nil,
		Type:             "movie",
	}

	content := []*models.Content{
		contentInst,
	}

	movies := []*models.Movie{
		&models.Movie{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	params := &models.ContentFilter{
		Year:     2001,
		Genre:    1,
		Country:  1,
		Actor:    1,
		Director: 1,
	}

	mocks.MockMovieRepoSelectByParamsReturnRows(mock, params, pgnt, userID, movies)
	dbMovies, err := moviePgRep.SelectByParams(params, pgnt, userID)
	assert.Equal(t, movies, dbMovies)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectLatest_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	movies := []*models.Movie{
		&models.Movie{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	mocks.MockMovieRepoSelectLatestReturnRows(mock, pgnt, userID, movies)
	dbMovies, err := moviePgRep.SelectLatest(pgnt, userID)
	assert.Equal(t, movies, dbMovies)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMoviePgRepository_SelectByRating_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	moviePgRep := NewMoviePgRepository(db)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	movies := []*models.Movie{
		&models.Movie{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	mocks.MockMovieRepoSelectByRatingReturnRows(mock, pgnt, userID, movies)
	dbMovies, err := moviePgRep.SelectByRating(pgnt, userID)
	assert.Equal(t, movies, dbMovies)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
