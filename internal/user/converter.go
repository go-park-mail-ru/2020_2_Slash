package user

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

func GrpcUserToModel(grpcUser *User) *models.User {
	return &models.User{
		ID:       grpcUser.ID,
		Nickname: grpcUser.Nickname,
		Email:    grpcUser.Email,
		Password: grpcUser.Password,
		Avatar:   grpcUser.Avatar,
		Role:     grpcUser.Role,
	}
}

func ModelUserToGrpc(modelUser *models.User) *User {
	return &User{
		ID:       modelUser.ID,
		Nickname: modelUser.Nickname,
		Email:    modelUser.Email,
		Password: modelUser.Password,
		Avatar:   modelUser.Avatar,
		Role:     modelUser.Role,
	}
}
