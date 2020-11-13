package delivery

import (
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
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
	"github.com/labstack/echo/v4/middleware"
	"io"
	"net/http"
	"os"
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
	e.POST("/api/v1/user/register", uh.RegisterUserHandler())
	e.GET("/api/v1/user/profile", uh.GetUserProfileHandler(), mw.CheckAuth)
	e.PUT("/api/v1/user/profile", uh.UpdateUserProfileHandler(), mw.CheckAuth, mw.CheckCSRF)
	e.POST("/api/v1/user/avatar", uh.UpdateAvatarHandler(), mw.CheckAuth, middleware.BodyLimit("10M"), mw.CheckCSRF)
}

func (uh *UserHandler) RegisterUserHandler() echo.HandlerFunc {
	type Request struct {
		Nickname         string `json:"nickname" validate:"omitempty,gte=3,lte=32"`
		Email            string `json:"email" validate:"required,email,lte=64"`
		Password         string `json:"password" validate:"required,gte=6,lte=32"`
		RepeatedPassword string `json:"repeated_password" validate:"eqfield=Password"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).ReadUser(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		user := &models.User{
			Nickname: req.Nickname,
			Email:    req.Email,
			Password: req.Password,
			Role:     User,
		}

		if err := uh.userUcase.Create(user); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		sess := models.NewSession(user.ID)
		if err := uh.sessUcase.Create(sess); err != nil {
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
		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"user": user,
			},
		})
	}
}

func (uh *UserHandler) GetUserProfileHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		userID := cntx.Get("userID").(uint64)
		user, err := uh.userUcase.GetByID(userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"user": user,
			},
		})
	}
}

func (uh *UserHandler) UpdateUserProfileHandler() echo.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname" validate:"omitempty,gte=3,lte=32"`
		Email    string `json:"email" validate:"omitempty,email,lte=64"`
		Password string `json:"password" validate:"omitempty,gte=6,lte=32"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).ReadUser(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		userData := &models.User{
			Nickname: req.Nickname,
			Email:    req.Email,
			Password: req.Password,
			Role:     User,
		}

		userID := cntx.Get("userID").(uint64)
		user, err := uh.userUcase.UpdateProfile(userID, userData)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"user": user,
			},
		})
	}
}

func (uh *UserHandler) UpdateAvatarHandler() echo.HandlerFunc {
	const avatarsDir = "/avatars/"

	return func(cntx echo.Context) error {
		image, customErr := reader.NewRequestReader(cntx).ReadImage("avatar")
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		imageFile, err := image.Open()
		if err != nil {
			logger.Error(err)
			customErr := errors.New(CodeBadRequest, err)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}
		defer imageFile.Close()

		fileExtension, err := helpers.GetImageExtension(image)
		if err != nil {
			logger.Error(err)
			customErr := errors.New(CodeBadRequest, err)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		userID := cntx.Get("userID").(uint64)
		newAvatarFileName := helpers.GetUniqFileName(userID, fileExtension)
		rltNewAvatarFilePath := avatarsDir + newAvatarFileName
		absNewAvatarFilePath := "." + rltNewAvatarFilePath
		fileMode := int(0777)

		// Save image to storage
		newAvatarFile, err := os.OpenFile(absNewAvatarFilePath, os.O_WRONLY|os.O_CREATE, os.FileMode(fileMode))
		if err != nil {
			logger.Error(err)
			customErr := errors.New(CodeInternalError, err)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}
		defer newAvatarFile.Close()

		if _, err := io.Copy(newAvatarFile, imageFile); err != nil {
			_ = os.Remove(absNewAvatarFilePath)
			logger.Error(err)
			customErr := errors.New(CodeInternalError, err)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		_, customErr = uh.userUcase.UpdateAvatar(userID, rltNewAvatarFilePath)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"avatar": newAvatarFileName,
			},
		})
	}
}
