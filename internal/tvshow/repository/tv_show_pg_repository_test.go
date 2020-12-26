package repository

import (
	"log"
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

func TestTVShowPgRepository_SelectByParams_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

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
		Type:             "tvshow",
	}

	content := []*models.Content{
		contentInst,
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	params := &models.ContentFilter{
		Year: []int{2001},
	}

	mocks.MockTVShowRepoSelectByParamsReturnRows(mock, params, pgnt, userID, tvshows)
	dbTVShows, err := tvshowPgRep.SelectByParams(params, pgnt, userID)
	log.Println(dbTVShows)
	assert.Equal(t, tvshows, dbTVShows)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_SelectLatest_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	mocks.MockTVShowRepoSelectLatestReturnRows(mock, pgnt, userID, tvshows)
	dbTVShows, err := tvshowPgRep.SelectLatest(pgnt, userID)
	assert.Equal(t, tvshows, dbTVShows)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTVShowPgRepository_SelectByRating_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	tvshowPgRep := NewTVShowPgRepository(db)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 1

	mocks.MockTVShowRepoSelectByRatingReturnRows(mock, pgnt, userID, tvshows)
	dbTVShows, err := tvshowPgRep.SelectByRating(pgnt, userID)
	assert.Equal(t, tvshows, dbTVShows)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
