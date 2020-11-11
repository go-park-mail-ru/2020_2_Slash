package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/rating"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type RatingHandler struct {
	ratingUseCase rating.RatingUsecase
}

func NewRatingHandler(ratingUseCase rating.RatingUsecase) *RatingHandler {
	return &RatingHandler{
		ratingUseCase: ratingUseCase,
	}
}

func (rh *RatingHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/rating/:cid", rh.CreateHandler(), mw.CheckAuth)
	e.PUT("/api/v1/rating/:cid", rh.ChangeHandler(), mw.CheckAuth)
	e.GET("api/v1/rating/:cid", rh.GetHandler(), mw.CheckAuth)
	e.DELETE("/api/v1/rating/:cid", rh.DeleteHandler(), mw.CheckAuth)
	e.GET("/api/v1/movies/:cid/rating", rh.GetContentRatingHandler())
}

func (rh *RatingHandler) CreateHandler() echo.HandlerFunc {
	type Request struct {
		Likes bool `json:"likes"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		contentID, err := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		userID := cntx.Get("userID").(uint64)
		rating := &models.Rating{
			UserID:    userID,
			ContentID: contentID,
			Likes:     req.Likes,
		}

		customErr := rh.ratingUseCase.Create(rating)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"rating": rating,
			},
		})
	}
}

func (rh *RatingHandler) DeleteHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {

		userID := cntx.Get("userID").(uint64)
		contentID, err := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		rating := &models.Rating{
			UserID:    userID,
			ContentID: contentID,
		}

		customErr := rh.ratingUseCase.Delete(rating)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{Message: "success"})
	}
}

func (rh *RatingHandler) ChangeHandler() echo.HandlerFunc {
	type Request struct {
		Likes bool `json:"likes"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		userID := cntx.Get("userID").(uint64)
		contentID, err := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}
		rating := &models.Rating{
			UserID:    userID,
			ContentID: contentID,
			Likes:     req.Likes,
		}

		customErr := rh.ratingUseCase.Change(rating)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"rating": rating,
			},
		})
	}
}

func (rh *RatingHandler) GetHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		userID := cntx.Get("userID").(uint64)
		contentID, err := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		rating, customErr := rh.ratingUseCase.GetByUserIDContentID(userID, contentID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"rating": rating,
			},
		})
	}
}

func (rh *RatingHandler) GetContentRatingHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		contentID, err := strconv.ParseUint(cntx.Param("cid"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		match, customErr := rh.ratingUseCase.GetContentRating(contentID)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"match": match,
			},
		})
	}
}
