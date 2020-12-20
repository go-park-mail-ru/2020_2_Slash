package delivery

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/episode"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EpisodeHandler struct {
	episodeUsecase episode.EpisodeUsecase
}

func NewEpisodeHandler(usecase episode.EpisodeUsecase) *EpisodeHandler {
	return &EpisodeHandler{episodeUsecase: usecase}
}

func (eh *EpisodeHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/episodes", eh.CreateHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.PUT("/api/v1/episodes/:eid", eh.ChangeHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.DELETE("/api/v1/episodes/:eid", eh.DeleteHandler(), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.GET("/api/v1/episodes/:eid", eh.GetHandler(), mw.GetAuth)
	e.PUT("/api/v1/episodes/:eid/poster", eh.UpdatePosterHandler(),
		middleware.BodyLimit("10M"), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
	e.PUT("/api/v1/episodes/:eid/video", eh.UpdateVideoHandler(),
		middleware.BodyLimit("1000M"), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
}

func (eh *EpisodeHandler) CreateHandler() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name" validate:"required,lte=128"`
		Number      int    `json:"number" validate:"required"`
		Description string `json:"description" validate:"required"`
		SeasonID    uint64 `json:"season_id" validate:"required"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		episode := &models.Episode{
			Name:        req.Name,
			Number:      req.Number,
			Description: req.Description,
			SeasonID:    req.SeasonID,
		}
		customErr := eh.episodeUsecase.Create(episode)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"episode": episode,
			},
		})
	}
}

func (eh *EpisodeHandler) ChangeHandler() echo.HandlerFunc {
	type Request struct {
		Name        string `json:"name" validate:"required,lte=128"`
		Number      int    `json:"number" validate:"required"`
		Description string `json:"description" validate:"required"`
		SeasonID    uint64 `json:"season_id" validate:"required"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		episodeID, err := strconv.ParseUint(cntx.Param("eid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		episode := &models.Episode{
			ID:          episodeID,
			Name:        req.Name,
			Number:      req.Number,
			Description: req.Description,
			SeasonID:    req.SeasonID,
		}
		customErr := eh.episodeUsecase.Change(episode)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"episode": episode,
			},
		})
	}
}

func (eh *EpisodeHandler) DeleteHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		episodeID, err := strconv.ParseUint(cntx.Param("eid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		customErr := eh.episodeUsecase.DeleteByID(episodeID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{Message: "success"})
	}
}

func (eh *EpisodeHandler) GetHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		episodeID, err := strconv.ParseUint(cntx.Param("eid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		episode, customErr := eh.episodeUsecase.GetByID(episodeID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"episode": episode,
			},
		})
	}
}

func (eh *EpisodeHandler) GetInfoToStoreMedia(episodeID uint64) (*models.Episode, *models.Content, int, *errors.Error) {
	episode, customErr := eh.episodeUsecase.GetByID(episodeID)
	if customErr != nil {
		return nil, nil, 0, customErr
	}

	content, customErr := eh.episodeUsecase.GetContentByEID(episodeID)
	if customErr != nil {
		return nil, nil, 0, customErr
	}

	seasonNumber, customErr := eh.episodeUsecase.GetSeasonNumber(episodeID)
	if customErr != nil {
		return nil, nil, 0, customErr
	}

	return episode, content, seasonNumber, nil
}

func buildSeasonDirPath(seasonNumber int, content *models.Content) string {
	// lowercaseorigin_cid
	contentName := helpers.GetContentDirTitle(content.OriginalName, content.ContentID)
	// lowercaseorigin_cid/season_number/
	seasonDirTitle := filepath.Join(contentName, strconv.Itoa(seasonNumber))
	return seasonDirTitle
}

func (eh *EpisodeHandler) UpdatePosterHandler() echo.HandlerFunc {
	const postersDirRoot = "/images/"
	const format = ".webp"

	return func(cntx echo.Context) error {
		posterImage, customErr := reader.NewRequestReader(cntx).ReadNotRequiredImage("poster")
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		if posterImage == nil {
			customErr := errors.Get(consts.CodeBadRequest)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		episodeID, err := strconv.ParseUint(cntx.Param("eid"), 10, 64)
		if err != nil {
			customErr := errors.Get(consts.CodeBadRequest)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		episode, content, seasonNumber, customErr := eh.GetInfoToStoreMedia(episodeID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		path, osErr := os.Getwd()
		if osErr != nil {
			err := errors.New(consts.CodeInternalError, osErr)
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Create posters directory
		// /images/name_cid/season_number/
		seasonDir := postersDirRoot + buildSeasonDirPath(seasonNumber, content)
		postersDirAbsPath := filepath.Join(path, seasonDir)
		helpers.InitTree(postersDirAbsPath)

		// Store poster
		posterName := strconv.Itoa(episode.Number) + format
		absPosterPath := filepath.Join(postersDirAbsPath, posterName)
		if err := helpers.StoreSmallImage(posterImage, absPosterPath); err != nil {
			if episode.Poster == "" {
				removeErr := os.RemoveAll(postersDirAbsPath)
				if removeErr != nil {
					logger.Error(removeErr)
				}
			}
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Update episode poster
		rltPosterPath := filepath.Join(seasonDir, posterName)
		if err := eh.episodeUsecase.UpdatePoster(episode, rltPosterPath); err != nil {
			if episode.Poster == "" {
				removeErr := os.RemoveAll(postersDirAbsPath)
				if removeErr != nil {
					logger.Error(removeErr)
				}
			}
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"poster": rltPosterPath,
			},
		})
	}
}

func (eh *EpisodeHandler) UpdateVideoHandler() echo.HandlerFunc {
	const videosDirRoot = "/videos/"
	const format = ".mp4"

	return func(cntx echo.Context) error {
		video, err := reader.NewRequestReader(cntx).ReadVideo("video")
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		episodeID, parseErr := strconv.ParseUint(cntx.Param("eid"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(consts.CodeInternalError, parseErr)
			logger.Error(customErr)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		episode, content, seasonNumber, customErr := eh.GetInfoToStoreMedia(episodeID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: err})
		}

		path, osErr := os.Getwd()
		if osErr != nil {
			customErr := errors.New(consts.CodeInternalError, osErr)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: err})
		}

		// Create videos directory
		// /videos/name_cid/season_number/
		seasonDir := videosDirRoot + buildSeasonDirPath(seasonNumber, content)
		videosDirPath := filepath.Join(path, seasonDir)
		helpers.InitTree(videosDirPath)

		// Store video
		videoName := strconv.Itoa(episode.Number) + format
		absVideoPath := filepath.Join(videosDirPath, videoName)
		if err := helpers.StoreFile(video, absVideoPath); err != nil {
			if episode.Video == "" {
				removeErr := os.RemoveAll(videosDirPath)
				if removeErr != nil {
					logger.Error(removeErr)
				}
			}
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		// Update episode
		rltVideoPath := filepath.Join(seasonDir, videoName)
		if err := eh.episodeUsecase.UpdateVideo(episode, rltVideoPath); err != nil {
			if episode.Video == "" {
				removeErr := os.RemoveAll(videosDirPath)
				if removeErr != nil {
					logger.Error(removeErr)
				}
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
