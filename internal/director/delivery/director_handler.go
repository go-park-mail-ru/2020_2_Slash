package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type DirectorHandler struct {
	directorUseCase director.DirectorUseCase
}

func NewDirectorHandler(directorUseCase director.DirectorUseCase) *DirectorHandler {
	return &DirectorHandler{
		directorUseCase: directorUseCase,
	}
}

func (dh *DirectorHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/directors", dh.CreateDirectorHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.PUT("/api/v1/directors/:id", dh.ChangeDirectorHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.GET("/api/v1/directors/:id", dh.GetDirectorHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.DELETE("/api/v1/directors/:id", dh.DeleteDirectorHandler(), mw.CheckAuth, mw.CheckAdmin)
}

func (dh *DirectorHandler) CreateDirectorHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		director := &models.Director{
			Name: req.Name,
		}
		err := dh.directorUseCase.Create(director)
		if err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"director": director,
			},
		})
	}
}

func (dh *DirectorHandler) ChangeDirectorHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		req := &Request{}
		if customErr := reader.NewRequestReader(cntx).Read(req); customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		director := &models.Director{
			ID:   id,
			Name: req.Name,
		}
		customErr := dh.directorUseCase.Change(director)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"director": director,
			},
		})
	}
}

func (ah *DirectorHandler) GetDirectorHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		director, customErr := ah.directorUseCase.Get(id)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"director": director,
			},
		})
	}
}

func (ah *DirectorHandler) DeleteDirectorHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		customErr := ah.directorUseCase.DeleteById(id)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}
