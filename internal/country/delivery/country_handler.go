package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
)

type CountryHandler struct {
	countryUcase country.CountryUsecase
}

func NewCountryHandler(countryUcase country.CountryUsecase) *CountryHandler {
	return &CountryHandler{
		countryUcase: countryUcase,
	}
}

func (ch *CountryHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/countries", ch.CreateCountryHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.PUT("/api/v1/countries/:cid", ch.UpdateCountryHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.DELETE("/api/v1/countries/:cid", ch.DeleteCountryHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.GET("/api/v1/countries", ch.GetCountriesListHandler())
}

func (ch *CountryHandler) CreateCountryHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,lte=64"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		country := &models.Country{
			Name: req.Name,
		}

		if err := ch.countryUcase.Create(country); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"country": country,
			},
		})
	}
}

func (ch *CountryHandler) UpdateCountryHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,lte=64"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		countryData := &models.Country{
			Name: req.Name,
		}

		countryID, parseErr := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(consts.CodeInternalError, parseErr)
			logger.Error(customErr)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		country, err := ch.countryUcase.UpdateByID(countryID, countryData)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"country": country,
			},
		})
	}
}

func (ch *CountryHandler) DeleteCountryHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		countryID, parseErr := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(consts.CodeInternalError, parseErr)
			logger.Error(customErr)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		if err := ch.countryUcase.DeleteByID(countryID); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (ch *CountryHandler) GetCountriesListHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		countries, err := ch.countryUcase.List()
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"countries": countries,
			},
		})
	}
}
