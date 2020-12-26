package delivery

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	actorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	countryMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
	directorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	genreMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	seasonMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/season/mocks"
	tvshowMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var countries = []*models.Country{
	&models.Country{
		ID:   1,
		Name: "США",
	},
}

var genres = []*models.Genre{
	&models.Genre{
		Name: "Мультфильм",
	},
	&models.Genre{
		Name: "Комедия",
	},
}

var actors = []*models.Actor{
	&models.Actor{
		Name: "Майк Майерс",
	},
	&models.Actor{
		Name: "Эдди Мёрфи",
	},
}

var directors = []*models.Director{
	&models.Director{
		Name: "Эндрю Адамсон",
	},
	&models.Director{
		Name: "Вики Дженсон",
	},
}

type Request struct {
	Name             string   `json:"name" validate:"required,lte=128"`
	OriginalName     string   `json:"original_name" validate:"required,lte=128"`
	Description      string   `json:"description" validate:"required"`
	ShortDescription string   `json:"short_description" validate:"required"`
	Year             int      `json:"year" validate:"required"`
	IsFree           *bool    `json:"is_free" validate:"required"`
	CountriesID      []uint64 `json:"countries" validate:"required"`
	GenresID         []uint64 `json:"genres" validate:"required"`
	ActorsID         []uint64 `json:"actors" validate:"required"`
	DirectorsID      []uint64 `json:"directors" validate:"required"`
}

