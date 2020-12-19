package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
)

type GenreHandler struct {
	genreUcase genre.GenreUsecase
}

func NewGenreHandler(genreUcase genre.GenreUsecase) *GenreHandler {
	return &GenreHandler{
		genreUcase: genreUcase,
	}
}

func (gh *GenreHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/genres", gh.CreateGenreHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.PUT("/api/v1/genres/:gid", gh.UpdateGenreHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.DELETE("/api/v1/genres/:gid", gh.DeleteGenreHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.GET("/api/v1/genres", gh.GetGenresListHandler())
}

func (gh *GenreHandler) CreateGenreHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,lte=64"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		genre := &models.Genre{
			Name: req.Name,
		}

		if err := gh.genreUcase.Create(genre); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"genre": genre,
			},
		})
	}
}

func (gh *GenreHandler) UpdateGenreHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,lte=64"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		genreData := &models.Genre{
			Name: req.Name,
		}

		genreID, parseErr := strconv.ParseUint(cntx.Param("gid"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(consts.CodeInternalError, parseErr)
			logger.Error(customErr)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}
		genre, err := gh.genreUcase.UpdateByID(genreID, genreData)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"genre": genre,
			},
		})
	}
}

func (gh *GenreHandler) DeleteGenreHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		genreID, parseErr := strconv.ParseUint(cntx.Param("gid"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(consts.CodeInternalError, parseErr)
			logger.Error(customErr)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		if err := gh.genreUcase.DeleteByID(genreID); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (gh *GenreHandler) GetGenresListHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		genres, err := gh.genreUcase.List()
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"genres": genres,
			},
		})
	}
}
