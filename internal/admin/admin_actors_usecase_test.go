package admin

import (
	"context"
	actorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActorUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := actorMocks.NewMockActorRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		actorsRep: actorRep,
	}

	actor := &Actor{
		Name: "Jamie Fox",
	}

	actorRep.
		EXPECT().
		Insert(ActorGRPCToModel(actor)).
		Return(nil)

	_, err := adminMicroservice.CreateActor(context.Background(), actor)
	assert.Equal(t, err, (error)(nil))
}

func TestActorUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := actorMocks.NewMockActorRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		actorsRep: actorRep,
	}

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

	_, err := adminMicroservice.DeleteActorByID(context.Background(),
		&ID{ID: actor.ID})
	assert.Equal(t, err, (error)(nil))
}

func TestActorUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorRep := actorMocks.NewMockActorRepository(ctrl)
	adminMicroservice := &AdminMicroservice{
		actorsRep: actorRep,
	}
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

	_, err := adminMicroservice.ChangeActor(context.Background(), ActorModelToGRPC(actor))
	assert.Equal(t, err, (error)(nil))
}
