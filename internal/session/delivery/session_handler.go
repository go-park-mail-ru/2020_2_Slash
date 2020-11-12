package delivery

import (
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"github.com/go-park-mail-ru/2020_2_Slash/tools"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/CSRFManager"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
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
	e.POST("/api/v1/session", sh.LoginHandler())
	e.DELETE("/api/v1/session", sh.LogoutHandler(), mw.CheckAuth, mw.CheckCSRF)
}

func (sh *SessionHandler) LoginHandler() echo.HandlerFunc {
	type Request struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,gte=6"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		dbUser, err := sh.userUcase.GetByEmail(req.Email)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		if err := sh.userUcase.CheckPassword(dbUser, req.Password); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		sess := models.NewSession(dbUser.ID)
		if err = sh.sessUcase.Create(sess); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		token, err := CSRFManager.CreateToken(sess)
		if err != nil {
			logger.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		cntx.Response().Header().Set("X-Csrf-Token", token)

		cookie := tools.CreateCookie(sess)
		cntx.SetCookie(cookie)
		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"user": dbUser,
			},
		})
	}
}

func (sh *SessionHandler) LogoutHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		session, hasCookie := cntx.Cookie(SessionName)

		if hasCookie == http.ErrNoCookie {
			err := errors.Get(CodeUserUnauthorized)
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		err := sh.sessUcase.Delete(session.Value)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		SetOverdueCookie(cntx, session)

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func SetOverdueCookie(cntx echo.Context, cookie *http.Cookie) {
	cookie.Path = "/"
	cookie.Expires = time.Now().AddDate(0, 0, -2)
	cntx.SetCookie(cookie)
}
