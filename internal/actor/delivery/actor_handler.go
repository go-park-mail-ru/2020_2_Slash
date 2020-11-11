package delivery

import (
	"github.com/go-park-mail-ru/2020_2_Slash/internal/actor"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/consts"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/helpers/errors"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	actorUseCase actor.ActorUseCase
}

func NewActorHandler(actorUseCase actor.ActorUseCase) *ActorHandler {
	return &ActorHandler{
		actorUseCase: actorUseCase,
	}
}

func (ah *ActorHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.POST("/api/v1/actors", ah.CreateActorHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.PUT("/api/v1/actors/:id", ah.ChangeActorHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.GET("/api/v1/actors/:id", ah.GetActorHandler(), mw.CheckAuth, mw.CheckAdmin)
	e.DELETE("/api/v1/actors/:id", ah.DeleteActorHandler(), mw.CheckAuth, mw.CheckAdmin)
}

func (ah *ActorHandler) CreateActorHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		actor := &models.Actor{
			Name: req.Name,
		}
		err := ah.actorUseCase.Create(actor)
		if err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		return cntx.JSON(http.StatusCreated, Response{
			Body: &Body{
				"actor": actor,
			},
		})
	}
}

func (ah *ActorHandler) ChangeActorHandler() echo.HandlerFunc {
	type Request struct {
		Name string `json:"name" validate:"required"`
	}

	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		req := &Request{}
		if customErr := reader.NewRequestReader(cntx).Read(req); customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		actor := &models.Actor{
			ID:   id,
			Name: req.Name,
		}
		customErr := ah.actorUseCase.Change(actor)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"actor": actor,
			},
		})
	}
}

func (ah *ActorHandler) GetActorHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		actor, customErr := ah.actorUseCase.Get(id)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"actor": actor,
			},
		})
	}
}

func (ah *ActorHandler) DeleteActorHandler() echo.HandlerFunc {
	return func(cntx echo.Context) error {
		id, err := strconv.ParseUint(cntx.Param("id"), 10, 64)
		if err != nil {
			customErr := errors.New(consts.CodeBadRequest, err)
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		customErr := ah.actorUseCase.DeleteById(id)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Message: "success",
		})
	}
}
