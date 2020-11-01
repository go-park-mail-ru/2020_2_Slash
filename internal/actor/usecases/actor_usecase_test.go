package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
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
