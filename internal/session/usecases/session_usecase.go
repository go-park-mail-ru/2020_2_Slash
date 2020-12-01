package usecases

import (
	"context"

	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	sessGRPC "github.com/go-park-mail-ru/2020_2_Slash/internal/session/delivery/grpc"
	"google.golang.org/grpc/status"
)

type SessionUsecase struct {
	sessBlockClient sessGRPC.SessionBlockClient
}

func NewSessionUsecase(client sessGRPC.SessionBlockClient) session.SessionUsecase {
	return &SessionUsecase{
		sessBlockClient: client,
	}
}

func (su *SessionUsecase) Create(sess *models.Session) *errors.Error {
	_, err := su.sessBlockClient.Create(context.Background(), sessGRPC.ModelSessionToGrpc(sess))
	if err != nil {
		customErr := errors.Get(ErrorCode(status.Code(err)))
		return customErr
	}
	return nil
}

func (su *SessionUsecase) Get(sessValue string) (*models.Session, *errors.Error) {
	sess, err := su.sessBlockClient.Get(context.Background(), &sessGRPC.SessionValue{Value: sessValue})
	if err != nil {
		customErr := errors.Get(ErrorCode(status.Code(err)))
		return nil, customErr
	}
	return sessGRPC.GrpcSessionToModel(sess), nil
}

func (su *SessionUsecase) Delete(sessValue string) *errors.Error {
	_, err := su.sessBlockClient.Delete(context.Background(), &sessGRPC.SessionValue{Value: sessValue})
	if err != nil {
		customErr := errors.Get(ErrorCode(status.Code(err)))
		return customErr
	}
	return nil
}

func (su *SessionUsecase) Check(sessValue string) (*models.Session, *errors.Error) {
	sess, err := su.sessBlockClient.Check(context.Background(), &sessGRPC.SessionValue{Value: sessValue})
	if err != nil {
		customErr := errors.Get(ErrorCode(status.Code(err)))
		return nil, customErr
	}
	return sessGRPC.GrpcSessionToModel(sess), nil
}