func TestTVShowHandler_CreateTVShowHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)

	isFree := true
	var contentInst *models.Content = &models.Content{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		IsFree:           &isFree,
		Countries:        countries,
		Genres:           genres,
		Actors:           actors,
		Directors:        directors,
		Type:             "tvshow",
	}

	countriesID := []uint64{1}
	directorsID := []uint64{1, 2}
	actorsID := []uint64{1, 2}
	genresID := []uint64{1, 2}

	requestData := &Request{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		IsFree:           &isFree,
		CountriesID:      countriesID,
		GenresID:         genresID,
		ActorsID:         actorsID,
		DirectorsID:      directorsID,
	}

	var tvshowInst *models.TVShow = &models.TVShow{
		Content: *contentInst,
	}

	tvshowJSON, err := converter.AnyToBytesBuffer(requestData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows", strings.NewReader(tvshowJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	tvshowHandler := NewTVShowHandler(tvshowUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase, seasonUseCase)
	handleFunc := tvshowHandler.CreateTVShowHandler()
	tvshowHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		ListByID(countriesID).
		Return(countries, nil)

	genreUseCase.
		EXPECT().
		ListByID(genresID).
		Return(genres, nil)

	actorUseCase.
		EXPECT().
		ListByID(actorsID).
		Return(actors, nil)

	directorUseCase.
		EXPECT().
		ListByID(directorsID).
		Return(directors, nil)

	tvshowUseCase.
		EXPECT().
		Create(tvshowInst).
		Return(nil)

	response := &response.Response{Body: &response.Body{"tvshow": tvshowInst}}

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

func TestTVShowHandler_UpdateTVShowHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)

	isFree := true
	var contentInst *models.Content = &models.Content{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		IsFree:           &isFree,
		Countries:        countries,
		Genres:           genres,
		Actors:           actors,
		Directors:        directors,
	}

	countriesID := []uint64{1}
	directorsID := []uint64{1, 2}
	actorsID := []uint64{1, 2}
	genresID := []uint64{1, 2}

	requestData := &Request{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		IsFree:           &isFree,
		CountriesID:      countriesID,
		GenresID:         genresID,
		ActorsID:         actorsID,
		DirectorsID:      directorsID,
	}

	var tvshowInst *models.TVShow = &models.TVShow{
		Content: *contentInst,
	}

	tvshowJSON, err := converter.AnyToBytesBuffer(requestData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(tvshowInst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows"+strId, strings.NewReader(tvshowJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tid")
	c.SetParamValues(strId)
	c.Set("userID", 3)

	tvshowHandler := NewTVShowHandler(tvshowUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase, seasonUseCase)
	handleFunc := tvshowHandler.UpdateTVShowHandler()
	tvshowHandler.Configure(e, nil)

	countryUseCase.
		EXPECT().
		ListByID(countriesID).
		Return(countries, nil)

	genreUseCase.
		EXPECT().
		ListByID(genresID).
		Return(genres, nil)

	actorUseCase.
		EXPECT().
		ListByID(actorsID).
		Return(actors, nil)

	directorUseCase.
		EXPECT().
		ListByID(directorsID).
		Return(directors, nil)

	tvshowUseCase.
		EXPECT().
		GetByID(tvshowInst.ID).
		Return(tvshowInst, nil)

	contentUseCase.
		EXPECT().
		UpdateByID(tvshowInst.ContentID, contentInst).
		Return(contentInst, nil)

	response := &response.Response{Body: &response.Body{"tvshow": tvshowInst}}

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

func TestTVShowHandler_DeleteTVShowHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)

	var contentInst *models.Content = &models.Content{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		Countries:        countries,
		Genres:           genres,
		Actors:           actors,
		Directors:        directors,
	}

	var tvshowInst *models.TVShow = &models.TVShow{
		Content: *contentInst,
	}

	e := echo.New()
	strId := strconv.Itoa(int(tvshowInst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tid")
	c.SetParamValues(strId)

	tvshowHandler := NewTVShowHandler(tvshowUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase, seasonUseCase)
	handleFunc := tvshowHandler.DeleteTVShowHandler()
	tvshowHandler.Configure(e, nil)

	tvshowUseCase.
		EXPECT().
		GetByID(tvshowInst.ID).
		Return(tvshowInst, nil)

	contentUseCase.
		EXPECT().
		DeleteByID(tvshowInst.ContentID).
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

func TestTVShowHandler_GetTVShowHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)

	var contentInst *models.Content = &models.Content{
		Name:             "Шрек",
		OriginalName:     "Shrek",
		Description:      "Полная сюрпризов сказка об ужасном болотном огре, который ненароком наводит порядок в Сказочной стране",
		ShortDescription: "Полная сюрпризов сказка об ужасном болотном огре",
		Year:             2001,
		Countries:        countries,
		Genres:           genres,
		Actors:           actors,
		Directors:        directors,
	}

	var tvshowInst *models.TVShow = &models.TVShow{
		Content: *contentInst,
	}
	var userID uint64 = 0

	e := echo.New()
	strId := strconv.Itoa(int(tvshowInst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tid")
	c.SetParamValues(strId)

	tvshowHandler := NewTVShowHandler(tvshowUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase, seasonUseCase)
	handleFunc := tvshowHandler.GetTVShowHandler()
	tvshowHandler.Configure(e, nil)

	tvshowUseCase.
		EXPECT().
		GetFullByID(tvshowInst.ID, userID).
		Return(tvshowInst, nil)

	response := &response.Response{Body: &response.Body{"tvshow": tvshowInst}}

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

func TestTVShowHandler_GetTVShowsHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)

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
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	tvshowHandler := NewTVShowHandler(tvshowUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase, seasonUseCase)
	handleFunc := tvshowHandler.GetTVShowsHandler()
	tvshowHandler.Configure(e, nil)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	tvshowUseCase.
		EXPECT().
		ListByParams(params, pgnt, userID).
		Return(tvshows, nil)

	response := &response.Response{Body: &response.Body{"tvshows": tvshows}}

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

func TestTVShowHandler_GetLatestTVShowsHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 0

	reqJSON, err := converter.AnyToBytesBuffer(pgnt)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows/latest", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	tvshowHandler := NewTVShowHandler(tvshowUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase, seasonUseCase)
	handleFunc := tvshowHandler.GetLatestTVShowsHandler()
	tvshowHandler.Configure(e, nil)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	tvshowUseCase.
		EXPECT().
		ListLatest(pgnt, userID).
		Return(tvshows, nil)

	response := &response.Response{Body: &response.Body{"tvshows": tvshows}}

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

func TestTVShowHandler_GetTVShowsByRatingHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tvshowUseCase := tvshowMocks.NewMockTVShowUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)
	seasonUseCase := seasonMocks.NewMockSeasonUsecase(ctrl)

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 0

	reqJSON, err := converter.AnyToBytesBuffer(pgnt)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tvshows/top", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	tvshowHandler := NewTVShowHandler(tvshowUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase, seasonUseCase)
	handleFunc := tvshowHandler.GetTopTVShowListHandler()
	tvshowHandler.Configure(e, nil)

	content := []*models.Content{
		&models.Content{
			Name: "Shrek",
		},
	}

	tvshows := []*models.TVShow{
		&models.TVShow{
			Content: *content[0],
		},
	}

	tvshowUseCase.
		EXPECT().
		ListByRating(pgnt, userID).
		Return(tvshows, nil)

	response := &response.Response{Body: &response.Body{"tvshows": tvshows}}

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
