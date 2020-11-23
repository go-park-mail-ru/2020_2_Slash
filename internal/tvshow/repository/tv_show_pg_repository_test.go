package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/mocks"
	"github.com/stretchr/testify/assert"
)

var contentInst *models.Content = &models.Content{
	ContentID: 1,
}

var tvshowInst *models.TVShow = &models.TVShow{
	ID:      1,
	Content: *contentInst,
}

func TestTVShowPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	mocks.MockTVShowRepoInsertReturnRows(mock, 1, tvshowInst)
	err = tvshowPgRep.Insert(tvshowInst)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_Insert_ContentAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	mocks.MockTVShowRepoInsertReturnRows(mock, 1, tvshowInst)
	err = tvshowPgRep.Insert(tvshowInst)
	assert.NoError(t, err)

	mocks.MockTVShowRepoInsertReturnErrNoUniq(mock, 2, tvshowInst)
	err = tvshowPgRep.Insert(tvshowInst)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	mocks.MockTVShowRepoSelectByIDReturnRows(mock, tvshowInst.ID, tvshowInst)
	dbTVShow, err := tvshowPgRep.SelectByID(tvshowInst.ID)
	assert.Equal(t, tvshowInst, dbTVShow)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_SelectById_NoTVShowWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	mocks.MockTVShowRepoSelectByIDReturnErrNoRows(mock, tvshowInst.ID)
	dbTVShow, err := tvshowPgRep.SelectByID(tvshowInst.ID)
	assert.Equal(t, dbTVShow, (*models.TVShow)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_SelectFullByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	var userID uint64 = 1
	mocks.MockTVShowRepoSelectFullByIDReturnRows(mock, tvshowInst.ID, userID, tvshowInst)
	dbTVShow, err := tvshowPgRep.SelectFullByID(tvshowInst.ID, userID)
	assert.Equal(t, tvshowInst, dbTVShow)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_SelectByContentID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	mocks.MockTVShowRepoSelectByContentIDReturnRows(mock, tvshowInst.ID, tvshowInst)
	dbTVShow, err := tvshowPgRep.SelectByContentID(tvshowInst.ContentID)
	assert.Equal(t, tvshowInst, dbTVShow)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_SelectById_NoTVShowWithThisContentID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	mocks.MockTVShowRepoSelectByContentIDReturnErrNoRows(mock, tvshowInst)
	dbTVShow, err := tvshowPgRep.SelectByContentID(tvshowInst.ContentID)
	assert.Equal(t, dbTVShow, (*models.TVShow)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
