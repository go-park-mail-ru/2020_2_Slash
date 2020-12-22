package usecases

import (
	"database/sql"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestActorUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	actorUseCase := NewActorUseCase(actorRep)

	actor := &models.Actor{
		Name: "Jamie Fox",
	}

	actorRep.
		EXPECT().
		Insert(gomock.Eq(actor)).
		Return(nil)

	err := actorUseCase.Create(actor)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestActorUseCase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	actorUseCase := NewActorUseCase(actorRep)

	actor := &models.Actor{
		ID:   3,
		Name: "Margo Robbie",
	}

	actorRep.
		EXPECT().
		SelectById(gomock.Eq(actor.ID)).
		Return(actor, nil)

	dbActor, err := actorUseCase.Get(actor.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbActor, actor)
}

func TestActorUseCase_Get_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	actorUseCase := NewActorUseCase(actorRep)

	actor := &models.Actor{
		ID:   3,
		Name: "Margo Robbie",
	}

	actorRep.
		EXPECT().
		SelectById(gomock.Eq(actor.ID)).
		Return(nil, sql.ErrNoRows)

	dbActor, err := actorUseCase.Get(actor.ID)
	assert.Equal(t, err, errors.Get(consts.CodeActorDoesNotExist))
	assert.Equal(t, dbActor, (*models.Actor)(nil))
}

func TestActorUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	actorUseCase := NewActorUseCase(actorRep)

	actor := &models.Actor{
		ID:   3,
		Name: "Margo Robbie",
	}

	actorRep.
		EXPECT().
		SelectById(gomock.Eq(actor.ID)).
		Return(actor, nil)

	actorRep.
		EXPECT().
		DeleteById(gomock.Eq(actor.ID)).
		Return(nil)

	err := actorUseCase.DeleteById(actor.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestActorUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	actorUseCase := NewActorUseCase(actorRep)

	actor := &models.Actor{
		ID:   3,
		Name: "Margo Robbie",
	}

	actorRep.
		EXPECT().
		SelectById(gomock.Eq(actor.ID)).
		Return(actor, nil)

	actorRep.
		EXPECT().
		Update(gomock.Eq(actor)).
		Return(nil)

	err := actorUseCase.Change(actor)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestActorUseCase_ListByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	actorUseCase := NewActorUseCase(actorRep)

	actors := []*models.Actor{
		&models.Actor{
			ID:   1,
			Name: "Margo Robbie",
		},
		&models.Actor{
			ID:   2,
			Name: "No Margo Robbie",
		},
	}

	actorsID := []uint64{1, 2}

	actorRep.
		EXPECT().
		SelectById(gomock.Eq(actorsID[0])).
		Return(actors[0], nil)

	actorRep.
		EXPECT().
		SelectById(gomock.Eq(actorsID[1])).
		Return(actors[1], nil)

	dbActors, err := actorUseCase.ListByID(actorsID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbActors, actors)
}

func TestActorUseCase_List_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	actorUseCase := NewActorUseCase(actorRep)

	actors := []*models.Actor{
		&models.Actor{
			ID:   1,
			Name: "Margo Robbie",
		},
		&models.Actor{
			ID:   2,
			Name: "No Margo Robbie",
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}

	actorRep.
		EXPECT().
		SelectAll(pgnt).
		Return(actors, nil)

	dbActors, err := actorUseCase.List(pgnt)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbActors, actors)
}
