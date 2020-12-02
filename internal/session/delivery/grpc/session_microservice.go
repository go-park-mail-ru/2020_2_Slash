package grpc

import (
	"database/sql"
	"time"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"

	"context"

	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type SessionBlockMicroservice struct {
	sessRepo session.SessionRepository
}

func NewSessionBlockMicroservice(repo session.SessionRepository) *SessionBlockMicroservice {
	return &SessionBlockMicroservice{
		sessRepo: repo,
	}
}

func (sm *SessionBlockMicroservice) Create(cntx context.Context, sess *Session) (*emptypb.Empty, error) {
	if err := sm.sessRepo.Insert(GrpcSessionToModel(sess)); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (sm *SessionBlockMicroservice) Get(cntx context.Context, sessValue *SessionValue) (*Session, error) {
	sess, err := sm.sessRepo.SelectByValue(sessValue.GetValue())
	switch {
	case err == sql.ErrNoRows:
		return nil, status.Error(codes.Code(consts.CodeUserUnauthorized), "")
	case err != nil:
		return nil, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return ModelSessionToGrpc(sess), nil
}

func (sm *SessionBlockMicroservice) Delete(cntx context.Context, sessValue *SessionValue) (*emptypb.Empty, error) {
	if !sm.isExist(sessValue.GetValue()) {
		return &emptypb.Empty{}, status.Error(codes.Code(consts.CodeSessionDoesNotExist), "")
	}
	if err := sm.sessRepo.DeleteByValue(sessValue.GetValue()); err != nil {
		return &emptypb.Empty{}, status.Error(codes.Code(consts.CodeInternalError), err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (sm *SessionBlockMicroservice) Check(cntx context.Context, sessValue *SessionValue) (*Session, error) {
	sess, err := sm.Get(cntx, sessValue)
	if err != nil {
		return nil, err
	}

	sessModel := GrpcSessionToModel(sess)
	if sessModel.ExpiresAt.Before(time.Now()) {
		_, err := sm.Delete(cntx, sessValue)
		if err != nil {
			return nil, err
		}
		return nil, status.Error(codes.Code(consts.CodeSessionExpired), "")
	}
	return sess, nil
}

func (sm *SessionBlockMicroservice) isExist(sessValue string) bool {
	_, err := sm.Get(context.Background(), &SessionValue{Value: sessValue})
	return err == nil
}
