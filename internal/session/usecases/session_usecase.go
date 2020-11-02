package usecases

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"time"
)

type SessionUsecase struct {
	sessRepo session.SessionRepository
}

func NewSessionUsecase(repo session.SessionRepository) session.SessionUsecase {
	return &SessionUsecase{
		sessRepo: repo,
	}
}

func (su *SessionUsecase) Create(sess *models.Session) *errors.Error {
	if err := su.sessRepo.Insert(sess); err != nil {
		return errors.New(CodeInternalError, err)
	}
	return nil
}

func (su *SessionUsecase) Get(sessValue string) (*models.Session, *errors.Error) {
	sess, err := su.sessRepo.SelectByValue(sessValue)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeUserUnauthorized)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return sess, nil
}

func (su *SessionUsecase) IsExist(sessValue string) bool {
	_, err := su.Get(sessValue)
	return err == nil
}

func (su *SessionUsecase) Delete(sessValue string) *errors.Error {
	if !su.IsExist(sessValue) {
		return errors.Get(CodeSessionDoesNotExist)
	}

	if err := su.sessRepo.DeleteByValue(sessValue); err != nil {
		return errors.New(CodeInternalError, err)
	}

	return nil
}

func (su *SessionUsecase) Check(sessValue string) (*models.Session, *errors.Error) {
	sess, customErr := su.Get(sessValue)
	if customErr != nil {
		return nil, customErr
	}
	if sess.ExpiresAt.Before(time.Now()) {
		deleteErr := su.Delete(sessValue)
		if deleteErr != nil {
			return nil, deleteErr
		}
		customErr = errors.Get(CodeSessionExpired)
		return nil, customErr
	}
	return sess, nil
}
