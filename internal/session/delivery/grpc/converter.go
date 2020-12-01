package grpc

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/golang/protobuf/ptypes"
)

func GrpcSessionToModel(grpcSess *Session) *models.Session {
	ExpiresAt, _ := ptypes.Timestamp(grpcSess.ExpiresAt)

	return &models.Session{
		ID:        grpcSess.ID,
		Value:     grpcSess.Value,
		UserID:    grpcSess.UserID,
		ExpiresAt: ExpiresAt,
	}
}

func ModelSessionToGrpc(modelSess *models.Session) *Session {
	ExpiresAt, _ := ptypes.TimestampProto(modelSess.ExpiresAt)

	return &Session{
		ID:        modelSess.ID,
		Value:     modelSess.Value,
		UserID:    modelSess.UserID,
		ExpiresAt: ExpiresAt,
	}
}
