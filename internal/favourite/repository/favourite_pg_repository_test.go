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
	result := []*models.Content{
		{
			ContentID: 2,
		},
		{
			ContentID: 5,
		},
	}

	favouritePgRep := NewFavouritePgRepository(db)

	mocks.MockSelectFavouriteContentReturnRows(mock, userID, result)
	dbFavourites, err := favouritePgRep.SelectFavouriteContent(userID)
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

	mocks.MockSelectFavouriteContentReturnErrNoRows(mock, userID)
	dbFavourites, err := favouritePgRep.SelectFavouriteContent(userID)
	assert.Equal(t, ([]*models.Content)(nil), dbFavourites)
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
