package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/content"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/favourite"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FavouriteHandler struct {
	favouriteUseCase favourite.FavouriteUsecase
	contentUseCase   content.ContentUsecase
}

func NewFavouriteHandler(favouriteUseCase favourite.FavouriteUsecase,
	contentUseCase content.ContentUsecase) *FavouriteHandler {
	return &FavouriteHandler{
		favouriteUseCase: favouriteUseCase,
		contentUseCase:   contentUseCase,
	}
}

func (fh *FavouriteHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/favourites", fh.CreateHandler(), mw.CheckAuth)
	e.DELETE("/api/v1/favourites", fh.DeleteHandler(), mw.CheckAuth)
	e.GET("/api/v1/favourites", fh.GetFavouritesHandler(), mw.CheckAuth)
}

func (fh *FavouriteHandler) CreateHandler() echo.HandlerFunc {
	type Request struct {
		ContentID uint64 `json:"content_id" validate:"required"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if customErr := reader.NewRequestReader(cntx).Read(req); customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		_, customErr := fh.contentUseCase.GetByID(req.ContentID)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		userID := cntx.Get("userID").(uint64)

		newFavourite := &models.Favourite{
			UserID:    userID,
			ContentID: req.ContentID,
			Created:   time.Now(),
		}

		customErr = fh.favouriteUseCase.Create(newFavourite)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"favourite": newFavourite,
			},
		})
	}
}

func (fh *FavouriteHandler) DeleteHandler() echo.HandlerFunc {
	type Request struct {
		ContentID uint64 `json:"content_id"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		userID := cntx.Get("userID").(uint64)

		newFavourite := &models.Favourite{
			UserID:    userID,
			ContentID: req.ContentID,
		}

		customErr := fh.favouriteUseCase.Delete(newFavourite)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{Message: "success"})
	}
}

func (fh *FavouriteHandler) GetFavouritesHandler() echo.HandlerFunc {
	type Request struct {
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logrus.Info(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		userID := cntx.Get("userID").(uint64)

		favouriteMovies, customErr := fh.favouriteUseCase.GetUserFavouriteMovies(userID, &req.Pagination)
		if customErr != nil {
			logrus.Info(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"favourites": favouriteMovies,
			},
		})
	}
}
