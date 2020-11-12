package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFavouritePgRepository_Insert_Success(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	fav := &models.Favourite{
		UserID:    3,
		ContentID: 2,
		Created:   time.Now(),
	}

	favouritePgRep := NewFavouritePgRepository(db)

	mocks.MockInsertSuccess(mock, fav)
	err = favouritePgRep.Insert(fav)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActorPgRepository_DeleteById_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	fav := &models.Favourite{
		UserID:    3,
		ContentID: 2,
		Created:   time.Now(),
	}

	favouritePgRep := NewFavouritePgRepository(db)

	mocks.MockDeleteSuccess(mock, fav)
	err = favouritePgRep.Delete(fav)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavouritePgRepository_SelectFavouritesById_Success(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var userID uint64 = 3
	content := models.Content{
		ContentID:        1,
		Name:             "content_1",
		OriginalName:     "content_1",
		Description:      "desc",
		ShortDescription: "short desc",
		Year:             2020,
		Images:           "images/content_1",
		Type:             "movie",
	}
	content2 := models.Content{
		ContentID:        2,
		Name:             "content_2",
		OriginalName:     "content_2",
		Description:      "desc",
		ShortDescription: "short desc",
		Year:             2020,
		Images:           "images/content_2",
		Type:             "movie",
	}

	movie := &models.Movie{
		ID:      1,
		Video:   "videos/movie_1.mp4",
		Content: content,
	}
	movie2 := &models.Movie{
		ID:      2,
		Video:   "videos/movie_2.mp4",
		Content: content2,
	}
	result := []*models.Movie{
		movie,
		movie2,
	}
	var limit uint64 = 2
	var	offset uint64 = 0

	favouritePgRep := NewFavouritePgRepository(db)

	mocks.MockSelectFavouriteMoviesReturnRows(mock, userID, result, limit, offset)
	dbFavourites, err := favouritePgRep.SelectFavouriteMovies(userID, limit, offset)
	assert.Equal(t, result, dbFavourites)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavouritePgRepository_SelectFavouritesById_NoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var userID uint64 = 3

	favouritePgRep := NewFavouritePgRepository(db)

	var limit, offset uint64 = 2, 0
	mocks.MockSelectFavouriteContentReturnErrNoRows(mock, userID, limit, offset)
	dbFavourites, err := favouritePgRep.SelectFavouriteMovies(userID, limit, offset)
	assert.Equal(t, ([]*models.Movie)(nil), dbFavourites)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavouritePgRepository_Select_NoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	fav := &models.Favourite{
		UserID:    3,
		ContentID: 2,
		Created:   time.Now(),
	}

	favouritePgRep := NewFavouritePgRepository(db)
	mocks.MockSelectReturnErrNoRows(mock, fav)
	err = favouritePgRep.Select(fav)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFavouritePgRepository_Select_Success(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	fav := &models.Favourite{
		UserID:    3,
		ContentID: 2,
		Created:   time.Now(),
	}

	favouritePgRep := NewFavouritePgRepository(db)
	mocks.MockSelectReturnRows(mock, fav, fav)
	err = favouritePgRep.Select(fav)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
