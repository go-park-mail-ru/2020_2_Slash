package delivery

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func TestDirectorHandler_CreateDirectorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorUseCase := mocks.NewMockDirectorUseCase(ctrl)

	director := &models.Director{
		Name: "Mike Nichols",
	}
	directorJSON, err := AnyToBytesBuffer(director)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/directors", strings.NewReader(directorJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	directorHandler := NewDirectorHandler(directorUseCase)
	handleFunc := directorHandler.CreateDirectorHandler()
	directorHandler.Configure(e, nil)

	directorUseCase.
		EXPECT().
		Create(director).
		Return(nil)

	response := &response.Response{Body: &response.Body{"director": director}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		expResBody, err := AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestDirectorHandler_ChangeDirectorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorUseCase := mocks.NewMockDirectorUseCase(ctrl)

	director := &models.Director{
		Name: "Mike Nichols",
	}
	directorJSON, err := AnyToBytesBuffer(director)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(director.ID))
	req := httptest.NewRequest(http.MethodPut, "/api/v1/directors/"+strId,
		strings.NewReader(directorJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	directorHandler := NewDirectorHandler(directorUseCase)
	handleFunc := directorHandler.ChangeDirectorHandler()
	directorHandler.Configure(e, nil)

	directorUseCase.
		EXPECT().
		Change(director).
		Return(nil)

	response := &response.Response{Body: &response.Body{"director": director}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestDirectorHandler_GetDirectorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorUseCase := mocks.NewMockDirectorUseCase(ctrl)

	director := &models.Director{
		ID:   3,
		Name: "Mike Nichols",
	}

	e := echo.New()
	strId := strconv.Itoa(int(director.ID))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/directors/"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	directorHandler := NewDirectorHandler(directorUseCase)
	handleFunc := directorHandler.GetDirectorHandler()
	e.PUT("/api/v1/directors/:id", handleFunc)

	directorUseCase.
		EXPECT().
		Get(director.ID).
		Return(director, nil)

	response := &response.Response{Body: &response.Body{"director": director}}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestDirectorHandler_GetDirectorHandler_NoDirector(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorUseCase := mocks.NewMockDirectorUseCase(ctrl)

	e := echo.New()
	var id uint64 = 3
	strId := strconv.Itoa(int(id))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/directors/"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	directorHandler := NewDirectorHandler(directorUseCase)
	handleFunc := directorHandler.GetDirectorHandler()
	directorHandler.Configure(e, nil)

	directorUseCase.
		EXPECT().
		Get(id).
		Return(nil, errors.Get(consts.CodeDirectorDoesNotExist))

	response := &response.Response{Error: errors.Get(consts.CodeDirectorDoesNotExist)}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		expResBody, err := AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestDirectorHandler_DeleteDirectorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorUseCase := mocks.NewMockDirectorUseCase(ctrl)

	director := &models.Director{
		Name: "Mike Nichols",
	}
	directorJSON, err := AnyToBytesBuffer(director)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(director.ID))
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/directors/"+strId,
		strings.NewReader(directorJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	directorHandler := NewDirectorHandler(directorUseCase)
	handleFunc := directorHandler.DeleteDirectorHandler()
	directorHandler.Configure(e, nil)

	directorUseCase.
		EXPECT().
		DeleteById(director.ID).
		Return(nil)

	response := &response.Response{Message: "success"}

	// Assertions
	if assert.NoError(t, handleFunc(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expResBody, err := AnyToBytesBuffer(response)
		if err != nil {
			t.Error(err)
			return
		}
		bytes, _ := ioutil.ReadAll(rec.Body)

		assert.JSONEq(t, expResBody.String(), string(bytes))
	}
}

func TestDirectorHandler_GetDirectorsListHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	directorUseCase := mocks.NewMockDirectorUseCase(ctrl)

	directors := []*models.Director{
		&models.Director{
			ID:   1,
			Name: "Mike Nichols",
		},
		&models.Director{
			ID:   2,
			Name: "No Mike Nichols",
		},
	}

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}

	reqJSON, err := converter.AnyToBytesBuffer(pgnt)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/directors/", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	directorHandler := NewDirectorHandler(directorUseCase)
	handleFunc := directorHandler.GetDirectorsListHandler()
	directorHandler.Configure(e, nil)

	directorUseCase.
		EXPECT().
		List(pgnt).
		Return(directors, nil)

	response := &response.Response{Body: &response.Body{"directors": directors}}

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
