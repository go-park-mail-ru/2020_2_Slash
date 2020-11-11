package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/country"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/director"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/genre"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/movie"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/uniq"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type MovieHandler struct {
	movieUcase    movie.MovieUsecase
	contentUcase  content.ContentUsecase
	countryUcase  country.CountryUsecase
	genreUcase    genre.GenreUsecase
	actorUcase    actor.ActorUseCase
	directorUcase director.DirectorUseCase
}

func NewMovieHandler(movieUcase movie.MovieUsecase, contentUcase content.ContentUsecase,
	countryUcase country.CountryUsecase, genreUcase genre.GenreUsecase,
	actorUcase actor.ActorUseCase, directorUcase director.DirectorUseCase) *MovieHandler {
	return &MovieHandler{
		movieUcase:    movieUcase,
		contentUcase:  contentUcase,
		countryUcase:  countryUcase,
		genreUcase:    genreUcase,
		actorUcase:    actorUcase,
		directorUcase: directorUcase,
	}
}

func (mh *MovieHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/movies", mh.CreateMovieHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.PUT("/api/v1/movies/:mid", mh.UpdateMovieHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.DELETE("/api/v1/movies/:mid", mh.DeleteMovieHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.GET("/api/v1/movies/:mid", mh.GetMovieHandler(), mw.GetAuth)
	e.PUT("/api/v1/movies/:mid/poster", mh.UpdateMoviePostersHandler(),
		middleware.BodyLimit("10M"), mw.CheckAuth, mw.CheckAdmin)
	e.PUT("/api/v1/movies/:mid/video", mh.UpdateMovieVideoHandler(),
		middleware.BodyLimit("1000M"), mw.CheckAuth, mw.CheckAdmin)
	e.GET("/api/v1/movies", mh.GetMoviesHandler(), mw.GetAuth)
	e.GET("/api/v1/movies/latest", mh.GetLatestMoviesHandler(), mw.GetAuth)
	e.GET("/api/v1/movies/top", mh.GetTopMovieListHandler(), mw.GetAuth)
}

