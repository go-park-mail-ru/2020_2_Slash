package delivery

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	actorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/actor/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	contentMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/content/mocks"
	countryMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/country/mocks"
	directorMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/director/mocks"
	genreMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/genre/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	movieMocks "github.com/go-park-mail-ru/2020_2_Slash/internal/movie/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
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
	CountriesID      []uint64 `json:"countries" validate:"required"`
	GenresID         []uint64 `json:"genres" validate:"required"`
	ActorsID         []uint64 `json:"actors" validate:"required"`
	DirectorsID      []uint64 `json:"directors" validate:"required"`
}

func TestMovieHandler_CreateMovieHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

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
		Type:             "movie",
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
		CountriesID:      countriesID,
		GenresID:         genresID,
		ActorsID:         actorsID,
		DirectorsID:      directorsID,
	}

	var movieInst *models.Movie = &models.Movie{
		Content: *contentInst,
	}

	movieJSON, err := converter.AnyToBytesBuffer(requestData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies", strings.NewReader(movieJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.CreateMovieHandler()
	movieHandler.Configure(e, nil)

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

	movieUseCase.
		EXPECT().
		Create(movieInst).
		Return(nil)

	response := &response.Response{Body: &response.Body{"movie": movieInst}}

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

func TestMovieHandler_UpdateMovieHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

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
		CountriesID:      countriesID,
		GenresID:         genresID,
		ActorsID:         actorsID,
		DirectorsID:      directorsID,
	}

	var movieInst *models.Movie = &models.Movie{
		Content: *contentInst,
	}

	movieJSON, err := converter.AnyToBytesBuffer(requestData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(movieInst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies"+strId, strings.NewReader(movieJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.UpdateMovieHandler()
	movieHandler.Configure(e, nil)

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

	movieUseCase.
		EXPECT().
		GetByID(movieInst.ID).
		Return(movieInst, nil)

	contentUseCase.
		EXPECT().
		Update(contentInst).
		Return(nil)

	response := &response.Response{Body: &response.Body{"movie": movieInst}}

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

func TestMovieHandler_DeleteMovieHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

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

	var movieInst *models.Movie = &models.Movie{
		Content: *contentInst,
	}

	e := echo.New()
	strId := strconv.Itoa(int(movieInst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.DeleteMovieHandler()
	movieHandler.Configure(e, nil)

	movieUseCase.
		EXPECT().
		GetByID(movieInst.ID).
		Return(movieInst, nil)

	movieUseCase.
		EXPECT().
		DeleteByID(movieInst.ID).
		Return(nil)

	contentUseCase.
		EXPECT().
		DeleteByID(movieInst.ContentID).
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

func TestMovieHandler_DeleteMovieHandler_NoMovie(t *testing.T) {
	t.Parallel()
	logger.DisableLogger()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

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

	var movieInst *models.Movie = &models.Movie{
		Content: *contentInst,
	}

	e := echo.New()
	strId := strconv.Itoa(int(movieInst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.DeleteMovieHandler()
	movieHandler.Configure(e, nil)

	movieUseCase.
		EXPECT().
		GetByID(movieInst.ID).
		Return(nil, errors.Get(consts.CodeMovieDoesNotExist))

	response := &response.Response{Error: errors.Get(consts.CodeMovieDoesNotExist)}

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

func TestMovieHandler_GetMovieHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

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

	var movieInst *models.Movie = &models.Movie{
		Content: *contentInst,
	}
	var userID uint64 = 0

	e := echo.New()
	strId := strconv.Itoa(int(movieInst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("mid")
	c.SetParamValues(strId)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.GetMovieHandler()
	movieHandler.Configure(e, nil)

	movieUseCase.
		EXPECT().
		GetFullByID(movieInst.ID, userID).
		Return(movieInst, nil)

	response := &response.Response{Body: &response.Body{"movie": movieInst}}

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

func TestMovieHandler_UpdateMovieVideoHandler(t *testing.T) {
	t.Parallel()
	// Setup
	logger.DisableLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

	e := echo.New()
	strId := strconv.Itoa(1)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.UpdateMovieVideoHandler()
	movieHandler.Configure(e, nil)

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

func TestMovieHandler_GetMoviesHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

	pgnt := &models.Pagination{
		From:  0,
		Count: 1,
	}
	var userID uint64 = 0

	params := &models.ContentFilter{
		Year:     2001,
		Genre:    1,
		Country:  1,
		Actor:    1,
		Director: 1,
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
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.GetMoviesHandler()
	movieHandler.Configure(e, nil)

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

	movieUseCase.
		EXPECT().
		ListByParams(params, pgnt, userID).
		Return(movies, nil)

	response := &response.Response{Body: &response.Body{"movies": movies}}

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

func TestMovieHandler_GetLatestMoviesHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

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
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies/latest", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.GetLatestMoviesHandler()
	movieHandler.Configure(e, nil)

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

	movieUseCase.
		EXPECT().
		ListLatest(pgnt, userID).
		Return(movies, nil)

	response := &response.Response{Body: &response.Body{"movies": movies}}

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

func TestMovieHandler_GetMoviesByRatingHandler(t *testing.T) {
	t.Parallel()
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

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
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies/top", strings.NewReader(reqJSON.String()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.GetTopMovieListHandler()
	movieHandler.Configure(e, nil)

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

	movieUseCase.
		EXPECT().
		ListByRating(pgnt, userID).
		Return(movies, nil)

	response := &response.Response{Body: &response.Body{"movies": movies}}

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
