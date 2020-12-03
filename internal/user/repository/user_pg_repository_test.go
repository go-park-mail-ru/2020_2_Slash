package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/stretchr/testify/assert"
)

var userInst = &models.User{
	Nickname: "Jhon",
	Email:    "jhon@gmail.com",
	Password: "hardpassword",
}

func TestUserPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userPgRep := NewUserPgRepository(db)

	mocks.MockUserRepoInsertReturnRows(mock, userInst)
	err = userPgRep.Insert(userInst)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_Insert_UserAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userPgRep := NewUserPgRepository(db)

	mocks.MockUserRepoInsertReturnRows(mock, userInst)
	err = userPgRep.Insert(userInst)
	assert.NoError(t, err)

	mocks.MockUserRepoInsertReturnErrNoUniq(mock, userInst)
	err = userPgRep.Insert(userInst)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userPgRep := NewUserPgRepository(db)

	mocks.MockUserRepoUpdateReturnResultOk(mock, userInst)
	err = userPgRep.Update(userInst)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_SelectByID_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userPgRep := NewUserPgRepository(db)

	mocks.MockUserRepoSelectByIDReturnRows(mock, userInst)
	dbUser, err := userPgRep.SelectByID(userInst.ID)
	assert.Equal(t, userInst, dbUser)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserPgRepository_SelectByEmail_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userPgRep := NewUserPgRepository(db)

	mocks.MockUserRepoSelectByEmailReturnRows(mock, userInst)
	dbUser, err := userPgRep.SelectByEmail(userInst.Email)
	assert.Equal(t, userInst, dbUser)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
