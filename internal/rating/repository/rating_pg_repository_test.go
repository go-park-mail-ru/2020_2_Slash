package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/rating/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRatingPgRepository_Insert_Success(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingPgRep := NewRatingPgRepository(db)

	mocks.MockInsertReturnRows(mock, rating)
	err = ratingPgRep.Insert(rating)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRatingPgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingPgRep := NewRatingPgRepository(db)

	mocks.MockUpdateReturnResultOK(mock, rating)
	err = ratingPgRep.Update(rating)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRatingPgRepository_Delete_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingPgRep := NewRatingPgRepository(db)

	mocks.MockDeleteReturnResultOK(mock, rating)
	err = ratingPgRep.Delete(rating)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRatingPgRepository_SelectByUserIDContentID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingPgRep := NewRatingPgRepository(db)

	mocks.MockSelectReturnRows(mock, rating)
	dbRating, err := ratingPgRep.SelectByUserIDContentID(rating.UserID, rating.ContentID)
	assert.Equal(t, rating, dbRating)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRatingPgRepository_SelectByUserIDContentID_Fail(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingPgRep := NewRatingPgRepository(db)

	mocks.MockSelectReturnErrNoRows(mock, rating)
	dbRating, err := ratingPgRep.SelectByUserIDContentID(rating.UserID, rating.ContentID)
	assert.Equal(t, dbRating, (*models.Rating)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRatingPgRepository_SelectLikesCount_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var contentID uint64 = 4
	likesCount := 5

	ratingPgRep := NewRatingPgRepository(db)

	mocks.MockSelectCountReturnRows(mock, contentID, likesCount)
	dbCount, err := ratingPgRep.SelectLikesCount(contentID)
	assert.Equal(t, dbCount, likesCount)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRatingPgRepository_SelectRatesCount_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var contentID uint64 = 4
	ratesCount := 5

	ratingPgRep := NewRatingPgRepository(db)

	mocks.MockSelectCountReturnRows(mock, contentID, ratesCount)
	dbCount, err := ratingPgRep.SelectRatesCount(contentID)
	assert.Equal(t, dbCount, ratesCount)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
