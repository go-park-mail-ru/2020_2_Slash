package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	e.POST("/api/v1/genres", gh.CreateGenreHandler())
	e.PUT("/api/v1/genres/:gid", gh.UpdateGenreHandler())
	e.DELETE("/api/v1/genres/:gid", gh.DeleteGenreHandler())
	e.GET("/api/v1/genres", gh.GetGenresListHandler())
}

func (gh *GenreHandler) CreateGenreHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,lte=64"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		genre := &models.Genre{
			Name: req.Name,
		}

		if err := gh.genreUcase.Create(genre); err != nil {
			logrus.Info(err.Message)
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
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		genreData := &models.Genre{
			Name: req.Name,
		}

		genreID, _ := strconv.ParseUint(cntx.Param("gid"), 10, 64)
		genre, err := gh.genreUcase.UpdateByID(genreID, genreData)
		if err != nil {
			logrus.Info(err.Message)
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
		genreID, _ := strconv.ParseUint(cntx.Param("gid"), 10, 64)

		if err := gh.genreUcase.DeleteByID(genreID); err != nil {
			logrus.Info(err.Message)
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
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"genres": genres,
			},
		})
	}
}
