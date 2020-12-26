package delivery

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	movieMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/mocks"
	tvshowMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestContentHandler_GetContentPostersHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 0

	params := &models.ContentFilter{
		Year:     []int{2001},
		Genre:    []int{1},
		Country:  []int{1},
		Actor:    []int{1},
		Director: []int{1},
	}

	type Request struct {
		models.ContentFilter
		models.Pagination
	}

	reqest := &Request{
		*params,
		*pgnt,
	}

	reqJSON, err := converter.AnyToBytesBuffer(reqest)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/content", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	contentHandler := NewContentHandler(contentUseCase, movieUseCase, tvshowUseCase)
	handleFunc := contentHandler.GetContentHandler()
	contentHandler.Configure(e, nil)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	movies := []*models.Movie{
		&models.Movie{
			Content: *content[0],
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	movieUseCase.
		EXPECT().
		ListByParams(params, pgnt, userID).
		Return(movies, nil)

	tvshowUseCase.
		EXPECT().
		ListByParams(params, pgnt, userID).
		Return(tvshows, nil)

	response := &response.Response{Body: &response.Body{
		"movies":  movies,
		"tvshows": tvshows,
	}}

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

func TestContentHandler_UpdateContentPostersHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)

	e := echo.New()
	strId := strconv.Itoa(1)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/content"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	contentHandler := NewContentHandler(contentUseCase, movieUseCase, tvshowUseCase)
	handleFunc := contentHandler.UpdatePostersHandler()
	contentHandler.Configure(e, nil)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"code":         101,
			"message":      "request Content-Type isn't multipart/form-data",
			"user_message": "Неверный формат запроса",
		},
	}

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
