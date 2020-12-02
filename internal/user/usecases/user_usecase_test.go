package usecases

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var userModel = &models.User{
	Nickname: "Jhon",
	Email:    "jhon@gmail.com",
	Password: "hardpassword",
}
var userInst = user.ModelUserToGrpc(userModel)

func TestUserUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userClient := mocks.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userClient.
		EXPECT().
		Create(context.Background(), userInst).
		Return(userInst, nil)

	err := userUseCase.Create(userModel)
	assert.Equal(t, err, (*errors.Error)(nil))
}

func TestUserUseCase_Update_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userClient := mocks.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userClient.
		EXPECT().
		UpdateProfile(context.Background(), userInst).
		Return(userInst, nil)

	dbUser, err := userUseCase.UpdateProfile(userModel)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userModel)
}

func TestUserUseCase_UpdateAvatar_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newAvatar := "/avatar"

	userClient := mocks.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userAvatar := &user.IdAvatar{
		Id:     &user.ID{ID: userModel.ID},
		Avatar: &user.Avatar{Avatar: newAvatar},
	}

	userClient.
		EXPECT().
		UpdateAvatar(context.Background(), userAvatar).
		Return(userInst, nil)

	dbUser, err := userUseCase.UpdateAvatar(userModel.ID, newAvatar)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userModel)
}

func TestUserUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userClient := mocks.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userClient.
		EXPECT().
		GetByID(context.Background(), &user.ID{ID: userModel.ID}).
		Return(userInst, nil)

	dbUser, err := userUseCase.GetByID(userModel.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userModel)
}

func TestUserUseCase_GetByEmail_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userClient := mocks.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userClient.
		EXPECT().
		GetByEmail(context.Background(), &user.Email{Email: userModel.Email}).
		Return(userInst, nil)

	dbUser, err := userUseCase.GetByEmail(userModel.Email)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userModel)
}