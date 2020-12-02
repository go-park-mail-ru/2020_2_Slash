package grpc

import (
	context "context"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var sessModel = models.NewSession(3)
var sess = ModelSessionToGrpc(sessModel)

func TestSessionUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionClient := NewSessionBlockMicroservice(sessionRep)

	sessionRep.
		EXPECT().
		Insert(gomock.Eq(GrpcSessionToModel(sess))).
		Return(nil)

	_, err := sessionClient.Create(context.Background(), sess)
	assert.Equal(t, err, (error)(nil))
}

func TestSessionUseCase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionClient := NewSessionBlockMicroservice(sessionRep)

	sessionRep.
		EXPECT().
		SelectByValue(gomock.Eq(sessModel.Value)).
		Return(sessModel, nil)

	dbSession, err := sessionClient.Get(context.Background(), &SessionValue{Value: sessModel.Value})
	assert.Equal(t, err, (error)(nil))
	assert.Equal(t, dbSession.Value, sessModel.Value)
}

func TestSessionUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionClient := NewSessionBlockMicroservice(sessionRep)

	sessionRep.
		EXPECT().
		SelectByValue(gomock.Eq(sessModel.Value)).
		Return(sessModel, nil)

	sessionRep.
		EXPECT().
		DeleteByValue(gomock.Eq(sessModel.Value)).
		Return(nil)

	_, err := sessionClient.Delete(context.Background(), &SessionValue{Value: sessModel.Value})
	assert.Equal(t, err, (error)(nil))
}

func TestSessionUsecase_Check_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionClient := NewSessionBlockMicroservice(sessionRep)

	sessionRep.
		EXPECT().
		SelectByValue(gomock.Eq(sessModel.Value)).
		Return(sessModel, nil)

	dbSession, err := sessionClient.Check(context.Background(), &SessionValue{Value: sessModel.Value})
	assert.Equal(t, err, (error)(nil))
	assert.Equal(t, dbSession.Value, sessModel.Value)
}

func TestSessionUsecase_Check_Expired(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionClient := NewSessionBlockMicroservice(sessionRep)

	sessModel := &models.Session{
		ID:        1,
		Value:     uuid.NewV4().String(),
		UserID:    3,
		ExpiresAt: time.Date(2019, 1, 1, 2, 3, 2, 3, time.UTC),
	}

	sessionRep.
		EXPECT().
		SelectByValue(gomock.Eq(sessModel.Value)).
		Return(sessModel, nil)

	sessionRep.
		EXPECT().
		DeleteByValue(gomock.Eq(sessModel.Value)).
		Return(nil)

	_, err := sessionClient.Delete(context.Background(), &SessionValue{Value: sessModel.Value})
	assert.Equal(t, err, (error)(nil))
}
