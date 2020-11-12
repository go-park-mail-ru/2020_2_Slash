package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
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

func TestGenreHandler_CreateGenreHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreUseCase := mocks.NewMockGenreUsecase(ctrl)

	genre := &models.Genre{
		Name: "USA",
	}
	genreJSON, err := converter.AnyToBytesBuffer(genre)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/genres", strings.NewReader(genreJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	genreHandler := NewGenreHandler(genreUseCase)
	handleFunc := genreHandler.CreateGenreHandler()
	genreHandler.Configure(e, nil)

	genreUseCase.
		EXPECT().
		Create(genre).
		Return(nil)

	response := &response.Response{Body: &response.Body{"genre": genre}}

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

func TestGenreHandler_CreateGenreHandler_NameAlreadyExists(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreUseCase := mocks.NewMockGenreUsecase(ctrl)

	genre := &models.Genre{
		Name: "USA",
	}
	genreJSON, err := converter.AnyToBytesBuffer(genre)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/genres", strings.NewReader(genreJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	genreHandler := NewGenreHandler(genreUseCase)
	handleFunc := genreHandler.CreateGenreHandler()
	genreHandler.Configure(e, nil)

	genreUseCase.
		EXPECT().
		Create(genre).
		Return(errors.Get(consts.CodeGenreNameAlreadyExists))

	response := &response.Response{Error: errors.Get(consts.CodeGenreNameAlreadyExists)}

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

func TestGenreHandler_UpdateGenreHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreUseCase := mocks.NewMockGenreUsecase(ctrl)

	genre := &models.Genre{
		Name: "USA",
	}

	newGenreData := &models.Genre{
		Name: "GB",
	}

	genreJSON, err := converter.AnyToBytesBuffer(newGenreData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(genre.ID))
	req := httptest.NewRequest(http.MethodPut, "/api/v1/genres/"+strId,
		strings.NewReader(genreJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	genreHandler := NewGenreHandler(genreUseCase)
	handleFunc := genreHandler.UpdateGenreHandler()
	genreHandler.Configure(e, nil)

	genreUseCase.
		EXPECT().
		UpdateByID(genre.ID, newGenreData).
		Return(newGenreData, nil)

	response := &response.Response{Body: &response.Body{"genre": newGenreData}}

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

func TestGenreHandler_UpdateGenreHandler_NameAlreadyExists(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreUseCase := mocks.NewMockGenreUsecase(ctrl)

	genre := &models.Genre{
		Name: "USA",
	}

	newGenreData := &models.Genre{
		Name: "GB",
	}

	genreJSON, err := converter.AnyToBytesBuffer(newGenreData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(genre.ID))
	req := httptest.NewRequest(http.MethodPut, "/api/v1/genres/"+strId,
		strings.NewReader(genreJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	genreHandler := NewGenreHandler(genreUseCase)
	handleFunc := genreHandler.UpdateGenreHandler()
	genreHandler.Configure(e, nil)

	genreUseCase.
		EXPECT().
		UpdateByID(genre.ID, newGenreData).
		Return(nil, errors.Get(consts.CodeGenreNameAlreadyExists))

	response := &response.Response{Error: errors.Get(consts.CodeGenreNameAlreadyExists)}

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

func TestGenreHandler_DeleteGenreHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreUseCase := mocks.NewMockGenreUsecase(ctrl)

	genre := &models.Genre{
		Name: "USA",
	}

	e := echo.New()
	strId := strconv.Itoa(int(genre.ID))
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/genres/"+strId,
		strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	genreHandler := NewGenreHandler(genreUseCase)
	handleFunc := genreHandler.DeleteGenreHandler()
	genreHandler.Configure(e, nil)

	genreUseCase.
		EXPECT().
		DeleteByID(genre.ID).
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

func TestGenreHandler_DeleteGenreHandler_NoGenre(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreUseCase := mocks.NewMockGenreUsecase(ctrl)

	var genreID uint64 = 0

	e := echo.New()
	strId := strconv.Itoa(int(genreID))
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/genres/"+strId,
		strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	genreHandler := NewGenreHandler(genreUseCase)
	handleFunc := genreHandler.DeleteGenreHandler()
	genreHandler.Configure(e, nil)

	genreUseCase.
		EXPECT().
		DeleteByID(genreID).
		Return(errors.Get(consts.CodeGenreDoesNotExist))

	response := &response.Response{Error: errors.Get(consts.CodeGenreDoesNotExist)}

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

func TestGenreHandler_GetGenresListHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	genreUseCase := mocks.NewMockGenreUsecase(ctrl)

	genres := []*models.Genre{
		&models.Genre{
			ID:   1,
			Name: "USA",
		},
		&models.Genre{
			ID:   2,
			Name: "GB",
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/genres/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	genreHandler := NewGenreHandler(genreUseCase)
	handleFunc := genreHandler.GetGenresListHandler()
	genreHandler.Configure(e, nil)

	genreUseCase.
		EXPECT().
		List().
		Return(genres, nil)

	response := &response.Response{Body: &response.Body{"genres": genres}}

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
