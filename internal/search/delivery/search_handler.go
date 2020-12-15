package delivery

import (
	"net/http"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/mwares"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/search"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	reader "github.com/go-park-mail-ru/2020_2_Slash/tools/request_reader"
	. "github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/labstack/echo/v4"
)

type SearchHandler struct {
	searchUsecase search.SearchUsecase
}

func NewSearchHandler(usecase search.SearchUsecase) *SearchHandler {
	return &SearchHandler{searchUsecase: usecase}
}

func (sh *SearchHandler) Configure(e *echo.Echo, mw *mwares.MiddlewareManager) {
	e.GET("/api/v1/search", sh.SearchHandler(), mw.GetAuth)
}

func (sh *SearchHandler) SearchHandler() echo.HandlerFunc {
	type Request struct {
		Query string `query:"q" validate:"required"`
		models.Pagination
	}

	return func(cntx echo.Context) error {
		req := &Request{}
		if err := reader.NewRequestReader(cntx).Read(req); err != nil {
			logger.Error(err.Message)
			return cntx.JSON(err.HTTPCode, Response{Error: err})
		}

		userID, _ := cntx.Get("userID").(uint64)
		result, customErr := sh.searchUsecase.Search(userID, req.Query, &req.Pagination)
		if customErr != nil {
			logger.Error(customErr.Message)
			return cntx.JSON(customErr.HTTPCode, Response{Error: customErr})
		}

		return cntx.JSON(http.StatusOK, Response{
			Body: &Body{
				"result": result,
			},
		})
	}
}
