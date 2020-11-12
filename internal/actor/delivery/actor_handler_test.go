package delivery

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
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

func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func TestActorHandler_CreateActorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorUseCase := mocks.NewMockActorUseCase(ctrl)

	actor := &models.Actor{
		Name: "Margo Robbie",
	}
	actorJSON, err := AnyToBytesBuffer(actor)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/actors", strings.NewReader(actorJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	actorHandler := NewActorHandler(actorUseCase)
	handleFunc := actorHandler.CreateActorHandler()
	actorHandler.Configure(e, nil)

	actorUseCase.
		EXPECT().
		Create(actor).
		Return(nil)

	response := &response.Response{Body: &response.Body{"actor": actor}}

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

func TestActorHandler_ChangeActorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorUseCase := mocks.NewMockActorUseCase(ctrl)

	actor := &models.Actor{
		Name: "Margo Robbie",
	}
	actorJSON, err := AnyToBytesBuffer(actor)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(actor.ID))
	req := httptest.NewRequest(http.MethodPut, "/api/v1/actors/"+strId,
		strings.NewReader(actorJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	actorHandler := NewActorHandler(actorUseCase)
	handleFunc := actorHandler.ChangeActorHandler()
	actorHandler.Configure(e, nil)

	actorUseCase.
		EXPECT().
		Change(actor).
		Return(nil)

	response := &response.Response{Body: &response.Body{"actor": actor}}

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

func TestActorHandler_GetActorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorUseCase := mocks.NewMockActorUseCase(ctrl)

	actor := &models.Actor{
		ID:   3,
		Name: "Margo Robbie",
	}

	e := echo.New()
	strId := strconv.Itoa(int(actor.ID))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/actors/"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	actorHandler := NewActorHandler(actorUseCase)
	handleFunc := actorHandler.GetActorHandler()
	e.PUT("/api/v1/actors/:id", handleFunc)

	actorUseCase.
		EXPECT().
		Get(actor.ID).
		Return(actor, nil)

	response := &response.Response{Body: &response.Body{"actor": actor}}

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

func TestActorHandler_GetActorHandler_NoActor(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorUseCase := mocks.NewMockActorUseCase(ctrl)

	e := echo.New()
	var id uint64 = 3
	strId := strconv.Itoa(int(id))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/actors/"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	actorHandler := NewActorHandler(actorUseCase)
	handleFunc := actorHandler.GetActorHandler()
	actorHandler.Configure(e, nil)

	actorUseCase.
		EXPECT().
		Get(id).
		Return(nil, errors.Get(consts.CodeActorDoesNotExist))

	response := &response.Response{Error: errors.Get(consts.CodeActorDoesNotExist)}

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

func TestActorHandler_DeleteActorHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	actorUseCase := mocks.NewMockActorUseCase(ctrl)

	actor := &models.Actor{
		Name: "Margo Robbie",
	}
	actorJSON, err := AnyToBytesBuffer(actor)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(actor.ID))
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/actors/"+strId,
		strings.NewReader(actorJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	actorHandler := NewActorHandler(actorUseCase)
	handleFunc := actorHandler.DeleteActorHandler()
	actorHandler.Configure(e, nil)

	actorUseCase.
		EXPECT().
		DeleteById(actor.ID).
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
