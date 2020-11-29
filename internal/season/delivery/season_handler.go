package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
)

type SeasonHandler struct {
	seasonUsecase season.SeasonUsecase
}

func NewSeasonHandler(usecase season.SeasonUsecase) *SeasonHandler {
	return &SeasonHandler{seasonUsecase: usecase}
}

func (sh *SeasonHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/seasons", sh.CreateHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.PUT("/api/v1/seasons/:id", sh.ChangeHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.DELETE("/api/v1/seasons/:id", sh.DeleteHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.GET("/api/v1/seasons/:id", sh.GetHandler(), mw.GetAuth)
}

func (sh *SeasonHandler) CreateHandler() echo.HandlerFunc {
	type Request struct {
		Number   int    `json:"number" validate:"required"`
		TVShowID uint64 `json:"tv_show_id" validate:"required"`
	}
	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		season := &models.Season{
			Number:   req.Number,
			TVShowID: req.TVShowID,
		}
		customErr := sh.seasonUsecase.Create(season)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"season": season,
			},
		})
	}
}

func (sh *SeasonHandler) ChangeHandler() echo.HandlerFunc {
	type Request struct {
		Number   int    `json:"number" validate:"required"`
		TVShowID uint64 `json:"tv_show_id" validate:"required"`
	}
	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		seasonID, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		season := &models.Season{
			ID:       seasonID,
			Number:   req.Number,
			TVShowID: req.TVShowID,
		}
		customErr := sh.seasonUsecase.Change(season)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"season": season,
			},
		})
	}
}

func (sh *SeasonHandler) DeleteHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		seasonID, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		customErr := sh.seasonUsecase.Delete(seasonID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{Message: "success"})
	}
}

func (sh *SeasonHandler) GetHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		seasonID, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		season, customErr := sh.seasonUsecase.Get(seasonID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		episodes, customErr := sh.seasonUsecase.GetEpisodes(seasonID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}
		season.Episodes = episodes

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"season": season,
			},
		})
	}
}
