package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	e.POST("/api/v1/countries", ch.CreateCountryHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.PUT("/api/v1/countries/:cid", ch.UpdateCountryHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.DELETE("/api/v1/countries/:cid", ch.DeleteCountryHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.GET("/api/v1/countries", ch.GetCountriesListHandler(), mw.CheckAuth, mw.CheckAdmin)
}

func (ch *CountryHandler) CreateCountryHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required,lte=64"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		country := &models.Country{
			Name: req.Name,
		}

		if err := ch.countryUcase.Create(country); err != nil {
			logrus.Info(err.Message)
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
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		countryData := &models.Country{
			Name: req.Name,
		}

		countryID, _ := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		country, err := ch.countryUcase.UpdateByID(countryID, countryData)
		if err != nil {
			logrus.Info(err.Message)
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
		countryID, _ := strconv.ParseUint(cntx.Param("cid"), 10, 64)

		if err := ch.countryUcase.DeleteByID(countryID); err != nil {
			logrus.Info(err.Message)
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
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"countries": countries,
			},
		})
	}
}
