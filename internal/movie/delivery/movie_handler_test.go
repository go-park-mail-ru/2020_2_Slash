package delivery

import (
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

	var content_inst *models.Content = &models.Content{
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

	var movie_inst *models.Movie = &models.Movie{
		Content: *content_inst,
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
		Create(movie_inst).
		Return(nil)

	response := &response.Response{Body: &response.Body{"movie": movie_inst}}

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

	var content_inst *models.Content = &models.Content{
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

	var movie_inst *models.Movie = &models.Movie{
		Content: *content_inst,
	}

	movieJSON, err := converter.AnyToBytesBuffer(requestData)
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	strId := strconv.Itoa(int(movie_inst.ID))
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
		GetByID(movie_inst.ID).
		Return(movie_inst, nil)

	contentUseCase.
		EXPECT().
		UpdateByID(movie_inst.ContentID, content_inst).
		Return(content_inst, nil)

	response := &response.Response{Body: &response.Body{"movie": movie_inst}}

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

	var content_inst *models.Content = &models.Content{
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

	var movie_inst *models.Movie = &models.Movie{
		Content: *content_inst,
	}

	e := echo.New()
	strId := strconv.Itoa(int(movie_inst.ID))
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
		GetByID(movie_inst.ID).
		Return(movie_inst, nil)

	movieUseCase.
		EXPECT().
		DeleteByID(movie_inst.ID).
		Return(nil)

	contentUseCase.
		EXPECT().
		DeleteByID(movie_inst.ContentID).
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
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	movieUseCase := movieMocks.NewMockMovieUsecase(ctrl)
	contentUseCase := contentMocks.NewMockContentUsecase(ctrl)
	countryUseCase := countryMocks.NewMockCountryUsecase(ctrl)
	genreUseCase := genreMocks.NewMockGenreUsecase(ctrl)
	actorUseCase := actorMocks.NewMockActorUseCase(ctrl)
	directorUseCase := directorMocks.NewMockDirectorUseCase(ctrl)

	var content_inst *models.Content = &models.Content{
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

	var movie_inst *models.Movie = &models.Movie{
		Content: *content_inst,
	}

	e := echo.New()
	strId := strconv.Itoa(int(movie_inst.ID))
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
		GetByID(movie_inst.ID).
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

	var content_inst *models.Content = &models.Content{
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

	var movie_inst *models.Movie = &models.Movie{
		Content: *content_inst,
	}

	e := echo.New()
	strId := strconv.Itoa(int(movie_inst.ID))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies"+strId, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strId)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.GetMovieHandler()
	movieHandler.Configure(e, nil)

	movieUseCase.
		EXPECT().
		GetFullByID(movie_inst.ID).
		Return(movie_inst, nil)

	response := &response.Response{Body: &response.Body{"movie": movie_inst}}

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

func TestMovieHandler_UpdateMoviePostersHandler(t *testing.T) {
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
	handleFunc := movieHandler.UpdateMoviePostersHandler()
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

func TestMovieHandler_UpdateMovieVideoHandler(t *testing.T) {
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

func TestMovieHandler_GetMovieListByGenreHandler(t *testing.T) {
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

	genre := &models.Genre{
		ID:   1,
		Name: "comedy",
	}

	e := echo.New()
	genreID := strconv.Itoa(1)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/movies?genre="+genreID, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamValues(genreID)

	movieHandler := NewMovieHandler(movieUseCase, contentUseCase,
		countryUseCase, genreUseCase, actorUseCase, directorUseCase)
	handleFunc := movieHandler.GetMovieListByGenreHandler()
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
		ListByGenre(genre.ID).
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
