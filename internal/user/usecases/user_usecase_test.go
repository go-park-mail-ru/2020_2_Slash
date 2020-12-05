package usecases

import (
	"context"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery/grpc"
	mocks2 "github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery/grpc/mocks"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var userModel = &models.User{
	Nickname: "Jhon",
	Email:    "jhon@gmail.com",
	Password: "hardpassword",
}
var userInst = grpc.ModelUserToGrpc(userModel)

func TestUserUseCase_Create_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userClient := mocks2.NewMockUserBlockClient(ctrl)
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

	userClient := mocks2.NewMockUserBlockClient(ctrl)
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

	userClient := mocks2.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userAvatar := &grpc.IdAvatar{
		Id:     &grpc.ID{ID: userModel.ID},
		Avatar: &grpc.Avatar{Avatar: newAvatar},
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

	userClient := mocks2.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userClient.
		EXPECT().
		GetByID(context.Background(), &grpc.ID{ID: userModel.ID}).
		Return(userInst, nil)

	dbUser, err := userUseCase.GetByID(userModel.ID)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userModel)
}

func TestUserUseCase_GetByEmail_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userClient := mocks2.NewMockUserBlockClient(ctrl)
	userUseCase := NewUserUsecase(userClient)

	userClient.
		EXPECT().
		GetByEmail(context.Background(), &grpc.Email{Email: userModel.Email}).
		Return(userInst, nil)

	dbUser, err := userUseCase.GetByEmail(userModel.Email)
	assert.Equal(t, err, (*errors.Error)(nil))
	assert.Equal(t, dbUser, userModel)
}
