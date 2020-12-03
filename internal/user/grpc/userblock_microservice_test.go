package grpc

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var userModel = &models.User{
	Nickname: "Jhon",
	Email:    "jhon@gmail.com",
	Password: "hardpassword",
	Avatar:   "/avatar",
}
var userInst = user.ModelUserToGrpc(userModel)

func TestUserUseCase_Create_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserblockMicroservice(userRep)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(userModel.Email)).
		Return(userModel, nil)

	_, err := userUseCase.Create(context.Background(), userInst)
	assert.Equal(t, err, status.Error(codes.Code(consts.CodeEmailAlreadyExists), ""))
}

func TestUserUseCase_Update_Fail(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserblockMicroservice(userRep)

	newUserModel := &models.User{
		Nickname: "Jhon1",
		Email:    "jhon1@gmail.com",
		Password: "veryhardpassword",
	}
	newUserInst := user.ModelUserToGrpc(newUserModel)

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(newUserInst.ID)).
		Return(userModel, nil)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(newUserInst.Email)).
		Return(userModel, nil)

	_, err := userUseCase.UpdateProfile(context.Background(), newUserInst)
	assert.Equal(t, err, status.Error(codes.Code(consts.CodeEmailAlreadyExists), ""))
}

func TestUserUseCase_UpdateAvatar_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserblockMicroservice(userRep)

	newAvatar := "/avatar"

	idAvatar := &user.IdAvatar{
		Id:     &user.ID{ID: userModel.ID},
		Avatar: &user.Avatar{Avatar: newAvatar},
	}

	newUserModel := userModel
	newUserModel.Avatar = newAvatar

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(userInst.ID)).
		Return(userModel, nil)

	userRep.
		EXPECT().
		Update(gomock.Eq(newUserModel)).
		Return(nil)

	_, err := userUseCase.UpdateAvatar(context.Background(), idAvatar)
	assert.NotEqual(t, err, (error)(nil))
}

func TestUserUseCase_GetByID_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserblockMicroservice(userRep)

	userRep.
		EXPECT().
		SelectByID(gomock.Eq(userInst.ID)).
		Return(userModel, nil)

	dbUser, err := userUseCase.GetByID(context.Background(), &user.ID{ID: userInst.ID})
	assert.Equal(t, err, (error)(nil))
	assert.Equal(t, dbUser, userInst)
}

func TestUserUseCase_GetByEmail_OK(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRep := mocks.NewMockUserRepository(ctrl)
	userUseCase := NewUserblockMicroservice(userRep)

	userRep.
		EXPECT().
		SelectByEmail(gomock.Eq(userInst.Email)).
		Return(userModel, nil)

	dbUser, err := userUseCase.GetByEmail(context.Background(), &user.Email{Email: userInst.Email})
	assert.Equal(t, err, (error)(nil))
	assert.Equal(t, dbUser, userInst)
}
