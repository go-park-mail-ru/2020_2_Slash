package session

import "github.com/go-park-mail-ru/2020_2_Slash/internal/models"

type SessionRepository interface {
	Insert(session *models.Session) error
}
