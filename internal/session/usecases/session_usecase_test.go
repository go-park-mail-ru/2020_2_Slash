package usecases

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	sessGRPC "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery/grpc"
	sessMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery/grpc/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

var sessModel = models.NewSession(3)
var sess = sessGRPC.ModelSessionToGrpc(sessModel)

func TestSessionUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionClient := sessMocks.NewMockSessionBlockClient(ctrl)
	sessionUseCase := NewSessionUsecase(sessionClient)

	sessionClient.
		EXPECT().
		Create(context.Background(), sess).
		Return(&emptypb.Empty{}, nil)

	err := sessionUseCase.Create(sessModel)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestSessionUseCase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionClient := sessMocks.NewMockSessionBlockClient(ctrl)
	sessionUseCase := NewSessionUsecase(sessionClient)

	sessionClient.
		EXPECT().
		Get(context.Background(), &sessGRPC.SessionValue{Value: sessModel.Value}).
		Return(sess, nil)

	dbSession, err := sessionUseCase.Get(sessModel.Value)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbSession.Value, sessModel.Value)
}

func TestSessionUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionClient := sessMocks.NewMockSessionBlockClient(ctrl)
	sessionUseCase := NewSessionUsecase(sessionClient)

	sessionClient.
		EXPECT().
		Delete(context.Background(), &sessGRPC.SessionValue{Value: sessModel.Value}).
		Return(&emptypb.Empty{}, nil)

	err := sessionUseCase.Delete(sessModel.Value)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestSessionUsecase_Check_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionClient := sessMocks.NewMockSessionBlockClient(ctrl)
	sessionUseCase := NewSessionUsecase(sessionClient)

	sessionClient.
		EXPECT().
		Check(context.Background(), &sessGRPC.SessionValue{Value: sessModel.Value}).
		Return(sess, nil)

	dbSession, err := sessionUseCase.Check(sessModel.Value)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbSession.Value, sessModel.Value)
}
