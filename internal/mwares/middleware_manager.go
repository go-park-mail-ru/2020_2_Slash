package mwares

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares/monitoring"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/session"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/user"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/CSRFManager"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
)

type MiddlewareManager struct {
	sessUcase session.SessionUsecase
	userUcase user.UserUsecase
	mntng     *monitoring.Monitoring
	origins   []string
}

func NewMiddlewareManager(sessUcase session.SessionUsecase,
	userUcase user.UserUsecase, mntng *monitoring.Monitoring) *MiddlewareManager {
	return &MiddlewareManager{
		sessUcase: sessUcase,
		userUcase: userUcase,
		mntng:     mntng,
		origins:   []string{"https://www.flicksbox.ru", "http://www.flicksbox.ru:3000"},
	}
}

func (m *MiddlewareManager) PanicRecovering(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		// nolint: errcheck
		defer func() error {
			if err := recover(); err != nil {
				status := strconv.Itoa(cntx.Response().Status)
				path := cntx.Request().URL.Path
				method := cntx.Request().Method

				re := regexp.MustCompile(`/\d+`)
				replacedPath := re.ReplaceAllString(path, "/*")

				m.mntng.Hits.WithLabelValues(status, replacedPath, method).Inc()
				m.mntng.Duration.WithLabelValues(status, replacedPath, method).Observe(0)

				logger.Warn(err)
				customErr := errors.Get(CodeInternalError)
				return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
			}
			return nil
		}()
		return next(cntx)
	}
}

func (m *MiddlewareManager) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		if cntx.Request().URL.String() == "/metrics" {
			return next(cntx)
		}

		logger.Info(cntx.Request().RemoteAddr, " ", cntx.Request().Method, " ", cntx.Request().URL)

		start := time.Now()
		err := next(cntx)
		end := time.Now()
		workTime := end.Sub(start)

		status := strconv.Itoa(cntx.Response().Status)
		path := cntx.Request().URL.Path
		method := cntx.Request().Method

		re := regexp.MustCompile(`/\d+`)
		replacedPath := re.ReplaceAllString(path, "/*")

		m.mntng.Hits.WithLabelValues(status, replacedPath, method).Inc()
		m.mntng.Duration.WithLabelValues(status, replacedPath, method).Observe(workTime.Seconds())

		logger.Info("Status: ", cntx.Response().Status, " Work time: ", workTime)
		logger.Println()

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
		res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Csrf-Token")
		res.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
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
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		sess, customErr := mw.sessUcase.Check(cookie.Value)
		if customErr != nil {
			logger.Error(customErr.Message)
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
		userID, ok := cntx.Get("userID").(uint64)
		if !ok {
			customErr := errors.Get(CodeGetFromContextError)
			logger.Error(customErr)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		isAdmin, customErr := mw.userUcase.IsAdmin(userID)
		if customErr != nil {
			logger.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		if !isAdmin {
			customErr := errors.Get(CodeAccessDenied)
			logger.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return next(cntx)
	}
}

func (mw *MiddlewareManager) CheckCSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(cntx echo.Context) error {
		switch cntx.Request().Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			return next(cntx)
		}

		token := cntx.Request().Header.Get("X-Csrf-Token")
		if token == "" {
			customErr := errors.Get(CodeCSRFTokenWasNotPassed)
			logger.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		cookie, err := cntx.Cookie(SessionName)
		if err != nil {
			logger.Info(err)
			customErr := errors.New(CodeUserUnauthorized, err)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		sess, customErr := mw.sessUcase.Get(cookie.Value)
		if customErr != nil {
			logger.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		customErr = CSRFManager.ValidateToken(sess, token)
		if customErr != nil {
			logger.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return next(cntx)
	}
}
