package usecases

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/admin"
	adminMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/admin/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActorUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	actorUseCase := NewActorUseCase(actorRep, adminPanelClient)

	actor := &models.Actor{
		Name: "Jamie Fox",
	}

	grpcActor := admin.ActorModelToGRPC(actor)

	adminPanelClient.
		EXPECT().
		CreateActor(context.Background(), grpcActor).
		Return(grpcActor, nil)

	err := actorUseCase.Create(actor)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestActorUseCase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	actorUseCase := NewActorUseCase(actorRep, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	actorUseCase := NewActorUseCase(actorRep, adminPanelClient)

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
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	actorUseCase := NewActorUseCase(actorRep, adminPanelClient)

	actor := &models.Actor{
		ID:   3,
		Name: "Margo Robbie",
	}

	adminPanelClient.
		EXPECT().
		DeleteActorByID(context.Background(), &admin.ID{ID: actor.ID}).
		Return(&empty.Empty{}, nil)

	err := actorUseCase.DeleteById(actor.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestActorUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	actorUseCase := NewActorUseCase(actorRep, adminPanelClient)

	actor := &models.Actor{
		ID:   3,
		Name: "Margo Robbie",
	}

	grpcActor := admin.ActorModelToGRPC(actor)
	adminPanelClient.
		EXPECT().
		ChangeActor(context.Background(), grpcActor).
		Return(&empty.Empty{}, nil)

	err := actorUseCase.Change(actor)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestActorUseCase_ListByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := mocks.NewMockActorRepository(ctrl)
	adminPanelClient := adminMocks.NewMockAdminPanelClient(ctrl)
	actorUseCase := NewActorUseCase(actorRep, adminPanelClient)

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
