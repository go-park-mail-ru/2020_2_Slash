package mwares

import (
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type MiddlewareManager struct {
	sessUcase session.SessionUsecase
	userUcase user.UserUsecase
	origins   []string
}

func NewMiddlewareManager(sessUcase session.SessionUsecase,
	userUcase user.UserUsecase) *MiddlewareManager {
	return &MiddlewareManager{
		sessUcase: sessUcase,
		userUcase: userUcase,
		origins:   []string{"http://www.flicksbox.ru"},
	}
}

func (m *MiddlewareManager) PanicRecovering(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logrus.Warn(err)
			}
		}()

		return next(cntx)
	}
}

func (m *MiddlewareManager) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		logrus.Info(cntx.Request().RemoteAddr, " ", cntx.Request().Method, " ", cntx.Request().URL)

		start := time.Now()
		err := next(cntx)
		end := time.Now()

		logrus.Info("Status: ", cntx.Response().Status, " Work time: ", end.Sub(start))
		logrus.Println()

		return err
	}
}

func (mw *MiddlewareManager) isAllowedOrigin(origin string) bool {
	for _, allowed := range mw.origins {
		if string(allowed) == origin {
			return true
		}
	}
	return false
}

func (mw *MiddlewareManager) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		origin := cntx.Request().Header.Get(echo.HeaderOrigin)

		allowOrigin := ""
		if mw.isAllowedOrigin(origin) {
			allowOrigin = origin
		}

		res := cntx.Response()
		res.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		res.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		res.Header().Set("Access-Control-Allow-Credentials", "true")

		if cntx.Request().Method == http.MethodOptions {
			return cntx.NoContent(http.StatusNoContent)
		}
		return next(cntx)
	}
}

func (mw *MiddlewareManager) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		cookie, err := cntx.Cookie(SessionName)
		if err != nil {
			customErr := errors.New(CodeUserUnauthorized, err)
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		sess, customErr := mw.sessUcase.Check(cookie.Value)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		cntx.Set("sessValue", sess.Value)
		cntx.Set("userID", sess.UserID)
		return next(cntx)
	}
}

func (mw *MiddlewareManager) GetAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		cookie, err := cntx.Cookie(SessionName)
		if err != nil {
			return next(cntx)
		}

		sess, customErr := mw.sessUcase.Check(cookie.Value)
		if customErr != nil {
			return next(cntx)
		}

		cntx.Set("sessValue", sess.Value)
		cntx.Set("userID", sess.UserID)
		return next(cntx)
	}
}

func (mw *MiddlewareManager) CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		userID := cntx.Get("userID").(uint64)
		isAdmin, customErr := mw.userUcase.IsAdmin(userID)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		if !isAdmin {
			customErr := errors.Get(CodeAccessDenied)
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return next(cntx)
	}
}
