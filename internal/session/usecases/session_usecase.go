package usecases

import (
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

func (su *SessionUsecase) Create(userID uint64) (*models.Session, error) {
	sess := models.NewSession(userID)
	err := su.sessRepo.Insert(sess)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
