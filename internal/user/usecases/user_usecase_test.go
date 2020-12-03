package usecases

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var userInst = &models.User{
	Nickname: "Jhon",
	Email:    "jhon@gmail.com",
	Password: "hardpassword",
}

func TestUserUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUsecase(userRep)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(userInst.Email)).
		Return(nil, sql.ErrNoRows)

	userRep.
		EXPECT().
		Insert(gomock.Eq(userInst)).
		Return(nil)

	err := userUseCase.Create(userInst)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestUserUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUsecase(userRep)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(userInst.Email)).
		Return(userInst, nil)

	err := userUseCase.Create(userInst)
	assert.Equal(t, err, errors.Get(consts.CodeEmailAlreadyExists))
}

func TestUserUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUsecase(userRep)

	var userInst = &models.User{
		Nickname: "Jhon",
		Email:    "jhon@gmail.com",
		Password: "hardpassword",
	}

	newUserData := &models.User{
		Nickname: "Jhon1",
		Email:    "jhon1@gmail.com",
		Password: "veryhardpassword",
	}

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(userInst.ID)).
		Return(userInst, nil)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(newUserData.Email)).
		Return(nil, sql.ErrNoRows)

	userRep.
		EXPECT().
		Update(gomock.Eq(userInst)).
		Return(nil)

	_, err := userUseCase.UpdateProfile(userInst.ID, newUserData)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestUserUseCase_UpdateAvatar_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUsecase(userRep)

	newAvatar := "avatar"

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(userInst.ID)).
		Return(userInst, nil)

	userRep.
		EXPECT().
		Update(gomock.Eq(userInst)).
		Return(nil)

	_, err := userUseCase.UpdateAvatar(userInst.ID, newAvatar)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestUserUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUsecase(userRep)

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(userInst.ID)).
		Return(userInst, nil)

	dbUser, err := userUseCase.GetByID(userInst.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userInst)
}

func TestUserUseCase_GetByEmail_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserUsecase(userRep)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(userInst.Email)).
		Return(userInst, nil)

	dbUser, err := userUseCase.GetByEmail(userInst.Email)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userInst)
}
