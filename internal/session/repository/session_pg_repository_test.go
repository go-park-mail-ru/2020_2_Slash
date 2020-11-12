package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessionPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	session := models.NewSession(3)

	sessionPgRepository := NewSessionPgRepository(db)

	mocks.MockInsertReturnRows(mock, session)
	err = sessionPgRepository.Insert(session)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionPgRepository_Delete_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	session := models.NewSession(3)

	sessionPgRepository := NewSessionPgRepository(db)

	mocks.MockDeleteReturnResultOk(mock, session.Value)
	err = sessionPgRepository.DeleteByValue(session.Value)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestSessionPgRepository_SelectByValue_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	session := models.NewSession(3)

	sessionPgRepository := NewSessionPgRepository(db)

	mocks.MockSelectReturnRows(mock, session)
	dbSession, err := sessionPgRepository.SelectByValue(session.Value)
	assert.Equal(t, session, dbSession)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionPgRepository_SelectByValue_NoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	session := models.NewSession(3)

	sessionPgRepository := NewSessionPgRepository(db)

	mocks.MockSelectReturnErrNoRows(mock, session.Value)
	dbSession, err := sessionPgRepository.SelectByValue(session.Value)
	assert.Equal(t, (*models.Session)(nil), dbSession)
	assert.Equal(t, sql.ErrNoRows, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}