package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSessionUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionUseCase := NewSessionUsecase(sessionRep)

	session := models.NewSession(3)

	sessionRep.
		EXPECT().
		Insert(gomock.Eq(session)).
		Return(nil)

	err := sessionUseCase.Create(session)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestSessionUseCase_Get_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionUseCase := NewSessionUsecase(sessionRep)

	session := models.NewSession(3)

	sessionRep.
		EXPECT().
		SelectByValue(session.Value).
		Return(session, nil)

	dbSession, err := sessionUseCase.Get(session.Value)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbSession, session)
}

func TestSessionUseCase_Get_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionUseCase := NewSessionUsecase(sessionRep)

	session := models.NewSession(3)

	sessionRep.
		EXPECT().
		SelectByValue(session.Value).
		Return(nil, sql.ErrNoRows)

	dbSession, err := sessionUseCase.Get(session.Value)
	assert.Equal(t, err, errors.Get(consts.CodeUserUnauthorized))
	assert.Equal(t, dbSession, (*models.Session)(nil))
}

func TestSessionUseCase_Delete_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionUseCase := NewSessionUsecase(sessionRep)

	session := models.NewSession(3)

	sessionRep.
		EXPECT().
		SelectByValue(session.Value).
		Return(session, nil)

	sessionRep.
		EXPECT().
		DeleteByValue(session.Value).
		Return(nil)

	err := sessionUseCase.Delete(session.Value)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestSessionUsecase_Check_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionUseCase := NewSessionUsecase(sessionRep)

	session := models.NewSession(3)

	sessionRep.
		EXPECT().
		SelectByValue(session.Value).
		Return(session, nil)

	dbSession, err := sessionUseCase.Check(session.Value)
	assert.Equal(t, session, dbSession)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestSessionUsecase_Check_Expired(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionRep := mocks.NewMockSessionRepository(ctrl)
	sessionUseCase := NewSessionUsecase(sessionRep)

	session := &models.Session{
		ID:        1,
		Value:     uuid.NewV4().String(),
		UserID:    3,
		ExpiresAt: time.Date(2019, 1, 1, 2, 3, 2, 3, time.UTC),
	}

	sessionRep.
		EXPECT().
		SelectByValue(session.Value).
		Return(session, nil).AnyTimes()
	sessionRep.
		EXPECT().
		DeleteByValue(session.Value).
		Return(nil)

	_, err := sessionUseCase.Check(session.Value)
	assert.Equal(t, err, errors.Get(consts.CodeSessionExpired))
}

