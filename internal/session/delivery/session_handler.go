package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"github.com/labstack/echo/v4"
)

type SessionHandler struct {
	sessUcase session.SessionUsecase
}

func NewSessionHandler(sessUcase session.SessionUsecase) *SessionHandler {
	return &SessionHandler{
		sessUcase: sessUcase,
	}
}

func (sh *SessionHandler) Configure(e *echo.Echo) {
}
