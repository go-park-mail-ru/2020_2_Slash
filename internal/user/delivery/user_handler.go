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

type UserHandler struct {
	userUcase user.UserUsecase
	sessUcase session.SessionUsecase
}

func NewUserHandler(userUcase user.UserUsecase, sessUcase session.SessionUsecase) *UserHandler {
	return &UserHandler{
		userUcase: userUcase,
		sessUcase: sessUcase,
	}
}

func (uh *UserHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/user/register", uh.registerUserHandler())
	e.GET("/api/v1/user/profile", uh.getUserProfileHandler(), mw.CheckAuth)
	e.PUT("/api/v1/user/profile", uh.updateUserProfileHandler(), mw.CheckAuth)
}

func (uh *UserHandler) registerUserHandler() echo.HandlerFunc {
	type Request struct {
		Nickname         string `json:"nickname"`
		Email            string `json:"email" validate:"required,email"`
		Password         string `json:"password" validate:"required,gte=6"`
		RepeatedPassword string `json:"repeated_password" validate:"eqfield=Password"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		user := &models.User{
			Nickname: req.Nickname,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := uh.userUcase.Create(user); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		sess := models.NewSession(user.ID)
		if err := uh.sessUcase.Create(sess); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		cookie := tools.CreateCookie(sess)
		cntx.SetCookie(cookie)
		return cntx.JSON(http.StatusOK, user.Sanitize())
	}
}

func (uh *UserHandler) getUserProfileHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		userID := cntx.Get("userID").(uint64)
		user, err := uh.userUcase.GetByID(userID)
		if err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}
		return cntx.JSON(http.StatusOK, user.Sanitize())
	}
}

func (uh *UserHandler) updateUserProfileHandler() echo.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"omitempty,gte=6"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}

		userData := &models.User{
			Nickname: req.Nickname,
			Email:    req.Email,
			Password: req.Password,
		}

		userID := cntx.Get("userID").(uint64)
		user, err := uh.userUcase.UpdateProfile(userID, userData)
		if err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, err)
		}
		return cntx.JSON(http.StatusOK, user.Sanitize())
	}
}