func (mh *MovieHandler) CreateMovieHandler() echo.HandlerFunc {
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
		countries, err := mh.countryUcase.ListByID(uniq.RemoveDuplicates(req.CountriesID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get genres
		genres, err := mh.genreUcase.ListByID(uniq.RemoveDuplicates(req.GenresID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get actors
		actors, err := mh.actorUcase.ListByID(uniq.RemoveDuplicates(req.ActorsID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get directors
		directors, err := mh.directorUcase.ListByID(uniq.RemoveDuplicates(req.DirectorsID))
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
			Type:             "movie",
		}

		movie := &models.Movie{
			Content: *content,
		}

		if err := mh.movieUcase.Create(movie); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"movie": movie,
			},
		})
	}
}

func (mh *MovieHandler) UpdateMovieHandler() echo.HandlerFunc {
	type Request struct {
		Name             string   `json:"name" validate:"lte=128"`
		OriginalName     string   `json:"original_name" validate:"lte=128"`
		Description      string   `json:"description"`
		ShortDescription string   `json:"short_description"`
		Year             int      `json:"year"`
		CountriesID      []uint64 `json:"countries"`
		GenresID         []uint64 `json:"genres"`
		ActorsID         []uint64 `json:"actors"`
		DirectorsID      []uint64 `json:"directors"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Get countries
		countries, err := mh.countryUcase.ListByID(uniq.RemoveDuplicates(req.CountriesID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get genres
		genres, err := mh.genreUcase.ListByID(uniq.RemoveDuplicates(req.GenresID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get actors
		actors, err := mh.actorUcase.ListByID(uniq.RemoveDuplicates(req.ActorsID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}
		// Get directors
		directors, err := mh.directorUcase.ListByID(uniq.RemoveDuplicates(req.DirectorsID))
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		contentData := &models.Content{
			Name:             req.Name,
			OriginalName:     req.OriginalName,
			Description:      req.Description,
			ShortDescription: req.ShortDescription,
			Year:             req.Year,
			Countries:        countries,
			Genres:           genres,
			Actors:           actors,
			Directors:        directors,
		}

		movieID, _ := strconv.ParseUint(cntx.Param("mid"), 10, 64)
		movie, err := mh.movieUcase.GetByID(movieID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		content, err := mh.contentUcase.UpdateByID(movie.ContentID, contentData)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		movie.Content = *content
		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"movie": movie,
			},
		})
	}
}

func (mh *MovieHandler) DeleteMovieHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		movieID, _ := strconv.ParseUint(cntx.Param("mid"), 10, 64)

		movie, err := mh.movieUcase.GetByID(movieID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		if err := mh.movieUcase.DeleteByID(movie.ID); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Other related data are deleted in CASCADE
		if err := mh.contentUcase.DeleteByID(movie.ContentID); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}

func (mh *MovieHandler) GetMovieHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		movieID, _ := strconv.ParseUint(cntx.Param("mid"), 10, 64)

		userID, _ := cntx.Get("userID").(uint64)
		movie, err := mh.movieUcase.GetFullByID(movieID, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"movie": movie,
			},
		})
	}
}

func (mh *MovieHandler) UpdateMoviePostersHandler() echo.HandlerFunc {
	const postersDirRoot = "/images/"
	const smallPosterName = "small.png"
	const largePosterName = "large.png"

	return func(cntx echo.Context) error {
		smallImage, err := reader.NewRequestReader(cntx).ReadNotRequiredImage("small_poster")
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		largeImage, err := reader.NewRequestReader(cntx).ReadNotRequiredImage("large_poster")
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Check for passing at least one image
		if smallImage == nil && largeImage == nil {
			err := errors.Get(CodeBadRequest)
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		movieID, _ := strconv.ParseUint(cntx.Param("mid"), 10, 64)
		movie, err := mh.movieUcase.GetWithContentByID(movieID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		path, osErr := os.Getwd()
		if osErr != nil {
			err := errors.New(CodeInternalError, osErr)
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Create posters directory
		postersDir := postersDirRoot + strconv.Itoa(int(movie.ContentID))
		postersDirPath := filepath.Join(path, postersDir)
		helpers.InitStorage(postersDirPath)

		// Store small poster
		if smallImage != nil {
			smallPosterPath := filepath.Join(postersDirPath, smallPosterName)
			if err := helpers.StoreFile(smallImage, smallPosterPath); err != nil {
				if movie.Content.Images == "" {
					os.RemoveAll(postersDirPath)
				}
				logger.Error(err.Message)
				return cntx.JSON(err.HTTPCode, Response{Error: err})
			}
		}

		// Store large poster
		if largeImage != nil {
			largePosterPath := filepath.Join(postersDirPath, largePosterName)
			if err := helpers.StoreFile(largeImage, largePosterPath); err != nil {
				if movie.Content.Images == "" {
					os.RemoveAll(postersDirPath)
				}
				logger.Error(err.Message)
				return cntx.JSON(err.HTTPCode, Response{Error: err})
			}
		}

		// Update content
		if err := mh.contentUcase.UpdatePosters(&movie.Content, postersDir); err != nil {
			if movie.Content.Images == "" {
				os.RemoveAll(postersDirPath)
			}
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"images": postersDir,
			},
		})
	}
}

func (mh *MovieHandler) UpdateMovieVideoHandler() echo.HandlerFunc {
	const videosDirRoot = "/videos/"
	const videoName = "movie.mp4"

	return func(cntx echo.Context) error {
		video, err := reader.NewRequestReader(cntx).ReadVideo("video")
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		movieID, _ := strconv.ParseUint(cntx.Param("mid"), 10, 64)
		movie, err := mh.movieUcase.GetByID(movieID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		path, osErr := os.Getwd()
		if osErr != nil {
			err := errors.New(CodeInternalError, osErr)
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Create videos directory
		videosDir := videosDirRoot + strconv.Itoa(int(movie.ContentID))
		videosDirPath := filepath.Join(path, videosDir)
		helpers.InitStorage(videosDirPath)

		// Store video
		absVideoPath := filepath.Join(videosDirPath, videoName)
		if err := helpers.StoreFile(video, absVideoPath); err != nil {
			if movie.Video == "" {
				os.RemoveAll(videosDirPath)
			}
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Update movie
		rltVideoPath := filepath.Join(videosDir, videoName)
		if err := mh.movieUcase.UpdateVideo(movie, rltVideoPath); err != nil {
			if movie.Video == "" {
				os.RemoveAll(videosDirPath)
			}
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"video": rltVideoPath,
			},
		})
	}
}

func (mh *MovieHandler) GetMoviesHandler() echo.HandlerFunc {
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
		movies, err := mh.movieUcase.ListByParams(&req.ContentFilter,
			&req.Pagination, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"movies": movies,
			},
		})
	}
}

func (mh *MovieHandler) GetLatestMoviesHandler() echo.HandlerFunc {
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
		movies, err := mh.movieUcase.ListLatest(&req.Pagination, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"movies": movies,
			},
		})
	}
}

func (mh *MovieHandler) GetTopMovieListHandler() echo.HandlerFunc {
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
		movies, err := mh.movieUcase.ListByRating(&req.Pagination, userID)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"movies": movies,
			},
		})
	}
}
