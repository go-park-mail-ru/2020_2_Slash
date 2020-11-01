package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActorPgRepository_Insert_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	actor := &models.Actor{
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoInsertReturnRows(mock, 0, actor.Name)
	err = actorPgRep.Insert(actor)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActorPgRepository_Insert_DBFail(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	actor := &models.Actor{
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoInsertReturnErrNoRows(mock, 0, actor.Name)
	err = actorPgRep.Insert(actor)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActorPgRepository_Insert_ActorAlreadyExist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	actor := &models.Actor{
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoInsertReturnRows(mock, 0, actor.Name)
	err = actorPgRep.Insert(actor)
	assert.NoError(t, err)

	mocks.MockActorRepoInsertReturnRows(mock, 1, actor.Name)
	err = actorPgRep.Insert(actor)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActorPgRepository_Update_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	actor := &models.Actor{
		ID:   3,
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoUpdateReturnResultOk(mock, actor.ID, actor.Name)
	err = actorPgRep.Update(actor)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActorPgRepository_Update_NoActorWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	actor := &models.Actor{
		ID:   3,
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoUpdateReturnResultZero(mock, actor.ID, actor.Name)
	err = actorPgRep.Update(actor)
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

	actor := &models.Actor{
		ID:   3,
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoDeleteReturnResultOk(mock, actor.ID, actor.Name)
	err = actorPgRep.DeleteById(actor.ID)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActorPgRepository_SelectById_OK(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	actor := &models.Actor{
		ID:   3,
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoSelectReturnRows(mock, actor.ID, actor.Name)
	dbActor, err := actorPgRep.SelectById(actor.ID)
	assert.Equal(t, actor, dbActor)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActorPgRepository_SelectById_NoActorWithThisID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	actor := &models.Actor{
		ID:   3,
		Name: "Brad Pitt",
	}

	actorPgRep := NewActorPgRepository(db)

	mocks.MockActorRepoSelectReturnErrNoRows(mock, actor.ID)
	dbActor, err := actorPgRep.SelectById(actor.ID)
	assert.Equal(t, dbActor, (*models.Actor)(nil))
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
