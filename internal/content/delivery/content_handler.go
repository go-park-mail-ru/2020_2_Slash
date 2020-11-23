package delivery

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ContentHandler struct {
	contentUcase content.ContentUsecase
}

func NewContentHandler(contentUcase content.ContentUsecase) *ContentHandler {
	return &ContentHandler{
		contentUcase: contentUcase,
	}
}

func (ch *ContentHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.PUT("/api/v1/content/:mid/poster", ch.UpdatePostersHandler(),
		middleware.BodyLimit("10M"), mw.CheckAuth, mw.CheckAdmin, mw.CheckCSRF)
}

func (ch *ContentHandler) UpdatePostersHandler() echo.HandlerFunc {
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

		contentID, _ := strconv.ParseUint(cntx.Param("mid"), 10, 64)
		content, err := ch.contentUcase.GetByID(contentID)
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
		postersDir := postersDirRoot + strconv.Itoa(int(contentID))
		postersDirPath := filepath.Join(path, postersDir)
		helpers.InitStorage(postersDirPath)

		// Store small poster
		if smallImage != nil {
			smallPosterPath := filepath.Join(postersDirPath, smallPosterName)
			if err := helpers.StoreFile(smallImage, smallPosterPath); err != nil {
				if content.Images == "" {
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
				if content.Images == "" {
					os.RemoveAll(postersDirPath)
				}
				logger.Error(err.Message)
				return cntx.JSON(err.HTTPCode, Response{Error: err})
			}
		}

		// Update content
		if err := ch.contentUcase.UpdatePosters(content, postersDir); err != nil {
			if content.Images == "" {
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
