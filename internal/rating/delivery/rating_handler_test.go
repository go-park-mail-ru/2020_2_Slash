package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/rating"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/rating/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type Request struct {
	Likes bool `json:"likes"`
}

func TestRatingHandler_CreateHandler_Success(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)

	likes := Request{Likes: false}

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     false,
	}

	ratingJSON, err := converter.AnyToBytesBuffer(likes)
	if err != nil {
		t.Fatal(err)
	}

	ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase,
		http.MethodPost, strings.NewReader(ratingJSON.String()))
	handleFunc := ratingHandler.CreateHandler()

	ratingUseCase.
		EXPECT().
		Create(rating).
		Return(nil)

	response := &response.Response{Body: &response.Body{"rating": rating}}

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

func TestRatingHandler_CreateHandler_Error(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)
	logger.DisableLogger()

	likes := Request{Likes: false}

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     false,
	}

	ratingJSON, err := converter.AnyToBytesBuffer(likes)
	if err != nil {
		t.Fatal(err)
	}

	errors := []*errors.Error{
		errors.Get(consts.CodeRatingAlreadyExist),
		errors.Get(consts.CodeContentDoesNotExist),
	}

	for _, customErr := range errors {
		ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase,
			http.MethodPost, strings.NewReader(ratingJSON.String()))
		handleFunc := ratingHandler.CreateHandler()

		ratingUseCase.
			EXPECT().
			Create(rating).
			Return(customErr)

		response := &response.Response{Error: customErr}
		// Assertions
		if assert.NoError(t, handleFunc(c)) {
			assert.Equal(t, customErr.HTTPCode, rec.Code)

			expResBody, err := converter.AnyToBytesBuffer(response)
			if err != nil {
				t.Error(err)
				return
			}
			bytes, _ := ioutil.ReadAll(rec.Body)

			assert.JSONEq(t, expResBody.String(), string(bytes))
		}
	}
}

func TestRatingHandler_ChangeHandler_Success(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)

	likes := Request{Likes: false}

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingJSON, err := converter.AnyToBytesBuffer(likes)
	if err != nil {
		t.Fatal(err)
	}
	ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase,
		http.MethodPut, strings.NewReader(ratingJSON.String()))
	handleFunc := ratingHandler.ChangeHandler()

	ratingUseCase.
		EXPECT().
		Change(rating).
		Return(nil)

	rating.Likes = false
	response := &response.Response{Body: &response.Body{"rating": rating}}

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

func TestRatingHandler_ChangeHandler_Error(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)
	logger.DisableLogger()

	likes := Request{Likes: true}

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     true,
	}

	ratingJSON, err := converter.AnyToBytesBuffer(likes)
	if err != nil {
		t.Fatal(err)
	}

	errors := []*errors.Error{
		errors.Get(consts.CodeRatingAlreadyExist),
		errors.Get(consts.CodeRatingDoesNotExist),
		errors.Get(consts.CodeContentDoesNotExist),
	}

	for _, customErr := range errors {
		ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase,
			http.MethodPut, strings.NewReader(ratingJSON.String()))
		handleFunc := ratingHandler.ChangeHandler()

		ratingUseCase.
			EXPECT().
			Change(rating).
			Return(customErr)

		response := &response.Response{Error: customErr}
		// Assertions
		if assert.NoError(t, handleFunc(c)) {
			assert.Equal(t, customErr.HTTPCode, rec.Code)

			expResBody, err := converter.AnyToBytesBuffer(response)
			if err != nil {
				t.Error(err)
				return
			}
			bytes, _ := ioutil.ReadAll(rec.Body)

			assert.JSONEq(t, expResBody.String(), string(bytes))
		}
	}
}

func TestRatingHandler_DeleteHandler_Success(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     false,
	}

	ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase, http.MethodDelete, strings.NewReader(""))
	handleFunc := ratingHandler.DeleteHandler()

	ratingUseCase.
		EXPECT().
		Delete(rating).
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

func TestRatingHandler_DeleteHandler_ContentDoesNotExist(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)
	logger.DisableLogger()

	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     false,
	}

	errors := []*errors.Error{
		errors.Get(consts.CodeRatingDoesNotExist),
		errors.Get(consts.CodeContentDoesNotExist),
	}

	for _, customErr := range errors {
		ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase,
			http.MethodDelete, strings.NewReader(""))
		handleFunc := ratingHandler.DeleteHandler()

		ratingUseCase.
			EXPECT().
			Delete(rating).
			Return(customErr)

		response := &response.Response{Error: customErr}
		// Assertions
		if assert.NoError(t, handleFunc(c)) {
			assert.Equal(t, customErr.HTTPCode, rec.Code)

			expResBody, err := converter.AnyToBytesBuffer(response)
			if err != nil {
				t.Error(err)
				return
			}
			bytes, _ := ioutil.ReadAll(rec.Body)

			assert.JSONEq(t, expResBody.String(), string(bytes))
		}
	}
}

func setupRatingHandler(rating *models.Rating,
	ratingUseCase rating.RatingUsecase,
	httpMethod string,
	stringifiedJSON io.Reader) (*RatingHandler, echo.Context,
	*httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(httpMethod, "/api/v1/rating",
		stringifiedJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", rating.UserID)
	c.SetParamNames("cid")
	c.SetParamValues(strconv.Itoa(int(rating.UserID)))
	ratingHandler := NewRatingHandler(ratingUseCase)
	ratingHandler.Configure(e, nil)
	return ratingHandler, c, rec
}

func TestRatingHandler_GetContentRatingHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)
	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     false,
	}
	percentage := 75

	ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase,
		http.MethodGet, strings.NewReader(""))
	handleFunc := ratingHandler.GetContentRatingHandler()

	ratingUseCase.
		EXPECT().
		GetContentRating(rating.ContentID).
		Return(percentage, nil)

	response := &response.Response{Body: &response.Body{"match": percentage}}

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

func TestRatingHandler_GetRatingHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ratingUseCase := mocks.NewMockRatingUsecase(ctrl)
	rating := &models.Rating{
		UserID:    3,
		ContentID: 3,
		Likes:     false,
	}

	ratingHandler, c, rec := setupRatingHandler(rating, ratingUseCase,
		http.MethodGet, strings.NewReader(""))
	handleFunc := ratingHandler.GetHandler()

	ratingUseCase.
		EXPECT().
		GetByUserIDContentID(rating.UserID, rating.ContentID).
		Return(rating, nil)

	response := &response.Response{Body: &response.Body{"rating": rating}}

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