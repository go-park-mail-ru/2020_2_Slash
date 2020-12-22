package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
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
	e.POST("/api/v1/directors", dh.CreateDirectorHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.PUT("/api/v1/directors/:id", dh.ChangeDirectorHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.GET("/api/v1/directors/:id", dh.GetDirectorHandler())
	e.DELETE("/api/v1/directors/:id", dh.DeleteDirectorHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.GET("/api/v1/directors", dh.GetDirectorsListHandler())
}

func (dh *DirectorHandler) CreateDirectorHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		director := &models.Director{
			Name: req.Name,
		}
		err := dh.directorUseCase.Create(director)
		if err != nil {
			logger.Error(err.Message)
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
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		req := &Request{}
		if customErr := reader.NewRequestReader(cntx).Read(req); customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		director := &models.Director{
			ID:   id,
			Name: req.Name,
		}
		customErr := dh.directorUseCase.Change(director)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"director": director,
			},
		})
	}
}

func (dh *DirectorHandler) GetDirectorHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		director, customErr := dh.directorUseCase.Get(id)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"director": director,
			},
		})
	}
}

func (dh *DirectorHandler) DeleteDirectorHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		customErr := dh.directorUseCase.DeleteById(id)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (dh *DirectorHandler) GetDirectorsListHandler() echo.HandlerFunc {
	type Request struct {
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if customErr := reader.NewRequestReader(cntx).Read(req); customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		directors, err := dh.directorUseCase.List(&req.Pagination)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"directors": directors,
			},
		})
	}
}
