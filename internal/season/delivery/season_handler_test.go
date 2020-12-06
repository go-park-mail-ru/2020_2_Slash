package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season/mocks"
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

func setupSeasonHandler(seasonUseCase season.SeasonUsecase, httpMethod string,
	stringifiedJSON string, seasonID uint64) (
	echo.Context, *SeasonHandler, *httptest.ResponseRecorder) {
	e := echo.New()
	strID := ""
	if seasonID != 0 {
		strID = "/" + strconv.FormatUint(seasonID, 10)
	}
	req := httptest.NewRequest(httpMethod, "/api/v1/season" + strID,
		strings.NewReader(stringifiedJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	seasonHandler := NewSeasonHandler(seasonUseCase)
	seasonHandler.Configure(e, nil)
	return c, seasonHandler, rec
}

var testSeason = &models.Season{
	ID:             1,
	Number:         3,
	EpisodesNumber: 8,
	TVShowID:       1,
	Episodes:       nil,
}

var createdSeason = &models.Season{
	ID:             0, // because id assigned in procedure style
	Number:         testSeason.Number,
	EpisodesNumber: 0,
	TVShowID:       testSeason.TVShowID,
	Episodes:       nil,
}

var updTestSeason = &models.Season{
	ID:             1,
	Number:         4,
	EpisodesNumber: 0,
	TVShowID:       1,
	Episodes:       nil,
}

var testEpisodes = []*models.Episode{
	&models.Episode{
		ID:          1,
		Name:        "Рикбег из Рикшенка",
		Number:      1,
		Video:       "/videos/rickandmorty_22/3/1",
		Description: "Саммер решает спасти Рика из тюрьмы.",
		Poster:      "/images/rickandmorty_22/3/1",
		SeasonID:    3,
	},
	&models.Episode{
		ID:          2,
		Name:        "Рикман с камнем",
		Number:      2,
		Video:       "/videos/rickandmorty_22/3/2",
		Description: "Рик, Морти и Саммер охотятся за новым источником энергии в постакалиптической версии Земли.",
		Poster:      "/images/rickandmorty_22/3/2",
		SeasonID:    3,
	},
	&models.Episode{
		ID:          2,
		Name:        "Огурчик Рик",
		Number:      3,
		Video:       "/videos/rickandmorty_22/3/3",
		Description: "Рик превращает себя в огурчик.",
		Poster:      "/images/rickandmorty_22/3/3",
		SeasonID:    3,
	},
}

type Request struct {
	Number   int    `json:"number" validate:"required"`
	TVShowID uint64 `json:"tv_show_id" validate:"required"`
}

func TestSeasonHandler_CreateHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	seasonUseCase := mocks.NewMockSeasonUsecase(ctrl)
	logger.DisableLogger()

	testRequest := &Request{
		Number: testSeason.Number,
		TVShowID: testSeason.TVShowID,
	}

	seasonJSON, err := converter.AnyToBytesBuffer(testRequest)
	if err != nil {
		t.Fatal(err)
	}
	c, seasonHandler, rec := setupSeasonHandler(seasonUseCase,
		http.MethodPost, seasonJSON.String(), 0)
	handleFunc := seasonHandler.CreateHandler()

	seasonUseCase.
		EXPECT().
		Create(createdSeason).
		Return(nil)

	response := &response.Response{Body: &response.Body{"season": createdSeason}}

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

func TestSeasonHandler_ChangeHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	seasonUseCase := mocks.NewMockSeasonUsecase(ctrl)
	logger.DisableLogger()

	testRequest := &Request{
		Number: updTestSeason.Number,
		TVShowID: updTestSeason.TVShowID,
	}

	seasonJSON, err := converter.AnyToBytesBuffer(testRequest)
	if err != nil {
		t.Fatal(err)
	}
	c, seasonHandler, rec := setupSeasonHandler(seasonUseCase,
		http.MethodPut, seasonJSON.String(), updTestSeason.ID)
	c.SetParamNames("id")
	strID := strconv.FormatUint(updTestSeason.ID, 10)
	c.SetParamValues(strID)
	handleFunc := seasonHandler.ChangeHandler()

	seasonUseCase.
		EXPECT().
		Change(updTestSeason).
		Return(nil)

	response := &response.Response{Body: &response.Body{"season": updTestSeason}}

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

func TestSeasonHandler_DeleteHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	seasonUseCase := mocks.NewMockSeasonUsecase(ctrl)
	logger.DisableLogger()

	c, seasonHandler, rec := setupSeasonHandler(seasonUseCase,
		http.MethodDelete, "", testSeason.ID)
	c.SetParamNames("id")
	strID := strconv.FormatUint(testSeason.ID, 10)
	c.SetParamValues(strID)
	handleFunc := seasonHandler.DeleteHandler()

	seasonUseCase.
		EXPECT().
		Delete(testSeason.ID).
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

func TestSeasonHandler_GetHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	seasonUseCase := mocks.NewMockSeasonUsecase(ctrl)
	logger.DisableLogger()

	c, seasonHandler, rec := setupSeasonHandler(seasonUseCase,
		http.MethodGet, "", testSeason.ID)
	c.SetParamNames("id")
	strID := strconv.FormatUint(updTestSeason.ID, 10)
	c.SetParamValues(strID)
	handleFunc := seasonHandler.GetHandler()

	seasonUseCase.
		EXPECT().
		Get(testSeason.ID).
		Return(testSeason, nil)

	seasonUseCase.
		EXPECT().
		GetEpisodes(testSeason.ID).
		Return(testEpisodes, nil)

	response := &response.Response{Body: &response.Body{"season": testSeason}}

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
