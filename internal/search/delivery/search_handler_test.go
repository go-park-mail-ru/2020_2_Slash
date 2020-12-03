package delivery

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/search"
	"github.com/go-park-mail-ru/2020_2_Slash/internal/search/mocks"
	"github.com/go-park-mail-ru/2020_2_Slash/pkg/converter"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/logger"
	"github.com/go-park-mail-ru/2020_2_Slash/tools/response"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type Request struct {
	Query string `query:"q"`
	models.Pagination
}

func setupSearchHandler(searchUseCase search.SearchUsecase, httpMethod string,
	stringifiedJSON string, searchID uint64) (
	echo.Context, *SearchHandler, *httptest.ResponseRecorder) {
	e := echo.New()
	strID := ""
	if searchID != 0 {
		strID = "/" + strconv.FormatUint(searchID, 10)
	}
	req := httptest.NewRequest(httpMethod, "/api/v1/search"+strID,
		strings.NewReader(stringifiedJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	searchHandler := NewSearchHandler(searchUseCase)
	searchHandler.Configure(e, nil)
	return c, searchHandler, rec
}

func TestSearchHandler_CreateHandler(t *testing.T) {
	// Setup
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	searchUseCase := mocks.NewMockSearchUsecase(ctrl)
	logger.DisableLogger()

	testRequest := &Request{
		Query:      "s",
		Pagination: models.Pagination{},
	}
	var curUserID uint64 = 0

	searchJSON, err := converter.AnyToBytesBuffer(testRequest)
	if err != nil {
		t.Fatal(err)
	}
	c, searchHandler, rec := setupSearchHandler(searchUseCase,
		http.MethodPost, searchJSON.String(), 0)
	handleFunc := searchHandler.SearchHandler()

	result := &models.SearchResult{}

	searchUseCase.
		EXPECT().
		Search(curUserID, testRequest.Query, &testRequest.Pagination).
		Return(result, nil)

	response := &response.Response{Body: &response.Body{"result": result}}

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
