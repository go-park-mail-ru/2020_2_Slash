package delivery

import (
	"bytes"
	"encoding/json"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestFavouriteHandler_Create_Success(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteUseCase := mocks.NewMockFavouriteUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)

	var userID uint64 = 3
	var contentID uint64 = 4
	favourite := &models.Favourite{
		UserID:    userID,
		ContentID: contentID,
	}
	type Request struct {
		ContentID uint64 `json:"content_id"`
	}

	favouriteJSON, err := AnyToBytesBuffer(Request{ContentID: contentID})
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/favourites", strings.NewReader(favouriteJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", userID)
	favouriteHandler := NewFavouriteHandler(favouriteUseCase, contentUseCase)
	handleFunc := favouriteHandler.CreateHandler()
	favouriteHandler.Configure(e, nil)

	contentUseCase.
		EXPECT().
		GetByID(contentID).
		Return(&models.Content{
			ContentID: contentID,
		}, nil)
	favouriteUseCase.
		EXPECT().
		Create(gomock.Any()).
		Return(nil)

	response := &response.Response{Body: &response.Body{"favourite": favourite}}

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

func TestFavouriteHandler_DeleteHandler_Success(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteUseCase := mocks.NewMockFavouriteUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)

	var userID uint64 = 3
	var contentID uint64 = 4
	favourite := &models.Favourite{
		UserID:    userID,
		ContentID: contentID,
	}
	type Request struct {
		ContentID uint64 `json:"content_id"`
	}

	favouriteJSON, err := AnyToBytesBuffer(Request{ContentID: contentID})
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/favourites", strings.NewReader(favouriteJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", userID)
	favouriteHandler := NewFavouriteHandler(favouriteUseCase, contentUseCase)
	handleFunc := favouriteHandler.DeleteHandler()
	favouriteHandler.Configure(e, nil)

	favouriteUseCase.
		EXPECT().
		Delete(favourite).
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

func TestFavouriteHandler_GetFavouritesHandler_Success(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	favouriteUseCase := mocks.NewMockFavouriteUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)

	var userID uint64 = 3
	expectReturn := []*models.Content{
		{
			ContentID: 2,
		},
		{
			ContentID: 4,
		},
		{
			ContentID: 421,
		},
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/favourites", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", userID)
	favouriteHandler := NewFavouriteHandler(favouriteUseCase, contentUseCase)
	handleFunc := favouriteHandler.GetFavouritesHandler()
	favouriteHandler.Configure(e, nil)

	favouriteUseCase.
		EXPECT().
		GetUserFavourites(userID).
		Return(expectReturn, nil)

	response := &response.Response{Body: &response.Body{"favourites": expectReturn}}

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
