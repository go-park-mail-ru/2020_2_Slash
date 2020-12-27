package usecases

import (
	"context"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user/delivery/grpc"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userBlockClient grpc.UserBlockClient
}

func NewUserUsecase(client grpc.UserBlockClient) *UserUsecase {
	return &UserUsecase{
		userBlockClient: client,
	}
}

func (uu *UserUsecase) Create(modelUser *models.User) *errors.Error {
	grpcUser, err := uu.userBlockClient.Create(context.Background(),
		grpc.ModelUserToGrpc(modelUser))
	if err != nil {
		customErr := errors.GetCustomErrFromStatus(err)
		return customErr
	}

	err = copier.Copy(modelUser, grpc.GrpcUserToModel(grpcUser))
	if err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (uu *UserUsecase) UpdatePassword(userID uint64, oldPassword, newPassword,
	repeatedNewPassword string) (*models.User, *errors.Error) {
	grpcUser, err := uu.userBlockClient.UpdatePassword(context.Background(),
		&grpc.UpdatePasswordMsg{
			Id:                  userID,
			OldPassword:         oldPassword,
			NewPassword:         newPassword,
			RepeatedNewPassword: repeatedNewPassword,
		})
	if err != nil {
		customErr := errors.GetCustomErrFromStatus(err)
		return nil, customErr
	}

	return grpc.GrpcUserToModel(grpcUser), nil
}

func (uu *UserUsecase) GetByEmail(email string) (*models.User, *errors.Error) {
	grpcUser, err := uu.userBlockClient.GetByEmail(context.Background(),
		&grpc.Email{Email: email})
	if err != nil {
		customErr := errors.GetCustomErrFromStatus(err)
		return nil, customErr
	}

	return grpc.GrpcUserToModel(grpcUser), nil
}

func (uu *UserUsecase) GetByID(userID uint64) (*models.User, *errors.Error) {
	grpcUser, err := uu.userBlockClient.GetByID(context.Background(),
		&grpc.ID{ID: userID})
	if err != nil {
		customErr := errors.GetCustomErrFromStatus(err)
		return nil, customErr
	}

	return grpc.GrpcUserToModel(grpcUser), nil
}

func (uu *UserUsecase) UpdateProfile(newUserData *models.User) (*models.User, *errors.Error) {
	grpcUser, err := uu.userBlockClient.UpdateProfile(context.Background(),
		grpc.ModelUserToGrpc(newUserData))
	if err != nil {
		customErr := errors.GetCustomErrFromStatus(err)
		return nil, customErr
	}

	return grpc.GrpcUserToModel(grpcUser), nil
}

func (uu *UserUsecase) UpdateAvatar(userID uint64, newAvatar string) (*models.User, *errors.Error) {
	grpcUser, err := uu.userBlockClient.UpdateAvatar(context.Background(),
		&grpc.IdAvatar{
			Id:     &grpc.ID{ID: userID},
			Avatar: &grpc.Avatar{Avatar: newAvatar},
		})
	if err != nil {
		customErr := errors.GetCustomErrFromStatus(err)
		return nil, customErr
	}

	return grpc.GrpcUserToModel(grpcUser), nil
}

func (uu *UserUsecase) CheckPassword(user *models.User, password string) *errors.Error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),
		[]byte(password)); err != nil {
		return errors.Get(CodeWrongPassword)
	}
	return nil
}

func (uu *UserUsecase) IsAdmin(userID uint64) (bool, *errors.Error) {
	user, err := uu.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.Role == Admin, nil
}
