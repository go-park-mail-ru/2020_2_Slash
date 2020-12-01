package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCountryHandler_CreateCountryHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryUseCase := mocks.NewMockCountryUsecase(ctrl)

	country := &models.Country{
		Name: "USA",
	}
	countryJSON, err := converter.AnyToBytesBuffer(country)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/countries", strings.NewReader(countryJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	countryHandler := NewCountryHandler(countryUseCase)
	handleFunc := countryHandler.CreateCountryHandler()
	countryHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		Create(country).
		Return(nil)

	response := &response.Response{Body: &response.Body{"country": country}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestCountryHandler_CreateCountryHandler_NameAlreadyExists(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryUseCase := mocks.NewMockCountryUsecase(ctrl)

	country := &models.Country{
		Name: "USA",
	}
	countryJSON, err := converter.AnyToBytesBuffer(country)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/countries", strings.NewReader(countryJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	countryHandler := NewCountryHandler(countryUseCase)
	handleFunc := countryHandler.CreateCountryHandler()
	countryHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		Create(country).
		Return(errors.Get(consts.CodeCountryNameAlreadyExists))

	response := &response.Response{Error: errors.Get(consts.CodeCountryNameAlreadyExists)}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestCountryHandler_UpdateCountryHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryUseCase := mocks.NewMockCountryUsecase(ctrl)

	country := &models.Country{
		Name: "USA",
	}

	newCountryData := &models.Country{
		ID:   country.ID,
		Name: "GB",
	}

	countryJSON, err := converter.AnyToBytesBuffer(newCountryData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(country.ID))
	req := httptest.NewRequest(http.MethodPut, "/api/v1/countries/"+strId,
		strings.NewReader(countryJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	countryHandler := NewCountryHandler(countryUseCase)
	handleFunc := countryHandler.UpdateCountryHandler()
	countryHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		Update(newCountryData).
		Return(nil)

	response := &response.Response{Body: &response.Body{"country": newCountryData}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestCountryHandler_UpdateCountryHandler_NameAlreadyExists(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryUseCase := mocks.NewMockCountryUsecase(ctrl)

	country := &models.Country{
		Name: "USA",
	}

	newCountryData := &models.Country{
		ID:   country.ID,
		Name: "GB",
	}

	countryJSON, err := converter.AnyToBytesBuffer(newCountryData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(country.ID))
	req := httptest.NewRequest(http.MethodPut, "/api/v1/countries/"+strId,
		strings.NewReader(countryJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	countryHandler := NewCountryHandler(countryUseCase)
	handleFunc := countryHandler.UpdateCountryHandler()
	countryHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		Update(newCountryData).
		Return(errors.Get(consts.CodeCountryNameAlreadyExists))

	response := &response.Response{Error: errors.Get(consts.CodeCountryNameAlreadyExists)}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestCountryHandler_DeleteCountryHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryUseCase := mocks.NewMockCountryUsecase(ctrl)

	country := &models.Country{
		Name: "USA",
	}

	e := echo.New()
	strId := strconv.Itoa(int(country.ID))
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/countries/"+strId,
		strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	countryHandler := NewCountryHandler(countryUseCase)
	handleFunc := countryHandler.DeleteCountryHandler()
	countryHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		DeleteByID(country.ID).
		Return(nil)

	response := &response.Response{Message: "success"}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestCountryHandler_DeleteCountryHandler_NoCountry(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryUseCase := mocks.NewMockCountryUsecase(ctrl)

	var countryID uint64 = 0

	e := echo.New()
	strId := strconv.Itoa(int(countryID))
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/countries/"+strId,
		strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	countryHandler := NewCountryHandler(countryUseCase)
	handleFunc := countryHandler.DeleteCountryHandler()
	countryHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		DeleteByID(countryID).
		Return(errors.Get(consts.CodeCountryDoesNotExist))

	response := &response.Response{Error: errors.Get(consts.CodeCountryDoesNotExist)}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestCountryHandler_GetCountriesListHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	countryUseCase := mocks.NewMockCountryUsecase(ctrl)

	countries := []*models.Country{
		&models.Country{
			ID:   1,
			Name: "USA",
		},
		&models.Country{
			ID:   2,
			Name: "GB",
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/countries/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	countryHandler := NewCountryHandler(countryUseCase)
	handleFunc := countryHandler.GetCountriesListHandler()
	countryHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		List().
		Return(countries, nil)

	response := &response.Response{Body: &response.Body{"countries": countries}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := converter.AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}
