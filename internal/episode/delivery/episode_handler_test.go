package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode/mocks"
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

func setupEpisodeHandler(seasonUseCase episode.EpisodeUsecase, httpMethod string,
	stringifiedJSON string, episodeID uint64) (
	echo.Context, *EpisodeHandler, *httptest.ResponseRecorder) {
	e := echo.New()
	strID := ""
	if episodeID != 0 {
		strID = "/" + strconv.FormatUint(episodeID, 10)
	}
	req := httptest.NewRequest(httpMethod, "/api/v1/episode"+strID,
		strings.NewReader(stringifiedJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	episodeHandler := NewEpisodeHandler(seasonUseCase)
	episodeHandler.Configure(e, nil)
	return c, episodeHandler, rec
}

var testEpisode = &models.Episode{
	ID:          0,
	Name:        "Огурчик Рик",
	Number:      3,
	Video:       "/videos/rickandmorty_22/3/3",
	Description: "Рик превращает себя в огурчик.",
	Poster:      "/images/rickandmorty_22/3/3",
	SeasonID:    3,
}

var createdEpisode = &models.Episode{
	Name:        "Огурчик Рик",
	Number:      3,
	Description: "Рик превращает себя в огурчик.",
	SeasonID:    3,
}

type Request struct {
	Name        string `json:"name" validate:"required,lte=128"`
	Number      int    `json:"number" validate:"required"`
	Description string `json:"description" validate:"required"`
	SeasonID    uint64 `json:"season_id" validate:"required"`
}

func TestEpisodeHandler_CreateHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	episodeUsecase := mocks.NewMockEpisodeUsecase(ctrl)
	logger.DisableLogger()

	testRequest := &Request{
		Number:      createdEpisode.Number,
		Name:        createdEpisode.Name,
		Description: createdEpisode.Description,
		SeasonID:    createdEpisode.SeasonID,
	}

	episodeJSON, err := converter.AnyToBytesBuffer(testRequest)
	if err != nil {
		t.Fatal(err)
	}
	c, episodeHandler, rec := setupEpisodeHandler(episodeUsecase,
		http.MethodPost, episodeJSON.String(), 0)
	handleFunc := episodeHandler.CreateHandler()

	episodeUsecase.
		EXPECT().
		Create(createdEpisode).
		Return(nil)

	response := &response.Response{Body: &response.Body{"episode": createdEpisode}}

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

func TestEpisodeHandler_ChangeHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	episodeUsecase := mocks.NewMockEpisodeUsecase(ctrl)
	logger.DisableLogger()

	testRequest := &Request{
		Number:      createdEpisode.Number,
		Name:        createdEpisode.Name,
		Description: createdEpisode.Description,
		SeasonID:    createdEpisode.SeasonID,
	}

	episodeJSON, err := converter.AnyToBytesBuffer(testRequest)
	if err != nil {
		t.Fatal(err)
	}
	c, episodeHandler, rec := setupEpisodeHandler(episodeUsecase,
		http.MethodPut, episodeJSON.String(), testEpisode.ID)
	c.SetParamNames("eid")
	strID := strconv.FormatUint(testEpisode.ID, 10)
	c.SetParamValues(strID)
	handleFunc := episodeHandler.ChangeHandler()

	episodeUsecase.
		EXPECT().
		Change(createdEpisode).
		Return(nil)

	response := &response.Response{Body: &response.Body{"episode": createdEpisode}}

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

func TestEpisodeHandler_DeleteHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	episodeUsecase := mocks.NewMockEpisodeUsecase(ctrl)
	logger.DisableLogger()

	testRequest := &Request{
		Number:      createdEpisode.Number,
		Name:        createdEpisode.Name,
		Description: createdEpisode.Description,
		SeasonID:    createdEpisode.SeasonID,
	}

	episodeJSON, err := converter.AnyToBytesBuffer(testRequest)
	if err != nil {
		t.Fatal(err)
	}
	c, episodeHandler, rec := setupEpisodeHandler(episodeUsecase,
		http.MethodDelete, episodeJSON.String(), testEpisode.ID)
	c.SetParamNames("eid")
	strID := strconv.FormatUint(testEpisode.ID, 10)
	c.SetParamValues(strID)
	handleFunc := episodeHandler.DeleteHandler()

	episodeUsecase.
		EXPECT().
		DeleteByID(createdEpisode.ID).
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
