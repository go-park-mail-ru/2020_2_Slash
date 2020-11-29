package delivery

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/season"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/tvshow"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/uniq"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
)

type TVShowHandler struct {
	tvshowUcase   tvshow.TVShowUsecase
	contentUcase  content.ContentUsecase
	countryUcase  country.CountryUsecase
	genreUcase    genre.GenreUsecase
	actorUcase    actor.ActorUseCase
	directorUcase director.DirectorUseCase
	seasonUcase   season.SeasonUsecase
}

func NewTVShowHandler(tvshowUcase tvshow.TVShowUsecase, contentUcase content.ContentUsecase,
	countryUcase country.CountryUsecase, genreUcase genre.GenreUsecase,
	actorUcase actor.ActorUseCase, directorUcase director.DirectorUseCase,
	seasonUcase season.SeasonUsecase) *TVShowHandler {
	return &TVShowHandler{
		tvshowUcase:   tvshowUcase,
		contentUcase:  contentUcase,
		countryUcase:  countryUcase,
		genreUcase:    genreUcase,
		actorUcase:    actorUcase,
		directorUcase: directorUcase,
		seasonUcase:   seasonUcase,
	}
}

func (th *TVShowHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/tvshows", th.CreateTVShowHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.DELETE("/api/v1/tvshows/:tid", th.DeleteTVShowHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.GET("/api/v1/tvshows/:tid", th.GetTVShowHandler(), mw.GetAuth)
	e.GET("/api/v1/tvshows/:tid/episodes", th.GetTVShowSeasonsHandler())
	e.GET("/api/v1/tvshows", th.GetTVShowsHandler(), mw.GetAuth)
	e.GET("/api/v1/tvshows/latest", th.GetLatestTVShowsHandler(), mw.GetAuth)
	e.GET("/api/v1/tvshows/top", th.GetTopTVShowListHandler(), mw.GetAuth)
}

func (th *TVShowHandler) CreateTVShowHandler() echo.HandlerFunc {
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

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Get countries
		countries, err := th.countryUcase.ListByID(uniq.RemoveDuplicates(req.CountriesID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get genres
		genres, err := th.genreUcase.ListByID(uniq.RemoveDuplicates(req.GenresID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get actors
		actors, err := th.actorUcase.ListByID(uniq.RemoveDuplicates(req.ActorsID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get directors
		directors, err := th.directorUcase.ListByID(uniq.RemoveDuplicates(req.DirectorsID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		content := &models.Content{
			Name:             req.Name,
			OriginalName:     req.OriginalName,
			Description:      req.Description,
			ShortDescription: req.ShortDescription,
			Year:             req.Year,
			Countries:        countries,
			Genres:           genres,
			Actors:           actors,
			Directors:        directors,
			Type:             "tvshow",
		}

		tvshow := &models.TVShow{
			Content: *content,
		}

		if err := th.tvshowUcase.Create(tvshow); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"tvshow": tvshow,
			},
		})
	}
}

func (th *TVShowHandler) DeleteTVShowHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		tvshowID, _ := strconv.ParseUint(cntx.Param("tid"), 10, 64)

		tvshow, err := th.tvshowUcase.GetByID(tvshowID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Other related data are deleted in CASCADE
		if err := th.contentUcase.DeleteByID(tvshow.ContentID); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (th *TVShowHandler) GetTVShowHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		tvshowID, _ := strconv.ParseUint(cntx.Param("tid"), 10, 64)

		userID, _ := cntx.Get("userID").(uint64)
		tvshow, err := th.tvshowUcase.GetFullByID(tvshowID, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"tvshow": tvshow,
			},
		})
	}
}

func (th *TVShowHandler) GetTVShowSeasonsHandler() echo.HandlerFunc {
	type ResponseData struct {
		ID      uint64           `json:"id"`
		Name    string           `json:"name"`
		Seasons []*models.Season `json:"seasons"`
	}

	return func(cntx echo.Context) error {
		tvshowID, _ := strconv.ParseUint(cntx.Param("tid"), 10, 64)

		tvshow, err := th.tvshowUcase.GetShortByID(tvshowID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		seasons, err := th.seasonUcase.ListByTVShow(tvshowID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		for _, season := range seasons {
			episodes, customErr := th.seasonUcase.GetEpisodes(season.ID)
			if customErr != nil {
				logger.Error(customErr.Message)
				return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
			}
			season.Episodes = episodes
		}

		res := &ResponseData{
			ID:      tvshow.ID,
			Name:    tvshow.Name,
			Seasons: seasons,
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"tvshow": res,
			},
		})
	}
}

func (th *TVShowHandler) GetTVShowsHandler() echo.HandlerFunc {
	type Request struct {
		models.ContentFilter
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		userID, _ := cntx.Get("userID").(uint64)
		tvshows, err := th.tvshowUcase.ListByParams(&req.ContentFilter,
			&req.Pagination, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"tvshows": tvshows,
			},
		})
	}
}

func (th *TVShowHandler) GetLatestTVShowsHandler() echo.HandlerFunc {
	type Request struct {
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if customErr := reader.NewRequestReader(cntx).Read(req); customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		userID, _ := cntx.Get("userID").(uint64)
		tvshows, err := th.tvshowUcase.ListLatest(&req.Pagination, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"tvshows": tvshows,
			},
		})
	}
}

func (th *TVShowHandler) GetTopTVShowListHandler() echo.HandlerFunc {
	type Request struct {
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if customErr := reader.NewRequestReader(cntx).Read(req); customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		userID, _ := cntx.Get("userID").(uint64)
		tvshows, err := th.tvshowUcase.ListByRating(&req.Pagination, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"tvshows": tvshows,
			},
		})
	}
}
