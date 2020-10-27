package usecases

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
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
	sess, err := su.sessRepo.SelectByID(sessValue)
	switch {
	case err == sql.ErrNoRows:
		return nil, errors.Get(CodeUserUnauthorized)
	case err != nil:
		return nil, errors.New(CodeInternalError, err)
	}
	return sess, nil
}
