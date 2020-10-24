package session

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type SessionUsecase interface {
	Create(userID uint64) (*models.Session, error)
}
