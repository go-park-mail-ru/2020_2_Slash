package session

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
)

type SessionUsecase interface {
	Create(sess *models.Session) *errors.Error
}
