package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"github.com/go-park-mail-ru/2020_2_Slash/tools"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SessionHandler struct {
	sessUcase session.SessionUsecase
	userUcase user.UserUsecase
}

func NewSessionHandler(sessUcase session.SessionUsecase,
	userUcase user.UserUsecase) *SessionHandler {
	return &SessionHandler{
		sessUcase: sessUcase,
		userUcase: userUcase,
	}
}

func (sh *SessionHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/session", sh.loginHandler())
}

func (sh *SessionHandler) loginHandler() echo.HandlerFunc {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,gte=6"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		dbUser, err := sh.userUcase.GetByEmail(req.Email)
		if err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		if err := sh.userUcase.CheckPassword(dbUser, req.Password); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		sess := models.NewSession(dbUser.ID)
		if err = sh.sessUcase.Create(sess); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		cookie := tools.CreateCookie(sess)
		cntx.SetCookie(cookie)
		return cntx.JSON(http.StatusOK, dbUser.Sanitize())
	}
}
