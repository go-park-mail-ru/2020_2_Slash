package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectorPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	director := &models.Director{
		Name: "Quentin Tarantino",
	}

	directorPgRepository := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoInsertReturnRows(mock, 0, director.Name)
	err = directorPgRepository.Insert(director)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDirectorPgRepository_Insert_DBFail(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	director := &models.Director{
		Name: "Quentin Tarantino",
	}

	directorPgRep := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoInsertReturnErrNoRows(mock, 0, director.Name)
	err = directorPgRep.Insert(director)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDirectorPgRepository_Insert_DirectorAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	director := &models.Director{
		Name: "Quentin Tarantino",
	}

	directorPgRep := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoInsertReturnRows(mock, 0, director.Name)
	err = directorPgRep.Insert(director)
	assert.NoError(t, err)

	mocks.MockDirectorRepoInsertReturnRows(mock, 1, director.Name)
	err = directorPgRep.Insert(director)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDirectorPgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	director := &models.Director{
		ID:   3,
		Name: "Quentin Tarantino",
	}

	directorPgRepository := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoUpdateReturnResultOk(mock, director.ID, director.Name)
	err = directorPgRepository.Update(director)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDirectorPgRepository_Update_NoDirectorWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	director := &models.Director{
		ID:   3,
		Name: "Quentin Tarantino",
	}

	directorPgRepository := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoUpdateReturnResultZero(mock, director.ID, director.Name)
	err = directorPgRepository.Update(director)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDirectorPgRepository_DeleteById_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	director := &models.Director{
		ID:   3,
		Name: "Quentin Tarantino",
	}

	directorPgRepository := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoDeleteReturnResultOk(mock, director.ID, director.Name)
	err = directorPgRepository.DeleteById(director.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDirectorPgRepository_SelectById_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	director := &models.Director{
		ID:   3,
		Name: "Quentin Tarantino",
	}

	directorPgRep := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoSelectReturnRows(mock, director.ID, director.Name)
	dbDirector, err := directorPgRep.SelectById(director.ID)
	assert.Equal(t, director, dbDirector)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDirectorPgRepository_SelectById_NoDirectorWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	director := &models.Director{
		ID:   3,
		Name: "Quentin Tarantino",
	}

	directorPgRepository := NewDirectorPgRepository(db)

	mocks.MockDirectorRepoSelectReturnErrNoRows(mock, director.ID)
	dbDirector, err := directorPgRepository.SelectById(director.ID)
	assert.Equal(t, dbDirector, (*models.Director)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
